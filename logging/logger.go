package logging

import (
	"io"
	"log"
	"os"
)

type LogWriter struct {
	dst io.Writer
}

func (w *LogWriter) Write(b []byte) (int, error) {
	return w.dst.Write(b)
}

func NewLogger() *log.Logger {
	w := &LogWriter{
		dst: os.Stdout,
	}
	logger := log.New(w, "", 0)
	return logger
}
