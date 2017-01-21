/*
Copyright Capgemini India. 2017 All Rights Reserved.
*/

package main

import (
	"errors"
	"fmt"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// Loyalty Program implementation
type LoyaltyProgramChaincode struct {
}

var merchantIndexTxStr = "_merchantIndexTxStr"

type MerchantData struct {
	MERCHANT_NAME string `json:"MERCHANT_NAME"`
	MERCHANT_CITY string `json:"MERCHANT_CITY"`
	MERCHANT_PHONE string `json:"MERCHANT_PHONE"`	
}


func (t *LoyaltyProgramChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	
	var err error
	// Initialize the chaincode
	
	fmt.Printf("Deployment of Loyalty Program is completed\n")
	
	
	// For Merchant Initialization
	var emptyMerchantDataTxs []MerchantData
	jsonAsBytes, _ := json.Marshal(emptyMerchantDataTxs)
	err = stub.PutState(merchantIndexTxStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}
	
	return nil, nil
}

// Add Merchant data in BLockChain
func (t *LoyaltyProgramChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	
	if function == "AddMerchant" {		
		return t.AddNewMerchantDetails(stub, args)
	}
	return nil, nil
}


func (t *LoyaltyProgramChaincode) AddNewMerchantDetails(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	
	var MerchantDataObj MerchantData
	var MerchantDataList []MerchantData
	var err error

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Need 14 arguments")
	}

	// Initialize the chaincode  for Merchant data
	MerchantDataObj.MERCHANT_NAME = args[0]
	MerchantDataObj.MERCHANT_CITY = args[1]
	MerchantDataObj.MERCHANT_PHONE = args[2]
	
	fmt.Printf("Input from user:%s\n", MerchantDataObj)
	
	merchantTxsAsBytes, err := stub.GetState(merchantIndexTxStr)
	if err != nil {
		return nil, errors.New("Failed to get consumer Transactions")
	}
	json.Unmarshal(merchantTxsAsBytes, &MerchantDataList)
	
	MerchantDataList = append(MerchantDataList, MerchantDataObj)
	jsonAsBytes, _ := json.Marshal(MerchantDataList)
	
	err = stub.PutState(merchantIndexTxStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// Query callback representing the query of a chaincode - for Merchant
func (t *LoyaltyProgramChaincode) Query(stub shim.ChaincodeStubInterface,function string, args []string) ([]byte, error) {
	
	var MerchantName string // Entities
	var err error
	var resAsBytes []byte

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the person to query")
	}

	MerchantName = args[0]
	
	resAsBytes, err = t.GetMerchantDetails(stub, MerchantName)
	
	fmt.Printf("Query Response:%s\n", resAsBytes)
	
	if err != nil {
		return nil, err
	}
	
	return resAsBytes, nil
}

func (t *LoyaltyProgramChaincode)  GetMerchantDetails(stub shim.ChaincodeStubInterface, MerchantName string) ([]byte, error) {
	
	//var requiredObj MerchantData
	var objFound bool
	MerchantTxsAsBytes, err := stub.GetState("read")
	if err != nil {
		return nil, errors.New("Failed to get Merchant Transactions")
	}
	var MerchantTxObjects []MerchantData
	var MerchantTxObjects1 []MerchantData
	json.Unmarshal(MerchantTxsAsBytes, &MerchantTxObjects)
	length := len(MerchantTxObjects)
	fmt.Printf("Output from chaincode: %s\n", MerchantTxsAsBytes)
	
	if MerchantName == "" {
		res, err := json.Marshal(MerchantTxObjects)
		if err != nil {
		return nil, errors.New("Failed to Marshal the required Obj")
		}
		return res, nil
	}
	
	objFound = false
	// iterate
	for i := 0; i < length; i++ {
		obj := MerchantTxObjects[i]
		if MerchantName == obj.MERCHANT_NAME {
			MerchantTxObjects1 = append(MerchantTxObjects1,obj)
			//requiredObj = obj
			objFound = true
		}
	}
	
	if objFound {
		res, err := json.Marshal(MerchantTxObjects1)
		if err != nil {
		return nil, errors.New("Failed to Marshal the required Obj")
		}
		return res, nil
	} else {
		res, err := json.Marshal("No Data found")
		if err != nil {
		return nil, errors.New("Failed to Marshal the required Obj")
		}
		return res, nil
	}
}



// #############################################################################


func main() {
	err := shim.Start(new(LoyaltyProgramChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
