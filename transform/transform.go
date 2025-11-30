// Package transform provides functions for transforming messages.
//
// Transforms are the core building blocks of Substation pipelines. Each transform
// implements the Transformer interface and processes messages by reading data,
// modifying it, and returning zero or more result messages.
//
// # Transform Categories
//
// Transforms are organized into the following categories:
//
//   - Aggregation: Combine or split messages (e.g., aggregate_to_array, aggregate_from_array)
//   - Array: Operations on JSON arrays (e.g., array_join, array_zip)
//   - Enrichment: Add data from external sources (e.g., enrich_http_get, enrich_dns_lookup)
//   - Format: Convert data formats (e.g., format_from_base64, format_to_gzip)
//   - Hash: Generate hashes (e.g., hash_md5, hash_sha256)
//   - Meta: Control flow and composition (e.g., meta_switch, meta_for_each)
//   - Network: Domain and IP operations (e.g., network_domain_subdomain)
//   - Number: Numeric operations (e.g., number_math_addition, number_maximum)
//   - Object: JSON object manipulation (e.g., object_copy, object_delete, object_insert)
//   - Send: Output data to destinations (e.g., send_stdout, send_aws_s3)
//   - String: String manipulation (e.g., string_replace, string_to_lower)
//   - Time: Time parsing and formatting (e.g., time_from_string, time_to_unix)
//   - Utility: Helper transforms (e.g., utility_drop, utility_delay)
//
// # Using Transforms
//
// Transforms are typically created using the New factory function with a
// configuration that specifies the transform type and settings:
//
//	cfg := config.Config{
//		Type: "object_copy",
//		Settings: map[string]interface{}{
//			"object": map[string]interface{}{
//				"source_key": "input.field",
//				"target_key": "output.field",
//			},
//		},
//	}
//	tf, err := transform.New(ctx, cfg)
//
// # Custom Transforms
//
// Custom transforms can be created by implementing the Transformer interface.
// The Transform method receives a message and returns zero or more result
// messages. Control messages (identified by IsControl()) should typically be
// passed through unchanged:
//
//	func (t *MyTransform) Transform(ctx context.Context, msg *message.Message) ([]*message.Message, error) {
//		if msg.IsControl() {
//			return []*message.Message{msg}, nil
//		}
//		// Process the message...
//		return []*message.Message{msg}, nil
//	}
package transform

import (
	"context"
	"fmt"
	"math"

	"github.com/brexhq/substation/v2/config"
	"github.com/brexhq/substation/v2/message"

	iconfig "github.com/brexhq/substation/v2/internal/config"
)

var errMsgInvalidObject = fmt.Errorf("message must be JSON object")

// Transformer is the interface implemented by all transforms and
// provides the ability to transform a message.
//
// The Transform method receives a context and a message, and returns
// a slice of messages and an error. Transforms may return:
//   - Zero messages: The message was dropped/filtered
//   - One message: The message was modified in place
//   - Multiple messages: The message was split into multiple messages
//
// Transforms should handle control messages appropriately, typically by
// passing them through unchanged or by using them to flush internal state.
type Transformer interface {
	Transform(context.Context, *message.Message) ([]*message.Message, error)
}

// Factory is a function type that creates Transformer instances from
// configuration. It can be used to implement custom transform factory
// functions that extend or replace the default transform types.
//
// Custom factories should typically delegate to the default New function
// for unknown transform types to maintain compatibility with built-in
// transforms.
type Factory func(context.Context, config.Config) (Transformer, error)

