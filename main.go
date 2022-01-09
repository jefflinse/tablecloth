package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/fatih/color"
)

type TableWriter struct {
	rowFormat string
	rows      []Row
	buf       *strings.Builder
	tw        *tabwriter.Writer
}

func NewTableWriter(w io.Writer, columns int, padding int) *TableWriter {
	buf := &strings.Builder{}
	return &TableWriter{
		rowFormat: strings.Repeat("%s\t", columns) + "\n",
		rows:      []Row{},

		buf: buf,
		tw:  tabwriter.NewWriter(buf, 0, 0, padding, ' ', tabwriter.Debug),
	}
}

func (t *TableWriter) AddRow(row Row) {
	t.rows = append(t.rows, row)
}

func (t *TableWriter) Render(w io.Writer, opts TableRenderOptions) error {
	for _, r := range t.rows {
		if _, err := fmt.Fprintf(t.tw, t.rowFormat, r.Formatted()...); err != nil {
			return err
		}
	}

	if err := t.tw.Flush(); err != nil {
		return err
	}

	// if lines, ok := w.nonTableLines[-1]; ok {
	// 	fmt.Fprintln(w.dest, strings.Join(lines, "\n"))
	// }

	s := bufio.NewScanner(strings.NewReader(t.buf.String()))
	i := 0
	for ; s.Scan(); i++ {
		fmt.Fprintln(w, s.Text())
		// if lines, ok := w.nonTableLines[i]; ok {
		// 	fmt.Fprintln(t.dest, strings.Join(lines, "\n"))
		// }
	}

	// if lines, ok := w.nonTableLines[i]; ok {
	// 	fmt.Fprintln(t.dest, strings.Join(lines, "\n"))
	// }

	return nil
}

type TableRenderOptions struct {
	ElasticColumnIndex int
	MaxWidth           int
}

type Row []Cell

func (r Row) Formatted() []interface{} {
	values := []interface{}{}
	for _, c := range r {
		values = append(values, c.Formatted())
	}

	return values
}

type Cell struct {
	Format string
	Values []ColorableCellValue
}

func (c Cell) Formatted() string {
	values := []interface{}{}
	for _, v := range c.Values {
		values = append(values, v.Formatted())
	}
	return fmt.Sprintf(c.Format, values...)
}

func (c Cell) Plain() string {
	values := []interface{}{}
	for _, v := range c.Values {
		values = append(values, v.Plain())
	}
	return fmt.Sprintf(c.Format, values...)
}

func (c Cell) NumColors() int {
	num := 0
	for _, v := range c.Values {
		num += len(v.Colors)
	}
	return num
}

// Overhead is the number of bytes attributable to the ANSI
// escape sequence(s) used to format the cell.
func (v Cell) Overhead() int {
	return len(v.Formatted()) - len(v.Plain())
}

type ColorableCellValue struct {
	Value  interface{}
	Colors []color.Attribute
}

func (v ColorableCellValue) Plain() string {
	return fmt.Sprint(v.Value)
}

func (v ColorableCellValue) Formatted() string {
	if len(v.Colors) == 0 {
		return v.Plain()
	}

	clr := color.New()
	clr.Add(v.Colors...)
	return clr.Sprint(v.Value)
}

func main() {
	t := NewTableWriter(os.Stdout, 3, 2)
	r := Row{
		Cell{
			Format: "Hello, %s It's %s degrees outside.",
			Values: []ColorableCellValue{
				{
					Value:  "World!",
					Colors: []color.Attribute{color.FgHiGreen},
				},
				{
					Value:  75,
					Colors: []color.Attribute{color.BgHiYellow, color.FgBlack},
				},
			},
		},
		Cell{
			Format: "No color here.",
		},
		Cell{
			Format: "Just %s color.",
			Values: []ColorableCellValue{
				{
					Value:  "one",
					Colors: []color.Attribute{color.FgHiGreen},
				},
			},
		},
	}
	t.AddRow(r)
	if err := t.Render(os.Stdout, TableRenderOptions{}); err != nil {
		log.Fatal(err)
	}

	// fmt.Printf("raw (%d):\n\t%s\n", len(c.Plain()), c.Plain())
	// fmt.Printf("fmt (%d (%d colors, %d overhead)):\n\t%s\n", len(c.Formatted()), c.NumColors(), c.Overhead(), c.Formatted())
}
