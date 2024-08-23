package ic78civCmd

import (
	"bytes"
	"fmt"
	"goRadio/serialDataExchange"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/albenik/bcd"
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
	requestFreque     []byte
	requestMode       []byte
	requestATT        []byte
	requestAFLevel    []byte
	requestRFLevel    []byte
	requestSQLLevel   []byte
}

func DataPollingGorutine(port serial.Port, serialAcces *sync.Mutex) {
	for {
		serialAcces.Lock()
		fmt.Println("HELLO")
		port.ResetInputBuffer()
		serialAcces.Unlock()
		time.Sleep(3 * time.Second)
	}
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
		requestFreque:     []byte{0xfe, 0xfe, transiverAddr, controllerAddr, 0x03, 0xfd},
		requestMode:       []byte{0xfe, 0xfe, transiverAddr, controllerAddr, 0x04, 0xfd},
		requestATT:        []byte{0xfe, 0xfe, transiverAddr, controllerAddr, 0x11, 0xfd},
		requestAFLevel:    []byte{0xfe, 0xfe, transiverAddr, controllerAddr, 0x14, 0x01, 0xfd},
		requestRFLevel:    []byte{0xfe, 0xfe, transiverAddr, controllerAddr, 0x14, 0x02, 0xfd},
	}
	return ic78civCommand
}

func addElementToFirstIndex(x []byte, y byte) []byte {
	x = append([]byte{y}, x...)
	return x
}

func printByte(data []byte) {
	for _, value := range data {
		fmt.Printf("%#x ", value)
	}
	fmt.Println()
}

func setFreque(freq int) {
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
	println(buf)

}

func bcdToInt(buff []byte) uint32 {
	if len(buff) > 2 {
		buffRevers := make([]byte, len(buff))
		j := 0
		for i := len(buff) - 1; i > -1; i-- {
			buffRevers[j] = buff[i]
			j++
		}
		return bcd.ToUint32(buffRevers)
	}
	return bcd.ToUint32(buff)
}

func civDataParser(request []byte, buff []byte) {

}

func requestTransiverAddr(port serial.Port, controllerAdr byte) byte {
	requestTAddres := []byte{0xfe, 0xfe, 0x00, controllerAdr, 0x19, 0x00, 0xfd}
	correctMsg := false
	buff := make([]byte, 30)
	var transiverAddr byte
	n := 0
	attempt := 0
	for !correctMsg {
		if attempt > 20 {
			break
		}
		attempt++
		port.ResetInputBuffer()
		time.Sleep(time.Duration(100) * time.Millisecond)
		serialDataExchange.WriteSerialPort(port, requestTAddres)
		time.Sleep(time.Duration(100) * time.Millisecond)
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
				buff = buff[idx+1 : len(buff)]
			} else {
				transiverAddr = buff[idx-1]
			}

		}
	}
	if attempt > 99 {
		fmt.Println("attempt>100!")
	}

	return transiverAddr
}

func requestMode(port serial.Port, p *civCommand) string {
	correctMsg := false
	buff := make([]byte, 30)
	var mode string
	var modeByte byte
	n := 0
	for !correctMsg {
		port.ResetInputBuffer()
		time.Sleep(time.Duration(100) * time.Millisecond)
		serialDataExchange.WriteSerialPort(port, p.requestMode)
		time.Sleep(time.Duration(100) * time.Millisecond)
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
			if bytes.Equal(buff[:idx+1], p.requestMode[:len(p.requestMode)]) {
				buff = buff[idx+1 : len(buff)]
			} else {
				modeByte = buff[idx-2]
			}

		}
	}

	switch modeByte {
	case 0x00:
		mode = "LSB"
	case 0x01:
		mode = "USB"
	case 0x02:
		mode = "AM"
	case 0x04:
		mode = "RTTY"
	case 0x07:
		mode = "CW"
	}

	return mode
}

func requestATT(port serial.Port, p *civCommand) string {
	correctMsg := false
	buff := make([]byte, 30)
	var att string
	var attByte byte
	n := 0
	for !correctMsg {
		port.ResetInputBuffer()
		time.Sleep(time.Duration(100) * time.Millisecond)
		serialDataExchange.WriteSerialPort(port, p.requestATT)
		time.Sleep(time.Duration(100) * time.Millisecond)
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
			if bytes.Equal(buff[:idx+1], p.requestATT[:len(p.requestATT)]) {
				buff = buff[idx+1 : len(buff)]
			} else {
				attByte = buff[idx-1]
			}

		}
	}
	switch attByte {
	case 0x00:
		att = "NO"
	case 0x20:
		att = "YES"
	}
	return att
}

func requestFreque(port serial.Port, p *civCommand) uint32 {
	correctMsg := false
	buff := make([]byte, 30)
	freque := make([]byte, 5)
	n := 0
	for !correctMsg {
		port.ResetInputBuffer()
		time.Sleep(time.Duration(100) * time.Millisecond)
		serialDataExchange.WriteSerialPort(port, p.requestFreque)
		time.Sleep(time.Duration(100) * time.Millisecond)
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
				buff = buff[idx+1 : len(buff)]
			} else {
				freque = buff[idx-5 : idx]
			}

		}
	}
	return bcdToInt(freque) / 1000
}

func requestAFLevel(port serial.Port, p *civCommand) uint32 {
	correctMsg := false
	buff := make([]byte, 30)
	level := make([]byte, 5)
	n := 0
	for !correctMsg {
		port.ResetInputBuffer()
		time.Sleep(time.Duration(100) * time.Millisecond)
		serialDataExchange.WriteSerialPort(port, p.requestAFLevel)
		time.Sleep(time.Duration(100) * time.Millisecond)
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
			if bytes.Equal(buff[:idx+1], p.requestAFLevel[:len(p.requestAFLevel)]) {
				buff = buff[idx+1 : len(buff)]
			} else {
				level = buff[idx-2 : idx]
			}

		}
	}
	return (bcdToInt(level) * 100) / 254
}

func IC78connect(port serial.Port, serialAcces *sync.Mutex) {
	serialAcces.Lock()
	fmt.Println("IC78 Connect")
	port.ResetInputBuffer()
	myic78civCommand := newIc78civCommand(0xe1, requestTransiverAddr(port, 0xe1))
	fmt.Printf("Transiver Addr: %#x \n", myic78civCommand.transiverAddr)
	fmt.Printf("Transiver Freque: %d Hz \n", requestFreque(port, myic78civCommand))
	fmt.Println("Transiver Mode:", requestMode(port, myic78civCommand))
	fmt.Println("ATT status:", requestATT(port, myic78civCommand))
	fmt.Printf("AF level: %d % \n", requestAFLevel(port, myic78civCommand))
	//setFreque(30569)
	serialAcces.Unlock()
}
