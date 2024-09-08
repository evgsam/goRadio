package ic78civCmd

import (
	"go.bug.st/serial"
)

func requestTransiverAddr(port serial.Port) {
	port.ResetInputBuffer()
	arg := []byte{byte(preambleCmd), byte(preambleCmd), 0x00, byte(controllerAddrCmd), byte(readAddrCmd), 0x00, byte(endMsgCmd)}
	sendDataToTransiver_(port, arg)
}