// New is a factory function for returning a configured Transformer based on
// the Type field in the configuration. The Settings field contains type-specific
// configuration options.
//
// Returns an error if the transform type is not recognized or if the
// configuration is invalid for the specified type.
func New(ctx context.Context, cfg config.Config) (Transformer, error) { //nolint: cyclop, gocyclo // ignore cyclomatic complexity
	switch cfg.Type {
	// Aggregation transforms.
	case "aggregate_from_array":
		return newAggregateFromArray(ctx, cfg)
	case "aggregate_to_array":
		return newAggregateToArray(ctx, cfg)
	case "aggregate_from_string":
		return newAggregateFromString(ctx, cfg)
	case "aggregate_to_string":
		return newAggregateToString(ctx, cfg)
	// Array transforms.
	case "array_join":
		return newArrayJoin(ctx, cfg)
	case "array_zip":
		return newArrayZip(ctx, cfg)
	// Enrichment transforms.
	case "enrich_aws_dynamodb_query":
		return newEnrichAWSDynamoDBQuery(ctx, cfg)
	case "enrich_aws_lambda":
		return newEnrichAWSLambda(ctx, cfg)
	case "enrich_dns_ip_lookup":
		return newEnrichDNSIPLookup(ctx, cfg)
	case "enrich_dns_domain_lookup":
		return newEnrichDNSDomainLookup(ctx, cfg)
	case "enrich_dns_text_lookup":
		return newEnrichDNSTxtLookup(ctx, cfg)
	case "enrich_http_get":
		return newEnrichHTTPGet(ctx, cfg)
	case "enrich_http_post":
		return newEnrichHTTPPost(ctx, cfg)
		// Deprecated: Use enrich_kv_store_item_get instead.
	case "enrich_kv_store_get":
		fallthrough
	case "enrich_kv_store_item_get":
		return newEnrichKVStoreItemGet(ctx, cfg)
		// Deprecated: Use enrich_kv_store_item_set instead.
	case "enrich_kv_store_set":
		fallthrough
	case "enrich_kv_store_item_set":
		return newEnrichKVStoreItemSet(ctx, cfg)
	case "enrich_kv_store_set_add":
		return newEnrichKVStoreSetAdd(ctx, cfg)
	// Format transforms.
	case "format_from_base64":
		return newFormatFromBase64(ctx, cfg)
	case "format_to_base64":
		return newFormatToBase64(ctx, cfg)
	case "format_from_gzip":
		return newFormatFromGzip(ctx, cfg)
	case "format_to_gzip":
		return newFormatToGzip(ctx, cfg)
	case "format_from_pretty_print":
		return newFormatFromPrettyPrint(ctx, cfg)
	case "format_from_zip":
		return newFormatFromZip(ctx, cfg)
	// Hash transforms.
	case "hash_md5":
		return newHashMD5(ctx, cfg)
	case "hash_sha256":
		return newHashSHA256(ctx, cfg)
	// Meta transforms.
	case "meta_err":
		return newMetaErr(ctx, cfg)
	case "meta_for_each":
		return newMetaForEach(ctx, cfg)
	case "meta_kv_store_lock":
		return newMetaKVStoreLock(ctx, cfg)
	case "meta_metric_duration":
		return newMetaMetricsDuration(ctx, cfg)
	case "meta_retry":
		return newMetaRetry(ctx, cfg)
	case "meta_switch":
		return newMetaSwitch(ctx, cfg)
	// Number transforms.
	case "number_maximum":
		return newNumberMaximum(ctx, cfg)
	case "number_minimum":
		return newNumberMinimum(ctx, cfg)
	case "number_math_addition":
		return newNumberMathAddition(ctx, cfg)
	case "number_math_division":
		return newNumberMathDivision(ctx, cfg)
	case "number_math_multiplication":
		return newNumberMathMultiplication(ctx, cfg)
	case "number_math_subtraction":
		return newNumberMathSubtraction(ctx, cfg)
	// Network transforms.
	case "network_domain_registered_domain":
		return newNetworkDomainRegisteredDomain(ctx, cfg)
	case "network_domain_subdomain":
		return newNetworkDomainSubdomain(ctx, cfg)
	case "network_domain_top_level_domain":
		return newNetworkDomainTopLevelDomain(ctx, cfg)
	// Object transforms.
	case "object_copy":
		return newObjectCopy(ctx, cfg)
	case "object_delete":
		return newObjectDelete(ctx, cfg)
	case "object_insert":
		return newObjectInsert(ctx, cfg)
	case "object_jq":
		return newObjectJQ(ctx, cfg)
	case "object_to_boolean":
		return newObjectToBoolean(ctx, cfg)
	case "object_to_float":
		return newObjectToFloat(ctx, cfg)
	case "object_to_integer":
		return newObjectToInteger(ctx, cfg)
	case "object_to_string":
		return newObjectToString(ctx, cfg)
	case "object_to_unsigned_integer":
		return newObjectToUnsignedInteger(ctx, cfg)
	// Send transforms.
	case "send_aws_dynamodb_put":
		return newSendAWSDynamoDBPut(ctx, cfg)
	case "send_aws_eventbridge":
		return newSendAWSEventBridge(ctx, cfg)
	case "send_aws_data_firehose":
		return newSendAWSDataFirehose(ctx, cfg)
	case "send_aws_kinesis_data_stream":
		return newSendAWSKinesisDataStream(ctx, cfg)
	case "send_aws_lambda":
		return newSendAWSLambda(ctx, cfg)
	case "send_aws_s3":
		return newSendAWSS3(ctx, cfg)
	case "send_aws_sns":
		return newSendAWSSNS(ctx, cfg)
	case "send_aws_sqs":
		return newSendAWSSQS(ctx, cfg)
	case "send_gcp_storage":
		return newSendGCPStorage(ctx, cfg)
	case "send_file":
		return newSendFile(ctx, cfg)
	case "send_http_post":
		return newSendHTTPPost(ctx, cfg)
	case "send_stdout":
		return newSendStdout(ctx, cfg)
	// String transforms.
	case "string_append":
		return newStringAppend(ctx, cfg)
	case "string_capture":
		return newStringCapture(ctx, cfg)
	case "string_to_lower":
		return newStringToLower(ctx, cfg)
	case "string_to_snake":
		return newStringToSnake(ctx, cfg)
	case "string_to_upper":
		return newStringToUpper(ctx, cfg)
	case "string_replace":
		return newStringReplace(ctx, cfg)
	case "string_split":
		return newStringSplit(ctx, cfg)
	case "string_uuid":
		return newStringUUID(ctx, cfg)
	// Test transforms.
	case "test_message":
		return newTestMessage(ctx, cfg)
	// Time transforms.
	case "time_from_string":
		return newTimeFromString(ctx, cfg)
	case "time_from_unix":
		return newTimeFromUnix(ctx, cfg)
	case "time_from_unix_milli":
		return newTimeFromUnixMilli(ctx, cfg)
	case "time_now":
		return newTimeNow(ctx, cfg)
	case "time_to_string":
		return newTimeToString(ctx, cfg)
	case "time_to_unix":
		return newTimeToUnix(ctx, cfg)
	case "time_to_unix_milli":
		return newTimeToUnixMilli(ctx, cfg)
	// Utility transforms.
	case "utility_control":
		return newUtilityControl(ctx, cfg)
	case "utility_delay":
		return newUtilityDelay(ctx, cfg)
	case "utility_drop":
		return newUtilityDrop(ctx, cfg)
	case "utility_err":
		return newUtilityErr(ctx, cfg)
	case "utility_metric_bytes":
		return newUtilityMetricBytes(ctx, cfg)
	case "utility_metric_count":
		return newUtilityMetricCount(ctx, cfg)
	case "utility_metric_freshness":
		return newUtilityMetricFreshness(ctx, cfg)
	case "utility_secret":
		return newUtilitySecret(ctx, cfg)
	default:
		return nil, fmt.Errorf("transform %s: %w", cfg.Type, iconfig.ErrInvalidFactoryInput)
	}
}

