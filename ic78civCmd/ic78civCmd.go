package ic78civCmd

import (
	"bytes"
	"fmt"
	"goRadio/serialDataExchange"
	"slices"
	"strconv"
	"time"

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
	requestFreque     []byte
}

func newIc78civCommand(controllerAddr byte, transiverAddr byte) *civCommand {
	ic78civCommand := &civCommand{
		preamble:          0xfe,
		transiverAddr:     transiverAddr,
		controllerAddr:    controllerAddr,
		setFrequeCommand:  0x05,
		readFrequeCommand: 0x03,
		readTransiverAddr: 0x19,
		endMsg:            0xfd,
		okCode:            0xfb,
		ngCode:            0xfa,
		requestTAddres:    []byte{0xfe, 0xfe, 0x00, controllerAddr, 0x19, 0x00, 0xfd},
		requestFreque:     []byte{0xfe, 0xfe, transiverAddr, controllerAddr, 0x03, 0xfd},
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

func requestTransiverAddr(port serial.Port, controllerAdr byte) byte {
	requestTAddres := []byte{0xfe, 0xfe, 0x00, controllerAdr, 0x19, 0x00, 0xfd}
	correctMsg := false
	buff := make([]byte, 30)
	var transiverAddr byte
	n := 0
	for !correctMsg {
		time.Sleep(time.Duration(100) * time.Millisecond)
		serialDataExchange.WriteSerialPort(port, requestTAddres)
		time.Sleep(time.Duration(10) * time.Millisecond)
		_ = serialDataExchange.ReadSerialPort(port, buff)
		for _, value := range buff {
			if value == 0xfd {
				n++
			}
		}
		if n < 2 {
			n = 0
			for i, _ := range buff {
				buff[i] = 0x00
			}
		} else {
			correctMsg = true
		}
	}
	for i := 0; i < n; i++ {
		idx := slices.Index(buff, 0xfd)
		if idx != -1 {
			if bytes.Equal(buff[:idx+1], requestTAddres[:len(requestTAddres)]) {
				fmt.Println("ECHO")
				buff = buff[idx+1 : len(buff)]
			} else {
				transiverAddr = buff[idx-1]
			}

		}
	}
	return transiverAddr
}

func requestFreque(port serial.Port, p *civCommand) []byte {
	correctMsg := false
	buff := make([]byte, 30)
	freque := make([]byte, 5)
	n := 0
	for !correctMsg {
		time.Sleep(time.Duration(100) * time.Millisecond)
		serialDataExchange.WriteSerialPort(port, p.requestFreque)
		time.Sleep(time.Duration(10) * time.Millisecond)
		_ = serialDataExchange.ReadSerialPort(port, buff)
		for _, value := range buff {
			if value == 0xfd {
				n++
			}
		}
		if n < 2 {
			n = 0
			for i, _ := range buff {
				buff[i] = 0x00
			}
		} else {
			correctMsg = true
		}
	}
	for i := 0; i < n; i++ {
		idx := slices.Index(buff, p.endMsg)
		if idx != -1 {
			if bytes.Equal(buff[:idx+1], p.requestFreque[:len(p.requestFreque)]) {
				fmt.Println("ECHO")
				buff = buff[idx+1 : len(buff)]
			} else {
				freque = buff[idx-5 : idx]
			}

		}
	}
	return freque
}

func IC78connect(port serial.Port) {
	myic78civCommand := newIc78civCommand(0xe1, requestTransiverAddr(port, 0xe1))
	fmt.Printf("Transiver Addr: %#x", myic78civCommand.transiverAddr)
	fmt.Println(requestFreque(port, myic78civCommand))
}
