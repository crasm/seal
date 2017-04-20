package seal

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func decodeHex(raw string) []byte {
	data, _ := hex.DecodeString(raw)
	return data
}

var goodCases = []struct {
	data, header string
	seal         *Seal
	bits         int
}{
	{
		data:   "",
		bits:   512,
		header: "SL%v0{cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e}\n",
		seal: &Seal{
			Magic:            `SL%v`,
			Version:          0,
			ClaimedSignature: decodeHex(`cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e`),
		},
	}, {
		data:   "seal!\n",
		bits:   512,
		header: "SL%v0{0d406d27279f9e9ff7dd349f49069c5dba677e013e5b5c9c1d857f9e560155bf02573d4e0275ee3ccbb60e2a7b84b6837d01152a995b3189fc5243b6ed471f94}\n",
		seal: &Seal{
			Magic:            `SL%v`,
			Version:          0,
			ClaimedSignature: decodeHex(`0d406d27279f9e9ff7dd349f49069c5dba677e013e5b5c9c1d857f9e560155bf02573d4e0275ee3ccbb60e2a7b84b6837d01152a995b3189fc5243b6ed471f94`),
		},
	}, {
		data:   "",
		bits:   8,
		header: "SL%v0{cf}\n",
		seal: &Seal{
			Magic:            `SL%v`,
			Version:          0,
			ClaimedSignature: decodeHex(`cf`),
		},
	}, {
		data:   "seal!\n",
		bits:   8,
		header: "SL%v0{0d}\n",
		seal: &Seal{
			Magic:            `SL%v`,
			Version:          0,
			ClaimedSignature: decodeHex(`0d`),
		},
	},
}

var badCases = []struct {
	data string
	bits int
}{
	{
		data: "seal!\n",
		bits: 0,
	}, {
		data: "seal!\n",
		bits: -1,
	}, {
		data: "seal!\n",
		bits: -8,
	}, {
		data: "seal!\n",
		bits: 9,
	}, {
		data: "seal!\n",
		bits: 513,
	}, {
		data: "seal!\n",
		bits: 520,
	},
}

func TestWrapBuffered(t *testing.T) {
	for _, c := range goodCases {
		file := bytes.NewBufferString(c.data)
		wrapped := &bytes.Buffer{}
		sl, err := WrapBufferedBits(file, wrapped, c.bits)
		require.Nil(t, err)

		assert.Equal(t, c.seal, sl)
		assert.Equal(t, c.header+c.data, wrapped.String())
	}
}

func TestWrapBufferedBad(t *testing.T) {
	for _, c := range badCases {
		file := bytes.NewBufferString(c.data)
		wrapped := &bytes.Buffer{}
		_, err := WrapBufferedBits(file, wrapped, c.bits)
		require.NotNil(t, err)
		assert.Empty(t, wrapped.String())
	}
}

func TestUnwrapGood(t *testing.T) {
	for _, c := range goodCases {
		file := bytes.NewBufferString(c.header + c.data)
		unwrapped := &bytes.Buffer{}
		usl, err := Unwrap(file, unwrapped)
		require.Nil(t, err)

		assert.Equal(t, c.seal, &usl.Seal)
		assert.Equal(t, c.data, unwrapped.String())
	}
}

func TestUnwrapBad(t *testing.T) {
	// TODO: wrong magic number, etc.
	t.Fatal("Not implemented")
}
