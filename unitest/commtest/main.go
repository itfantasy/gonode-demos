package main

import (
	"fmt"

	"github.com/itfantasy/gonode/components/logger"
)

func main() {
	testTheInterfaceWithNil()
}

type testAStructContainsAnInterface struct {
	netReporter logger.INetReporter
}

func testTheInterfaceWithNil() {
	var myStruct *testAStructContainsAnInterface = new(testAStructContainsAnInterface)
	fmt.Println(myStruct.netReporter == nil)
}
