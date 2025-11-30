// Package config provides structures for building configurations.
//
// The Config type is the fundamental building block used throughout Substation
// to configure transforms, conditions, and other components. It follows a
// consistent pattern where Type identifies the component type and Settings
// contains type-specific configuration options.
//
// # Usage Pattern
//
// All Substation components use the same configuration pattern:
//
//	cfg := config.Config{
//		Type: "transform_type",
//		Settings: map[string]interface{}{
//			"option1": "value1",
//			"option2": 42,
//		},
//	}
//
// This pattern allows configurations to be easily serialized to and from JSON,
// making them suitable for use with Jsonnet configuration files.
package config

import (
	"encoding/json"
)

// Config is a template used by Substation interface factories to produce new
// instances. Type refers to the type of instance (e.g., "object_copy", "string_match")
// and Settings contains options used to configure the instance.
//
// Examples of Config usage can be found in the condition and transform packages.
//
// When used with Jsonnet, configurations are typically generated using the
// helper functions in substation.libsonnet.
type Config struct {
	// Type identifies the component type to be created. This corresponds
	// to the registered name in the relevant factory function.
	Type string `json:"type"`
	// Settings contains type-specific configuration options as key-value pairs.
	// The expected keys and value types depend on the component Type.
	Settings map[string]interface{} `json:"settings"`
}

// String returns a JSON representation of the Config. This is useful for
// logging and debugging configurations.
func (c Config) String() string {
	b, _ := json.Marshal(c)
	return string(b)
}