// Apply executes one or more transform functions on one or more messages in
// sequence. Each transform processes all messages from the previous step
// before passing results to the next transform.
//
// If a transform returns an error, Apply immediately returns that error
// along with a nil message slice. Transforms that need error handling
// should use the meta_err transform.
//
// Control messages are passed through each transform and can be used to
// trigger flushing of internal state in stateful transforms.
func Apply(ctx context.Context, tf []Transformer, msgs ...*message.Message) ([]*message.Message, error) {
	resultMsgs := make([]*message.Message, len(msgs))
	copy(resultMsgs, msgs)

	for i := 0; len(resultMsgs) > 0 && i < len(tf); i++ {
		var nextResultMsgs []*message.Message
		for _, m := range resultMsgs {
			rMsgs, err := tf[i].Transform(ctx, m)
			if err != nil {
				// We immediately return if a transform hits an unrecoverable
				// error on a message.
				return nil, err
			}
			nextResultMsgs = append(nextResultMsgs, rMsgs...)
		}
		resultMsgs = nextResultMsgs
	}

	return resultMsgs, nil
}

func bytesToValue(b []byte) message.Value {
	msg := message.New()
	_ = msg.SetValue("_", b)

	return msg.GetValue("_")
}

func anyToBytes(v any) []byte {
	msg := message.New()
	_ = msg.SetValue("_", v)

	return msg.GetValue("_").Bytes()
}

// truncateTTL truncates the time-to-live (TTL) value from any precision greater
// than seconds (e.g., milliseconds, nanoseconds) to seconds.
//
// For example:
//   - 1696482368492 -> 1696482368
//   - 1696482368492290 -> 1696482368
func truncateTTL(v message.Value) int64 {
	if len(v.String()) <= 10 {
		return v.Int()
	}

	l := len(v.String()) - 10
	return v.Int() / int64(math.Pow10(l))
}
