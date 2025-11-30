# message

The message package provides the core data structure used by conditions and transforms in Substation pipelines.

## Overview

Messages are the fundamental data structure in Substation. They flow through the transformation pipeline and can contain JSON text, binary data, or both (as data and metadata).

## Message Types

There are two types of messages:

- **Data messages**: Contain data to be processed by transforms. Data can be JSON text (accessed via GetValue/SetValue) or binary (accessed via Data/SetData).

- **Control messages**: Special messages used for flow control, typically to signal the end of a batch or to flush stateful transforms. Created using `AsControl()`.

## Usage

### Creating Messages

```go
// Create a message with JSON data
msg := message.New().SetData([]byte(`{"name": "Alice", "age": 30}`))

// Create a control message
ctrl := message.New().AsControl()
```

### Working with JSON Data

For JSON messages, the `GetValue`, `SetValue`, and `DeleteValue` methods provide a convenient way to read and modify specific fields using dot notation:

```go
msg := message.New().SetData([]byte(`{"user": {"name": "Alice"}}`))

// Read a nested value
name := msg.GetValue("user.name").String() // "Alice"

// Set a new value
msg.SetValue("user.age", 30)

// Delete a value
msg.DeleteValue("user.name")
```

### Working with Metadata

Each message can carry metadata, which is accessed by prefixing keys with "meta ":

```go
// Set metadata
msg.SetValue("meta source", "api")
msg.SetValue("meta timestamp", 1234567890)

// Read metadata
source := msg.GetValue("meta source").String()
```

Binary metadata is accessed using the `Metadata` and `SetMetadata` methods:

```go
msg.SetMetadata([]byte(`{"key": "value"}`))
meta := msg.Metadata()
```

### Checking Message Types

```go
if msg.IsControl() {
    // Handle control message
    return []*message.Message{msg}, nil
}
// Process data message
```

## Value Type

The `Value` type is returned by `GetValue` and provides methods to convert JSON values to Go types:

```go
val := msg.GetValue("field")

// Check if value exists
if val.Exists() {
    str := val.String()    // Get as string
    num := val.Int()       // Get as int64
    flt := val.Float()     // Get as float64
    b := val.Bool()        // Get as bool
    arr := val.Array()     // Get as []Value
}
```

## Best Practices

1. **Always check for control messages** in transforms and handle them appropriately (usually by passing them through).

2. **Use method chaining** for creating and configuring messages:
   ```go
   msg := message.New().SetData(data).SetMetadata(meta)
   ```

3. **Check if values exist** before using them to avoid operating on empty data:
   ```go
   if val := msg.GetValue("key"); val.Exists() {
       // Use val
   }
   ```
