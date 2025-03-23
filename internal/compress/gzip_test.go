package compress

import (
	"compress/gzip"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGzipCompressAndDecompress(t *testing.T) {
	inputData := "hello world"
	input := strings.NewReader(inputData)

	comp := &GzipCompressor{}
	reader, err := Compress(input, comp)
	require.NoError(t, err)
	require.NotNil(t, reader)

	gzReader, err := gzip.NewReader(reader)
	require.NoError(t, err)
	defer gzReader.Close()

	outputData, err := io.ReadAll(gzReader)
	require.NoError(t, err)
	require.Equal(t, inputData, string(outputData))
}

type BrokenCompressor struct{}

func (b *BrokenCompressor) WrapWriter(w io.Writer) (io.WriteCloser, error) {
	return nil, fmt.Errorf("forced error")
}

func (b *BrokenCompressor) Extension() string {
	return "broken"
}

func TestCompress_CompressorError(t *testing.T) {
	input := strings.NewReader("test")
	comp := &BrokenCompressor{}

	reader, err := Compress(input, comp)
	require.NoError(t, err)

	_, err = io.ReadAll(reader)
	require.Error(t, err)
	require.Contains(t, err.Error(), "forced error")
}

type ErrorReader struct{}

func (e *ErrorReader) Read(p []byte) (int, error) {
	return 0, fmt.Errorf("reader exploded")
}

func TestCompress_IOCopyError(t *testing.T) {
	input := &ErrorReader{}
	comp := &GzipCompressor{}

	reader, err := Compress(input, comp)
	require.NoError(t, err)

	_, err = io.ReadAll(reader)
	require.Error(t, err)
	require.Contains(t, err.Error(), "reader exploded")
}
