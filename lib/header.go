package seal

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"strconv"
)

// Returns the number of bytes in a header.
func headerLen(bytes int) int {
	// Hex encoding requires two bytes per byte, so we multiply by 2.
	return IdentLen + len("{}\n") + bytes*2
}

// Parses the header of a seal file. Does not read beyond the
// header.
func parseHeader(in *bufio.Reader) (*Seal, error) {

	// TODO: Stop reading and error out if we pass headerLen(maxBytes).
	header, err := in.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	sl := &Seal{}

	sl.Magic = string(header[:len(Magic)]) // example: `SL%v`
	if sl.Magic != Magic {
		return nil, fmt.Errorf("seal: bad magic number: %q", sl.Magic)
	}

	uint64Version, err := strconv.ParseUint(
		string(header[len(Magic)]), 10, 8) // example: `0`
	if err != nil {
		return nil, fmt.Errorf("seal: couldn't parse version: %v", err)
	}

	sl.Version = int(uint64Version)
	if sl.Version != Version {
		return nil, fmt.Errorf("seal: unsupported version: %v", sl.Version)
	}

	sig := header[IdentLen : len(header)-1]

	if len(sig) <= 2 || (sig[0] != '{' && sig[len(sig)-1] != '}') {
		return nil, fmt.Errorf("seal: invalid signature")
	}
	sig = sig[1 : len(sig)-1]

	sl.ClaimedSignature, err = hex.DecodeString(string(sig))
	if err != nil {
		return nil, fmt.Errorf("seal: couldn't decode signature: %v", err)
	}

	return sl, err
}
