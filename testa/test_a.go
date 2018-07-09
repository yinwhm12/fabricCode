package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"fmt"
	"strconv"
)

type SimpleChaincode struct {

}

func (t *SimpleChaincode)Init(stud shim.ChaincodeStubInterface) pb.Response  {
	fmt.Println("hello SimpleChaincode")
	_, args := stud.GetFunctionAndParameters()
	var A string
	var Aval int
	var err error
	if len(args) != 2{
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	A = args[0]
	Aval, err = strconv.Atoi(args[1])
	if err != nil{
		return shim.Error("Expecting integer value for asset holding")
	}

	fmt.Printf("Aval = %d\n", Aval)
	err = stud.PutState(A,[]byte(strconv.Itoa(Aval)))
	if err != nil{
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (t *SimpleChaincode)Invoke(stud shim.ChaincodeStubInterface)pb.Response  {
	fmt.Println("SimpleChaincode Invoke noding")
	fun, args := stud.GetFunctionAndParameters()
	if fun == "query"{
		return t.query(stud,args)
	}else{
		channelName := ""
		if len(args) >= 3{
			channelName = args[2]
		}
		return stud.InvokeChaincode(fun,toChaincodeArgs(args[1]),channelName)
	}
}

func toChaincodeArgs(args ...string) [][]byte {
	bargs := make([][]byte, len(args))
	for i, arg := range args {
		bargs[i] = []byte(arg)
	}
	return bargs
}



func (t *SimpleChaincode)query(stub shim.ChaincodeStubInterface, args []string)pb.Response  {
	var A string // Entities
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	A = args[0]

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success(Avalbytes)

}

func main()  {
	err := shim.Start(new(SimpleChaincode))
	if err != nil{
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}

}