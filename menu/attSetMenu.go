/******************************************/
//Регулировка уровня ATT
/*****************************************/

package menu

func attSetMenu(chDataSet chan map[byte]string) error {

	chDataSet <- map[byte]string{
		byte(att): "+",
	}

	return nil
}
