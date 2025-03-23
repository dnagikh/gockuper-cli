package compress

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNoneCompressor(t *testing.T) {
	data := "information without compression"
	input := strings.NewReader(data)

	comp := &NoneCompressor{}
	reader, err := Compress(input, comp)
	require.NoError(t, err)

	result, err := io.ReadAll(reader)
	require.NoError(t, err)
	require.Equal(t, data, string(result))
}
