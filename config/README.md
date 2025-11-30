# config

The config package provides structures for building configurations in Substation.

## Overview

The `Config` type is the fundamental building block used throughout Substation to configure transforms, conditions, and other components. It follows a consistent pattern where `Type` identifies the component type and `Settings` contains type-specific configuration options.

## Usage

### Basic Configuration

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
```

### JSON Representation

Configurations can be easily serialized to and from JSON:

```json
{
    "type": "object_copy",
    "settings": {
        "object": {
            "source_key": "input.field",
            "target_key": "output.field"
        }
    }
}
```

### Jsonnet Configuration

In practice, configurations are typically written in Jsonnet using the `substation.libsonnet` library:

```jsonnet
local sub = import 'substation.libsonnet';

{
    transforms: [
        sub.tf.obj.cp({
            object: {
                source_key: 'input.field',
                target_key: 'output.field',
            },
        }),
    ],
}
```

## Factory Pattern

The `Config` type is used with factory functions throughout Substation:

- `transform.New(ctx, cfg)` creates transforms
- `condition.New(ctx, cfg)` creates conditions
- Other factories for KV stores, secrets, metrics destinations, etc.

## Common Settings

While settings vary by component type, some common patterns include:

| Setting | Description |
|---------|-------------|
| `id` | Unique identifier for the component |
| `object` | Configuration for reading from/writing to JSON objects |
| `batch` | Configuration for batching multiple messages |
| `transforms` | Nested transforms (for meta transforms) |
| `condition` | A condition to evaluate (for conditional transforms) |

### Object Settings

```go
"object": map[string]interface{}{
    "source_key": "path.to.source",   // Where to read data from
    "target_key": "path.to.target",   // Where to write data to
    "batch_key":  "path.to.batch",    // Key for batching
}
```

### Batch Settings

```go
"batch": map[string]interface{}{
    "count":    1000,    // Maximum number of messages
    "size":     1000000, // Maximum size in bytes
    "duration": "1m",    // Maximum duration
}
```
