package main

import "github.com/klauspost/compress/zstd"

func decompressZstd(compressedData []byte) (decompressedData []byte, resultCode int) {
	decoder, err := zstd.NewReader(nil)
	if err != nil {
		return nil, handleError([2]byte{0x00, 0x03}, err)
	}
	defer decoder.Close()
	decompressedData, err = decoder.DecodeAll(compressedData, nil)
	if err != nil {
		return nil, handleError([2]byte{0x00, 0x03}, err)
	}
	return decompressedData, 0
}
