// Package substation provides a toolkit for routing, normalizing, and enriching
// security event and audit logs.
//
// Substation is designed to be extensible and can be used to build data pipeline
// systems and microservices using out-of-the-box applications and 100+ data
// transformation functions, or custom transforms written in Go.
//
// # Core Concepts
//
// The package is built around three main concepts:
//
//   - Messages: The data structure that flows through the pipeline. Messages can
//     contain JSON text or binary data, along with metadata.
//   - Transforms: Functions that process messages, such as copying fields,
//     formatting data, or sending data to external services.
//   - Conditions: Functions that evaluate messages and return true or false,
//     used for conditional routing and filtering.
//
// # Basic Usage
//
// To use Substation, create a configuration with transforms and instantiate
// a new Substation instance:
//
//	ctx := context.Background()
//	cfg := substation.Config{
//		Transforms: []config.Config{
//			{Type: "object_copy", Settings: map[string]interface{}{
//				"object": map[string]interface{}{
//					"source_key": "a",
//					"target_key": "b",
//				},
//			}},
//		},
//	}
//
//	sub, err := substation.New(ctx, cfg)
//	if err != nil {
//		// Handle error
//	}
//
//	msg := message.New().SetData([]byte(`{"a":"hello"}`))
//	results, err := sub.Transform(ctx, msg)
//
// # Configuration
//
// Substation configurations are typically written in Jsonnet using the provided
// library (substation.libsonnet), which provides a fluent API for building
// configurations. The Library variable contains the embedded Jsonnet library.
//
// # Custom Transforms
//
// Custom transforms can be implemented by creating types that satisfy the
// transform.Transformer interface and using WithTransformFactory to provide
// a custom factory function.
package substation

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/brexhq/substation/v2/config"
	"github.com/brexhq/substation/v2/message"
	"github.com/brexhq/substation/v2/transform"
)

// Library contains the embedded Jsonnet library (substation.libsonnet) that
// provides helper functions and default configurations for building Substation
// configurations. This can be used with Jsonnet tools to compile configurations.
//
//go:embed substation.libsonnet
var Library string

var errNoTransforms = fmt.Errorf("no transforms configured")

// Config is the core configuration for the application. Custom applications
// should embed this struct and add additional configuration options.
//
// Example of embedding Config in a custom configuration:
//
//	type MyConfig struct {
//		substation.Config
//		APIKey string `json:"api_key"`
//	}
type Config struct {
	// Transforms contains a list of data transformations that are executed
	// in sequence. Each transform processes messages and passes the results
	// to the next transform in the pipeline.
	Transforms []config.Config `json:"transforms"`
}

// Substation provides access to data transformation functions. It holds
// the configuration and compiled transforms, and provides methods for
// processing messages through the transformation pipeline.
//
// Substation instances are safe for concurrent use after creation.
type Substation struct {
	cfg Config

	factory transform.Factory
	tforms  []transform.Transformer
}

// New returns a new Substation instance configured with the provided Config.
// The context is used during transform initialization and may be used by
// transforms that require external resources.
//
// Options can be provided to customize the Substation instance, such as
// using a custom transform factory with WithTransformFactory.
//
// Returns an error if no transforms are configured or if any transform
// fails to initialize.
func New(ctx context.Context, cfg Config, opts ...func(*Substation)) (*Substation, error) {
	if cfg.Transforms == nil {
		return nil, errNoTransforms
	}

	sub := &Substation{
		cfg:     cfg,
		factory: transform.New,
	}

	for _, o := range opts {
		o(sub)
	}

	// Create transforms from the configuration.
	for _, c := range cfg.Transforms {
		t, err := sub.factory(ctx, c)
		if err != nil {
			return nil, err
		}

		sub.tforms = append(sub.tforms, t)
	}

	return sub, nil
}

// WithTransformFactory returns an option that configures a Substation instance
// to use a custom transform factory. This allows applications to add custom
// transforms or override the default transform implementations.
//
// The custom factory should handle unknown transform types by delegating to
// the default transform.New function:
//
//	func customFactory(ctx context.Context, cfg config.Config) (transform.Transformer, error) {
//		switch cfg.Type {
//		case "my_custom_transform":
//			return &myCustomTransform{}, nil
//		default:
//			return transform.New(ctx, cfg)
//		}
//	}
func WithTransformFactory(fac transform.Factory) func(*Substation) {
	return func(s *Substation) {
		s.factory = fac
	}
}

// Transform runs the configured data transformation functions on the
// provided messages.
//
// This is safe to use concurrently.
func (s *Substation) Transform(ctx context.Context, msg ...*message.Message) ([]*message.Message, error) {
	return transform.Apply(ctx, s.tforms, msg...)
}

// String returns a JSON representation of the configuration.
func (s *Substation) String() string {
	b, err := json.Marshal(s.cfg)
	if err != nil {
		return fmt.Sprintf("substation: %v", err)
	}

	return string(b)
}
