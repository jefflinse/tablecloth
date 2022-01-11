package tablecloth

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
	"text/tabwriter"
)

// A Table is a set of rows each containing cells.
type Table struct {
	columns       []ColumnDefinition
	rowFormat     string
	rows          []Row
	overheads     []int
	nonTableLines map[int][]string
	currentLine   int
	buf           *strings.Builder
	tw            *tabwriter.Writer

	Debug bool
}

// NewTableWithColumns creates a new table with the specified columns.
func NewTableWithColumns(columns []ColumnDefinition) *Table {
	buf := &strings.Builder{}
	return &Table{
		columns:       columns,
		rowFormat:     strings.Repeat("%s\t", len(columns)) + "\n",
		rows:          []Row{},
		overheads:     make([]int, len(columns)),
		nonTableLines: map[int][]string{},
		currentLine:   -1,
		buf:           buf,
		tw:            tabwriter.NewWriter(buf, 0, 0, 2, ' ', tabwriter.Debug),
	}
}

// NewTable creates a new table with the specified number of columns.
func NewTable(columns int) *Table {
	return NewTableWithColumns(make([]ColumnDefinition, columns))
}

// A ColumnDefinition defines a column in a table.
type ColumnDefinition struct {
	Name      string
	MinLength int
}

// AddRow adds a row to the table.
func (t *Table) AddRow(row ...Cell) {
	rendered := Row(row).Render()
	for i := range rendered {
		if rendered[i].overhead > t.overheads[i] {
			t.overheads[i] = rendered[i].overhead
		}
	}
	t.rows = append(t.rows, row)
	t.currentLine++
}

// AddLine adds a spanning line to the table.
//
// Spanning lines can appear between table rows and are rendered as-is.
func (t *Table) AddLine(line string) {
	t.nonTableLines[t.currentLine] = append(t.nonTableLines[t.currentLine], line)
}

// Write writes the table to the given io.Writer.
func (t *Table) Write(w io.Writer) error {
	rows := [][]interface{}{}

	for r := range t.rows {
		rowValues := []interface{}{}
		cells := t.rows[r].Render()
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
		fmt.Fprintln(w, strings.Join(lines, "\n"))
	}

	s := bufio.NewScanner(strings.NewReader(t.buf.String()))
	i := 0
	for ; s.Scan(); i++ {
		str := s.Text()
		// chars := []string{}
		// for c := 0; c < len(str); c++ {
		// 	chars = append(chars, string(str[c]))
		// }
		fmt.Fprintln(w, str)
		// fmt.Println(strings.Join(chars, " "))
		if lines, ok := t.nonTableLines[i]; ok {
			fmt.Fprintln(w, strings.Join(lines, "\n"))
		}
	}

	if lines, ok := t.nonTableLines[i]; ok {
		fmt.Fprintln(w, strings.Join(lines, "\n"))
	}

	return nil
}

// A Row is a set of Cells.
type Row []Cell

// Render returns a set of RenderedCells for the row.
func (r Row) Render() []RenderedCell {
	rendered := make([]RenderedCell, len(r))
	for i, col := range r {
		rendered[i] = col.Render()
	}

	return rendered
}

// A Cell is a single table cell whose value is comprised of a format string
// and zero or more colorable values to be formatted into the string.
type Cell struct {
	Format string
	Values []FormattableCellValue
}

// Render returns the string representation of the cell with any colors
// applied, and the total overhead in bytes added by the ANSI escape sequences.
func (c Cell) Render() RenderedCell {
	values := []interface{}{}
	totalOverhead := 0

	value := ""
	overhead := 0
	for _, v := range c.Values {
		value, overhead = v.Render()
		values = append(values, value)
		totalOverhead += overhead
	}

	format := "\x1b[0m" + c.Format
	return RenderedCell{
		value:    fmt.Sprintf(format, values...),
		overhead: totalOverhead,
	}
}

// FormattableCellValue is a value that can be formatted with color.
type FormattableCellValue struct {
	Format func(...interface{}) string
	Value  interface{}
}

// Render returns the string representation of the cell value with any colors
// applied, and the total overhead in bytes added by the ANSI escape sequences.
func (v *FormattableCellValue) Render() (string, int) {
	unformatted := fmt.Sprint(v.Value)
	formatted := v.Format(unformatted)
	overhead := len(formatted) - len(unformatted)
	return formatted, overhead
}

// A RenderedCell is the result of rendering a Cell and contains the rendered
// string value and the total overhead in bytes added by the ANSI escape sequences.
type RenderedCell struct {
	value    string
	overhead int
}

// AdjustOverhead adjusts the overhead of the rendered cell by the specified
// amount. It does so by padding with zeros the first escape sequence found in the cell's
// rendered value.
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
