package ic78civCmd

import (
	"bytes"
	"slices"

	"go.bug.st/serial"
)

func requestTransiverAddr(port serial.Port) (byte, error) {
	addr, er := commandSend_(port, nil, taddr)
	if er != nil {
		return 0x00, er
	} else {
		return addr[1], nil
	}

}

func commandSend(port serial.Port, p *civCommand, c commandName) ([]byte, error) {
	readBuff := make([]byte, maxReadBuff)
	dataBuff := make([]byte, 7)
	var arg []byte
	var cmd byte
	switch c {
	case freqRead:
		arg = p.requestFreque
		cmd = byte(readFreqCmd)
	case taddr:
		arg = []byte{byte(preambleCmd), byte(preambleCmd), 0x00, byte(controllerAddrCmd), byte(readAddrCmd), 0x00, byte(endMsgCmd)}
		cmd = byte(readAddrCmd)
	case mode:
		arg = p.requestMode
		cmd = byte(readModeCmd)
	case att:
		arg = p.requestATT
		cmd = byte(attCmd)
	case af:
		arg = p.requestAFLevel
		cmd = byte(afrfsqlCmd)
	case rf:
		arg = p.requestRFLevel
		cmd = byte(afrfsqlCmd)
	case sql:
		arg = p.requestSQLLevel
		cmd = byte(afrfsqlCmd)
	case preamp:
		arg = p.requestPreamp
		cmd = byte(preampCmd)
	}

	n, readBuff, err := sendDataToTransiver(port, arg)
	if err != nil {
		return readBuff, err
	}
	for i := 0; i < n; i++ {
		idxCmd := slices.Index(readBuff, cmd)
		idxEnd := slices.Index(readBuff, byte(endMsgCmd))
		if idxCmd > idxEnd {
			readBuff = readBuff[idxCmd+1 : len(readBuff)]
		}
		if idxEnd != -1 {
			if bytes.Equal(readBuff[:idxEnd+1], arg[:len(arg)]) {
				readBuff = readBuff[idxEnd+1 : len(readBuff)]
			} else {
				dataBuff = readBuff[idxCmd:idxEnd]
				dataBuff = append(make([]byte, 0), dataBuff[1:]...)
			}
		}
	}
	return dataBuff, nil
}
