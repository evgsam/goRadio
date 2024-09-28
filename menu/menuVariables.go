package menu

type commandName int

const (
	freqRead commandName = iota
	taddr
	mode
	att
	af
	rf
	sql
	preamp
	status
	mainViews
)

var (
	viewsNames      = make(map[byte]string)
	infoViewArray   = make([]viewsStruct, 0)
	hotkeyViewArray = make([]viewsStruct, 0)
	inputViewArray  = make([]viewsStruct, 0)
	spMenuActive    bool
	freqMenuActive  bool
)

type viewsStruct struct {
	cmd            commandName
	name           string
	x0, y0, x1, y1 int
	value          string
}

func requireValidator(value string) bool {
	if value == "" {
		return false
	}
	return true
}
