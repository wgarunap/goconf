//go:generate mockgen -source=register.go -destination=mocks/register_mock.go -package=mocks

package goconf

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"

	"github.com/olekukonko/tablewriter"
)

// SensitiveDataMaskString is the default mask used to hide sensitive configuration values
const SensitiveDataMaskString = "***************"

// OutputFormat defines the format for configuration output
type OutputFormat string

const (
	// OutputFormatTable outputs configuration as a table (default)
	OutputFormatTable OutputFormat = "table"
	// OutputFormatJSON outputs configuration as JSON
	OutputFormatJSON OutputFormat = "json"
)

var currentOutputFormat = OutputFormatTable

// SetOutputFormat sets the output format for configuration printing
func SetOutputFormat(format OutputFormat) {
	currentOutputFormat = format
}

// Configer interface must be implemented by configuration structs
// to enable environment variable registration
type Configer interface {
	Register() error
}

// Validater interface can be implemented to enable configuration validation
type Validater interface {
	Validate() error
}

// Printer interface can be implemented to enable configuration output
type Printer interface {
	Print() interface{}
}

// Load registers, validates, and prints one or more configuration objects
func Load(configs ...Configer) error {
	for _, c := range configs {
		err := c.Register()
		if err != nil {
			return err
		}

		v, ok := c.(Validater)
		if ok {
			err = v.Validate()
			if err != nil {
				return err
			}
		}

		p, ok := c.(Printer)
		if ok {
			switch currentOutputFormat {
			case OutputFormatJSON:
				if err := printJSON(p); err != nil {
					return err
				}
			default:
				if err := printTable(p); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// extractFields recursively extracts fields from a struct and returns them as table rows
func extractFields(prefix string, values reflect.Value) [][]string {
	var data [][]string

	for i := 0; i < values.NumField(); i++ {
		field := values.Field(i)
		structField := values.Type().Field(i)

		fieldName := structField.Name
		if prefix != "" {
			fieldName = prefix + "." + fieldName
		}

		// Check if field is marked as secret
		secretTag, ok := structField.Tag.Lookup("secret")
		if ok && secretTag == "true" {
			data = append(data, []string{fieldName, SensitiveDataMaskString})
			continue
		}

		// Handle different field types
		switch field.Kind() {
		case reflect.Struct:
			// Recursively process nested structs
			nestedData := extractFields(fieldName, field)
			data = append(data, nestedData...)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			data = append(data, []string{fieldName, strconv.FormatInt(field.Int(), 10)})
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			data = append(data, []string{fieldName, strconv.FormatUint(field.Uint(), 10)})
		case reflect.Float32, reflect.Float64:
			data = append(data, []string{fieldName, strconv.FormatFloat(field.Float(), 'f', -1, 64)})
		case reflect.Bool:
			data = append(data, []string{fieldName, strconv.FormatBool(field.Bool())})
		case reflect.String:
			data = append(data, []string{fieldName, field.String()})
		default:
			// For other types, use string representation
			data = append(data, []string{fieldName, fmt.Sprintf("%v", field.Interface())})
		}
	}

	return data
}

func printTable(p Printer) error {
	table := tablewriter.NewWriter(os.Stdout)

	printer := p.Print()

	values := reflect.ValueOf(printer)
	if values.Kind() == reflect.Ptr {
		values = values.Elem()
	}

	if values.Kind() == reflect.Interface {
		values = values.Elem()
	}

	data := extractFields("", values)

	table.Header("Config", "Value")

	if err := table.Bulk(data); err != nil {
		return fmt.Errorf("failed to add table data: %w", err)
	}

	if err := table.Render(); err != nil {
		return fmt.Errorf("failed to render table: %w", err)
	}

	return nil
}

// extractJSONFields recursively extracts fields from a struct and returns them as a map for JSON marshaling
func extractJSONFields(values reflect.Value) map[string]interface{} {
	configMap := make(map[string]interface{})

	for i := 0; i < values.NumField(); i++ {
		field := values.Field(i)
		structField := values.Type().Field(i)

		// Check if field is marked as secret
		secretTag, ok := structField.Tag.Lookup("secret")
		if ok && secretTag == "true" {
			configMap[structField.Name] = SensitiveDataMaskString
			continue
		}

		// Handle nested structs recursively
		if field.Kind() == reflect.Struct {
			configMap[structField.Name] = extractJSONFields(field)
		} else {
			configMap[structField.Name] = field.Interface()
		}
	}

	return configMap
}

func printJSON(p Printer) error {
	printer := p.Print()

	values := reflect.ValueOf(printer)
	if values.Kind() == reflect.Ptr {
		values = values.Elem()
	}

	if values.Kind() == reflect.Interface {
		values = values.Elem()
	}

	configMap := extractJSONFields(values)

	jsonData, err := json.MarshalIndent(configMap, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config to JSON: %w", err)
	}

	// Create a logger that writes to stdout with timestamp
	logger := log.New(os.Stdout, "", log.LstdFlags)
	logger.Println(string(jsonData))

	return nil
}
