package ic78civCmd

import (
	"goRadio/serialDataExchange"
	"sync"

	"go.bug.st/serial"
)

func splitByFD(data []byte) [][]byte {
	var result [][]byte
	start := 0
	buffer := make([]byte, 0)
	if data[0] == 0xfe && data[1] == 0xfe { // Проверка на начало пакета
		start = 2
	}
	for _, b := range data {
		if b == 0xfd {
			break
		} else if b == 0xfe && data[start] == 0xfe {
			start += 1
		}
		buffer = append(buffer, b)
	}
	result = append(result, buffer[4:])
	return result
}

func CivCmdParser(port serial.Port, serialAcces *sync.Mutex) {
	_, _ = requestTransiverAddr(port)
	ch := serialDataExchange.SerialPortPoller(port, serialAcces)
	for msg := range ch {
		d := splitByFD(msg)
		printByte(d[0])
	}

}
