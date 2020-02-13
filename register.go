package goconf

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/pickme-go/log/v2"
	"gopkg.in/oleiade/reflections.v1"
	"os"
)

type Configer interface {
	Register()
}
type Validater interface {
	Validate()
}
type Printer interface {
	Print() interface{}
}

func Load(configs ...Configer) {
	for _, c := range configs {
		c.Register()

		v, ok := c.(Validater)
		if ok {
			v.Validate()
		}

		p, ok := c.(Printer)
		if ok {
			printTable(p)
		}
	}
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
