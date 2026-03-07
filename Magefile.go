//go:build mage
// +build mage

package main

import (
	"os"

	"github.com/magefile/mage/sh"

	// mage:import
	build "github.com/grafana/grafana-plugin-sdk-go/build"
)

// Default configures the default target.
var Default = build.BuildAll

func TestIntegration() error {
	timeout := "120s"
	if t := os.Getenv("TEST_INTEGRATION_TIMEOUT"); t != "" {
		timeout = t
	}
	return sh.RunV("go", "test",
		"-tags", "integration",
		"-v",
		"-timeout", timeout,
		"-run", "TestIntegration",
		"./pkg/plugin/...",
	)
}
