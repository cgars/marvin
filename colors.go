package main

import (
	"bytes"
	"fmt"

	irc "github.com/fluffle/goirc/client"
)

const (
	cPref = '\u0003'
)

type ColorCode byte

const (
	White = ColorCode(iota)
	Black
	Blue
	Green
	LightRed
	Brown
	Purple
	Orange
	Yellow
	LightGreen
	Cyan
	LightCyan
	LightBlue
	Pink
	Grey
	LightGrey
)

type TextBuffer struct {
	buf     *bytes.Buffer
	closing string
}

func (b *TextBuffer) String() string {
	return b.buf.String()
}

func Text(format string, args ...interface{}) *TextBuffer {
	buf := bytes.NewBufferString(fmt.Sprintf(format, args...))
	return &TextBuffer{buf: buf}
}

func (b *TextBuffer) Text(format string, args ...interface{}) *TextBuffer {
	b.buf.WriteString(fmt.Sprintf(format, args...))
	b.maybeFinishClosing()
	return b
}

func (b *TextBuffer) Fg(c ColorCode, format string, args ...interface{}) *TextBuffer {
	b.buf.WriteRune(cPref)
	ccStr := fmt.Sprintf("%02d", int(c))
	b.buf.WriteString(ccStr)
	b.Text(format, args...)
	b.buf.WriteRune(cPref)
	b.maybeFinishClosing()
	return b
}

func (b *TextBuffer) maybeFinishClosing() {
	if b.closing != "" {
		b.buf.WriteString(b.closing)
		b.closing = ""
	}
}

func (b *TextBuffer) S(s, e string) *TextBuffer {
	b.buf.WriteString(s)
	b.closing = e
	return b
}

func (b *TextBuffer) Send(c *irc.Conn, t string) {
	b.maybeFinishClosing()
	b.buf.WriteString(b.closing)
	b.buf.WriteString("\n")
	c.Privmsg(t, b.String())
}
