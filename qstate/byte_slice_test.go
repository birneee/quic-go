package qstate

import (
	"crypto/rand"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestByteSlice_MarshalJSON(t *testing.T) {
	bs := make(HexByteSlice, 20)
	rand.Read(bs)
	j, err := json.Marshal(&bs)
	require.NoError(t, err)
	var bs2 HexByteSlice
	err = json.Unmarshal(j, &bs2)
	require.NoError(t, err)
	require.Equal(t, bs, bs2)
}