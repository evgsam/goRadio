package ic78civCmd

import (
	"bytes"
	"fmt"
	"goRadio/serialDataExchange"
	"strconv"

	"go.bug.st/serial"
)

type civCommand struct {
	preamble          byte
	transiverAddr     byte
	controllerAddr    byte
	setFrequeCommand  byte
	readFrequeCommand byte
	readTransiverAddr byte
	endMsg            byte
	okCode            byte
	ngCode            byte
}

func newIc78civCommand(controllerAddr byte) *civCommand {
	ic78civCommand := &civCommand{
		preamble:          0xfe,
		controllerAddr:    controllerAddr,
		setFrequeCommand:  0x05,
		readFrequeCommand: 0x03,
		readTransiverAddr: 0x19,
		endMsg:            0xfd,
		okCode:            0xfb,
		ngCode:            0xfa,
	}
	return ic78civCommand
}

func addElementToFirstIndex(x []byte, y byte) []byte {
	x = append([]byte{y}, x...)
	return x
}

func SetFreque(freq int) {
	buf := make([]byte, 5)
	arr := make([]byte, len(strconv.Itoa(freq)), 10)
	for i := len(arr) - 1; freq > 0; i-- {
		arr[i] = byte(freq % 10)
		freq = int(freq / 10)
	}
	for len(arr) < 10 {
		arr = addElementToFirstIndex(arr, 0)
	}
	dig := 5
	for i := 0; i < 10; i = i + 2 {
		dig--
		buf[dig] = (arr[i] * 10) + arr[i+1]
	}

}

func requestTransmitterAddr(port serial.Port, p *civCommand) {
	buff := make([]byte, 20)
	sendBuf := []byte{p.preamble, p.preamble, 0x00, p.controllerAddr, p.readTransiverAddr, 0x00, p.endMsg}
	serialDataExchange.WriteSerialPort(port, sendBuf)
	_ = serialDataExchange.ReadSerialPort(port, buff)

	if bytes.Equal(buff[:len(sendBuf)], sendBuf[:len(sendBuf)]) {
		fmt.Println("OK")
	} else {
		fmt.Println("ERROR")
	}

	fmt.Println(buff[len(sendBuf)+1])
	fmt.Println(buff[len(sendBuf)+2])
	fmt.Println(buff[len(sendBuf)+3])

	if buff[len(sendBuf)+1] == p.preamble && buff[len(sendBuf)+2] == p.preamble && buff[len(sendBuf)+3] == p.controllerAddr {
		for i := len(sendBuf) + 1; i < len(buff); i++ {
			if buff[i] == p.endMsg {
				p.transiverAddr = buff[i-1]
			}
		}
	}

}

func requestFreque(port serial.Port, p *civCommand) {
	buff := make([]byte, 100)
	serialDataExchange.WriteSerialPort(port, []byte{p.preamble, p.preamble, p.transiverAddr, p.controllerAddr, p.readFrequeCommand, p.endMsg})
	n := serialDataExchange.ReadSerialPort(port, buff)
	if n > 0 {
		for _, value := range buff {
			fmt.Printf("%#x ", value)
		}
	}

}

func IC78connect(port serial.Port) {
	myic78civCommand := newIc78civCommand(0xe1)
	requestTransmitterAddr(port, myic78civCommand)

	//requestFreque(port, myic78civCommand)
}
