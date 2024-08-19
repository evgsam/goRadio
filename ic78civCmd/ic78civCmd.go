package ic78civCmd

import (
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
	endMsg            byte
	okCode            byte
	ngCode            byte
}

func newIc78civCommand(transiverAddr byte, controllerAddr byte) *civCommand {
	ic78civCommand := &civCommand{
		preamble:          0xfe,
		transiverAddr:     transiverAddr,
		controllerAddr:    controllerAddr,
		setFrequeCommand:  0x05,
		readFrequeCommand: 0x03,
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

func IC78connect(port serial.Port) {
	buff := make([]byte, 11)
	var nmbrByteRead int
	myic78civCommand := newIc78civCommand(0x62, 0xe1)
	serialDataExchange.WriteSerialPort(port, []byte{myic78civCommand.preamble, myic78civCommand.preamble, myic78civCommand.transiverAddr,
		myic78civCommand.controllerAddr, myic78civCommand.readFrequeCommand, myic78civCommand.endMsg})
	nmbrByteRead = serialDataExchange.ReadSerialPort(port, buff)

	if nmbrByteRead == 0 {
		fmt.Println("\nEOF")
	}
	if nmbrByteRead > 0 {
		for _, value := range buff {
			fmt.Printf("%#x ", value)
		}
	}
}
