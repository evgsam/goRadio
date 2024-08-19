package main

import (
	"fmt"
	"goRadio/ic78civCmd"
	"goRadio/serialDataExchange"

	"go.bug.st/serial"
)

func printByte(data []byte) {
	for _, value := range data {
		fmt.Printf("%#x ", value)
	}
	fmt.Println()
}

func main() {
	var port serial.Port

	/*fmt.Println(ic78civCmd.GetTransiverAddr(myic78civCommand))
	receiveOk := false
	var nmbrByteRead int
	var attemptСount int


	connectCommand := []byte{0xfe, 0xfe, 0x62, 0xe1, 0x19, 0x00, 0xfd}
	frequeCommand := []byte{0xfe, 0xfe, 0x62, 0xe1, 0x00, 0x50, 0x34, 0x12, 0x05, 0x00, 0xfd}
	//frequeCommand2 := []byte{0xfe, 0xfe, 0x62, 0xe1, 0x15, 0x02, 0xfd}
	answerOk := []byte{0xfe, 0xfe, 0xe1, 0x62, 0xfb, 0xfd}
	*/
	port = serialDataExchange.OpenSerialPort(19200, 8)
	fmt.Println(port)
	ic78civCmd.IC78connect(port)
	/*for attemptСount <= 100 {
		//	writeSerialPort(port, []byte{myic78civCommand.preamble[0], myic78civCommand.preamble[1], myic78civCommand.transiverAddr, myic78civCommand.controllerAddr, 0x19, 0x00, myic78civCommand.endMsg})
		fmt.Print("TX:")
		printByte(connectCommand)
		fmt.Print("RX:")
		buff := make([]byte, 7)
		for {
			nmbrByteRead = readSerialPort(port, buff)
			if nmbrByteRead == 0 {
				fmt.Println("\nEOF")
				break
			}
			if nmbrByteRead > 0 {
				if bytes.Equal(buff[:4], answerOk[:4]) {
					receiveOk = true
				}
				printByte(buff)
				break
			}
		}
		if receiveOk {
			//var freque *[]byte

			//writeSerialPort(port, []byte{myic78civCommand.preamble[0], myic78civCommand.preamble[1], myic78civCommand.transiverAddr, myic78civCommand.controllerAddr})
			//writeSerialPort(port, setFreque(536978))
			//writeSerialPort(port, []byte{myic78civCommand.endMsg})
			time.Sleep(time.Duration(10) * time.Millisecond)
			//writeSerialPort(port, []byte{myic78civCommand.preamble[0], myic78civCommand.preamble[1], myic78civCommand.transiverAddr, myic78civCommand.controllerAddr, 0x15, 0x02, myic78civCommand.endMsg})
			fmt.Print("Freque set conmmand TX:")
			printByte(frequeCommand)
			nmbrByteRead = readSerialPort(port, buff)
			fmt.Print("Freque set conmmand RX:")
			if nmbrByteRead == 0 {
				fmt.Println("\nEOF")
				break
			}
			if nmbrByteRead > 0 {
				printByte(buff)
			}
			break
		}
		attemptСount++
		time.Sleep(time.Duration(20) * time.Millisecond)
		nmbrByteRead = readSerialPort(port, buff)
		fmt.Print("Freque set conmmand RX:")
		if nmbrByteRead == 0 {
			fmt.Println("\nEOF")
			break
		}
		if nmbrByteRead > 0 {
			printByte(buff)
		}
	}
	fmt.Printf("\nreceiveOk=%v, attemptСount=%v \n", receiveOk, attemptСount)
	*/
}
