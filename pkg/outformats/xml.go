package outformats

import (
	"bufio"
	"encoding/xml"
	"errors"
	"io"
)

// XMLWriter XML format writer
type XMLWriter struct {
	fields []string
	w      *bufio.Writer
}

// Write writes a single JSON record to w one line.
// A record is a slice of strings with each string being one field.
// Writes are buffered, so Flush must eventually be called to ensure
// that the record is written to the underlying io.Writer.
func (w *XMLWriter) Write(records []string) error {
	if len(records) != len(w.fields) {
		return errors.New("records length is not equal to fields")
	}

	m := make(result)
	for i := range w.fields {
		m[w.fields[i]] = records[i]
	}
	d, err := xml.Marshal(m)
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
func (w *XMLWriter) WriteAll(records [][]string) error {
	for _, record := range records {
		err := w.Write(record)
		if err != nil {
			return err
		}
	}
	return w.w.Flush()
}

func (w *XMLWriter) Flush() {
	w.w.Flush()
}

// StringMap is a map[string]string.
type result map[string]string

// StringMap marshals into XML.
func (s result) MarshalXML(e *xml.Encoder, start xml.StartElement) error {

	tokens := []xml.Token{start}

	for key, value := range s {
		t := xml.StartElement{Name: xml.Name{"", key}}
		tokens = append(tokens, t, xml.CharData(value), xml.EndElement{t.Name})
	}

	tokens = append(tokens, xml.EndElement{start.Name})

	for _, t := range tokens {
		err := e.EncodeToken(t)
		if err != nil {
			return err
		}
	}

	// flush to ensure tokens are written
	return e.Flush()
}

// NewXMLWriter generate xml writer
// fields are key field
func NewXMLWriter(w io.Writer, fields []string) *XMLWriter {
	return &XMLWriter{
		w:      bufio.NewWriter(w),
		fields: fields,
	}
}
