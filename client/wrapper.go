package main

import (
	"io"
)

//ReaderWrapper is a reader...
type ReaderWrapper struct {
	io.Reader
	reader io.Reader
}

//NewReaderWrapper constructor...
func NewReaderWrapper(reader io.Reader) *ReaderWrapper {
	r := new(ReaderWrapper)

	r.reader = reader

	return r
}

func (w *ReaderWrapper) Read(p []byte) (n int, err error) {
	n, err = w.reader.Read(p)

	for i := range p[:n] {
		p[i] = ((0x0F & p[i]) << 4) | ((0xF0 & p[i]) >> 4)
	}

	// log.Printf("[ReaderWrapper.Read] : %s\n", hex.EncodeToString(p[:n]))

	return n, err
}

//WriterWrapper is a writer...
type WriterWrapper struct {
	io.Writer
	writer io.Writer
}

//NewWriterWrapper constructor...
func NewWriterWrapper(writer io.Writer) *WriterWrapper {
	r := new(WriterWrapper)

	r.writer = writer

	return r
}

func (w *WriterWrapper) Write(p []byte) (n int, err error) {
	for i := range p {
		p[i] = ((0x0F & p[i]) << 4) | ((0xF0 & p[i]) >> 4)
	}

	n, err = w.writer.Write(p)

	// log.Printf("[WriterWrapper.Write] : %s\n", hex.EncodeToString(p[:n]))

	return n, err
}
