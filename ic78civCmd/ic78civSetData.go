package ic78civCmd

import (
	"errors"
	"slices"
	"strconv"

	"go.bug.st/serial"
)

func IC78civCmdSet(port serial.Port, ch chan map[byte]string) {
	transiverAddr := 0x62
	for {
		m := <-ch
		for key, val := range m {
			switch key {
			case byte(mode):
				setMode(port, byte(transiverAddr), val)
			case byte(att):
			case byte(preamp):
				setPreamp(port, byte(transiverAddr), val)
			case byte(freqRead):
				freq, _ := strconv.ParseUint(val, 10, 32)
				setFreque(port, byte(transiverAddr), int(freq))
			case byte(af):
				level, _ := strconv.ParseUint(val, 10, 32)
				setAfRfSql(port, byte(transiverAddr), af, int(level))
			case byte(rf):
				level, _ := strconv.ParseUint(val, 10, 32)
				setAfRfSql(port, byte(transiverAddr), af, int(level))
			case byte(sql):
				level, _ := strconv.ParseUint(val, 10, 32)
				setAfRfSql(port, byte(transiverAddr), af, int(level))
			}
		}
	}

}

func setAfRfSql(port serial.Port, trAddr byte, c commandName, level int) error {
	levelBuf := intToArr((level*255)/100, 4)
	for len(levelBuf) < 4 {
		levelBuf = addElementToFirstIndex(levelBuf, 0)
	}
	var subcmd byte
	switch c {
	case af:
		subcmd = byte(afSubCmd)
	case rf:
		subcmd = byte(rfSubCmd)
	case sql:
		subcmd = byte(sqlSubCmd)
	}
	buf := make([]byte, 7)
	buf = []byte{byte(preambleCmd), byte(preambleCmd), trAddr, byte(controllerAddrCmd), byte(afrfsqlCmd),
		subcmd, byteArrToBCD(levelBuf, 2)[1], byteArrToBCD(levelBuf, 2)[0], byte(endMsgCmd)}
	_, readBuff, err := sendDataToTransiver(port, buf)
	if err != nil {
		return err
	}
	if slices.Index(readBuff, byte(ngCode)) > 0 {
		return errors.New("transceiver sent NG")
	}
	return nil
}

func setMode(port serial.Port, trAddr byte, mode string) error {
	var arg byte
	switch mode {
	case "LSB":
		arg = 0x00
	case "USB":
		arg = 0x01
	case "AM":
		arg = 0x02
	case "RTTY":
		arg = 0x04
	case "CW":
		arg = 0x07
	}
	buf := make([]byte, 7)
	buf = []byte{byte(preambleCmd), byte(preambleCmd), trAddr, byte(controllerAddrCmd), byte(setModeCmd), arg, byte(endMsgCmd)}
	_, readBuff, err := sendDataToTransiver(port, buf)
	if err != nil {
		return err
	}
	if slices.Index(readBuff, byte(ngCode)) > 0 {
		return errors.New("transceiver sent NG")
	}
	return nil
}

func setFreque(port serial.Port, trAddr byte, freq int) error {
	freqBuf := intFreqToBcdArr(freq)
	buf := make([]byte, 11)
	buf = []byte{byte(preambleCmd), byte(preambleCmd), trAddr, byte(controllerAddrCmd), byte(setFreqCmd)}
	for i := 0; i < len(freqBuf); i++ {
		buf = append(buf, freqBuf[i])
	}
	buf = append(buf, byte(endMsgCmd))
	_, readBuff, err := sendDataToTransiver(port, buf)
	if err != nil {
		return err
	}
	if slices.Index(readBuff, byte(ngCode)) > 0 {
		return errors.New("transceiver sent NG")
	}
	return nil
}

func setPreamp(port serial.Port, trAddr byte, preamp string) error {
	var cmd byte
	switch preamp {
	case "OFF":
		cmd = 0x00
	case "P.AMP":
		cmd = 0x01
	}
	buf := make([]byte, 7)
	buf = []byte{byte(preambleCmd), byte(preambleCmd), trAddr, byte(controllerAddrCmd), byte(preampCmd), byte(preampSubCmd), cmd, byte(endMsgCmd)}
	_, readBuff, err := sendDataToTransiver(port, buf)
	if err != nil {
		return err
	}
	if slices.Index(readBuff, byte(ngCode)) > 0 {
		return errors.New("transceiver sent NG")
	}
	return nil

}
