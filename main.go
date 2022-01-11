package main

import (
	"log"
	"os"

	"github.com/fatih/color"
)

func main() {
	t := NewTableWriter(os.Stdout, 2)
	t.Debug = true
	t.AddRow(Row{
		Cell{Format: "%s color", Values: []ColorableCellValue{
			{Value: "000", Colors: []color.Attribute{}},
		}},
		Cell{Format: "No color here."},
	})
	t.AddRow(Row{
		Cell{Format: "%s color", Values: []ColorableCellValue{
			{Value: "111", Colors: []color.Attribute{color.BgYellow}},
		}},
		Cell{Format: "No color here."},
	})
	t.AddRow(Row{
		Cell{Format: "%s color", Values: []ColorableCellValue{
			{Value: "222", Colors: []color.Attribute{color.FgHiGreen, color.Underline}},
		}},
		Cell{Format: "No color here."},
	})
	t.AddRow(Row{
		Cell{Format: "%s color", Values: []ColorableCellValue{
			{Value: "333", Colors: []color.Attribute{color.FgHiGreen, color.Underline, color.BgRed}},
		}},
		Cell{Format: "No color here."},
	})
	t.AddRow(Row{
		Cell{Format: "%s color", Values: []ColorableCellValue{
			{Value: "444", Colors: []color.Attribute{color.FgHiGreen, color.Underline, color.BgRed, color.CrossedOut}},
		}},
		Cell{Format: "No color here."},
	})

	if err := t.Flush(); err != nil {
		log.Fatal(err)
	}

	// fmt.Printf("raw (%d):\n\t%s\n", len(c.Plain()), c.Plain())
	// fmt.Printf("fmt (%d (%d colors, %d overhead)):\n\t%s\n", len(c.Formatted()), c.NumColors(), c.Overhead(), c.Formatted())
}
