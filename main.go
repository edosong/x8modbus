package main

// The {OA, Bias/Gap} x 8 ch.
// Data: 32bit float,
// Data (each byte) order: LittleEndian -> need to change to BigEdian

import (
	"flag"
	"fmt"
	"log"
	"math"
	"time"
)

const (
	StartRegAddress uint16 = 42101
	defaultHost            = "192.168.179.5"
	defaultPort            = "502" //standard
)

func main() {

	// get device host (url or ip address) and port from the command line
	var (
		host string
		port string
	)

	flag.StringVar(&host, "host", defaultHost, "Slave device host (url or ip address)")
	flag.StringVar(&port, "port", defaultPort, fmt.Sprintf("Slave device port (the default is %s)", defaultPort))
	flag.Parse()
	// client := modbus.TCPClient(host + ":" + port)
	// mbHandler := modbus.NewTCPClientHandler(host + ":" + port)
	// mbHandler.Timeout = 10 * time.Second
	// mbHandler.SlaveId = 1

	// var err error

	// if err = mbHandler.Connect(); err != nil {
	// 	log.Fatal("Connect error:", err)
	// }
	// defer mbHandler.Close()

	// client := modbus.NewClient(mbHandler)
	for { // read data in loop
		// x8Data, err := client.ReadHoldingRegisters(StartRegAddress, 32)
		// x8Data := []byte{96, 113, 95, 60, 197, 145, 20, 61, 21, 212, 91, 60, 2, 175, 3, 61, 48, 33, 90, 60, 95, 248, 1, 61, 193, 145, 95, 60, 87, 87, 201, 60, 138, 56, 174, 60, 230, 139, 217, 60, 92, 52, 108, 60, 119, 208, 195, 60, 195, 64, 118, 60, 128, 83, 208, 60, 66, 8, 134, 60, 60, 59, 188, 60}
		x8Data := []byte{226, 73, 111, 60, 26, 234, 103, 61, 166, 72, 90, 60, 24, 39, 238, 60, 206, 32, 95, 60, 111, 251, 235, 60, 172, 91, 94, 60, 159, 59, 234, 60, 38, 122, 198, 60, 40, 62, 228, 60, 220, 217, 139, 60, 14, 33, 236, 60, 54, 76, 141, 60, 38, 149, 232, 60, 52, 13, 65, 158, 186, 84, 71, 22}

		// if err != nil {
		// 	log.Fatal("Read Register err:", err)
		// }
		processX8Data(x8Data)
		time.Sleep(10 * time.Second)
	}
}

func processX8Data(data []byte) {
	if len(data) != 64 {
		log.Println("Data length error.")
	}

	fmt.Printf("Read:%v\n-----------\n", time.Now())
	for i := 0; i < 8; i++ {
		fmt.Printf("Channel:%d.", i+1)
		getOnePair(data[i*8 : i*8+8])
	}
}

func getOnePair(d []byte) {
	for i := 0; i < 2; i++ {
		// fmt.Printf("..%d(%d, %08b)", i, d[i], d[i])
		dout := uint32(d[0+i*4]) | uint32(d[1+i*4])<<8 | uint32(d[2+i*4])<<16 | uint32(d[3+i*4])<<24
		fout := math.Float32frombits(dout)
		// s := fmt.Sprintf("%.2v",fout)
		// fvalue := strconv.ParseFloat()
		if i == 0 {
			fmt.Printf("OA: %.2f, ", fout)
		} else {
			fmt.Printf("Bias/Gap: %.2f\n", fout)
		}
	}

}
