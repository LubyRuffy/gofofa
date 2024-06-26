package outformats

import (
	"encoding/csv"
	"io"
)

// CSVWriter CSV format writer
type CSVWriter struct {
	*csv.Writer
}

func (w *CSVWriter) Flush() {
	w.Writer.Flush()
}

// NewCSVWriter generate CSVWriter
func NewCSVWriter(w io.Writer) *CSVWriter {
	return &CSVWriter{
		Writer: csv.NewWriter(w),
	}
}
