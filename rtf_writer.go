package main

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

// rtfWriter converts ANSI-colored text into a minimal RTF document preserving colors.
type rtfWriter struct {
	writer  io.Writer
	started bool
	closed  bool
	color   int
	bold    bool
}

func newRTFWriter(w io.Writer) *rtfWriter {
	return &rtfWriter{writer: w}
}

func (r *rtfWriter) Write(p []byte) (int, error) {
	if r.closed {
		return 0, fmt.Errorf("write on closed rtf writer")
	}

	if !r.started {
		if err := r.start(); err != nil {
			return 0, err
		}
	}

	content := string(p)
	var builder strings.Builder

	for i := 0; i < len(content); {
		if content[i] == 0x1b && i+1 < len(content) && content[i+1] == '[' {
			end := strings.IndexByte(content[i:], 'm')
			if end == -1 {
				i++
				continue
			}

			sequence := content[i+2 : i+end]
			r.handleEscape(&builder, sequence)
			i += end + 1
			continue
		}

		runeValue, size := utf8.DecodeRuneInString(content[i:])
		r.writeEscapedRune(&builder, runeValue)
		i += size
	}

	_, err := r.writer.Write([]byte(builder.String()))
	if err != nil {
		return 0, err
	}

	return len(p), nil
}

func (r *rtfWriter) Close() error {
	if r.closed {
		return nil
	}

	r.closed = true

	if !r.started {
		if err := r.start(); err != nil {
			return err
		}
	}

	_, err := r.writer.Write([]byte("\\line\n}"))
	return err
}

func (r *rtfWriter) start() error {
	r.started = true

	header := "{\\rtf1\\ansi\\deff0\n{\\colortbl ;\\red255\\green0\\blue0;\\red0\\green128\\blue0;\\red255\\green255\\blue0;\\red0\\green0\\blue255;\\red255\\green0\\blue255;\\red0\\green255\\blue255;\\red255\\green255\\blue255;}\n"

	_, err := r.writer.Write([]byte(header))
	return err
}

func (r *rtfWriter) handleEscape(builder *strings.Builder, sequence string) {
	for _, part := range strings.Split(sequence, ";") {
		switch part {
		case "0":
			r.reset(builder)
		case "1":
			r.setBold(builder, true)
		case "31", "32", "33", "34", "35", "36", "37":
			colorIndex := int(part[1] - '0')
			r.setColor(builder, colorIndex)
		}
	}
}

func (r *rtfWriter) reset(builder *strings.Builder) {
	if r.bold {
		builder.WriteString("\\b0 ")
	}
	r.bold = false
	r.color = 0
	builder.WriteString("\\cf0 ")
}

func (r *rtfWriter) setBold(builder *strings.Builder, enable bool) {
	if r.bold == enable {
		return
	}

	r.bold = enable
	if enable {
		builder.WriteString("\\b ")
	} else {
		builder.WriteString("\\b0 ")
	}
}

func (r *rtfWriter) setColor(builder *strings.Builder, color int) {
	if r.color == color {
		return
	}

	r.color = color
	builder.WriteString(fmt.Sprintf("\\cf%d ", color))
}

func (r *rtfWriter) writeEscapedRune(builder *strings.Builder, rne rune) {
	switch rne {
	case '\\':
		builder.WriteString("\\\\")
	case '{':
		builder.WriteString("\\{")
	case '}':
		builder.WriteString("\\}")
	case '\n':
		builder.WriteString("\\line\n")
	case '\r':
		// Ignore carriage returns
	default:
		builder.WriteRune(rne)
	}
}

// dualWriter mirrors every write to both the console and the RTF writer.
type dualWriter struct {
	console io.Writer
	rtf     *rtfWriter
}

func (d dualWriter) Write(p []byte) (int, error) {
	if d.console != nil {
		if _, err := d.console.Write(p); err != nil {
			return 0, err
		}
	}

	if d.rtf != nil {
		if _, err := d.rtf.Write(p); err != nil {
			return 0, err
		}
	}

	return len(p), nil
}
