package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
)

var (
	port       string
	udp        UDP
	bytearr    []byte
	bytesArray []byte
	byteMatrix [][]byte
	blkcount   int
	blksize    int = 512 //blk Size of
)

const (
	maxUDP int = 65535
	RRQ        = uint16(1) //Op Codes
	WRQ        = uint16(2)
	DATA       = uint16(3)
	ACK        = uint16(4)
	ERROR      = uint16(5)
	OACK       = uint16(6)
)

/* TFTP Formats

   Type   Op #     Format without header

          2 bytes    string   1 byte     string   1 byte
          -----------------------------------------------
   RRQ/  | 01/02 |  Filename  |   0  |    Mode    |   0  |
   WRQ    -----------------------------------------------
          2 bytes    2 bytes       n bytes
          ---------------------------------
   DATA  | 03    |   Block #  |    Data    |
          ---------------------------------
          2 bytes    2 bytes
          -------------------
   ACK   | 04    |   Block #  |
          --------------------
          2 bytes  2 bytes        string    1 byte
          ----------------------------------------
   ERROR | 05    |  ErrorCode |   ErrMsg   |   0  |
          ----------------------------------------
*/

type UDP struct {
	OpCode  []byte
	opCode2 []byte
	reader  Reader
	writer  Writer
	data    Data
	ack     Ack
	errors  Error
}
type Reader struct {
	FileName string
	Mode     string
}
type Writer struct {
	FileName string
	Mode     string
}
type Data struct {
	blocknumber uint16
	data        []byte
}
type Ack struct {
	blocknumber []byte
}
type Error struct {
	errorcode []byte
	errmsg    string
}

//Combines ip Address with port for connection w/ UDP
func SetServer(ipAddress string, port string) (udpAddress string) {
	return string(ipAddress + ":" + port)
}

func udpParser(buffer []byte) {
	//If its the read or Write command it grabs that
	if int(buffer[1]) == 1 || int(buffer[1]) == 2 {
		udp.OpCode = buffer[:2]
	}
	//Grabs current Opcode
	udp.opCode2 = buffer[:2]
	println(string(udp.OpCode))
	switch uint16(udp.opCode2[1]) {
	case WRQ:
		log.Println(" WRQ Packet recieved ")
		//skips OP CODE and Splices by the 0x00
		buff := bytes.Split(buffer[2:], []byte{0})
		//grabs FileName from buff
		udp.writer.FileName = string(buff[0])
		//grabs Mode from buff
		udp.writer.Mode = string(buff[1])
		//tries to change blkSize if option is appended
		if buff[2] != nil && string(buff[2]) == "blksize" {
			blksize, _ = strconv.Atoi(string(buff[3]))
			log.Println(" blksize Changed = ", blksize)
		}
		fmt.Println(" Found the Mode and FileName " + udp.writer.Mode + " " + udp.writer.FileName)
	case RRQ:
		log.Println(" RRQ Packet recieved ")
		//skips OP CODE and Splices by the 0x00
		buff := bytes.Split(buffer[2:], []byte{0})
		//grabs FileName from buff
		udp.reader.FileName = string(buff[0])
		//grabs Mode from buff
		udp.reader.Mode = string(buff[1])

		//tries to change blkSize if option is appended
		if buff[2] != nil && string(buff[2]) == "blksize" {
			blksize, _ = strconv.Atoi(string(buff[3]))
			log.Println(" blksize Changed = ", blksize)
		}
		fmt.Println(" Found the Mode and FileName " + udp.reader.Mode + " " + udp.reader.FileName)
	case DATA:
		log.Println(" Data Packet recieved ")
		//skips OP CODE and BlockNumber and Splices by the 0x00
		buff := bytes.Split(buffer[4:], []byte{0})
		//grabs block number from buff
		udp.data.blocknumber = uint16(buffer[3])
		//grabs data from spliced buffer
		udp.data.data = buff[0]
	case ACK:
		log.Println(" ACK Packet recieved ")
		//grabs block number from buff
		udp.ack.blocknumber = buffer[2:4]
		log.Println(udp.ack.blocknumber)
	case ERROR:
		log.Println(" ERROR Packet recieved ")
		buff := bytes.Split(buffer[2:], []byte{0})
		//grabs errOpCode from buffer
		udp.errors.errorcode = buffer[2:4]
		//grabs error msg from buff
		udp.errors.errmsg = string(buff[1])
		fmt.Print(udp.errors.errmsg)
	case OACK:
		log.Println(" OACK Packet recieved ")
	}
}

//Creates the OS.File to Read From ./files/
func readRRQ(filepath string) (reader *os.File) {
	name := string("./files/" + filepath)
	reader, err := os.Open(name)
	if err != nil {
		log.Print(err)

	}
	return
}

