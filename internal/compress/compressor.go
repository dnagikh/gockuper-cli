package compress

import (
	"fmt"
	"io"
	"strings"
)

type Compressor interface {
	WrapWriter(dst io.Writer) (io.WriteCloser, error)
	Extension() string
}

func FromString(name string) (Compressor, error) {
	switch strings.ToLower(name) {
	case "", "none":
		return &NoneCompressor{}, nil
	case "gzip":
		return &GzipCompressor{}, nil
	default:
		return nil, fmt.Errorf("unknown compressor: %s", name)
	}
}

func Compress(reader io.Reader, c Compressor) (io.Reader, error) {
	pr, pw := io.Pipe()

	go func() {
		defer pw.Close()
		writer, err := c.WrapWriter(pw)
		if err != nil {
			pw.CloseWithError(fmt.Errorf("could not compress pg_dump: %w", err))
		}
		defer writer.Close()

		_, err = io.Copy(writer, reader)
		if err != nil {
			pw.CloseWithError(fmt.Errorf("could not compress pg_dump: %w", err))
		}
	}()

	return pr, nil
}
