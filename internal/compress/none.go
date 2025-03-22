package compress

import "io"

type nopWriteCloser struct {
	io.Writer
}

func (n nopWriteCloser) Close() error { return nil }

type NoneCompressor struct {
	io.Writer
}

func (c *NoneCompressor) WrapWriter(dst io.Writer) (io.WriteCloser, error) {
	return nopWriteCloser{dst}, nil
}

func (c *NoneCompressor) Extension() string {
	return "dump"
}
