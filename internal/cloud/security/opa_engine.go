package security

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

// OPAPolicyResult represents Open Policy Agent Rego evaluation results
type OPAPolicyResult struct {
	ResourceName string
	ProviderType string
	PolicyRule   string
	Passed       bool
	Violation    string
}

var (
	regoDenyRuleRegex = regexp.MustCompile(`(?s)(?:deny|violation)\s*\[\s*([a-zA-Z0-9_"]+)\s*\]\s*(?:if\s*)?\{\s*(.*?)\s*\}`)
	regoAssignMsgRe   = regexp.MustCompile(`^\s*[a-zA-Z0-9_]+\s*:=\s*"([^"]+)"\s*$`)
	regoEqRe          = regexp.MustCompile(`^\s*input\.([a-zA-Z0-9_.]+)\s*==\s*(.+?)\s*$`)
	regoNeqRe         = regexp.MustCompile(`^\s*input\.([a-zA-Z0-9_.]+)\s*!=\s*(.+?)\s*$`)
	regoNotRe         = regexp.MustCompile(`^\s*not\s+input\.([a-zA-Z0-9_.]+)\s*$`)
	regoCountZeroRe   = regexp.MustCompile(`^\s*count\(\s*input\.([a-zA-Z0-9_.]+)\s*\)\s*==\s*0\s*$`)
)

type regoCondOp uint8

const (
	regoOpInvalid regoCondOp = iota
	regoOpCountZero
	regoOpNot
	regoOpEq
	regoOpNeq
)

type compiledRegoCond struct {
	op   regoCondOp
	path []string
	rhs  string
}

type compiledRegoRule struct {
	violationMsg string
	conditions   []compiledRegoCond
}

type cachedPolicyFile struct {
	modTime time.Time
	size    int64
	content string
}

var (
	regoModuleCache sync.Map // map[string][]compiledRegoRule
	policyFileCache sync.Map // map[string]cachedPolicyFile
)

func loadCachedPolicyFile(rawPath string) (string, error) {
	cleanPath := filepath.Clean(rawPath)
	/* #nosec G304 G703 */
	info, err := os.Stat(cleanPath)
	if err != nil {
		return "", err
	}
	if cached, ok := policyFileCache.Load(cleanPath); ok {
		entry := cached.(cachedPolicyFile)
		if entry.modTime.Equal(info.ModTime()) && entry.size == info.Size() {
			return entry.content, nil
		}
	}
	/* #nosec G304 G703 */
	data, err := os.ReadFile(cleanPath)
	if err != nil {
		return "", err
	}
	content := string(data)
	policyFileCache.Store(cleanPath, cachedPolicyFile{
		modTime: info.ModTime(),
		size:    info.Size(),
		content: content,
	})
	return content, nil
}

// EvaluateOPARegoPolicy evaluates unified resource attributes against Open Policy Agent (OPA) Rego rules
func EvaluateOPARegoPolicy(providerType string, resourceType string, resourceName string, regoRule string) OPAPolicyResult {
	return EvaluateOPARegoPolicyWithAttrs(providerType, resourceType, resourceName, regoRule, nil)
}

