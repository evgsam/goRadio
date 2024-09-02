package ic78civCmd

import "go.bug.st/serial"

func requestTransiverAddr(port serial.Port) (byte, error) {
	addr, er := commandSend(port, nil, taddr)
	if er != nil {
		return 0x00, er
	} else {
		return addr[1], nil
	}

}

func requestMode(port serial.Port, p *civCommand) (string, error) {
	buff, err := commandSend(port, p, mode)
	var mode string
	if err != nil {
		return "error", err
	}
	switch buff[0] {
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
	return mode, err
}

func requestPreamp(port serial.Port, p *civCommand) (string, error) {
	buff, err := commandSend(port, p, preamp)
	var preamp string
	if err != nil {
		return "error", err
	}
	buff = append(make([]byte, 0), buff[1:]...)
	switch buff[0] {
	case 0x00:
		preamp = "OFF"
	case 0x01:
		preamp = "P.AMP"
	}
	return preamp, err
}

func requestATT(port serial.Port, p *civCommand) (string, error) {
	buff, err := commandSend(port, p, att)
	var att string
	if err != nil {
		return "error", err
	}
	switch buff[0] {
	case 0x00:
		att = "NO"
	case 0x20:
		att = "YES"
	}
	return att, err
}

func requestFreque(port serial.Port, p *civCommand) (uint32, error) {
	buff, err := commandSend(port, p, freqRead)
	if err != nil {
		return 0, err
	}
	buffRevers := make([]byte, len(buff))
	j := 0
	for i := len(buff) - 1; i > -1; i-- {
		buffRevers[j] = buff[i]
		j++
	}
	return bcdToInt(buffRevers) / 1000, err
}

func requestAFLevel(port serial.Port, p *civCommand) (uint32, error) {
	buff, err := commandSend(port, p, af)
	if err != nil {
		return 0, err
	}
	buff = append(make([]byte, 0), buff[1:]...)
	return (bcdToInt(buff) * 100) / 254, err
}

func requestSQLLevel(port serial.Port, p *civCommand) (uint32, error) {
	buff, err := commandSend(port, p, sql)
	if err != nil {
		return 0, err
	}
	buff = append(make([]byte, 0), buff[1:]...)
	return (bcdToInt(buff) * 100) / 254, err
}

func requestRFLevel(port serial.Port, p *civCommand) (uint32, error) {
	buff, err := commandSend(port, p, rf)
	if err != nil {
		return 0, err
	}
	buff = append(make([]byte, 0), buff[1:]...)
	return (bcdToInt(buff) * 100) / 254, err
}
