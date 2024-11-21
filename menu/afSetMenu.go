/******************************************/
//Регулировка уровня AF
/*****************************************/

package menu

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
