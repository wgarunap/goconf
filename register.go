package goconf

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/tryfix/log"
	"gopkg.in/oleiade/reflections.v1"
	"os"
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

	var data = [][]string{}

	pr := p.Print()
	var fields []string
	fields, _ = reflections.Fields(pr)

	for _, field := range fields {
		value, err := reflections.GetField(pr, field)
		if err != nil {
			log.Error("error printing the goconf table", err)
		}
		data = append(data, []string{field, fmt.Sprint(value)})
	}

	table.SetHeader([]string{"Config", "Value"})
	table.AppendBulk(data)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	table.Render()
}
