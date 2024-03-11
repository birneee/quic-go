//go:generate msgp
package qstate

import (
	"encoding/hex"
	"encoding/json"
)

type ByteSlice []byte

func (s *ByteSlice) MarshalJSON() ([]byte, error) {
	return json.Marshal(hex.EncodeToString(*s))
}

func (s *ByteSlice) UnmarshalJSON(b []byte) error {
	*s = make(ByteSlice, (len(b)-2)/2)
	_, err := hex.Decode(*s, b[1:len(b)-1])
	if err != nil {
		return err
	}
	return nil
}
