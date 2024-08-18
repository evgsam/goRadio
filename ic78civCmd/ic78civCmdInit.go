package ic78civCmd

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
