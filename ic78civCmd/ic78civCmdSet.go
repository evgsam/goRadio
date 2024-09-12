package ic78civCmd

import (
	"goRadio/serialDataExchange"
	"strconv"
	"time"

	"go.bug.st/serial"
)

func sendData(port serial.Port, arg []byte) {
	port.ResetInputBuffer()
	time.Sleep(time.Duration(100) * time.Millisecond)
	serialDataExchange.WriteSerialPort(port, arg)
}

func modeByte(val string) byte {
	var cmd byte
	switch val {
	case "LSB":
		cmd = 0x00
	case "USB":
		cmd = 0x01
	case "AM":
		cmd = 0x02
	case "CW":
		cmd = 0x03
	case "RTTY":
		cmd = 0x04
	}
	return cmd
}

func modePreamp(val string) byte {
	var cmd byte
	switch val {
	case "LSB":
		cmd = 0x00
	case "USB":
		cmd = 0x01
	case "AM":
		cmd = 0x02
	case "CW":
		cmd = 0x03
	case "RTTY":
		cmd = 0x04
	}
	return cmd
}

func freqBuf(val string) []byte {
	freq, _ := strconv.ParseUint(val, 10, 32)
	return intFreqToBcdArr(int(freq))
}

func IC78civCmdSet(port serial.Port, ch chan map[byte]string) {
	var arg []byte
	transiverAddr := 0x62
	for {
		m := <-ch
		for key, val := range m {
			switch key {
			case byte(mode):
				arg = []byte{byte(preambleCmd), byte(preambleCmd), byte(transiverAddr), byte(controllerAddrCmd), byte(setModeCmd), modeByte(val), byte(endMsgCmd)}
				sendData(port, arg)
			case byte(att):
				arg = []byte{byte(preambleCmd), byte(preambleCmd), byte(transiverAddr), byte(controllerAddrCmd), byte(sendFreqCmd), byte(endMsgCmd)}
				sendData(port, arg)
			case byte(preamp):

			case byte(freqRead):
				cmd := freqBuf(val)
				arg = []byte{byte(preambleCmd), byte(preambleCmd), byte(transiverAddr), byte(controllerAddrCmd), byte(sendFreqCmd),
					cmd[0], cmd[1], cmd[2], cmd[3], cmd[4], byte(endMsgCmd)}
				sendData(port, arg)
			case byte(af):
			case byte(rf):
			case byte(sql):
			}
		}
	}

}
