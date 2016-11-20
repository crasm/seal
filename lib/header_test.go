package seal

import (
	"bufio"
	"bytes"
	"reflect"
	"testing"
)

var cafeClaim = []byte{0xca, 0xfe}
var cafeHeader = []byte("SL%v0{cafe}\n")
var cafeSeal = &Seal{
	Magic:            `SL%v`,
	Version:          0,
	ClaimedSignature: cafeClaim,
}

func TestHeaderLen(t *testing.T) {
	result := headerLen(2)
	expected := len(cafeHeader)
	if result != expected {
		t.Errorf("expected a header length of %v but got %v", expected, result)
	}
}

func TestCreateHeader(t *testing.T) {
	result := createHeader(cafeClaim)
	expected := cafeHeader

	if !bytes.Equal(result, expected) {
		t.Errorf("expected %q but got %q", expected, result)
	}
}

func TestParseHeaderValid(t *testing.T) {
	cases := []struct {
		header []byte
		seal   *Seal
	}{
		{
			header: cafeHeader,
			seal:   cafeSeal,
		},
	}

	for _, c := range cases {
		generated, err := parseHeader(makeReader(c.header))
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(c.seal, generated) {
			t.Errorf("expected %q but got %q", c.seal, generated)
		}
	}
}

func TestParseHeaderInvalid(t *testing.T) {
	headers := []string{
		"SL%v0{cafe}",       // no newline
		"SL%q\n",            // bad magic number
		"SL%v!\n",           // bad version
		"SL%v8\n",           // unsupported version
		"SL%v0\n",           // missing claim
		"SL%v0cafe\n",       // missing braces
		"SL%v0}cafe{\n",     // backwards braces
		"SL%v0{cafe\n",      // unmatched brace
		"SL%v0{  cafe  }\n", // spaces
	}

	for i, h := range headers {
		sl, err := parseHeader(makeReader([]byte(h)))
		if sl != nil {
			t.Errorf("header %d: erroneously returned a Seal: %q", i, sl)
		}
		if err == nil {
			t.Errorf("header %d: no error returned, but should have", i)
		}
	}

}

func makeReader(data []byte) *bufio.Reader {
	return bufio.NewReader(bytes.NewReader(data))
}
