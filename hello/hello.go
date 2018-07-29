package main

import (
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"fmt"
)

type Hello struct {

}

func (t *Hello)Init(stud shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (t *Hello)Invoke(stud shim.ChaincodeStubInterface)peer.Response  {
	fmt.Println("TestA Invoke noding")
	fun, args := stud.GetFunctionAndParameters()
	if fun == "hello"{
		return sayHi(stud,args[0])
	}else {
		return shim.Success([]byte("fuck----------"))
	}
}

func sayHi(stud shim.ChaincodeStubInterface, str string) peer.Response {
	return shim.Success([]byte(str))
}
