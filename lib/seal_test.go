package seal

import (
	"bytes"
	"testing"
)

var cases = []struct {
	data, seal string
}{
	{
		data: "",
		seal: "SL%v0{cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e}\n",
	}, {
		data: "seal!\n",
		seal: "SL%v0{0d406d27279f9e9ff7dd349f49069c5dba677e013e5b5c9c1d857f9e560155bf02573d4e0275ee3ccbb60e2a7b84b6837d01152a995b3189fc5243b6ed471f94}\n",
	},
}

func TestGenerateBasic(t *testing.T) {
	for _, c := range cases {
		s, err := Generate(bytes.NewBufferString(c.data))
		if err != nil {
			t.Fatal(err)
		}

		if c.seal != s.String() {
			t.Errorf("Expected %s, got %s", c.seal, s.String())
		}
	}
}

func TestValidateBasic(t *testing.T) {
	// TODO: create Seal structs and use them to validate against data.
	// May require an API update.
	t.Fatal("Not implemented")
}
