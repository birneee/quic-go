//go:generate msgp
package qstate

import (
	"encoding/hex"
	"encoding/json"
)

type StatelessResetToken [16]byte

func (s *StatelessResetToken) MarshalJSON() ([]byte, error) {
	return json.Marshal(hex.EncodeToString((*s)[:]))
}

func (s *StatelessResetToken) UnmarshalJSON(b []byte) error {
	_, err := hex.Decode(s[:], b[1:len(b)-1])
	if err != nil {
		return err
	}
	return nil
}
