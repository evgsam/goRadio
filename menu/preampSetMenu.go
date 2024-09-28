package menu

func preampSetMenu(chDataSet chan map[byte]string) error {

	chDataSet <- map[byte]string{
		byte(preamp): "+",
	}

	return nil
}
