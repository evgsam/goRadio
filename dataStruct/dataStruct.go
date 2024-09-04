package datastruct

type RadioSettings struct {
	Err    error
	Status string
	Mode   string
	ATT    string
	Preamp string
	Freque uint32
	AF     uint32
	RF     uint32
	SQL    uint32
	TrAddr byte
}
