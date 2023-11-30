package model

import (
	"bytes"
	"hash/crc32"
	"io"
)

type Image struct {
	Blob  []byte
	CRC32 uint32
}

func NewImage(r io.Reader) (*Image, error) {
	icon := &Image{}

	raw, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	icon.Blob = raw
	icon.CRC32 = crc32.ChecksumIEEE(raw)

	return icon, nil
}

func (i *Image) Reader() io.Reader {
	return bytes.NewReader(i.Blob)
}
