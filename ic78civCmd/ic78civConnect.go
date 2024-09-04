package ic78civCmd

import (
	"fmt"
	"sync"

	"go.bug.st/serial"
)

func newIc78civCommand(transiverAddr byte) *civCommand {
	ic78civCommand := &civCommand{
		transiverAddr:   transiverAddr,
		requestFreque:   []byte{byte(preambleCmd), byte(preambleCmd), transiverAddr, byte(controllerAddrCmd), byte(readFreqCmd), byte(endMsgCmd)},
		requestMode:     []byte{byte(preambleCmd), byte(preambleCmd), transiverAddr, byte(controllerAddrCmd), byte(readModeCmd), byte(endMsgCmd)},
		requestATT:      []byte{byte(preambleCmd), byte(preambleCmd), transiverAddr, byte(controllerAddrCmd), byte(attCmd), byte(endMsgCmd)},
		requestAFLevel:  []byte{byte(preambleCmd), byte(preambleCmd), transiverAddr, byte(controllerAddrCmd), byte(afrfsqlCmd), byte(afSubCmd), byte(endMsgCmd)},
		requestRFLevel:  []byte{byte(preambleCmd), byte(preambleCmd), transiverAddr, byte(controllerAddrCmd), byte(afrfsqlCmd), byte(rfSubCmd), byte(endMsgCmd)},
		requestSQLLevel: []byte{byte(preambleCmd), byte(preambleCmd), transiverAddr, byte(controllerAddrCmd), byte(afrfsqlCmd), byte(sqlSubCmd), byte(endMsgCmd)},
		requestPreamp:   []byte{byte(preambleCmd), byte(preambleCmd), transiverAddr, byte(controllerAddrCmd), byte(preampCmd), byte(preampSubCmd), byte(endMsgCmd)},
	}
	return ic78civCommand
}

func IC78connect(port serial.Port, serialAcces *sync.Mutex) error {
	serialAcces.Lock()
	//fmt.Println("IC78 Connect")
	port.ResetInputBuffer()
	var myic78civCommand *civCommand
	addr, err := requestTransiverAddr(port)
	if err != nil {
		serialAcces.Unlock()
		return err
	} else {
		myic78civCommand = newIc78civCommand(addr)
		fmt.Printf("Transiver Addr: %#x \n", myic78civCommand.transiverAddr)
	}
	freq, err := requestFreque(port, myic78civCommand)
	if err != nil {
		serialAcces.Unlock()
		return err
	} else {
		fmt.Printf("Transiver Freque: %d Hz \n", freq)
	}
	mode, err := requestMode(port, myic78civCommand)
	if err != nil {
		serialAcces.Unlock()
		return err
	} else {
		fmt.Println("Transiver Mode:", mode)
	}
	att, err := requestATT(port, myic78civCommand)
	if err != nil {
		serialAcces.Unlock()
		return err
	} else {
		fmt.Println("ATT status:", att)
	}
	afLevel, err := requestAFLevel(port, myic78civCommand)
	if err != nil {
		serialAcces.Unlock()
		return err
	} else {
		fmt.Printf("AF level: %d % \n", afLevel)
	}
	rfLevel, err := requestRFLevel(port, myic78civCommand)
	if err != nil {
		serialAcces.Unlock()
		return err
	} else {
		fmt.Printf("RF level: %d % \n", rfLevel)
	}
	sqlLevel, err := requestSQLLevel(port, myic78civCommand)
	if err != nil {
		serialAcces.Unlock()
		return err
	} else {
		fmt.Printf("SQL level: %d % \n", sqlLevel)
	}
	preamp, err := requestPreamp(port, myic78civCommand)
	if err != nil {
		serialAcces.Unlock()
		return err
	} else {
		fmt.Println("Preamp status:", preamp)
	}
	err = setFreque(port, myic78civCommand, 3501)
	if err != nil {
		serialAcces.Unlock()
		return err
	} else {
		fmt.Println("freque set")
	}
	err = setMode(port, myic78civCommand, "AM")
	if err != nil {
		serialAcces.Unlock()
		return err
	} else {
		fmt.Println("mode set")
	}
	err = setAfRfSql(port, myic78civCommand, af, 93)
	if err != nil {
		serialAcces.Unlock()
		return err
	} else {
		fmt.Println("af level set")
	}

	err = setAfRfSql(port, myic78civCommand, rf, 99)
	if err != nil {
		serialAcces.Unlock()
		return err
	} else {
		fmt.Println("rf level set")
	}
	err = setAfRfSql(port, myic78civCommand, sql, 69)
	if err != nil {
		serialAcces.Unlock()
		return err
	} else {
		fmt.Println("sql level set")
	}

	err = setPreamp(port, myic78civCommand, "P.AMP")
	if err != nil {
		serialAcces.Unlock()
		return err
	} else {
		fmt.Println("P.AMP set")
	}

	serialAcces.Unlock()
	return nil

}
