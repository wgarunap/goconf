package goconf

import (
	"github.com/olekukonko/tablewriter"
	"os"
	"reflect"
)

const (
	defaultLength    = 10
	defaultAddLength = 5
)

type Configer interface {
	Register() error
}
type Validater interface {
	Validate() error
}
type Printer interface {
	Print() interface{}
}

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
			printTable(p)
		}
	}
	return nil
}

func printTable(p Printer) {
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

		_, ok := structField.Tag.Lookup("secret")
		if ok {
			data = append(data, []string{structField.Name, mask(field.String())})
		} else {
			data = append(data, []string{structField.Name, field.String()})
		}
	}

	table.SetHeader([]string{"Config", "Value"})
	table.AppendBulk(data)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	table.Render()
}

func mask(value string) string {
	length := defaultLength
	if len(value) > defaultLength {
		length = len(value) + defaultAddLength
	}
	runes := make([]rune, length)

	for i := 0; i < length; i++ {
		runes[i] = '*'
	}

	return string(runes)
}
