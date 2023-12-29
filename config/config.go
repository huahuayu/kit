package config

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Loader[T any] interface {
	LoadDefault(fn func() (T, error)) (T, error)
	Load(configPath string, fileType FileType) (T, error)
}

type FileType int

const (
	YAML FileType = iota
	TOML
	JSON
	DOTENV
)

type configLoader[T any] struct{}

func NewConfigLoader[T any]() Loader[T] {
	return &configLoader[T]{}
}

func (c *configLoader[T]) LoadDefault(fn func() (T, error)) (T, error) {
	return fn()
}

// Load loads configurations from a file
// support env tag to mapping env variable to struct field
// priority: env > config file
func (c *configLoader[T]) Load(configPath string, fileType FileType) (T, error) {
	var cfg T
	// First, load from file
	err := c.loadFromFile(&cfg, configPath, fileType)
	if err != nil {
		return cfg, err
	}

	// Override with env variables
	err = overrideWithEnv(&cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}

func (c *configLoader[T]) loadFromFile(cfg *T, configPath string, fileType FileType) error {
	switch fileType {
	case YAML:
		return loadYAML(configPath, cfg)
	case TOML:
		return loadTOML(configPath, cfg)
	case JSON:
		return loadJSON(configPath, cfg)
	case DOTENV:
		return loadDotEnv(configPath)
	default:
		return fmt.Errorf("unsupported file type: %v", fileType)
	}
}

func overrideWithEnv[T any](cfg *T) error {
	v := reflect.ValueOf(cfg).Elem()
	return overrideWithEnvRecursive(v)
}

func overrideWithEnvRecursive(v reflect.Value) error {
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.Struct {
			// Recursively handle nested structs
			if err := overrideWithEnvRecursive(field); err != nil {
				return err
			}
		} else if tag, ok := t.Field(i).Tag.Lookup("env"); ok && field.CanSet() {
			if envVal, exists := os.LookupEnv(tag); exists {
				log.Printf("config field %s is overridden by env variable %s\n", t.Field(i).Name, tag)
				if err := setField(field, envVal); err != nil {
					return fmt.Errorf("error setting field '%s' with env variable '%s': %w", t.Field(i).Name, tag, err)
				}
			}
		}
	}
	return nil
}

func loadYAML[T any](configPath string, cfg *T) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, cfg)
	return err
}

func loadTOML[T any](configPath string, cfg *T) error {
	_, err := toml.DecodeFile(configPath, cfg)
	return err
}

func loadJSON[T any](configPath string, cfg *T) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, cfg)
	return err
}

func loadDotEnv(configPath string) error {
	file, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid line in .env file: %s", line)
		}
		key := parts[0]
		value := parts[1]
		os.Setenv(key, value)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func setField(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intVal, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(intVal)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintVal, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(uintVal)

	case reflect.Bool:
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(boolVal)

	case reflect.Float32, reflect.Float64:
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		field.SetFloat(floatVal)

	case reflect.Slice:
		slice, err := parseSlice(field.Type().Elem(), value)
		if err != nil {
			return fmt.Errorf("error parsing slice for field: %w", err)
		}
		field.Set(slice)

	case reflect.Map:
		mapValue, err := parseMap(field.Type().Key(), field.Type().Elem(), value)
		if err != nil {
			return fmt.Errorf("error parsing map for field: %w", err)
		}
		field.Set(mapValue)

	default:
		return errors.New("unsupported field type")
	}
	return nil
}

func parseSlice(elemType reflect.Type, value string) (reflect.Value, error) {
	values := strings.Split(value, ",")
	slice := reflect.MakeSlice(reflect.SliceOf(elemType), len(values), len(values))
	for i, v := range values {
		v = strings.TrimSpace(v)
		switch elemType.Kind() {
		case reflect.String:
			slice.Index(i).SetString(v)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			intVal, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return reflect.ValueOf(nil), err
			}
			slice.Index(i).SetInt(intVal)
		// Add cases for other types as needed
		default:
			return reflect.ValueOf(nil), fmt.Errorf("unsupported slice element type: %v", elemType)
		}
	}
	return slice, nil
}

func parseMap(keyType, valueType reflect.Type, value string) (reflect.Value, error) {
	pairs := strings.Split(value, ",")
	result := reflect.MakeMapWithSize(reflect.MapOf(keyType, valueType), len(pairs))
	for _, pair := range pairs {
		kv := strings.SplitN(pair, ":", 2)
		if len(kv) != 2 {
			return reflect.ValueOf(nil), fmt.Errorf("invalid map value")
		}
		key := reflect.New(keyType).Elem()
		if err := setField(key, kv[0]); err != nil {
			return reflect.ValueOf(nil), err
		}
		val := reflect.New(valueType).Elem()
		if err := setField(val, kv[1]); err != nil {
			return reflect.ValueOf(nil), err
		}
		result.SetMapIndex(key, val)
	}
	return result, nil
}
