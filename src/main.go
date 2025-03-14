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
	Hexcode      string
	Lat          float32
	Long         float32
	Callsign     string
	AcType       string
	AcReg        string
	Squawk       string
	DbFlag       uint8
	AcCat        uint8
	AcSrc        uint8
	HbarFt       uint16
	HgeoFt       uint16
	SpdGroundKts float32
	Timestamp    uint64
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

	if resultCode != 0 {
		return resultCode
	}

	if configMap["OutputType"] != "csv" {
		return handleError([2]byte{0x00, 0x05}, errors.New("for now only saving to csv is available"))
	}

	return writeToCsv(aircrafts, configMap["OutputDirectory"])
}

func main() {

	mainLog = log.New(os.Stdout, "Ailrun: ", log.LstdFlags|log.Lmicroseconds)
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
