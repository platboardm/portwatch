package sink_test

import (
	"bytes"
	"context"
	"errors"
	"log"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/sink"
)

func newLogger(buf *bytes.Buffer) *log.Logger {
	return log.New(buf, "", 0)
}

func TestNewMulti_PanicsOnNilSink(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on nil Sink")
		}
	}()
	sink.NewMulti(nil, newLogger(new(bytes.Buffer)))
}

func TestNewMulti_PanicsOnNilLogger(t *testing.T) {
	s := sink.New(map[string]sink.Sender{"a": &stubSender{}})
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on nil logger")
		}
	}()
	sink.NewMulti(s, nil)
}

func TestMulti_Send_LogsErrors(t *testing.T) {
	var buf bytes.Buffer
	s := sink.New(map[string]sink.Sender{
		"failing": &stubSender{err: errors.New("net down")},
	})
	m := sink.NewMulti(s, newLogger(&buf))
	m.Send(context.Background(), makeAlert())

	got := buf.String()
	if !strings.Contains(got, "failing") {
		t.Errorf("expected log to contain sender name, got: %q", got)
	}
	if !strings.Contains(got, "net down") {
		t.Errorf("expected log to contain error text, got: %q", got)
	}
}

func TestMulti_Send_SilentOnSuccess(t *testing.T) {
	var buf bytes.Buffer
	s := sink.New(map[string]sink.Sender{
		"ok": &stubSender{},
	})
	m := sink.NewMulti(s, newLogger(&buf))
	m.Send(context.Background(), makeAlert())

	if buf.Len() != 0 {
		t.Errorf("expected no log output on success, got: %q", buf.String())
	}
}
