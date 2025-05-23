package fastjson

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"
	"unicode/utf8"
)

func TestWriterReset(t *testing.T) {
	var w Writer
	w.String("foo")
	capBefore := cap(w.Bytes())
	w.Reset()
	if len(w.Bytes()) != 0 {
		t.Fatalf("expected w.Bytes() to be reset")
	}
	if cap(w.Bytes()) != capBefore {
		t.Fatalf("expected w.Bytes() capacity to be retained")
	}
}

func TestWriterRewind(t *testing.T) {
	var w Writer
	w.String("foo")
	before := w.Size()
	w.RawByte(',')
	w.String("bar")
	assertEncoded(t, &w, `"foo","bar"`)

	w.Rewind(before)
	assertEncoded(t, &w, `"foo"`)
}

func TestWriterTime(t *testing.T) {
	var w Writer
	w.Time(time.Unix(0, 0).UTC(), time.RFC1123Z)
	assertEncoded(t, &w, `Thu, 01 Jan 1970 00:00:00 +0000`)
}

func TestWriterString(t *testing.T) {
	var i rune
	for i = 0; i < utf8.MaxRune; i++ {
		s := string(i)
		var w Writer
		w.String(s)

		expected, err := json.Marshal(s)
		if err != nil {
			t.Errorf("json.Marshal returned unexpected error: %v", err)
		}

		got := w.Bytes()
		if !bytes.Equal(expected, got) {
			t.Errorf("rune %d: expected '%s', got '%s'", i, expected, got)
		}
	}
}

func TestWriterStringEscapes(t *testing.T) {
	var w Writer
	w.StringContents("\t\r\n\\\"\x00")
	assertEncoded(t, &w, `\t\r\n\\\"\u0000`)

	w.Reset()
	w.StringContents("\u2028\u2029")
	assertEncoded(t, &w, `\u2028\u2029`)

	w.Reset()
	w.StringContents("世界")
	assertEncoded(t, &w, "世界")

	w.Reset()
	w.StringContents(string([]byte{255}))
	assertEncoded(t, &w, `\ufffd`)
}

func BenchmarkString(b *testing.B) {
	w := Writer{
		buf: make([]byte, 0, 1024),
	}
	b.Run("string", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			w.Reset()
			w.String("foobarstring-go-fastjson")
		}
	})
}

func BenchmarkFloat(b *testing.B) {
	w := Writer{
		buf: make([]byte, 0, 1024),
	}
	b.Run("string", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			w.Reset()
			Marshal(&w, 3.14+float64(i))
		}
	})
}
