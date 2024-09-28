package menu

func modeSetMenu(chDataSet chan map[byte]string) error {

	chDataSet <- map[byte]string{
		byte(mode): "+",
	}

	return nil
}
