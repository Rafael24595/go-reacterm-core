package local

import (
	"io"

	"github.com/Rafael24595/go-log/log"
	"github.com/Rafael24595/go-log/log/record"
)

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

func Log(f func() error) {
	Writer(log.WriterFromCategory(record.ERROR), f)
}
