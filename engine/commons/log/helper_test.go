package local

import (
	"bytes"
	"errors"
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

type failingWriter struct{}

func (f failingWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("writer failure")
}

func TestWriter_NoError(t *testing.T) {
	var buf bytes.Buffer
	executed := false

	Writer(&buf, func() error {
		executed = true
		return nil
	})

	assert.True(t, executed)
	assert.Equal(t, 0, buf.Len())
}

func TestWriter_WritesErrorToWriter(t *testing.T) {
	var buf bytes.Buffer

	Writer(&buf, func() error {
		return errors.New("something went wrong")
	})

	assert.Inside(t, "something went wrong", buf.String())
}

func TestWriter_FallbackOnWriterFailure(t *testing.T) {
	assert.NotPanic(t, func() {
		Writer(failingWriter{}, func() error {
			return errors.New("primary error")
		})
	})
}
