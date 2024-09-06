package serialDataExchange

import (
	"sync"
	"time"

	"go.bug.st/serial"
)

func SerialPortPoller(port serial.Port, serialAcces *sync.Mutex) <-chan []byte {
	readBuff := make([]byte, 255)
	ch := make(chan []byte, 255) // 1 - so we keep at least one message
	go func() {
		for {
			serialAcces.Lock()
			_ = ReadSerialPort(port, readBuff)
			select {
			case ch <- readBuff:
			default:
				time.Sleep(200 * time.Millisecond)
			}
			serialAcces.Unlock()
			time.Sleep(200 * time.Millisecond)
		}
	}()
	return ch
}
