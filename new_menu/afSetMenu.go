package newmenu

func afPlusMenu(chDataSet chan map[byte]string) error {

	chDataSet <- map[byte]string{
		byte(af): "+",
	}

	return nil
}

func afMinusMenu(chDataSet chan map[byte]string) error {

	chDataSet <- map[byte]string{
		byte(af): "-",
	}

	return nil
}
