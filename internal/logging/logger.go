package logging

import (
	"io"
	"log"
	"os"
)

func NewLogger() *log.Logger {
	return New(os.Stdout)
}

func NewSilentLogger() *log.Logger {
	return New(io.Discard)
}

func New(w io.Writer) *log.Logger {
	return log.New(w, "", log.Ldate|log.Ltime)
}
