package ic78civCmd

import (
	"errors"
	"goRadio/serialDataExchange"
	"time"

	"go.bug.st/serial"
)

func sendDataToTransiver_(port serial.Port, arg []byte) {
	port.ResetInputBuffer()
	time.Sleep(time.Duration(100) * time.Millisecond)
	serialDataExchange.WriteSerialPort(port, arg)
}

func sendDataToTransiver(port serial.Port, arg []byte) (int, []byte, error) {
	n := 0
	attempt := 0
	readBuff := make([]byte, maxReadBuff)
	correctMsg := false
	for !correctMsg {
		if attempt < 10 {
			attempt++
			port.ResetInputBuffer()
			time.Sleep(time.Duration(100) * time.Millisecond)
			serialDataExchange.WriteSerialPort(port, arg)
			time.Sleep(time.Duration(100) * time.Millisecond)
			_ = serialDataExchange.ReadSerialPort(port, readBuff)
			for _, value := range readBuff {
				if value == 0xfd {
					n++
				}
			}
			if n < 2 {
				n = 0
				for i, _ := range readBuff {
					readBuff[i] = 0x00
				}
			} else {
				correctMsg = true
			}
		} else {
			return n, readBuff, errors.New("can't connect Transiver")
		}
	}
	return n, readBuff, nil
}

func commandSend(port serial.Port, p *civCommand, c commandName, value string) {
	var arg []byte
	switch c {
	case freqRead:
		arg = p.requestFreque
	case taddr:
		arg = []byte{byte(preambleCmd), byte(preambleCmd), 0x00, byte(controllerAddrCmd), byte(readAddrCmd), 0x00, byte(endMsgCmd)}
	case mode:
		arg = p.requestMode
	case att:
		arg = p.requestATT
	case af:
		arg = p.requestAFLevel
	case rf:
		arg = p.requestRFLevel
	case sql:
		arg = p.requestSQLLevel
	case preamp:
		arg = p.requestPreamp
	}
	port.ResetInputBuffer()
	time.Sleep(time.Duration(100) * time.Millisecond)
	serialDataExchange.WriteSerialPort(port, arg)
}
