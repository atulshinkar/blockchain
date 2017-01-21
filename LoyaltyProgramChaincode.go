/*
Copyright Capgemini India. 2016 All Rights Reserved.
*/

package main

import (
	"errors"
	"fmt"
	"encoding/json"
	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// Loyalty Program implementation
type LoyaltyProgramChaincode struct {
}

var merchantIndexTxStr = "_merchantIndexTxStr"
var userIndexTxStr = "_userIndexTxStr"
var loginIndexTxStr = "_loginIndexTxStr"
var transferIndexTxStr = "_transferIndexTxStr"

type MerchantData struct {
	MERCHANT_NAME string `json:"MERCHANT_NAME"`
	MERCHANT_CITY string `json:"MERCHANT_CITY"`
	MERCHANT_PHONE string `json:"MERCHANT_PHONE"`	
}

type UserData struct {
	NAME string `json:"NAME"`
	PHONENO string `json:"PHONENO"`
	USERNAME string `json:"USERNAME"`
	PASSWORD string `json:"PASSWORD"`
	MERCHANTNAME string `json:"MERCHANTNAME"`
	POINTS float64 `json:"POINTS"`
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
	
	//For User Initialization
	var emptyUserTxs []UserData
	jsonAsBytes2, _ := json.Marshal(emptyUserTxs)
	err = stub.PutState(userIndexTxStr, jsonAsBytes2)
	if err != nil {
		return nil, err
	}
	
	return nil, nil
}

// Add Merchant data in BLockChain
func (t *LoyaltyProgramChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	
	if function == "AddMerchant" {		
		return t.AddNewMerchantDetails(stub, args)
	} else if function == "AddUser" {		
		return t.RegisterUser(stub, args)
	} else if function == "Login" {		
		return t.Login(stub, args)
	} else if function == "Transfer" {		
		return t.Transfer(stub, args)
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

func (t *LoyaltyProgramChaincode) RegisterUser(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	
	var User UserData
	var UserDataList []UserData
	var err error

	if len(args) != 5 {
		return nil, errors.New("Incorrect number of arguments. Need 5 arguments")
	}
	
	// Initialize the chaincode  for User data
	User.NAME = args[0]
	User.PHONENO = args[1]
	User.USERNAME = args[2]
	User.PASSWORD = args[3]
	User.MERCHANTNAME = args[4]
	 if User.MERCHANTNAME == "KMT" {
		User.POINTS = 100
	 } else if User.MERCHANTNAME == "SMC" {
		User.POINTS = 150
	 } 
	
	fmt.Printf("Input from user:%s\n", User)
	
	userTxsAsBytes, err := stub.GetState(userIndexTxStr)
	if err != nil {
		return nil, errors.New("Failed to get user data")
	}
	json.Unmarshal(userTxsAsBytes, &UserDataList)
	
	UserDataList = append(UserDataList, User)
	jsonAsBytes2, _ := json.Marshal(UserDataList)
	
	err = stub.PutState(userIndexTxStr, jsonAsBytes2)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (t *LoyaltyProgramChaincode) Login(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	
	var username string
	var password string
	var users UserData
	var err error

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Need 2 arguments")
	}
	
	username = args[0]
	password = args[1]
	
	users, err = t.GetUserDetails(stub, username) 

	if users.USERNAME == username {
		if users.PASSWORD == password {
				return nil, nil
			}
	}	
	
	fmt.Printf("Output from chaincode: %s\n", err)
	return nil , errors.New("Incorrect Username or Password")
	
}

func (t *LoyaltyProgramChaincode) Transfer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	
	var touser string
	var pointstotransfer float64
	var err error
	var currentuser string
	var user UserData
	var user2 UserData
	var pt1 float64
	var pt2 float64

	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Need 4 arguments")
	}
	
	currentuser  = args[0]
	touser  = args[1]
	pointstotransfer,err  = strconv.ParseFloat(args[3], 64)
	
	user, err = t.GetUserDetails(stub, currentuser) 
	
	if user.POINTS<=pointstotransfer {
		user.POINTS = user.POINTS - pointstotransfer;
		if user.MERCHANTNAME == "KMT" {
		pt1 = user.POINTS * 0.1 
		} else if user.MERCHANTNAME == "SMC" {
			pt1 = user.POINTS * 0.15 
		}
	}

	
	user2, err = t.GetUserDetails(stub, touser) 
	
	if user2.MERCHANTNAME == "KMT" {
		pt2 = user2.POINTS * 0.1 
		pt2 = pt2 + pt1
		user2.POINTS = pt2 / 0.1
	} else if user2.MERCHANTNAME == "SMC" {
		pt2 = user2.POINTS * 0.15 
		pt2 = pt2 + pt1
		user2.POINTS = pt2 / 0.15
	}

	res,err := json.Marshal(user2)
	err = stub.PutState(userIndexTxStr, res)
	if err != nil {
		return nil, err
	}
	return nil, nil
	
}


// Query callback representing the query of a chaincode - for Merchant
func (t *LoyaltyProgramChaincode) Query(stub shim.ChaincodeStubInterface,function string, args []string) ([]byte, error) {
	
	var MerchantName string // Entities
	var UserName string // Entities
	var err error
	var resAsBytes []byte

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the person to query")
	}

	if function == "GetMerchantDetails" {		
		MerchantName = args[0]	
		resAsBytes, err = t.GetMerchantDetails(stub, MerchantName)
	} else if function == "GetUserDetails" {		
		UserName = args[0]	
		resAsBytes, err := json.Marshal(t.GetUserDetails(stub, UserName))
	} 
	
	fmt.Printf("Query Response:%s\n", resAsBytes)
	
	if err != nil {
		return nil, err
	}
	
	return resAsBytes, nil
}

