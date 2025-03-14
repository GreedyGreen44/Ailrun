package main

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

func readConfigFile(fileName string) (configMap map[string]string, resultCode int) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, handleError([2]byte{0x01, 0x02}, errors.New("failed to open config file"))
	}
	defer file.Close()
	configMap = make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		param := strings.Split(line, "=")
		if len(param) != 2 {
			handleError([2]byte{0x00, 0x01}, errors.New("invalid config string format"))
			continue
		}
		configMap[param[0]] = param[1]
	}

	resultCode = checkMap(configMap)

	return configMap, resultCode
}

func checkMap(configMap map[string]string) (resultCode int) {
	if _, ok := configMap["Host"]; !ok {
		return handleError([2]byte{0x01, 0x03}, errors.New("no Host parameter in config"))
	}
	if _, ok := configMap["BoxBot"]; !ok {
		return handleError([2]byte{0x01, 0x03}, errors.New("no BoxBot parameter in config"))
	}
	if _, ok := configMap["BoxTop"]; !ok {
		return handleError([2]byte{0x01, 0x03}, errors.New("no BoxTop parameter in config"))
	}
	if _, ok := configMap["BoxLeft"]; !ok {
		return handleError([2]byte{0x01, 0x03}, errors.New("no BoxLeft parameter in config"))
	}
	if _, ok := configMap["BoxRight"]; !ok {
		return handleError([2]byte{0x01, 0x03}, errors.New("no BoxRight parameter in config"))
	}
	if _, ok := configMap["OutputType"]; !ok {
		return handleError([2]byte{0x01, 0x03}, errors.New("no OutputType parameter in config"))
	}
	if _, ok := configMap["OutputDirectory"]; !ok {
		return handleError([2]byte{0x01, 0x03}, errors.New("no OutputDirectory parameter in config"))
	}
	if _, ok := configMap["TimerValueSecs"]; !ok {
		return handleError([2]byte{0x01, 0x03}, errors.New("no TimerValueSecs parameter in config"))
	}
	return 0
}
