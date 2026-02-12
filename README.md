# YAMLConfig

![YAMLConfig](.github/logo.png?raw=true)

## What is YAMLConfig?

`YAMLConfig` is a Go package designed to simplify the loading, decoding, and validation of YAML configuration files into Go structs. It supports deep validation of nested structs, ensuring that your application's configuration is correctly loaded and validated with minimal boilerplate code.

## Features

- Easy loading and decoding of YAML files into Go structs.
- Deep validation of nested structs to ensure all required configuration items are present and correctly formatted.
- Optional fields support via a custom `yamlconfig:"omitempty"` tag, allowing certain configuration fields to be omitted.
- Custom error messages for missing or invalid configuration items.
- Support for a wide range of field types within configuration structs.

## Install

To install `YAMLConfig`, use the `go get` command:

```bash
go get github.com/sculley/yamlconfig
```

## Getting Started

### Support Struct Types

The following types are supported in your Config struct and will be validated against the fields in your YAML config file.

- String
- Int/Int8/Int64/Int32/Int64
- Uint/Uint8/Uint32/Uint64
- Float32/Float64
- Bool
- Slice/Array/Map
- Struct

### Optional Fields

By default, YAMLConfig expects all fields to be present and non-empty. However, you may have optional fields that you want to allow missing or empty values for. To mark a field as optional, annotate it with the yamlconfig:"omitempty" tag. If a field is empty and is marked as omitempty, it will not produce a validation error.

#### Example

```go
type Config struct {
    Service string            `yaml:"service"`
    Image   string            `yaml:"image"`
    Servers []string          `yaml:"servers"`
    Env     map[string]string `yaml:"env" yamlconfig:"omitempty"`
    Volumes []string          `yaml:"volumes" yamlconfig:"omitempty"`
    Ports   []string          `yaml:"ports" yamlconfig:"omitempty"`
}
```

In this example, Env, Volumes, and Ports fields are optional. If your YAML file omits these fields or leaves them empty, YAMLConfig will not return an error during validation.

#### Optional Nested Structs (Mutually Exclusive Sections)

For optional nested structs where you need "if present, validate children" (e.g. backend-specific config), use a **pointer to struct** with `omitempty`. When the section is absent, child validation is skipped. When present, required children are validated.

```go
type Config struct {
    Backend   string `yaml:"backend"`
    S3        *struct { Bucket string `yaml:"bucket"` }  `yaml:"s3" yamlconfig:"omitempty"`
    FileSystem *struct { Path string `yaml:"path"` }     `yaml:"file-system" yamlconfig:"omitempty"`
}
```

- `backend: s3` with `s3:` absent → valid (no S3 validation)
- `backend: s3` with `s3:` present but `bucket` missing → validation error
- `backend: s3` with `s3: { bucket: "my-bucket" }` → valid

With value structs (not pointers), absent and present-but-empty are indistinguishable, so child validation is skipped when the struct is empty.

### Creating a Configuration File

Define your configuration in a YAML file as follows:

```yaml
server:
  address: 10.0.0.1
  port: 8080
database:
  user: sculley
  password: SUPERSECUREPASSWORD
  database: yamlconfig
```

### Implementing YAMLConfig

- Integrate YAMLConfig into your Go application with the following steps:

```go
package main

import (
    "log"
    "github.com/sculley/yamlconfig"
)

type Config struct {
    Server struct {
        Address string `yaml:"address"`
        Port    int    `yaml:"port"`
    } `yaml:"server"`
    Database struct {
        User     string `yaml:"user"`
        Password string `yaml:"password"`
        Database string `yaml:"database"`
    } `yaml:"database"`
}
```

- Load and Validate Configuration:

```go
func main() {
    cfg := Config
    err := yamlconfig.LoadConfig("path/to/your/config.yml", &cfg)
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }
}
```

- Accessing Values

```go
addr := cfg.Server.Address
```

`YAMLConfig` automatically validates the YAML configuration against your defined Go structs. It ensures all required fields are present and correctly typed, returning an error for any discrepancies.

By following these steps, you can leverage YAMLConfig to efficiently manage your application's configuration, focusing more on your core application logic and less on configuration management.

## Development

Run the test suite:

```shell
make test
make coverage
```

Run linters:

```shell
make lint
```

Some linter violations can automatically be fixed:

```shell
make fmt
```

## License

The project is licensed under the [MIT License](LICENSE).
