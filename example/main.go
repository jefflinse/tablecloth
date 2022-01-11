package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/jefflinse/tablecloth"
)

func main() {
	t := tablecloth.NewTableWithColumns([]tablecloth.ColumnDefinition{
		{Name: "First", MinLength: 10},
		{Name: "Second"},
		{Name: "Third", MinLength: 10},
		{Name: "Fourth"},
	})
	t.AddRow(
		tablecloth.Cell{Format: "%s color me some text", Values: []tablecloth.FormattableCellValue{
			{Value: "000", Format: fmt.Sprint},
		}},
		tablecloth.Cell{Format: "No color here."},
		tablecloth.Cell{Format: "color %s text", Values: []tablecloth.FormattableCellValue{
			{Value: "me some more", Format: color.New().SprintFunc()},
		}},
		tablecloth.Cell{Format: "No color here."},
	)
	t.AddRow(
		tablecloth.Cell{Format: "%s color me some text ", Values: []tablecloth.FormattableCellValue{
			{Value: "111", Format: color.New(color.BgYellow).SprintFunc()},
		}},
		tablecloth.Cell{Format: "No color here."},
		tablecloth.Cell{Format: "color %s text", Values: []tablecloth.FormattableCellValue{
			{Value: "me some more", Format: color.New(color.BgCyan).SprintFunc()},
		}},
		tablecloth.Cell{Format: "No color here."},
	)
	t.AddRow(
		tablecloth.Cell{Format: "%s color me some text ", Values: []tablecloth.FormattableCellValue{
			{Value: "222", Format: color.New(color.FgHiGreen, color.Underline).SprintFunc()},
		}},
		tablecloth.Cell{Format: "No color here."},
		tablecloth.Cell{Format: "color %s text", Values: []tablecloth.FormattableCellValue{
			{Value: "me some more", Format: color.New(color.BgCyan, color.FgHiRed).SprintFunc()},
		}},
		tablecloth.Cell{Format: "No color here."},
	)
	t.AddRow(
		tablecloth.Cell{Format: "%s color me some text ", Values: []tablecloth.FormattableCellValue{
			{Value: "333", Format: color.New(color.FgHiGreen, color.Underline, color.BgRed).SprintFunc()},
		}},
		tablecloth.Cell{Format: "No color here."},
		tablecloth.Cell{Format: "color %s text", Values: []tablecloth.FormattableCellValue{
			{Value: "me some more", Format: color.New(color.BgCyan, color.FgHiRed, color.Underline).SprintFunc()},
		}},
		tablecloth.Cell{Format: "No color here."},
	)
	t.AddRow(
		tablecloth.Cell{Format: "%s color me some text ", Values: []tablecloth.FormattableCellValue{
			{Value: "444", Format: color.New(color.FgHiGreen, color.Underline, color.BgRed, color.CrossedOut).SprintFunc()},
		}},
		tablecloth.Cell{Format: "No color here."},
		tablecloth.Cell{Format: "color %s text", Values: []tablecloth.FormattableCellValue{
			{Value: "me some more", Format: color.New(color.BgCyan, color.FgHiRed, color.Underline, color.CrossedOut).SprintFunc()},
		}},
		tablecloth.Cell{Format: "No color here."},
	)

	if err := t.Write(os.Stdout); err != nil {
		log.Fatal(err)
	}

	// fmt.Printf("raw (%d):\n\t%s\n", len(c.Plain()), c.Plain())
	// fmt.Printf("fmt (%d (%d colors, %d overhead)):\n\t%s\n", len(c.Formatted()), c.NumColors(), c.Overhead(), c.Formatted())
}
