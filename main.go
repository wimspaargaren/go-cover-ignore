// Package main provides a simple tool removing ignored files from the code coverage report.
package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/caarlos0/env/v9"
	"gopkg.in/yaml.v3"
)

// Config represents the configuration of the tool.
type Config struct {
	IgnoreSpecPath      string `env:"GO_COVER_IGNORE_SPEC_PATH" envDefault:".coverage-ignore.yaml"`
	CoverageProfilePath string `env:"GO_COVER_IGNORE_COVER_PROFILE_PATH" envDefault:"cover.out"`
}

type ignoreSpec struct {
	Module      string   `yaml:"module"`
	IgnoreRules []string `yaml:"ignore_rules"`
}

func newConfig() (*Config, error) {
	cfg := Config{}
	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func main() {
	config, err := newConfig()
	if err != nil {
		log.Fatalf("unable to load go cover ignore config: %v", err)
	}
	ignoreData, err := os.ReadFile(config.IgnoreSpecPath)
	if err != nil {
		log.Fatalf("unable to read ignore spec yaml: %v at path: %s", err, config.IgnoreSpecPath)
	}
	t := ignoreSpec{}
	err = yaml.Unmarshal(ignoreData, &t)
	if err != nil {
		log.Fatalf("unable to read ignore spec yaml: %v", err)
	}
	covData, err := os.ReadFile(config.CoverageProfilePath)
	if err != nil {
		log.Fatalf("unable to read coverage profile output file: %v at path: %s", err, config.CoverageProfilePath)
	}
	result, err := filterCoverage(t, covData)
	if err != nil {
		log.Fatalf("unable to apply filter rules for ignoring coverage: %v", err)
	}
	err = os.WriteFile(config.CoverageProfilePath, []byte(strings.Join(result, "\n")), 0o644)
	if err != nil {
		log.Fatalf("unable to safe filtered coverage profile: %v", err)
	}
}

// filterCoverage takes in the original coverage report, applies
// the specified filter rules/regexes and returns the resulting lines.
func filterCoverage(t ignoreSpec, covData []byte) ([]string, error) {
	regexMatches := []*regexp.Regexp{}

	for _, ignoreRule := range t.IgnoreRules {
		regularExpression, err := regexp.Compile(ignoreRule)
		if err != nil {
			return nil, fmt.Errorf("unable to parse regular expression: %s, %w", ignoreRule, err)
		}
		regexMatches = append(regexMatches, regularExpression)
	}

	covDataLines := strings.Split(string(covData), "\n")
	result := []string{}
	for _, covDataLine := range covDataLines {
		matched := false
		for _, regexMatch := range regexMatches {
			if regexMatch.MatchString(covDataLine) {
				matched = true
				break
			}
		}
		if !matched {
			result = append(result, covDataLine)
		}
	}
	return result, nil
}
