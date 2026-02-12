package yamlconfig

import (
	"fmt"
	"os"
	"reflect"

	"gopkg.in/yaml.v3"
)

// LoadConfig loads a YAML configuration file from the provided path and decodes it
// into the provided struct pointer. It also validates the loaded configuration.
//
// Parameters:
//
// path: The path to the configuration file.
// config: A pointer to the struct to decode the configuration into.
//
// Returns:
// error: An error if the configuration file could not be loaded or decoded.
//
// Example:
//
// cfg := config.Config{}
// err := config.LoadConfig("config.yml", &cfg)
//
//	if err != nil {
//	    log.Fatal(err)
//	}
func LoadConfig(path string, config interface{}) error {
	// Open the configuration file
	file, fileErr := os.Open(path)
	if fileErr != nil {
		return fmt.Errorf("failed to load config file: %w", fileErr)
	}
	defer file.Close()

	// Create a new YAML decoder for the file
	d := yaml.NewDecoder(file)

	// Decode the YAML content into the provided struct pointer
	if yamlDecodeErr := d.Decode(config); yamlDecodeErr != nil {
		return fmt.Errorf("failed to decode config file: %w", yamlDecodeErr)
	}

	// Validate the loaded configuration
	if validateConfigErr := validateConfig(config); validateConfigErr != nil {
		return fmt.Errorf(("failed to load the config: %w"), validateConfigErr)
	}

	return nil
}

// validateConfig function checks if the provided configuration is valid. It
// ensures that all required fields are present and non-empty.
func validateConfig(config interface{}) error {
	val := reflect.ValueOf(config)

	// Check if the config is a pointer and points to a struct
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("expected a pointer to a struct, please ensure the input is a struct pointer")
	}

	// Recursively validate the struct
	return validateStruct(val.Elem())
}

// validateStruct function recursively validates a struct and its fields.
// It checks if all required fields are present and non-empty.
// A field is considered required if it does not have the yamlconfig tag "omitempty".
//
// Optional nested structs (omitempty): Validation of child fields is skipped only when
// the parent section is absent from the config. When the parent is present, its
// required children are validated. Use a pointer to struct (e.g. *struct{...}) for
// omitempty nested sections so the decoder can distinguish absent (nil) from present
// (non-nil); with value structs, absent and empty are indistinguishable.
func validateStruct(val reflect.Value) error {
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		typ := val.Type().Field(i)

		// Check for the yamlconfig tag
		yamlConfigTag := typ.Tag.Get("yamlconfig")
		isOmitEmpty := yamlConfigTag == "omitempty"

		// If the field is required (no omitempty) and empty, return an error
		if !isOmitEmpty && isEmpty(field) {
			return fmt.Errorf("missing required config item: %s", typ.Name)
		}

		// Recursively validate nested structs when present. Skip only when the struct
		// is absent: for omitempty pointers, nil means absent; for value structs,
		// empty means we can't distinguish absent from present, so we skip.
		if nested := getNestedStruct(field); nested.IsValid() && isStructPresent(field) {
			if err := validateStruct(nested); err != nil {
				return err
			}
		}
	}

	return nil
}

// getNestedStruct returns the struct to validate for a field that may be a struct
// or a pointer to struct. Returns zero Value if the field is not a nested struct.
func getNestedStruct(field reflect.Value) reflect.Value {
	switch field.Kind() {
	case reflect.Struct:
		return field
	case reflect.Ptr:
		if !field.IsNil() && field.Elem().Kind() == reflect.Struct {
			return field.Elem()
		}
	}
	return reflect.Value{}
}

// isStructPresent returns true if a struct/ptr-to-struct field is present in the
// config. For pointers: non-nil means the key was in YAML. For value structs:
// non-empty means at least one field was set (we cannot distinguish absent from
// present-but-empty for value structs).
func isStructPresent(field reflect.Value) bool {
	if field.Kind() == reflect.Ptr {
		return !field.IsNil()
	}
	return !isEmpty(field)
}

// isEmpty function checks if a value is empty. It is used to validate the
// configuration values.
func isEmpty(v reflect.Value) bool {
	switch v.Kind() { //nolint:exhaustive // We don't need to handle all types
	case reflect.Ptr:
		return v.IsNil() || isEmpty(v.Elem())
	case reflect.String:
		return v.String() == ""
	case reflect.Slice, reflect.Map, reflect.Array:
		return v.Len() == 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Bool:
		return false
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			if !isEmpty(v.Field(i)) {
				return false
			}
		}
		return true
	}

	return false
}
