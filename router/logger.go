package router

import (
	"io"
	"log"
	"os"
)

type Logger interface {
	Print(...interface{})
	Printf(string, ...interface{})
	Println(...interface{})
}

func NewLogger(out io.Writer) Logger {
	if out == nil {
		out = os.Stderr
	}
	return log.New(out, "", log.LstdFlags)
}
