package serialDataExchange

import (
	"goRadio/menu"
	"log"

	"go.bug.st/serial"
)

func OpenSerialPort(BaudRate int, DataBits uint8) serial.Port {
	serialListCh := make(chan []string)
	go menu.SerialPortSelectMenu(serialListCh)
	ports, err := serial.GetPortsList()
	var portsnum int
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		log.Fatal("No serial ports found!")
	}
	serialListCh <- ports
	myPort := <-serialListCh
	close(serialListCh)
	for i, val := range ports {
		if val == myPort[0] {
			portsnum = i
			break
		}
	}
	mode := &serial.Mode{
		BaudRate: BaudRate,
		Parity:   serial.NoParity,
		DataBits: int(DataBits),
		StopBits: serial.OneStopBit,
	}
	port, err := serial.Open(ports[portsnum], mode)
	if err != nil {
		log.Fatal(err)
	}
	return port
}

func ReadSerialPort(port serial.Port, buff []byte) int {
	n, err := port.Read(buff)
	if err != nil {
		log.Fatal(err)
	}
	return n
}

func WriteSerialPort(port serial.Port, buff []byte) {
	_, err := port.Write(buff)
	if err != nil {
		log.Fatal(err)
	}
}
