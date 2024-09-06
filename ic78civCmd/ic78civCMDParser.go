package ic78civCmd

import (
	datastruct "goRadio/dataStruct"
	"goRadio/menu"
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
		return make([]byte, 0)
	}
	return append(make([]byte, 0), buffer[2:]...)
}

func parser(buffer []byte, ch chan *datastruct.RadioSettings) {
	switch buffer[0] {
	case byte(sendFreqCmd):
		buffer = append(make([]byte, 0), buffer[1:]...)
		buffRevers := make([]byte, len(buffer))
		j := 0
		for i := len(buffer) - 1; i > -1; i-- {
			buffRevers[j] = buffer[i]
			j++
		}
		freq_data = bcdToInt(buffRevers) / 1000
	}
	ch <- &datastruct.RadioSettings{
		Err:    nil,
		Status: "Connect",
		Mode:   mode_data,
		ATT:    att_data,
		Preamp: preamp_data,
		Freque: freq_data,
		AF:     af_data,
		RF:     rf_data,
		SQL:    sql_data,
		TrAddr: adr_data,
	}
}
func CivCmdParser(port serial.Port, serialAcces *sync.Mutex) {
	port.ResetInputBuffer()
	adr_data, err := requestTransiverAddr(port)
	if err != nil {
		for err != nil {
			adr_data, err = requestTransiverAddr(port)
			time.Sleep(50 * time.Millisecond)
		}
	}
	myic78civCommand = newIc78civCommand(adr_data)
	mode_data, _ := requestMode(port, myic78civCommand)
	att_data, _ := requestATT(port, myic78civCommand)
	preamp_data, _ := requestPreamp(port, myic78civCommand)
	freq_data, _ := requestFreque(port, myic78civCommand)
	af_data, _ := requestAFLevel(port, myic78civCommand)
	rf_data, _ := requestRFLevel(port, myic78civCommand)
	sql_data, _ := requestSQLLevel(port, myic78civCommand)

	chRadioSettings := make(chan *datastruct.RadioSettings, 30)
	go menu.Menu(chRadioSettings)

	chRadioSettings <- &datastruct.RadioSettings{
		Err:    err,
		Status: "Connect",
		Mode:   mode_data,
		ATT:    att_data,
		Preamp: preamp_data,
		Freque: freq_data,
		AF:     af_data,
		RF:     rf_data,
		SQL:    sql_data,
		TrAddr: adr_data,
	}
	port.ResetInputBuffer()
	chSerialData := serialDataExchange.SerialPortPoller(port, serialAcces)
	for msg := range chSerialData {
		parser(splitByFD(adr_data, msg), chRadioSettings)
	}

}
