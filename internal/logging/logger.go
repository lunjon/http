package logging

import (
	"io"
	"log"
	"os"
)

func NewLogger() *log.Logger {
	return newLogger(os.Stdout)
}

func NewSilentLogger() *log.Logger {
	return newLogger(io.Discard)
}

func newLogger(w io.Writer) *log.Logger {
	return log.New(w, "", log.Ldate|log.Ltime)
}
