package qstate

import (
	"crypto/rand"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStatelessResetToken_MarshalJSON(t *testing.T) {
	var srt StatelessResetToken
	rand.Read(srt[:])
	j, err := json.Marshal(&srt)
	require.NoError(t, err)
	var srt2 StatelessResetToken
	err = json.Unmarshal(j, &srt2)
	require.NoError(t, err)
	require.Equal(t, srt, srt2)
}
