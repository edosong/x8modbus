package main

// The {OA, Bias/Gap} x 8 ch.
// Data: 32bit float,
// Data (each byte) order: LittleEndian -> need to change to BigEdian
// Date:20230131-0206

import (
	"flag"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/goburrow/modbus"
)

const (
	StartOAAddress uint16 = 0xA475 // 42101
	RegQuantity           = 32

	DefaultHost     = "192.168.179.5"
	DefaultPort     = "502" //standard
	DefaultRate     = 60    //second
	MinRate         = 10
	DefaultDataNums = 10 //
	MaxDataNums     = 1000
)

func main() {

	// get device host (url or ip address) and port from the command line
	var (
		host     string
		port     string
		rate     int64
		dataNums int64
	)

	flag.StringVar(&host, "host", DefaultHost, "Slave device host (url or ip address)")
	flag.StringVar(&port, "port", DefaultPort, fmt.Sprintf("Slave device port (the default is %s)", DefaultPort))
	flag.Int64Var(&rate, "rate", DefaultRate, "Data collection rate in Second. > 10 required.")
	flag.Int64Var(&rate, "nums", DefaultDataNums, "Max number of data to collect.")

	flag.Parse()
	if rate < MinRate {
		rate = MinRate
	}
	if dataNums > MaxDataNums {
		dataNums = MaxDataNums
	}
	mbHandler := modbus.NewTCPClientHandler(host + ":" + port)
	mbHandler.Timeout = 10 * time.Second
	mbHandler.SlaveId = 1

	var err error

	if err = mbHandler.Connect(); err != nil {
		log.Fatal("Connect error:", err)
	}
	defer mbHandler.Close()

	client := modbus.NewClient(mbHandler)
	printOACsvHeader()

	//-- processing 8 chan per specified rate
	// x8Data := make([]byte, 64)
	for i := 0; i < 100; i++ {
		readX8OA(client)

		fmt.Println()
		time.Sleep(time.Duration(rate) * time.Second)
	}
}

func readX8OA(client modbus.Client) {
	x8Data, err := client.ReadHoldingRegisters(StartOAAddress, RegQuantity)

	if err != nil {
		fmt.Println("Read holding reg error.", err)
	}
	if len(x8Data) != 64 {
		fmt.Println("x8Data length error.", len(x8Data))
	}

	fmt.Printf("%s,", time.Now().Format("2006-01-02 15:04:05"))

	for i := 0; i < 8; i++ { //8 chan
		for j := 0; j < 2; j++ {
			// dout := uint32(d[0+i*4]) | uint32(d[1+i*4])<<8 | uint32(d[2+i*4])<<16 | uint32(d[3+i*4])<<24
			dout := uint32(x8Data[0+j*4+i*8]) | uint32(x8Data[1+j*4+i*8])<<8 | uint32(x8Data[2+j*4+i*8])<<16 | uint32(x8Data[3+j*4+i*8])<<24
			fout := math.Float32frombits(dout)
			fmt.Printf("%.2f, ", fout)
		}
	}
}

func printOACsvHeader() {
	fmt.Printf("X8II OA, Bias/Gap Modbus Data\n----------------------------\n")
	fmt.Print("Time,")
	for i := 1; i <= 8; i++ {
		fmt.Printf("OA(ch%d), Bias/Gap(ch%d),", i, i)
	}
	fmt.Println()
}