// EvaluateOPARegoPolicyWithAttrs evaluates unified resource attributes against Open Policy Agent (OPA) Rego rules or inline Rego modules
func EvaluateOPARegoPolicyWithAttrs(providerType string, resourceType string, resourceName string, regoRule string, attributes map[string]interface{}) OPAPolicyResult {
	p := strings.ToUpper(providerType)

	// Support inline Rego module or external .rego policy file
	if strings.Contains(regoRule, "package ") || strings.Contains(regoRule, "deny[") || strings.Contains(regoRule, "violation[") {
		return EvaluateRegoModule(providerType, resourceType, resourceName, regoRule, attributes)
	}
	if policyPath := os.Getenv("OPA_POLICY_PATH"); policyPath != "" {
		content, err := loadCachedPolicyFile(policyPath)
		if err != nil || len(strings.TrimSpace(content)) == 0 {
			return OPAPolicyResult{
				ResourceName: resourceName,
				ProviderType: providerType,
				PolicyRule:   "opa_policy_path",
				Passed:       false,
				Violation:    fmt.Sprintf("[%s OPA Engine] Failed to load required Rego policy from OPA_POLICY_PATH", p),
			}
		}
		res := EvaluateRegoModule(providerType, resourceType, resourceName, content, attributes)
		if !res.Passed {
			return res
		}
	}

	switch regoRule {
	case "must_have_tags":
		if attributes != nil {
			if tags, exists := attributes["tags"]; exists {
				switch t := tags.(type) {
				case map[string]string:
					if len(t) == 0 {
						return OPAPolicyResult{
							ResourceName: resourceName,
							ProviderType: providerType,
							PolicyRule:   regoRule,
							Passed:       false,
							Violation:    fmt.Sprintf("[%s OPA Engine] Rego Policy Violation ('%s'): Resource '%s' has an empty tags map.", p, regoRule, resourceName),
						}
					}
				case map[string]interface{}:
					if len(t) == 0 {
						return OPAPolicyResult{
							ResourceName: resourceName,
							ProviderType: providerType,
							PolicyRule:   regoRule,
							Passed:       false,
							Violation:    fmt.Sprintf("[%s OPA Engine] Rego Policy Violation ('%s'): Resource '%s' has an empty tags map.", p, regoRule, resourceName),
						}
					}
				case nil:
					return OPAPolicyResult{
						ResourceName: resourceName,
						ProviderType: providerType,
						PolicyRule:   regoRule,
						Passed:       false,
						Violation:    fmt.Sprintf("[%s OPA Engine] Rego Policy Violation ('%s'): Resource '%s' is missing required tags.", p, regoRule, resourceName),
					}
				}
			}
		}

	case "disallow_public_ip":
		if strings.Contains(resourceType, "virtual_machine") {
			if attributes != nil {
				if pub, ok := attributes["associate_public_ip"].(bool); ok && !pub {
					break
				}
				if pub, ok := attributes["is_public"].(bool); ok && !pub {
					break
				}
			}
			return OPAPolicyResult{
				ResourceName: resourceName,
				ProviderType: providerType,
				PolicyRule:   regoRule,
				Passed:       false,
				Violation:    fmt.Sprintf("[%s OPA Engine] Rego Policy Violation ('%s'): Virtual Machine '%s' cannot allocate a public IP in production.", p, regoRule, resourceName),
			}
		}

	case "require_encryption":
		if attributes != nil {
			if enc, ok := attributes["encryption_enabled"].(bool); ok && !enc {
				return OPAPolicyResult{
					ResourceName: resourceName,
					ProviderType: providerType,
					PolicyRule:   regoRule,
					Passed:       false,
					Violation:    fmt.Sprintf("[%s OPA Engine] Rego Policy Violation ('%s'): Resource '%s' must enable encryption at rest.", p, regoRule, resourceName),
				}
			}
		}

	case "require_multi_az":
		if strings.Contains(resourceType, "db_instance") && attributes != nil {
			if multiAZ, ok := attributes["multi_az"].(bool); ok && !multiAZ {
				return OPAPolicyResult{
					ResourceName: resourceName,
					ProviderType: providerType,
					PolicyRule:   regoRule,
					Passed:       false,
					Violation:    fmt.Sprintf("[%s OPA Engine] Rego Policy Violation ('%s'): Database '%s' must enable Multi-AZ high availability.", p, regoRule, resourceName),
				}
			}
		}
	}

	return OPAPolicyResult{
		ResourceName: resourceName,
		ProviderType: providerType,
		PolicyRule:   regoRule,
		Passed:       true,
		Violation:    "",
	}
}

func compileRegoCondition(expr string) compiledRegoCond {
	if m := regoCountZeroRe.FindStringSubmatch(expr); len(m) == 2 {
		return compiledRegoCond{op: regoOpCountZero, path: strings.Split(m[1], ".")}
	}
	if m := regoNotRe.FindStringSubmatch(expr); len(m) == 2 {
		return compiledRegoCond{op: regoOpNot, path: strings.Split(m[1], ".")}
	}
	if m := regoEqRe.FindStringSubmatch(expr); len(m) == 3 {
		return compiledRegoCond{op: regoOpEq, path: strings.Split(m[1], "."), rhs: strings.Trim(m[2], `"`)}
	}
	if m := regoNeqRe.FindStringSubmatch(expr); len(m) == 3 {
		return compiledRegoCond{op: regoOpNeq, path: strings.Split(m[1], "."), rhs: strings.Trim(m[2], `"`)}
	}
	return compiledRegoCond{op: regoOpInvalid}
}

