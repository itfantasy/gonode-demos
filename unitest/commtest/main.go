package main

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/itfantasy/gonode/components/logger"
	"github.com/itfantasy/gonode/gnbuffers"
)

func main() {
	//testTheInterfaceWithNil()
	//testTheByteBuffer()
	//testTheGnBuffers()
	testTheGnBuffersObject()
}

type testAStructContainsAnInterface struct {
	netReporter logger.INetReporter
}

func testTheInterfaceWithNil() {
	var myStruct *testAStructContainsAnInterface = new(testAStructContainsAnInterface)
	fmt.Println(myStruct.netReporter == nil)
}

func testTheByteBuffer() {
	bytess := make([]byte, 1024)
	bytesBuffer := bytes.NewBuffer(bytess)
	bytesBuffer.Reset()
	err := binary.Write(bytesBuffer, binary.BigEndian, (int32)(16))
	if err != nil {
		fmt.Println(err)
	}
	buffer := ([]byte)("威猛的王大侠")
	binary.Write(bytesBuffer, binary.BigEndian, (int32)(len(buffer)))
	err2 := binary.Write(bytesBuffer, binary.BigEndian, (buffer))
	if err2 != nil {
		fmt.Println(err2)
	}
	fmt.Println(bytesBuffer.Bytes())
	fmt.Println(bytess)

	// -------------- test the reader
	var value int32
	binary.Read(bytesBuffer, binary.BigEndian, &value)
	fmt.Println(value)
	var lenth int32
	binary.Read(bytesBuffer, binary.BigEndian, &lenth)
	fmt.Println(lenth)
	var bytes []byte = make([]byte, lenth)
	binary.Read(bytesBuffer, binary.BigEndian, &bytes)
	fmt.Println(bytes)
	fmt.Println(string(bytes))
}

func testTheGnBuffers() {
	gnbuffer, err := gnbuffers.BuildBuffer(1024)
	if err != nil {
		fmt.Println(err.Error())
	}
	gnbuffer.PushInt(32)
	gnbuffer.PushLong(0xAA5555AA)
	gnbuffer.PushString("威猛的王大侠,     犀利的东哥")

	gnparser := gnbuffers.BuildParser(gnbuffer.Bytes(), 0)
	if val, err := gnparser.Int(); err == nil {
		fmt.Println(val)
	}
	if val, err := gnparser.Long(); err == nil {
		fmt.Println(val)
	}
	if val, err := gnparser.String(); err == nil {
		fmt.Println(val)
	}
}

func testTheGnBuffersObject() {
	gnbuffer, err := gnbuffers.BuildBuffer(1024)
	if err != nil {
		fmt.Println(err.Error())
	}
	if err := gnbuffer.PushObject(float32(32)); err != nil {
		fmt.Println(err.Error())
	}

	gnbuffer.PushObject(int64(0xAA5555AA))
	gnbuffer.PushObject("威猛的王大侠,     犀利的东哥")

	gnparser := gnbuffers.BuildParser(gnbuffer.Bytes(), 0)
	if val, err := gnparser.Object(); err == nil {
		fmt.Println("A")
		fmt.Println(val)
	} else {
		fmt.Println("AE")
		fmt.Println(err.Error())
	}
	if val, err := gnparser.Object(); err == nil {
		fmt.Println("B")
		fmt.Println(val)
	} else {
		fmt.Println("BE")
		fmt.Println(err.Error())
	}
	if val, err := gnparser.Object(); err == nil {
		fmt.Println("C")
		fmt.Println(val)
	} else {
		fmt.Println("CE")
		fmt.Println(err.Error())
	}
}
