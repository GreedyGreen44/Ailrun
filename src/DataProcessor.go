package main

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
)

func processData(decompressedData []byte) (aircrafts []AircraftInfo, resultCode int) {
	const frameSize = 112

	dataSize := len(decompressedData)
	mainLog.Printf("Data size: %v\n", dataSize)

	if dataSize%frameSize != 0 {
		return nil, handleError([2]byte{0x00, 0x04}, errors.New("invalid size of response data11"))
	}
	unixTimestamp := binary.LittleEndian.Uint64(decompressedData[0:8])
	mainLog.Printf("Current timestamp: %v\n", unixTimestamp)

	framesCount := binary.LittleEndian.Uint16(decompressedData[32:34])
	mainLog.Printf("Aircrafts recived: %v\n", framesCount)
	if int(framesCount) != dataSize/frameSize-1 {
		return nil, handleError([2]byte{0x00, 0x04}, errors.New("frames count and observed aircrafts mismatch"))
	}

	var offset int
	offset = frameSize

	aircrafts = make([]AircraftInfo, dataSize/frameSize-1)

	i := 0
	for offset+frameSize <= dataSize {
		aircrafts[i] = AircraftInfo{
			Hexcode:         hex.EncodeToString([]byte{decompressedData[offset+2], decompressedData[offset+1], decompressedData[offset]}),
			Lat:             float32(int32(binary.LittleEndian.Uint32(decompressedData[offset+8:offset+12]))) / 1000000,
			Long:            float32(int32(binary.LittleEndian.Uint32(decompressedData[offset+12:offset+16]))) / 1000000,
			Callsign:        string(decompressedData[offset+78 : offset+86]),
			AcType:          string(decompressedData[offset+88 : offset+92]),
			AcReg:           string(decompressedData[offset+92 : offset+99]),
			Squawk:          hex.EncodeToString([]byte{decompressedData[offset+33], decompressedData[offset+32]}),
			DbFlag:          decompressedData[offset+86],
			AcCat:           decompressedData[offset+64],
			AcSrc:           uint8(decompressedData[offset+66]) & 240 >> 4,
			HbarFt:          binary.LittleEndian.Uint16(decompressedData[offset+20 : offset+22]),
			HgeoFt:          binary.LittleEndian.Uint16(decompressedData[offset+22 : offset+24]),
			SpdGroundKtsX10: int16(binary.LittleEndian.Uint16(decompressedData[offset+34 : offset+36])),
			Timestamp:       unixTimestamp - uint64(binary.LittleEndian.Uint16(decompressedData[offset+6:offset+8]))*100,
		}
		i++
		offset += frameSize
	}

	return aircrafts, 0
}
