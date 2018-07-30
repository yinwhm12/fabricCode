package main

import (
	"github.com/hyperledger/fabric/bccsp"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/hyperledger/fabric/core/chaincode/shim/ext/entities"
	"fmt"
	"encoding/json"
	"errors"
	"github.com/hyperledger/fabric/bccsp/factory"
)

type EncCC struct {
	baccspInst bccsp.BCCSP
}

type ResumeData struct {
	ID string
	Year string
	WorkPlace string
	Career string
}

const DECKEY = "DECKEY"
const ENCKEY = "ENCKEY"
const IV = "IV"

func (t *EncCC)Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// encryptAndPutState encrypts the supplied value using the
// supplied entity and puts it to the ledger associated to
// the supplied KVS key
func encryptAndPutState(stub shim.ChaincodeStubInterface, ent entities.Encrypter, key string, value []byte) error {
	// at first we use the supplied entity to encrypt the value
	ciphertext, err := ent.Encrypt(value)
	if err != nil {
		return err
	}
	return stub.PutState(key, ciphertext)
}

func getStateAndDecrypt(stub shim.ChaincodeStubInterface, ent entities.Encrypter, key string) ([]byte, error) {
	// at first we retrieve the ciphertext from the ledger
	ciphertext, err := stub.GetState(key)
	if err != nil {
		return nil, err
	}

	// GetState will return a nil slice if the key does not exist.
	// Note that the chaincode logic may want to distinguish between
	// nil slice (key doesn't exist in state db) and empty slice
	// (key found in state db but value is empty). We do not
	// distinguish the case here
	if len(ciphertext) == 0 {
		return nil, errors.New("no ciphertext to decrypt")
	}

	return ent.Decrypt(ciphertext)
}

func (t *EncCC)EncRecord(stub shim.ChaincodeStubInterface,args []string,encKey, IV []byte)pb.Response  {
	ent, err := entities.NewAES256EncrypterEntity("ID",t.baccspInst,encKey,IV)
	if err != nil{
		return shim.Error(fmt.Sprintf("entities.NewAES256EncrypterEntity failed, err %s", err))
	}

	if len(args) != 4 {
		return shim.Error("Expected 4 parameters to function EncRecord")
	}
	key := args[0]+args[1]
	err, existed := checkExistedUniqueByID(stub,key)
	if err != nil{
		return shim.Error(fmt.Sprintf("get key = %s,state err %s",key, err))
	}
	if existed {
		return shim.Error(fmt.Sprintf("can not add this record with the key = %s,because had existed",key))
	}
	resumeData := ResumeData{
		ID:args[0],
		Year:args[1],
		WorkPlace:args[2],
		Career:args[3],
	}
	value, err := json.Marshal(&resumeData)
	if err != nil{
		return shim.Error(fmt.Sprintf("Json Marshal failed, err %s", err))
	}

	err = encryptAndPutState(stub,ent,key,value)
	if err != nil{
		return shim.Error(fmt.Sprintf("encryptAndPutState failed, err %+v", err))
	}
	return shim.Success(nil)
}

func (t *EncCC)DecRecord(stub shim.ChaincodeStubInterface, args []string, decKey, IV []byte)pb.Response  {
	ent, err := entities.NewAES256EncrypterEntity("ID",t.baccspInst,decKey,IV)
	if err != nil {
		return shim.Error(fmt.Sprintf("entities.NewAES256EncrypterEntity failed, err %s", err))
	}
	if len(args) != 2 {
		return shim.Error("Expected 2 parameters to function Decrypter")
	}
	key := args[0]+args[1]
	valueByte, err := getStateAndDecrypt(stub,ent,key)
	if err != nil{
		return shim.Error(fmt.Sprintf("getStateAndDecrypt failed, err %+v", err))
	}
	var resumeData ResumeData
	err = json.Unmarshal(valueByte,&resumeData)
	if err != nil{
		return shim.Error(fmt.Sprintf("Json unmarshal failed, err %s", err))
	}
	return shim.Success([]byte(resumeData.WorkPlace))
}

// 检查是否 已存在账本 存在返回(nil,true) 不存在返回(nil,false) 错误返回(err,false)
func checkExistedUniqueByID(stub shim.ChaincodeStubInterface,id string) (error,bool ){
	vbyte, err := stub.GetState(id)
	if err != nil{
		return err,false
	}
	if vbyte == nil{
		return nil,false
	}
	return nil, true
}

