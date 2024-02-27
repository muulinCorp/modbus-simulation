package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"time"
	"unsafe"

	"github.com/tbrandon/mbserver"
)

var (
	confpath = flag.String("confpath", "./conf/", "write memory profile to `file`")
	v        = flag.Bool("v", false, "version")

	Version   = "1.0.0"
	BuildTime = "2000-01-01T00:00:00+0800"
)

func main() {
	flag.Parse()

	if *v {
		fmt.Println(Version)
		return
	}

	serv := mbserver.NewServer()
	serv.RegisterFunctionHandler(3,
		func(s *mbserver.Server, frame mbserver.Framer) ([]byte, *mbserver.Exception) {
			fmt.Println("bytes: ", frame.Bytes())
			fmt.Println("data: ", frame.GetData())
			fmt.Println("func: ", frame.GetFunction())
			data := frame.GetData()
			register := int(binary.BigEndian.Uint16(data[0:2]))
			numRegs := int(binary.BigEndian.Uint16(data[2:4]))
			endRegister := register + numRegs
			if register != 40000 {
				return mbserver.ReadHoldingRegisters(s, frame)
			}
			fmt.Println(register, numRegs, endRegister)
			if endRegister > 65536 {
				return []byte{}, &mbserver.IllegalDataAddress
			}
			fmt.Println(time.Now().Unix())
			return append([]byte{byte(numRegs * 2)}, Uint32ToByteArray(uint32(time.Now().Unix()))...), &mbserver.Success
		})
	fmt.Println("modbus server listen on 502")
	err := serv.ListenTCP("0.0.0.0:502")
	if err != nil {
		log.Printf("%v\n", err)
	}
	defer serv.Close()

	// Wait forever
	for {
		time.Sleep(1 * time.Second)
	}
}

func Uint32ToByteArray(num uint32) []byte {
	size := int(unsafe.Sizeof(num))
	arr := make([]byte, size)
	for i := 0; i < size; i++ {
		byt := *(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&num)) + uintptr(i)))
		arr[i] = byt
	}
	return arr
}
