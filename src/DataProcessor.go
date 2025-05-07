package main

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"regexp"
)

func processData(decompressedData []byte) (aircrafts []AircraftInfo, resultCode int) {
	const frameSize = 112

	dataSize := len(decompressedData)

	if dataSize%frameSize != 0 {
		return nil, handleError([2]byte{0x00, 0x04}, errors.New("invalid size of response data11"))
	}
	unixTimestamp := binary.LittleEndian.Uint64(decompressedData[0:8])

	framesCount := binary.LittleEndian.Uint16(decompressedData[32:34])
	if int(framesCount) != dataSize/frameSize-1 {
		return nil, handleError([2]byte{0x00, 0x04}, errors.New("frames count and observed aircrafts mismatch"))
	}

	var offset int
	offset = frameSize

	aircrafts = make([]AircraftInfo, dataSize/frameSize-1)

	re := regexp.MustCompile("[^0-9A-Za-z\\-_]*")

	i := 0
	for offset < dataSize {
		aircrafts[i] = AircraftInfo{
			Hexcode:      hex.EncodeToString([]byte{decompressedData[offset+2], decompressedData[offset+1], decompressedData[offset]}),
			Lat:          float32(int32(binary.LittleEndian.Uint32(decompressedData[offset+12:offset+16]))) / 1000000,
			Long:         float32(int32(binary.LittleEndian.Uint32(decompressedData[offset+8:offset+12]))) / 1000000,
			Callsign:     re.ReplaceAllString(string(decompressedData[offset+78:offset+86]), ""),
			AcType:       re.ReplaceAllString(string(decompressedData[offset+88:offset+92]), ""),
			AcReg:        re.ReplaceAllString(string(decompressedData[offset+92:offset+99]), ""),
			Squawk:       hex.EncodeToString([]byte{decompressedData[offset+33], decompressedData[offset+32]}),
			DbFlag:       decompressedData[offset+86],
			AcCat:        decompressedData[offset+64],
			AcSrc:        uint8(decompressedData[offset+66]) & 240 >> 4,
			HbarFt:       binary.LittleEndian.Uint16(decompressedData[offset+20:offset+22]) * 25,
			HgeoFt:       binary.LittleEndian.Uint16(decompressedData[offset+22:offset+24]) * 25,
			SpdGroundKts: float32(binary.LittleEndian.Uint16(decompressedData[offset+34:offset+36])) / 10,
			Timestamp:    unixTimestamp - uint64(binary.LittleEndian.Uint16(decompressedData[offset+6:offset+8]))*100,
		}
		i++
		offset += frameSize
	}

	return aircrafts, 0
}
