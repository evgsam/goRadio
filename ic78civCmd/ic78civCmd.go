package ic78civCmd

import (
	"bytes"
	"fmt"
	"goRadio/serialDataExchange"
	"slices"
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
	requestTAddres    []byte
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
		requestTAddres:    []byte{0xfe, 0xfe, 0x00, controllerAddr, 0x19, 0x00, 0xfd},
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

func civDataParser(request []byte, buff []byte) {

}

func requestTransmitterAddr(port serial.Port, p *civCommand) {
	correctMsg := false
	buff := make([]byte, 30)
	cmdBuff := make([][]byte, 3)
	for i := range cmdBuff {
		cmdBuff[i] = make([]byte, 20)
	}
	n := 0
	for !correctMsg {
		serialDataExchange.WriteSerialPort(port, p.requestTAddres)
		_ = serialDataExchange.ReadSerialPort(port, buff)

		for _, value := range buff {
			if value == 0xfd {
				n++
			}
		}
		if n < 2 {
			for i, _ := range buff {
				buff[i] = 0x00
			}
		} else {
			correctMsg = true
		}
		fmt.Println(n)
	}
	for i := 0; i < n; i++ {
		idx := slices.Index(buff, p.endMsg)
		idx2 := slices.Index(buff, p.preamble)
		fmt.Println(idx + 1)
		fmt.Println(len(p.requestTAddres))
		if idx != -1 && idx2 != -1 {
			if bytes.Equal(buff[:idx+1], p.requestTAddres[:len(p.requestTAddres)]) {
				fmt.Println("ECHO")
				buff = buff[idx+1 : len(buff)]
			} else {
				//buff1 = buff[0 : idx+1]
				cmdBuff[i] = buff[0 : idx+1]
				//buff = buff[idx+1 : len(buff)]
			}

		}
	}
	println("END")

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
