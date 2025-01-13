//go:generate mockgen -source=register.go -destination=mocks/register_mock.go -package=mocks

package goconf

import (
	"os"
	"reflect"
	"strconv"

	"github.com/olekukonko/tablewriter"
)

const SensitiveDataMaskString = "***************"

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

	table.SetHeader([]string{"Config", "Value"})
	table.AppendBulk(data)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	table.Render()
}
