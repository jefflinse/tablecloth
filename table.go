package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"regexp"
	"strings"
	"text/tabwriter"

	"github.com/fatih/color"
)

type TableWriter struct {
	rowFormat     string
	rows          []Row
	overheads     []int
	nonTableLines map[int][]string
	currentLine   int
	buf           *strings.Builder
	tw            *tabwriter.Writer
	dest          io.Writer

	Debug bool
}

func NewTableWriter(w io.Writer, columns int) *TableWriter {
	buf := &strings.Builder{}
	return &TableWriter{
		rowFormat:     strings.Repeat("%s\t", columns) + "\n",
		rows:          []Row{},
		overheads:     make([]int, columns),
		nonTableLines: map[int][]string{},
		currentLine:   -1,
		buf:           buf,
		tw:            tabwriter.NewWriter(buf, 0, 0, 2, ' ', tabwriter.Debug),
		dest:          w,
	}
}

func (t *TableWriter) AddRow(row Row) {
	rendered := row.Render(true, 0)
	for i := range rendered {
		if rendered[i].overhead > t.overheads[i] {
			t.overheads[i] = rendered[i].overhead
		}
	}
	t.rows = append(t.rows, row)
	t.currentLine++
}

func (t *TableWriter) AddLine(line string) {
	t.nonTableLines[t.currentLine] = append(t.nonTableLines[t.currentLine], line)
}

func (t *TableWriter) Flush() error {
	rows := [][]interface{}{}

	for r := range t.rows {
		rowValues := []interface{}{}
		cells := t.rows[r].Render(true, 0)
		for c := range cells {
			if cells[c].overhead < t.overheads[c] {
				cells[c] = cells[c].AdjustOverhead(t.overheads[c] - cells[c].overhead)
			}
			rowValues = append(rowValues, cells[c].value)
		}

		rows = append(rows, rowValues)
	}

	for _, row := range rows {
		if _, err := fmt.Fprintf(t.tw, t.rowFormat, row...); err != nil {
			return err
		}
	}

	if err := t.tw.Flush(); err != nil {
		return err
	}

	if lines, ok := t.nonTableLines[-1]; ok {
		fmt.Fprintln(t.dest, strings.Join(lines, "\n"))
	}

	s := bufio.NewScanner(strings.NewReader(t.buf.String()))
	i := 0
	for ; s.Scan(); i++ {
		str := s.Text()
		// chars := []string{}
		// for c := 0; c < len(str); c++ {
		// 	chars = append(chars, string(str[c]))
		// }
		fmt.Fprintln(t.dest, str)
		// fmt.Println(strings.Join(chars, " "))
		if lines, ok := t.nonTableLines[i]; ok {
			fmt.Fprintln(t.dest, strings.Join(lines, "\n"))
		}
	}

	if lines, ok := t.nonTableLines[i]; ok {
		fmt.Fprintln(t.dest, strings.Join(lines, "\n"))
	}

	return nil
}

type Row []Cell

func (r Row) Render(formatted bool, truncate int) []RenderedCell {
	if truncate < 0 {
		truncate = 0
	}

	rendered := make([]RenderedCell, len(r))
	for i, col := range r {
		if i == 0 {
			rendered[i] = col.Render(formatted, truncate)
		} else {
			rendered[i] = col.Render(formatted, 0)
		}
	}

	return rendered
}

type Cell struct {
	Format string
	Values []ColorableCellValue
}

// Render returns the string representation of the cell with any colors
// applied, and the total overhead in bytes added by the ANSI escape sequences.
func (c Cell) Render(formatted bool, truncate int) RenderedCell {
	values := []interface{}{}
	totalOverhead := 0
	trimmed := false

	value := ""
	overhead := 0
	for _, v := range c.Values {
		if truncate > 0 && !trimmed {
			value, overhead = v.Render(formatted, truncate)
			trimmed = true
		} else {
			value, overhead = v.Render(formatted, 0)
		}

		values = append(values, value)
		totalOverhead += overhead
	}

	return RenderedCell{
		value:    fmt.Sprintf(c.Format, values...),
		overhead: totalOverhead,
	}
}

type RenderedCell struct {
	value    string
	overhead int
}

func (c RenderedCell) AdjustOverhead(delta int) RenderedCell {
	if delta == 0 {
		return c
	}

	v := c.value
	// find the first escape sequence
	regex := regexp.MustCompile(`^(.*?)(\x1b\[)([0-9;]*m)(.*)$`)

	// pad with zeros as needed
	repStr := fmt.Sprintf("${1}${2}%s${3}${4}", strings.Repeat("0", delta))
	v = regex.ReplaceAllString(v, repStr)

	cell := RenderedCell{
		value:    v,
		overhead: c.overhead + int(delta),
	}

	return cell
}

type ColorableCellValue struct {
	Value     interface{}
	Colors    []color.Attribute
	Trimmable bool
}

// Render returns the string representation of the cell value with any colors
// applied, and the total overhead in bytes added by the ANSI escape sequences.
func (v *ColorableCellValue) Render(formatted bool, trim int) (string, int) {
	unformatted := fmt.Sprint(v.Value)

	if trim > 0 {
		trimAmount := math.Min(float64(trim), float64(len(unformatted)))
		unformatted = unformatted[:len(unformatted)-int(trimAmount)]
	}

	if formatted {
		if len(v.Colors) == 0 {
			v.Colors = append(v.Colors, color.Reset)
		}

		formatted := color.New(v.Colors...).Sprint(unformatted)
		overhead := len(formatted) - len(unformatted)
		return formatted, overhead
	}

	return unformatted, 0
}
