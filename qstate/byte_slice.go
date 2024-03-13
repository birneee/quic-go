//go:generate msgp
package qstate

import (
	"encoding/hex"
	"encoding/json"
)

type HexByteSlice []byte

func (s *HexByteSlice) MarshalJSON() ([]byte, error) {
	return json.Marshal(hex.EncodeToString(*s))
}

func (s *HexByteSlice) UnmarshalJSON(b []byte) error {
	*s = make(HexByteSlice, (len(b)-2)/2)
	_, err := hex.Decode(*s, b[1:len(b)-1])
	if err != nil {
		return err
	}
	return nil
}
