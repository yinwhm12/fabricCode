package testb

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"fmt"
	"strconv"
)

type TestB struct {

}

func (t *TestB)Init(stud shim.ChaincodeStubInterface) peer.Response  {
	fmt.Println("hello TestA")
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

func (t *TestB)Invoke(stud shim.ChaincodeStubInterface)peer.Response  {
	fmt.Println("TestA Invoke noding")
	fun, args := stud.GetFunctionAndParameters()
	if fun == "query"{
		return t.query(stud,args)
	}else{
		//return stud.InvokeChaincode(fun,toChaincodeArgs(args[0]),args[1])
		//return shim.Success([]byte("no any more functions"))
		if len(args) < 1{
			args[0] = "testB"
		}
		return t.hello(stud,args[0])
	}
}

func (t *TestB)hello(stud shim.ChaincodeStubInterface,str string)peer.Response  {
	return shim.Success([]byte(str))

}

func (t *TestB)query(stub shim.ChaincodeStubInterface, args []string)peer.Response  {
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