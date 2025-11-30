// Package condition provides functions for evaluating messages using success
// or failure criteria.
//
// Conditions are used to make decisions about message processing, such as
// filtering messages, routing them to different transforms, or triggering
// specific behaviors based on message content.
//
// # Condition Categories
//
// Conditions are organized into the following categories:
//
//   - Meta: Combine multiple conditions (e.g., all, any, none)
//   - Format: Check data formats (e.g., format_json, format_mime)
//   - Network: Validate IP addresses (e.g., network_ip_private, network_ip_valid)
//   - Number: Numeric comparisons (e.g., number_equal_to, number_greater_than)
//   - String: String comparisons (e.g., string_equal_to, string_contains, string_match)
//   - Utility: Helper conditions (e.g., utility_random)
//
// # Using Conditions
//
// Conditions are typically created using the New factory function with a
// configuration that specifies the condition type and settings:
//
//	cfg := config.Config{
//		Type: "string_equal_to",
//		Settings: map[string]interface{}{
//			"object": map[string]interface{}{
//				"source_key": "status",
//			},
//			"value": "active",
//		},
//	}
//	cond, err := condition.New(ctx, cfg)
//
//	// Evaluate the condition
//	ok, err := cond.Condition(ctx, msg)
//	if ok {
//		// Condition matched
//	}
//
// # Meta Conditions
//
// Meta conditions combine multiple conditions using logical operators:
//   - all (meta_all): All conditions must be true
//   - any (meta_any): At least one condition must be true
//   - none (meta_none): No conditions can be true
package condition

import (
	"context"
	"fmt"

	"github.com/brexhq/substation/v2/config"
	"github.com/brexhq/substation/v2/message"

	iconfig "github.com/brexhq/substation/v2/internal/config"
)

// Conditioner is the interface implemented by all conditions. It provides
// the ability to evaluate a message and return a boolean result.
//
// The Condition method receives a context and a message, and returns true
// if the condition is satisfied, false otherwise. Control messages typically
// return false.
type Conditioner interface {
	Condition(context.Context, *message.Message) (bool, error)
}

// New is a factory function for returning a configured Conditioner based on
// the Type field in the configuration. The Settings field contains type-specific
// configuration options.
//
// Returns an error if the condition type is not recognized or if the
// configuration is invalid for the specified type.
func New(ctx context.Context, cfg config.Config) (Conditioner, error) { //nolint: cyclop, gocyclo // ignore cyclomatic complexity
	switch cfg.Type {
	// Meta inspectors.
	case "all", "meta_all":
		return newMetaAll(ctx, cfg)
	case "any", "meta_any":
		return newMetaAny(ctx, cfg)
	case "none", "meta_none":
		return newMetaNone(ctx, cfg)
	// Format inspectors.
	case "format_mime":
		return newFormatMIME(ctx, cfg)
	case "format_json":
		return newFormatJSON(ctx, cfg)
	// Network inspectors.
	case "network_ip_global_unicast":
		return newNetworkIPGlobalUnicast(ctx, cfg)
	case "network_ip_link_local_multicast":
		return newNetworkIPLinkLocalMulticast(ctx, cfg)
	case "network_ip_link_local_unicast":
		return newNetworkIPLinkLocalUnicast(ctx, cfg)
	case "network_ip_loopback":
		return newNetworkIPLoopback(ctx, cfg)
	case "network_ip_multicast":
		return newNetworkIPMulticast(ctx, cfg)
	case "network_ip_private":
		return newNetworkIPPrivate(ctx, cfg)
	case "network_ip_unicast":
		return newNetworkIPUnicast(ctx, cfg)
	case "network_ip_unspecified":
		return newNetworkIPUnspecified(ctx, cfg)
	case "network_ip_valid":
		return newNetworkIPValid(ctx, cfg)
	// Number inspectors.
	case "number_equal_to":
		return newNumberEqualTo(ctx, cfg)
	case "number_less_than":
		return newNumberLessThan(ctx, cfg)
	case "number_greater_than":
		return newNumberGreaterThan(ctx, cfg)
	case "number_bitwise_and":
		return newNumberBitwiseAND(ctx, cfg)
	case "number_bitwise_or":
		return newNumberBitwiseOR(ctx, cfg)
	case "number_bitwise_xor":
		return newNumberBitwiseXOR(ctx, cfg)
	case "number_bitwise_not":
		return newNumberBitwiseNOT(ctx, cfg)
	case "number_length_less_than":
		return newNumberLengthLessThan(ctx, cfg)
	case "number_length_greater_than":
		return newNumberLengthGreaterThan(ctx, cfg)
	case "number_length_equal_to":
		return newNumberLengthEqualTo(ctx, cfg)
	// String inspectors.
	case "string_contains":
		return newStringContains(ctx, cfg)
	case "string_ends_with":
		return newStringEndsWith(ctx, cfg)
	case "string_equal_to":
		return newStringEqualTo(ctx, cfg)
	case "string_greater_than":
		return newStringGreaterThan(ctx, cfg)
	case "string_less_than":
		return newStringLessThan(ctx, cfg)
	case "string_starts_with":
		return newStringStartsWith(ctx, cfg)
	case "string_match":
		return newStringMatch(ctx, cfg)
	// Utility inspectors.
	case "utility_random":
		return newUtilityRandom(ctx, cfg)
	default:
		return nil, fmt.Errorf("condition %s: %w", cfg.Type, iconfig.ErrInvalidFactoryInput)
	}
}
