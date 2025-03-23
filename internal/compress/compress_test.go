package compress

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFromString(t *testing.T) {
	tests := []struct {
		name      string
		expectErr bool
		expectExt string
	}{
		{"gzip", false, "gz"},
		{"GZIP", false, "gz"},
		{"none", false, "dump"},
		{"", false, "dump"},
		{"something", true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comp, err := FromString(tt.name)
			if tt.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, comp)
			require.Equal(t, tt.expectExt, comp.Extension())
		})
	}
}
