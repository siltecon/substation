# condition

The condition package contains interfaces and methods for evaluating data using success or failure criteria.

## Overview

Conditions are used to make decisions about message processing in Substation pipelines. They evaluate messages and return a boolean result, enabling conditional routing, filtering, and other decision-based operations.

## Condition Categories

Conditions are organized into the following categories:

| Category | Description | Examples |
|----------|-------------|----------|
| **Meta** | Combine multiple conditions | `all` (meta_all), `any` (meta_any), `none` (meta_none) |
| **Format** | Check data formats | `format_json`, `format_mime` |
| **Network** | Validate IP addresses | `network_ip_private`, `network_ip_valid`, `network_ip_loopback` |
| **Number** | Numeric comparisons | `number_equal_to`, `number_greater_than`, `number_less_than` |
| **String** | String comparisons | `string_equal_to`, `string_contains`, `string_match`, `string_starts_with` |
| **Utility** | Helper conditions | `utility_random` |

## Usage

Conditions are typically created using the `New` factory function with a configuration:

```go
cfg := config.Config{
    Type: "string_equal_to",
    Settings: map[string]interface{}{
        "object": map[string]interface{}{
            "source_key": "status",
        },
        "value": "active",
    },
}
cond, err := condition.New(ctx, cfg)

// Evaluate the condition
ok, err := cond.Condition(ctx, msg)
if ok {
    // Condition matched
}
```

## Meta Conditions

Meta conditions combine multiple conditions using logical operators:

- **all (meta_all)**: All conditions must be true
- **any (meta_any)**: At least one condition must be true
- **none (meta_none)**: No conditions can be true

Example of combining conditions:

```go
cfg := config.Config{
    Type: "meta_all",
    Settings: map[string]interface{}{
        "conditions": []interface{}{
            map[string]interface{}{
                "type": "string_equal_to",
                "settings": map[string]interface{}{
                    "object": map[string]interface{}{
                        "source_key": "type",
                    },
                    "value": "user",
                },
            },
            map[string]interface{}{
                "type": "number_greater_than",
                "settings": map[string]interface{}{
                    "object": map[string]interface{}{
                        "source_key": "age",
                    },
                    "value": 18,
                },
            },
        },
    },
}
```

## Control Messages

Control messages typically return `false` when evaluated by conditions. This allows conditions to effectively pass through control messages in conditional transforms like `meta_switch`.
