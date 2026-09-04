package local

import (
	"io"

	"github.com/Rafael24595/go-log/log"
	"github.com/Rafael24595/go-log/log/record"
)

// Writer executes the function f and writes any returned error to the provided io.Writer.
// If writing to the writer also fails, it falls back to standard print output (println).
func Writer(w io.Writer, f func() error) {
	err := f()
	if err == nil {
		return
	}

	_, err = w.Write([]byte(err.Error()))
	if err == nil {
		return
	}

	println(err.Error())
}

// Log executes the function f and writes any returned error to the ERROR category logger writer.
func Log(f func() error) {
	Writer(log.WriterFromCategory(record.ERROR), f)
}
