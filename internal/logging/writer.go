package logging

import (
	"io"
	"os"
)

type LogWriter struct {
	dst io.Writer
}

func NewStdoutWriter() *LogWriter {
	return &LogWriter{
		dst: os.Stdout,
	}
}

func (w *LogWriter) Write(b []byte) (int, error) {
	return w.dst.Write(b)
}
