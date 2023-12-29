# Config

The `config` package provides a simple and efficient way to manage configuration in Go applications.

## Features

- **Easy to use**: The `config` package is designed to be simple and intuitive to use.
- **Flexible**: It supports multiple configuration formats such as JSON, YAML, and TOML.
- **Environment variable support**: You can easily override configuration values using environment variables.

## Usage

Import the `config` package:

```go
import "github.com/huahuayu/kit/config"
```

Create a new config instance and load a configuration file:

```go
loader := NewConfigLoader[sampleConfig]()
cfg, err := loader.Load("testdata/sample.yaml", config.YAML)
```