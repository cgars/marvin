package main

import (
	"bytes"
	"fmt"
	"strconv"

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
	buf *bytes.Buffer
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
	return b
}

func (b *TextBuffer) Fg(c ColorCode, format string, args ...interface{}) *TextBuffer {
	b.buf.WriteRune(cPref)
	b.buf.WriteString(strconv.Itoa(int(c)))
	b.Text(format, args...)
	b.buf.WriteRune(cPref)
	return b
}

func (b *TextBuffer) Send(c *irc.Conn, t string) {
	c.Privmsg(t, b.String())
}
