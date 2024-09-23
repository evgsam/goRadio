package ic78civCmd

import (
	"fmt"
	datastruct "goRadio/dataStruct"
	"goRadio/serialDataExchange"
	"sync"
	"time"

	"go.bug.st/serial"
)

var (
	myRadiosettings *datastruct.RadioSettings
	mode_data       string
	att_data        string
	preamp_data     string
	freq_data       uint32
	af_data         uint32
	rf_data         uint32
	sql_data        uint32
	adr_data        byte
)

func dataRequest(port serial.Port, myic78civCommand *civCommand) {
	for {
		commandSend(port, myic78civCommand, freqRead)
		time.Sleep(300 * time.Millisecond)
		commandSend(port, myic78civCommand, mode)
		time.Sleep(300 * time.Millisecond)
		commandSend(port, myic78civCommand, att)
		time.Sleep(300 * time.Millisecond)
		commandSend(port, myic78civCommand, af)
		time.Sleep(300 * time.Millisecond)
		commandSend(port, myic78civCommand, rf)
		time.Sleep(300 * time.Millisecond)
		commandSend(port, myic78civCommand, sql)
		time.Sleep(300 * time.Millisecond)
		commandSend(port, myic78civCommand, preamp)
		time.Sleep(300 * time.Millisecond)
	}

}

func splitByFD(adr byte, data []byte) []byte {
	start := 0
	buffer := make([]byte, 0)
	if data[0] == 0xfe && data[1] == 0xfe { // Проверка на начало пакета
		start = 2
	}
	for _, b := range data {
		if b == 0xfd {
			break
		} else if b == 0xfe && data[start] == 0xfe {
			start += 1
		}
		buffer = append(buffer, b)
	}
	buffer = append(make([]byte, 0), buffer[2:]...)
	if buffer[0] == adr {
		return []byte{0x99}
	}
	return append(make([]byte, 0), buffer[2:]...)
}

func switchMode(data byte) string {
	switch data {
	case 0x00:
		return "LSB"
	case 0x01:
		return "USB"
	case 0x02:
		return "AM"
	case 0x03:
		return "CW"
	case 0x04:
		return "RTTY"
	}
	return ""
}

func switchATT(data byte) string {
	switch data {
	case 0x00:
		return "NO"
	case 0x20:
		return "YES"
	}
	return ""
}

func detectionAFRFSQL(buffer []byte) string {
	return fmt.Sprintf("%d%%", (bcdToInt(buffer)*100)/254)
}

func detectionFreque(buffer []byte) string {
	buffer = append(make([]byte, 0), buffer[1:]...)
	buffRevers := make([]byte, len(buffer))
	j := 0
	for i := len(buffer) - 1; i > -1; i-- {
		buffRevers[j] = buffer[i]
		j++
	}
	return fmt.Sprintf("%d Hz", bcdToInt(buffRevers)/1000)
}

func parser(buffer []byte, ch chan map[byte]string) {
	switch buffer[0] {
	case 0x99:
		ch <- map[byte]string{
			byte(status): "disconnect",
		}
	case byte(sendFreqCmd):
		ch <- map[byte]string{
			byte(freqRead): detectionFreque(buffer),
		}
	case byte(readFreqCmd):
		ch <- map[byte]string{
			byte(freqRead): detectionFreque(buffer),
		}
	case byte(sendModeCmd):
		ch <- map[byte]string{
			byte(mode): switchMode(buffer[1]),
		}
	case byte(readModeCmd):
		ch <- map[byte]string{
			byte(mode): switchMode(buffer[1]),
		}
	case byte(attCmd):
		ch <- map[byte]string{
			byte(att): switchATT(buffer[1]),
		}
	case byte(preampCmd):
		switch buffer[2] {
		case 0x00:
			ch <- map[byte]string{
				byte(preamp): "OFF",
			}
		case 0x01:
			ch <- map[byte]string{
				byte(preamp): "P.AMP",
			}
		}
	case byte(afrfsqlCmd):
		switch buffer[1] {
		case byte(afSubCmd):
			ch <- map[byte]string{
				byte(af): detectionAFRFSQL(buffer[2:4]),
			}
		case byte(rfSubCmd):
			ch <- map[byte]string{
				byte(rf): detectionAFRFSQL(buffer[2:4]),
			}
		case byte(sqlSubCmd):
			ch <- map[byte]string{
				byte(sql): detectionAFRFSQL(buffer[2:4]),
			}
		}
	}
}

func CivCmdParser(port serial.Port, serialAcces *sync.Mutex, chRadioSettings chan map[byte]string) {
	adr_data = 0x62
	myic78civCommand = newIc78civCommand(adr_data)
	go dataRequest(port, myic78civCommand)
	chSerialData := serialDataExchange.SerialPortPoller(port, serialAcces)
	for msg := range chSerialData {
		parser(splitByFD(adr_data, msg), chRadioSettings)
	}

}
