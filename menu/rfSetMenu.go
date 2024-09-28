package menu

func rfPlusMenu(chDataSet chan map[byte]string) error {
	chDataSet <- map[byte]string{
		byte(rf): "+",
	}
	return nil
}

func rfMinusMenu(chDataSet chan map[byte]string) error {
	chDataSet <- map[byte]string{
		byte(rf): "-",
	}
	return nil
}
