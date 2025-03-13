package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type AircraftInfo struct {
	Hexcode         string
	Lat             float32
	Long            float32
	Callsign        string
	AcType          string
	AcReg           string
	Squawk          string
	DbFlag          uint8
	AcCat           uint8
	AcSrc           uint8
	HbarFt          uint16
	HgeoFt          uint16
	SpdGroundKtsX10 int16
	Timestamp       uint64
}

var (
	mainLog    *log.Logger
	errorLog   *log.Logger
	warningLog *log.Logger
	done       chan bool
)

func initialCheck(clArgs []string) (map[string]string, int) {
	if len(clArgs) != 1 {
		return nil, handleError([2]byte{0x01, 0x01}, errors.New("incorrect number of arguments. You only need to provide path to config file"))
	}
	configMap, resultCode := readConfigFile(clArgs[0])

	return configMap, resultCode
}

func timerTick(configMap map[string]string) (resultCode int) {
	compressedData, resultCode := handleRequest(configMap)
	if resultCode != 0 {
		return resultCode
	}
	decompressedData, resultCode := decompressZstd(compressedData)
	if resultCode != 0 {
		return resultCode
	}
	aircrafts, resultCode := processData(decompressedData)

	for i, aircraft := range aircrafts {
		fmt.Println(i)
		fmt.Println(aircraft.Hexcode)
		fmt.Println(aircraft.Lat)
		fmt.Println(aircraft.Long)
		fmt.Println(aircraft.Callsign)
		fmt.Println(aircraft.AcType)
		fmt.Println(aircraft.AcReg)
		fmt.Println(aircraft.Squawk)
		fmt.Println(aircraft.DbFlag)
		fmt.Println(aircraft.AcCat)
		fmt.Println(aircraft.AcSrc)
		fmt.Println(aircraft.HbarFt)
		fmt.Println(aircraft.HgeoFt)
		fmt.Println(aircraft.SpdGroundKtsX10)
		fmt.Println(aircraft.Timestamp)
		fmt.Println("")
	}
	return resultCode
}

func main() {

	mainLog = log.New(os.Stdout, "Main: ", log.LstdFlags|log.Lmicroseconds)
	errorLog = log.New(os.Stderr, "Error: ", log.LstdFlags|log.Lmicroseconds)
	warningLog = log.New(os.Stdout, "Warning: ", log.LstdFlags|log.Lmicroseconds)

	mainLog.Println("Starting Ailrun...")

	clArgs := os.Args[1:]
	configMap, resultCode := initialCheck(clArgs)
	if resultCode == 2 {
		return
	}
	tickerInterval, err := strconv.Atoi(configMap["TimerValueSecs"])
	if err != nil {
		handleError([2]byte{0x00, 0x01}, err)
	}

	ticker := time.NewTicker(time.Duration(tickerInterval) * time.Second)
	done = make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				if timerTick(configMap) == 2 {
					os.Exit(0)
				}
			}
		}
	}()
	var inputCommand string
	fmt.Println("Press enter to stop Ailrun...")
	_, err = fmt.Scanln(&inputCommand)
	ticker.Stop()
	done <- true
}
