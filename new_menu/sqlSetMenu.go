package newmenu

func sqlPlusMenu(chDataSet chan map[byte]string) error {
	chDataSet <- map[byte]string{
		byte(sql): "+",
	}
	return nil
}

func sqlMinusMenu(chDataSet chan map[byte]string) error {
	chDataSet <- map[byte]string{
		byte(sql): "-",
	}
	return nil
}
