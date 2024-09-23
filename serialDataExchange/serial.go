package serialDataExchange

import (
	"log"

	"go.bug.st/serial"
)

func GetSerialPortList() []string {
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		log.Fatal("No serial ports found!")
	}
	return ports
}

func OpenSerialPort(BaudRate int, DataBits uint8, ports string) serial.Port {
	mode := &serial.Mode{
		BaudRate: BaudRate,
		Parity:   serial.NoParity,
		DataBits: int(DataBits),
		StopBits: serial.OneStopBit,
	}
	port, err := serial.Open(ports, mode)
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
