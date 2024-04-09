package outformats

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
)

// JSONWriter JSON format writer
type JSONWriter struct {
	fields []string
	w      *bufio.Writer
}

// Write writes a single JSON record to w one line.
// A record is a slice of strings with each string being one field.
// Writes are buffered, so Flush must eventually be called to ensure
// that the record is written to the underlying io.Writer.
func (w *JSONWriter) Write(records []string) error {
	if len(records) != len(w.fields) {
		return errors.New("records length is not equal to fields")
	}

	m := make(map[string]string)
	for i := range w.fields {
		m[w.fields[i]] = records[i]
	}
	d, err := json.Marshal(m)
	if err != nil {
		return err
	}

	if _, err := w.w.Write(d); err != nil {
		return err
	}
	if _, err := w.w.WriteString("\n"); err != nil {
		return err
	}

	return nil
}

// WriteAll writes multiple json records to w using Write and then calls Flush,
// returning any error from the Flush.
func (w *JSONWriter) WriteAll(records [][]string) error {
	for _, record := range records {
		err := w.Write(record)
		if err != nil {
			return err
		}
	}
	return w.w.Flush()
}

func (w *JSONWriter) Flush() {
	w.w.Flush()
}

// NewJSONWriter generate json writer
// fields are key field
func NewJSONWriter(w io.Writer, fields []string) *JSONWriter {
	return &JSONWriter{
		w:      bufio.NewWriter(w),
		fields: fields,
	}
}