func compileRegoModule(regoModule string) []compiledRegoRule {
	if cached, ok := regoModuleCache.Load(regoModule); ok {
		return cached.([]compiledRegoRule)
	}

	matches := regoDenyRuleRegex.FindAllStringSubmatch(regoModule, -1)
	rules := make([]compiledRegoRule, 0, len(matches))
	for _, m := range matches {
		violationMsg := strings.Trim(m[1], `"`)
		lines := strings.Split(m[2], "\n")
		conds := make([]compiledRegoCond, 0, len(lines))
		for _, rawLine := range lines {
			line := strings.TrimSpace(rawLine)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			if assign := regoAssignMsgRe.FindStringSubmatch(line); len(assign) == 2 {
				violationMsg = assign[1]
				continue
			}
			conds = append(conds, compileRegoCondition(line))
		}
		rules = append(rules, compiledRegoRule{
			violationMsg: violationMsg,
			conditions:   conds,
		})
	}

	regoModuleCache.Store(regoModule, rules)
	return rules
}

// EvaluateRegoModule parses and evaluates Rego deny/violation rules against a resource input document
func EvaluateRegoModule(providerType, resourceType, resourceName, regoModule string, attributes map[string]interface{}) OPAPolicyResult {
	input := map[string]interface{}{
		"provider":      strings.ToLower(providerType),
		"resource_type": resourceType,
		"resource_name": resourceName,
		"attributes":    attributes,
	}

	rules := compileRegoModule(regoModule)
	for _, rule := range rules {
		allTrue := true
		for _, cond := range rule.conditions {
			if !evalCompiledCondition(cond, input) {
				allTrue = false
				break
			}
		}
		if allTrue {
			return OPAPolicyResult{
				ResourceName: resourceName,
				ProviderType: providerType,
				PolicyRule:   "rego_module",
				Passed:       false,
				Violation:    fmt.Sprintf("[%s OPA Engine] %s", strings.ToUpper(providerType), rule.violationMsg),
			}
		}
	}

	return OPAPolicyResult{
		ResourceName: resourceName,
		ProviderType: providerType,
		PolicyRule:   "rego_module",
		Passed:       true,
	}
}

func resolveInputSegments(input map[string]interface{}, parts []string) (interface{}, bool) {
	var curr interface{} = input
	for _, p := range parts {
		m, ok := curr.(map[string]interface{})
		if !ok || m == nil {
			return nil, false
		}
		curr, ok = m[p]
		if !ok {
			if attrs, hasAttrs := m["attributes"].(map[string]interface{}); hasAttrs {
				if val, found := attrs[p]; found {
					curr = val
					continue
				}
			}
			return nil, false
		}
	}
	return curr, true
}

func evalCompiledCondition(cond compiledRegoCond, input map[string]interface{}) bool {
	switch cond.op {
	case regoOpCountZero:
		val, ok := resolveInputSegments(input, cond.path)
		if !ok || val == nil {
			return true
		}
		switch v := val.(type) {
		case map[string]string:
			return len(v) == 0
		case map[string]interface{}:
			return len(v) == 0
		case []interface{}:
			return len(v) == 0
		}
		return false

	case regoOpNot:
		val, ok := resolveInputSegments(input, cond.path)
		if !ok || val == nil {
			return true
		}
		if b, ok := val.(bool); ok {
			return !b
		}
		return false

	case regoOpEq:
		val, ok := resolveInputSegments(input, cond.path)
		if !ok {
			return false
		}
		return fmt.Sprintf("%v", val) == cond.rhs

	case regoOpNeq:
		val, ok := resolveInputSegments(input, cond.path)
		if !ok {
			return true
		}
		return fmt.Sprintf("%v", val) != cond.rhs

	default:
		return false
	}
}