func (t *LoyaltyProgramChaincode)  GetMerchantDetails(stub shim.ChaincodeStubInterface, MerchantName string) ([]byte, error) {
	
	//var requiredObj MerchantData
	var objFound bool
	MerchantTxsAsBytes, err := stub.GetState(merchantIndexTxStr)
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

func (t *LoyaltyProgramChaincode)  GetUserDetails(stub shim.ChaincodeStubInterface, username string) (UserData, error) {
	
	var objFound bool
	var UserTxObjects []UserData
	var UserTxObjects1 []UserData
	var currentuser UserData
	UserTxsAsBytes, err := stub.GetState(userIndexTxStr)
	if err != nil {
		return  currentuser,errors.New("Failed to get Merchant Transactions")
	}
	
	json.Unmarshal(UserTxsAsBytes, &UserTxObjects)
	length := len(UserTxObjects)
	fmt.Printf("Output from chaincode: %s\n", UserTxsAsBytes)
	
	if username == "" {
		res, err := json.Marshal(UserTxObjects)
		if err != nil {
		return currentuser, errors.New("Failed to Marshal the required Obj")
		}
		fmt.Printf("Output from chaincode: %s\n", res)
		return currentuser, nil
	}
	
	objFound = false
	// iterate
	for i := 0; i < length; i++ {
		obj := UserTxObjects[i]
		if username == obj.USERNAME {
			UserTxObjects1 = append(UserTxObjects1,obj)
			//requiredObj = obj
			objFound = true
			currentuser = obj
		}
	}
	
	if objFound {
		if err != nil {
		return currentuser, errors.New("Failed to Marshal the required Obj")
		}
		return currentuser, nil
	} else {
		res, err := json.Marshal("No Data found")
		if err != nil {
		return currentuser, errors.New("Failed to Marshal the required Obj")
		}
		fmt.Printf("Output from chaincode: %s\n", res)
		return currentuser, nil
	}
}


// #############################################################################


func main() {
	err := shim.Start(new(LoyaltyProgramChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
