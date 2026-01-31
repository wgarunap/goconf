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

func printTable(p Printer) error {
	table := tablewriter.NewWriter(os.Stdout)

	var data [][]string

	printer := p.Print()

	values := reflect.ValueOf(printer)
	if values.Kind() == reflect.Ptr {
		values = values.Elem()
	}

	if values.Kind() == reflect.Interface {
		values = values.Elem()
	}

	for i := 0; i < values.NumField(); i++ {
		field := values.Field(i)
		structField := values.Type().Field(i)

		secretTag, ok := structField.Tag.Lookup("secret")
		if ok && secretTag == "true" {
			data = append(data, []string{structField.Name, SensitiveDataMaskString})
			continue
		}

		if field.Kind() == reflect.Int {
			data = append(data, []string{structField.Name, strconv.Itoa(int(field.Int()))})
			continue
		}

		data = append(data, []string{structField.Name, field.String()})
	}

	table.Header("Config", "Value")

	if err := table.Bulk(data); err != nil {
		return fmt.Errorf("failed to add table data: %w", err)
	}

	if err := table.Render(); err != nil {
		return fmt.Errorf("failed to render table: %w", err)
	}

	return nil
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

	configMap := make(map[string]interface{})

	for i := 0; i < values.NumField(); i++ {
		field := values.Field(i)
		structField := values.Type().Field(i)

		secretTag, ok := structField.Tag.Lookup("secret")
		if ok && secretTag == "true" {
			configMap[structField.Name] = SensitiveDataMaskString
			continue
		}

		configMap[structField.Name] = field.Interface()
	}

	jsonData, err := json.MarshalIndent(configMap, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config to JSON: %w", err)
	}

	// Create a logger that writes to stdout with timestamp
	logger := log.New(os.Stdout, "", log.LstdFlags)
	logger.Println(string(jsonData))

	return nil
}