func (t *EncCC)AddRecord(stub shim.ChaincodeStubInterface, args []string)pb.Response  {
	if len(args) != 4 {
		return shim.Error("Expected 4 parameters to function AddRecord")
	}
	resumeData := ResumeData{
		ID:args[0],
		Year:args[1],
		WorkPlace:args[2],
		Career:args[3],
	}
	key := args[0]+args[1]
	err,existed := checkExistedUniqueByID(stub,key)
	if err != nil{
		return shim.Error(fmt.Sprintf("get key = %s,state err %s",key, err))
	}
	//不存在 可以操作
	if !existed{
		vbyte, err := json.Marshal(resumeData)
		if err != nil{
			return shim.Error(fmt.Sprintf("Json Marshal failed, err %s", err))
		}
		err = stub.PutState(key,vbyte)
		return shim.Success(nil)
	}
	return shim.Error(fmt.Sprintf("can not add this record with the key = %s,because had existed",key))
}

func (t *EncCC)GetRecord(stub shim.ChaincodeStubInterface, args []string)pb.Response  {
	if len(args) != 2 {
		return shim.Error("Expected 2 parameters to function AddRecord")
	}
	key := args[0]+args[1]
	vbyte, err := stub.GetState(key)
	if err != nil{
		return shim.Error(fmt.Sprintf("get key = %s,state err %s",key, err))
	}
	if vbyte == nil{
		return shim.Error(fmt.Sprintf("can not find the key = %s",key))
	}
	var resumeData ResumeData
	err = json.Unmarshal(vbyte,&resumeData)
	if err != nil{
		return shim.Error(fmt.Sprintf("Json unmarshal failed, err %s", err))
	}
	return shim.Success([]byte(resumeData.WorkPlace))
}

func (t *EncCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	f, args := stub.GetFunctionAndParameters()
	var tMap map[string][]byte
	var err error
	if f == "encRecord" || f == "decRecord"{
		tMap, err = stub.GetTransient()
		if err != nil {
			return shim.Error(fmt.Sprintf("Could not retrieve transient, err %s", err))
		}
	}
	switch f {
	case "addRecord":
		return t.AddRecord(stub,args)
	case "getRecord":
		return t.GetRecord(stub,args)
	case "enckey":
		//if tMap[ENCKEY] == nil || tMap[IV] == nil{
		//	return shim.Error("Please input ENCKEY and IV values")
		//}
		//fmt.Println("enckey:----------------------------------",tMap[ENCKEY],"---------------------len:",len(tMap[ENCKEY]))
		//fmt.Println("iv :----------",tMap[IV],"----------len:",len(tMap[IV]))
		return shim.Success(tMap[ENCKEY])
		//return t.EncRecord(stub,args,tMap[ENCKEY],tMap[IV])
		//for k, v := range tMap{
		//	if len(strings.TrimSpace(k))<1 || v == nil{
		//
		//	}
		//	return t.EncRecord(stub,args,[]byte(k),v)
		//}
		//return shim.Error("Please input ENCKEY and IV values")
	case "iv":
		return shim.Success(tMap[IV])
		//if tMap[DECKEY] == nil || tMap[IV] == nil{
		//	return shim.Error("Please input ENCKEY and IV values")
		//}
		//return t.DecRecord(stub,args,tMap[DECKEY],tMap[IV])
		//return t.EncRecord(stub,args,tMap[ENCKEY],tMap[IV])
		//for k, v := range tMap{
		//	if len(strings.TrimSpace(k))<1 || v == nil{
		//		return shim.Error("Please input ENCKEY and IV values")
		//	}
		//	return t.DecRecord(stub,args,[]byte(k),v)
		//}
		//return shim.Error("Please input ENCKEY and IV values")
	case "deckey":
		return shim.Success(tMap[DECKEY])

	default:
		return shim.Error(fmt.Sprintf("Unsupported function %s", f))
	}
}

func main()  {
	factory.InitFactories(nil)
	err := shim.Start(&EncCC{factory.GetDefault()})
	if err != nil {
		fmt.Printf("Error starting EncCC chaincode: %s", err)
	}
}