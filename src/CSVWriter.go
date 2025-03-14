package main

import (
	"encoding/csv"
	"errors"
	"os"
	"strconv"
	"time"
)

func transformStruct(aircraft AircraftInfo) (aircraftStr []string) {
	aircraftStr = make([]string, 0)
	aircraftStr = append(aircraftStr, aircraft.Hexcode)
	aircraftStr = append(aircraftStr, strconv.FormatFloat(float64(aircraft.Lat), 'f', -1, 32))
	aircraftStr = append(aircraftStr, strconv.FormatFloat(float64(aircraft.Long), 'f', -1, 32))
	aircraftStr = append(aircraftStr, aircraft.Callsign)
	aircraftStr = append(aircraftStr, aircraft.AcReg)
	aircraftStr = append(aircraftStr, aircraft.AcType)
	aircraftStr = append(aircraftStr, aircraft.Squawk)
	aircraftStr = append(aircraftStr, strconv.FormatUint(uint64(aircraft.DbFlag), 16))
	aircraftStr = append(aircraftStr, strconv.FormatUint(uint64(aircraft.AcCat), 16))
	aircraftStr = append(aircraftStr, strconv.FormatUint(uint64(aircraft.AcSrc), 16))
	aircraftStr = append(aircraftStr, strconv.FormatUint(uint64(aircraft.HbarFt), 10))
	aircraftStr = append(aircraftStr, strconv.FormatUint(uint64(aircraft.HgeoFt), 10))
	aircraftStr = append(aircraftStr, strconv.FormatFloat(float64(aircraft.SpdGroundKts), 'f', -1, 32))
	aircraftStr = append(aircraftStr, strconv.FormatUint(aircraft.Timestamp, 10))

	return aircraftStr

}

func writeToCsv(aircrafts []AircraftInfo, outputDirectory string) (resultCode int) {
	if len(aircrafts) == 0 {
		return handleError([2]byte{0x00, 0x05}, errors.New("no aircrafts found"))
	}

	writer, file, resultCode := createCSVWriter(outputDirectory)
	if resultCode != 0 {
		return resultCode
	}
	defer file.Close()

	badResultCode := 0
	writtenRecords := 0
	for _, aircraft := range aircrafts {
		resultCode = writeCSVRecord(writer, transformStruct(aircraft))
		if resultCode != 0 {
			badResultCode = resultCode
		} else {
			writtenRecords++
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return handleError([2]byte{0x00, 0x05}, err)
	}
	mainLog.Printf("Written %v records to csv file\n", writtenRecords)
	return badResultCode
}

func createCSVWriter(outputDirectory string) (writer *csv.Writer, file *os.File, resultCode int) {
	t := time.Now()
	fileName := outputDirectory + "/AilrunOut_" + t.Format("20060102_15") + ".csv"

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, nil, handleError([2]byte{0x00, 0x05}, err)
	}

	writer = csv.NewWriter(file)
	writer.Comma = ';'
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, nil, handleError([2]byte{0x00, 0x05}, err)
	}
	if fileInfo.Size() == 0 {
		header := []string{"HexCode", "Lat", "Lon", "Callsign", "Reg", "Type", "Squawk", "Mil", "Category", "Source", "Alt.baro, ft", "Alt.geo, ft", "Gr Speed, kts", "Timestamp"}
		writeCSVRecord(writer, header)
	}
	return writer, file, 0
}

func writeCSVRecord(writer *csv.Writer, aircraft []string) (resultCode int) {
	err := writer.Write(aircraft)
	if err != nil {
		return handleError([2]byte{0x00, 0x05}, err)
	}
	return 0
}
