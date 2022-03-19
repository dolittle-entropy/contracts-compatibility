package main

import "github.com/coreos/go-semver/semver"

// Compatibility represents a structure containing compatible versions by Runtime and SDKs versions
type Compatibility struct {
	Runtime map[string]map[string]semver.Versions `json:"runtime"`
	SDKs    map[string]map[string]semver.Versions `json:"sdk"`
}
