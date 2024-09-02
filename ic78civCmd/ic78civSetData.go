package ic78civCmd

import (
	"errors"
	"slices"

	"go.bug.st/serial"
)

func setAfRfSql(port serial.Port, p *civCommand, c commandName, level int) error {
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
	buf = []byte{byte(preambleCmd), byte(preambleCmd), p.transiverAddr, byte(controllerAddrCmd), byte(afrfsqlCmd),
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

func setMode(port serial.Port, p *civCommand, mode string) error {
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
	buf = []byte{byte(preambleCmd), byte(preambleCmd), p.transiverAddr, byte(controllerAddrCmd), byte(setModeCmd), arg, byte(endMsgCmd)}
	_, readBuff, err := sendDataToTransiver(port, buf)
	if err != nil {
		return err
	}
	if slices.Index(readBuff, byte(ngCode)) > 0 {
		return errors.New("transceiver sent NG")
	}
	return nil
}

func setFreque(port serial.Port, p *civCommand, freq int) error {
	freqBuf := intFreqToBcdArr(freq)
	buf := make([]byte, 11)
	buf = []byte{byte(preambleCmd), byte(preambleCmd), p.transiverAddr, byte(controllerAddrCmd), byte(setFreqCmd)}
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

func setPreamp(port serial.Port, p *civCommand, preamp string) error {
	var cmd byte
	switch preamp {
	case "OFF":
		cmd = 0x00
	case "P.AMP":
		cmd = 0x01
	}
	buf := make([]byte, 7)
	buf = []byte{byte(preambleCmd), byte(preambleCmd), p.transiverAddr, byte(controllerAddrCmd), byte(preampCmd), byte(preampSubCmd), cmd, byte(endMsgCmd)}
	_, readBuff, err := sendDataToTransiver(port, buf)
	if err != nil {
		return err
	}
	if slices.Index(readBuff, byte(ngCode)) > 0 {
		return errors.New("transceiver sent NG")
	}
	return nil

}
