package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

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

	switch configMap["OutputType"] {
	case "csv":
		return writeToCsv(aircrafts, configMap["OutputDirectory"])
	case "json":
		return writeToJson(aircrafts, configMap["OutputDirectory"])
	default:
		return handleError([2]byte{0x01, 0x05}, errors.New("output type is not supported"))
	}
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
