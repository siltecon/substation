# transform

The transform package contains interfaces and methods for transforming data in Substation pipelines.

## Overview

Transforms are the core building blocks of Substation data pipelines. Each transform implements the `Transformer` interface and processes messages by reading data, modifying it, and returning zero or more result messages.

## Transform Categories

Transforms are organized into the following categories:

| Category | Description | Examples |
|----------|-------------|----------|
| **Aggregation** | Combine or split messages | `aggregate_to_array`, `aggregate_from_array`, `aggregate_to_string` |
| **Array** | Operations on JSON arrays | `array_join`, `array_zip` |
| **Enrichment** | Add data from external sources | `enrich_http_get`, `enrich_http_post`, `enrich_dns_lookup`, `enrich_aws_dynamodb_query` |
| **Format** | Convert data formats | `format_from_base64`, `format_to_gzip`, `format_from_zip` |
| **Hash** | Generate hashes | `hash_md5`, `hash_sha256` |
| **Meta** | Control flow and composition | `meta_switch`, `meta_for_each`, `meta_err`, `meta_retry` |
| **Network** | Domain and IP operations | `network_domain_subdomain`, `network_domain_registered_domain` |
| **Number** | Numeric operations | `number_math_addition`, `number_maximum`, `number_minimum` |
| **Object** | JSON object manipulation | `object_copy`, `object_delete`, `object_insert`, `object_jq` |
| **Send** | Output data to destinations | `send_stdout`, `send_aws_s3`, `send_aws_kinesis_data_stream` |
| **String** | String manipulation | `string_replace`, `string_to_lower`, `string_capture`, `string_uuid` |
| **Time** | Time parsing and formatting | `time_from_string`, `time_to_unix`, `time_now` |
| **Utility** | Helper transforms | `utility_drop`, `utility_delay`, `utility_control`, `utility_secret` |

## Usage

Transforms are typically created using the `New` factory function with a configuration:

```go
cfg := config.Config{
    Type: "object_copy",
    Settings: map[string]interface{}{
        "object": map[string]interface{}{
            "source_key": "input.field",
            "target_key": "output.field",
        },
    },
}
tf, err := transform.New(ctx, cfg)
```

Multiple transforms can be applied in sequence using the `Apply` function:

```go
msgs, err := transform.Apply(ctx, transforms, messages...)
```

## Custom Transforms

Custom transforms can be created by implementing the `Transformer` interface:

```go
type MyTransform struct {
    // Configuration fields
}

func (t *MyTransform) Transform(ctx context.Context, msg *message.Message) ([]*message.Message, error) {
    if msg.IsControl() {
        return []*message.Message{msg}, nil
    }
    // Process the message...
    return []*message.Message{msg}, nil
}
```

## Control Messages

Control messages are special messages used for flow control. Transforms should typically pass control messages through unchanged, or use them to trigger flushing of internal state in stateful transforms.