//Creates the OS.File to Write From ./files/
func WriteWRQ(filepath string) (writer *os.File) {
	name := string("./files/" + filepath)
	writer, err := os.Create(name)
	if err != nil {
		log.Fatal(err)
	}
	return
}

//Grabs Bytes from File and Puts into chunks of size BLKsize
func getChunks() (byteMatrix [][]byte) {
	file := readRRQ(udp.reader.FileName)

	buffer := bytes.NewBuffer(bytesArray)

	buffer.ReadFrom(file)

	for len(buffer.Next(1)) < buffer.Len() {
		buffer.UnreadByte()
		byteMatrix = append(byteMatrix, buffer.Next(blksize))
	}

	return
}

//Uses Bytes from Client and Puts into chunks to send to outputFile
func putChunks(udpcon *net.UDPConn, cAddress *net.UDPAddr) {
	writer := WriteWRQ(udp.writer.FileName)

	//create first ACK packet
	headerBytes := make([]byte, 4)
	headerBytes[1] = byte(ACK)
	//sends Ack to Client to recieve Data
	udpcon.WriteToUDP(headerBytes, cAddress)
	//Sets current blknumber to 1 the first Data blknum should be
	currentblknum := uint16(1)
	for {
		//initilizes new buffer everyloop to not have old info in array
		buffer := make([]byte, maxUDP)

		_, _, err := udpcon.ReadFromUDP(buffer)

		if err != nil {
			log.Print(err)
		}
		//Parses AckPacket
		udpParser(buffer)
		log.Println("blknumber : ", udp.data.blocknumber, "currentblknum : ", currentblknum)
		//If data.blknum didn't update then Resends the Ack Packet
		if udp.data.blocknumber == currentblknum-1 {
			headerBytes := make([]byte, 4)
			headerBytes[1] = byte(ACK)
			udpcon.WriteToUDP(headerBytes, cAddress)
		}

		//writes the Data to file
		writer.Write(udp.data.data)
		//Updates data.blknum to be currentblk becuase of the Write
		udp.data.blocknumber = currentblknum
		//Updates Client with new ACK
		headerBytes := make([]byte, 4)
		headerBytes[1] = byte(ACK)
		headerBytes[3] = byte(currentblknum)
		udpcon.WriteToUDP(headerBytes, cAddress)
		//Exit condition is when Length of DATA is < BLKSize
		size := len(udp.data.data)
		log.Println("Size of Data : ", size)
		if size < blksize {
			return
		}
		//Deletes Data From Data.Data to make sure nothing is left inside before Parse is Called again
		udp.data.data = nil
		//Increments currBLKNum
		currentblknum++
	}
}

func main() {

	ipAddress := "localhost"
	port = "69"

	//tries to connect to a Client
	bindAdd, err := net.ResolveUDPAddr("udp", SetServer(ipAddress, port))
	if err != nil {
		log.Fatal(err)
	}

	udpcon, err := net.ListenUDP("udp", bindAdd)
	if err != nil {
		log.Fatal(err)
	}

	//boolean
	first := true
	for {
		// for UDP to server
		bytearr = make([]byte, maxUDP)
		//Reads UDP from Client
		_, cAddress, err := udpcon.ReadFromUDP(bytearr)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Connected to a UDP Client")
		//Parse UDP to find whether WRQ/RRQ
		udpParser(bytearr)
		//Either Goes to RRQ or WRQ depending on Client UDP parsed information
		switch uint16(udp.OpCode[1]) {

		case WRQ:
			//Passes connections so that putChunks can call WriteToUDP()
			putChunks(udpcon, cAddress)
			udpcon.Close()
			return
		case RRQ:
			//Creates New Header every Loop
			headerBytes := make([]byte, 4)
			headerBytes[1] = byte(DATA)
			//Grabs entire File only when first Loop
			if first {
				byteMatrix = getChunks()
				first = false
			}
			//increments BLKCount which is currentblknumber
			blkcount++
			//Sets Header to have Value of currentblknumber
			headerBytes[3] = byte(blkcount)
			//appends the correct chunk from file with current Header
			byteMatrix[blkcount-1] = append(headerBytes, byteMatrix[blkcount-1]...)
			//Writes to Client the current Data tftp
			udpcon.WriteToUDP(byteMatrix[blkcount-1], cAddress)
			/*Exit condition is when the currentBlkNum is = to the length of
			the amount of bytes of byteMatrix(blk sized chunks of File bytes)*/
			if blkcount == len(byteMatrix) {
				udpcon.Close()
			}
		}
	}
}
