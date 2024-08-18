package ic78civCmd

import "strconv"

type civCommand struct {
	preamble          [2]byte
	transiverAddr     byte
	controllerAddr    byte
	setFrequeCommand  byte
	readFrequeCommand byte
	subcodeNumber     byte
	endMsg            byte
	okCode            byte
	ngCode            byte
}

func NewIc78civCommand(transiverAddr byte, controllerAddr byte) *civCommand {
	ic78civCommand := &civCommand{
		preamble:         [2]byte{0xfe, 0xfe},
		transiverAddr:    transiverAddr,
		controllerAddr:   controllerAddr,
		setFrequeCommand: 0x05,
		endMsg:           0xfd,
		okCode:           0xfb,
		ngCode:           0xfa,
	}
	return ic78civCommand
}

func GetCivPreamble(p *civCommand) []byte {
	return p.preamble[:]
}

func GetTransiverAddr(p *civCommand) byte {
	return p.transiverAddr
}

func addElementToFirstIndex(x []byte, y byte) []byte {
	x = append([]byte{y}, x...)
	return x
}

func SetFreque(freq int) {
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

}
