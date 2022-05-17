package outformats

import (
	"encoding/csv"
	"io"
)

type CSVWriter struct {
	*csv.Writer
}

func NewCSVWriter(w io.Writer) *CSVWriter {
	return &CSVWriter{
		Writer: csv.NewWriter(w),
	}
}
