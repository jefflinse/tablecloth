package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
)

func main() {
	t := NewTableWithColumns([]ColumnDefinition{
		{Name: "First", MinLength: 10},
		{Name: "Second"},
		{Name: "Third", MinLength: 10},
		{Name: "Fourth"},
	})
	t.Debug = true
	t.AddRow(Row{
		Cell{Format: "%s color me some text", Values: []FormattableCellValue{
			{Value: "000", Format: fmt.Sprint},
		}},
		Cell{Format: "No color here."},
		Cell{Format: "color %s text", Values: []FormattableCellValue{
			{Value: "me some more", Format: color.New().SprintFunc()},
		}},
		Cell{Format: "No color here."},
	})
	t.AddRow(Row{
		Cell{Format: "%s color me some text ", Values: []FormattableCellValue{
			{Value: "111", Format: color.New(color.BgYellow).SprintFunc()},
		}},
		Cell{Format: "No color here."},
		Cell{Format: "color %s text", Values: []FormattableCellValue{
			{Value: "me some more", Format: color.New(color.BgCyan).SprintFunc()},
		}},
		Cell{Format: "No color here."},
	})
	t.AddRow(Row{
		Cell{Format: "%s color me some text ", Values: []FormattableCellValue{
			{Value: "222", Format: color.New(color.FgHiGreen, color.Underline).SprintFunc()},
		}},
		Cell{Format: "No color here."},
		Cell{Format: "color %s text", Values: []FormattableCellValue{
			{Value: "me some more", Format: color.New(color.BgCyan, color.FgHiRed).SprintFunc()},
		}},
		Cell{Format: "No color here."},
	})
	t.AddRow(Row{
		Cell{Format: "%s color me some text ", Values: []FormattableCellValue{
			{Value: "333", Format: color.New(color.FgHiGreen, color.Underline, color.BgRed).SprintFunc()},
		}},
		Cell{Format: "No color here."},
		Cell{Format: "color %s text", Values: []FormattableCellValue{
			{Value: "me some more", Format: color.New(color.BgCyan, color.FgHiRed, color.Underline).SprintFunc()},
		}},
		Cell{Format: "No color here."},
	})
	t.AddRow(Row{
		Cell{Format: "%s color me some text ", Values: []FormattableCellValue{
			{Value: "444", Format: color.New(color.FgHiGreen, color.Underline, color.BgRed, color.CrossedOut).SprintFunc()},
		}},
		Cell{Format: "No color here."},
		Cell{Format: "color %s text", Values: []FormattableCellValue{
			{Value: "me some more", Format: color.New(color.BgCyan, color.FgHiRed, color.Underline, color.CrossedOut).SprintFunc()},
		}},
		Cell{Format: "No color here."},
	})

	if err := t.Write(os.Stdout); err != nil {
		log.Fatal(err)
	}

	// fmt.Printf("raw (%d):\n\t%s\n", len(c.Plain()), c.Plain())
	// fmt.Printf("fmt (%d (%d colors, %d overhead)):\n\t%s\n", len(c.Formatted()), c.NumColors(), c.Overhead(), c.Formatted())
}
