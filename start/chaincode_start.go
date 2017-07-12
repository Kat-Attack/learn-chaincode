/*
Copyright IBM Corp 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type Task struct {
	Uid         string `json:"id"`
	User        string `json:"email"`
	Amount      int    `json:"amount"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Submissions string `json:"submissions"`
	CompletedBy string `json:"completed_by"` // either null or a user's email
}

type Marketplace struct {
	Tasks []string `json:"tasks"`
}

type CompletedTasks struct { // all tasks here shoud have Completed by not null
	Tasks []string `json:"tasks"`
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// ============================================================================================================================
// Init - reset all the things
// ============================================================================================================================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var Aval int
	var err error

	// if len(args) != 1 {
	// 	return nil, errors.New("Incorrect number of arguments. Expecting 1")
	// }

	// Initialize the chaincode
	Aval, err = strconv.Atoi(args[0])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding")
	}

	// Write the state to the ledger
	err = stub.PutState("abc", []byte(strconv.Itoa(Aval))) //making a test var "abc", I find it handy to read/write to it right away to test the network
	if err != nil {
		return nil, err
	}

	// var empty []string
	// jsonAsBytes, _ := json.Marshal(empty)								//marshal an emtpy array of strings to clear the index
	// err = stub.PutState(marbleIndexStr, jsonAsBytes)
	// if err != nil {
	// 	return nil, err
	// }

	// var trades AllTrades
	// jsonAsBytes, _ = json.Marshal(trades)								//clear the open trade struct
	// err = stub.PutState(openTradesStr, jsonAsBytes)
	// if err != nil {
	// 	return nil, err
	// }

	return nil, nil
}

// ============================================================================================================================
// Invoke - Our entry point to invoke a chaincode function (eg. write, createAccount, etc)
// ============================================================================================================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.write(stub, args)
	}
	// else if function == "createAccount" {
	// 	return t.CreateAccount(stub, args)
	// } else if function == "createProduct" {
	// 	return t.CreateProduct(stub, args)
	// } else if function == "purchaseProduct" { 						// for rewards (includes turing 100 savings into 100 spendings)
	// 	return t.PurchaseProduct(stub, args)
	// } else if function == "addAllowance" {							// transactions from admin panel
	// 	return t.AddAllowance(stub, args)
	// } else if function == "exchange" {
	// 	return t.Exchange(stub, args)
	// } else if function == "deposit" {
	// 	return t.Deposit(stub, args)
	// } else if function == "set_user" {										//change owner of a marble
	// 	res, err := t.set_user(stub, args)											//lets make sure all open trades are still valid
	// 	return res, err
	// }
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// ============================================================================================================================
// Query - Our entry point for Queries
// ============================================================================================================================
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

// ============================================================================================================================
// Write - Invoke function to write
// ============================================================================================================================
func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, value string
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	key = args[0] //rename for funsies
	value = args[1]
	err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// ============================================================================================================================
// Read - read a variable from chaincode state
// ============================================================================================================================

func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}

//from marbles and our cc:
// func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
// 	var name, jsonResp string
// 	var err error

// 	if len(args) != 1 {
// 		return nil, errors.New("Incorrect number of arguments. Expecting name of the var to query")
// 	}

// 	name = args[0]
// 	valAsbytes, err := stub.GetState(name)									//get the var from chaincode state
// 	if err != nil {
// 		jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
// 		return nil, errors.New(jsonResp)
// 	}

// 	return valAsbytes, nil													//send it onward
// }

func (t *SimpleChaincode) add_task(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	//   0       1       2         3           4              5              6
	// "uid", "user", "amount", "title", "description", "submissions", "completed_by" (completed_by + submissions expected to be null)
	if len(args) != 6 {
		return nil, errors.New("Incorrect number of arguments. Expecting 6")
	}

	fmt.Println("- create and add task")

	uid := args[0]
	user := args[1]
	amount := args[2]
	title := args[3]
	description := args[4]
	submissions := args[5]
	completed_by := args[6]

	//check if task already exists
	taskAsBytes, err := stub.GetState(uid)
	if err != nil {
		return nil, errors.New("Failed to get task id")
	}
	res := Task{}
	json.Unmarshal(taskAsBytes, &res)
	if res.Uid == uid {
		fmt.Println("This task arleady exists: id = " + uid)
		fmt.Println(res)
		return nil, errors.New("This task arleady exists") //all stop a task by this name exists
	}

	//build the marble json string manually
	str := `{"uid": "` + uid + `", "user": "` + user + `", "amount": ` + strconv.Itoa(amount) + `, "title": "` + title + `",
			"description": "` + description + `", "submissions": "` + submissions + `", "completed_by": "` + completed - by + `"}`
	err = stub.PutState(uid, []byte(str)) //store task with id as key
	if err != nil {
		return nil, err
	}

}
