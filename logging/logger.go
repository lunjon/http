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
	return log.New(os.Stdout, "", 0)
}
