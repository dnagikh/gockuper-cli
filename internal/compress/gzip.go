package compress

import (
	"compress/gzip"
	"io"
)

type GzipCompressor struct{}

func (c *GzipCompressor) WrapWriter(dst io.Writer) (io.WriteCloser, error) {
	return gzip.NewWriter(dst), nil
}

func (c *GzipCompressor) Extension() string {
	return "gz"
}
