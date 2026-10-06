package main

import (
	"os"
	"testing"
)

func TestGenCoverageMain(t *testing.T) {
	main()
	defer func() { _ = os.Remove("coverage.out") }()
}

func TestParseCoverageProfile(t *testing.T) {
	sample := `mode: set
github.com/abmarcum/multi-cloud-provider/main.go:10.15,14.2 4 1
github.com/abmarcum/multi-cloud-provider/internal/provider/provider.go:20.1,30.2 6 0
`
	total, covered := parseCoverageProfile(sample)
	if total != 10 || covered != 4 {
		t.Errorf("parseCoverageProfile returned total=%d covered=%d; want total=10 covered=4", total, covered)
	}
}
