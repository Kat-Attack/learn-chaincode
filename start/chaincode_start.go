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

var MarketplaceStr = "_marketplace"        // name for the key/value that will store all open tasks
var CompletedTasksStr = "_completedTasks " // name for the key/value that will store all completed tasks

type Task struct {
	Uid         string `json:"id"`
	User        string `json:"email"` // users are defined by their emails
	Amount      int    `json:"amount"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Submissions string `json:"submissions"`
	CompletedBy string `json:"completed_by"` // either null or a user's email
}

type Marketplace struct {
	Tasks []Task `json:"marketplace_tasks"`
}

type CompletedTasks struct { // all tasks here shoud have Completed by not null
	Tasks []Task `json:"tasks"`
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

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

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

	var cTasks CompletedTasks
	jsonAsBytes, _ := json.Marshal(cTasks) //clearr the CompletedTasks struct
	err = stub.PutState(CompletedTasksStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}

	var mplace Marketplace
	jsonAsBytes, _ = json.Marshal(mplace) //clear the Marketplace struct
	err = stub.PutState(MarketplaceStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}

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
	} else if function == "add_task" {
		return t.add_task(stub, args)
	} else if function == "update_task" {
		return t.update_task(stub, args)
	}
	// else if function == "purchaseProduct" { 						// for rewards (includes turing 100 savings into 100 spendings)
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
// Read - read a variable from chaincode state (used by Query)
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
// add_task - creates a task and adds it to Marketplace struct
// ============================================================================================================================

func (t *SimpleChaincode) add_task(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	fmt.Println("Number of args: ")
	fmt.Println(len(args))
	fmt.Println(args[0])
	fmt.Println(args[1])
	fmt.Println(args[2])
	fmt.Println(args[3])
	fmt.Println(args[4])
	fmt.Println(args[5])
	fmt.Println(args[6])
	//   0       1       2         3           4              5              6
	// "uid", "user", "amount", "title", "description", "submissions", "completed_by" (completed_by + submissions expected to be null)
	if len(args) < 5 {
		return nil, errors.New("Incorrect number of arguments. Expecting 5 to 7 (last two should be empty)")
	}

	fmt.Println("- create and add task")

	amount, err := strconv.Atoi(args[2])
	if err != nil {
		return nil, errors.New("3rd argument (amount) must be a numeric string")
	}

	var task = Task{}
	task.Uid = args[0]
	task.User = args[1]
	task.Amount = amount
	task.Title = args[3]
	task.Description = args[4]
	task.Submissions = args[5]
	task.CompletedBy = args[6]

	fmt.Println("below is task: ")
	fmt.Println(task)

	/////////////// for debugging: query for _debug1 to see your parameters ///////
	jsonAsBytes, _ := json.Marshal(task)
	err = stub.PutState("_debug1", jsonAsBytes)
	//////////////// end for debugging ////////////////////////////////////////////

	//get the open trade struct
	MarketplaceAsBytes, err := stub.GetState(MarketplaceStr)
	if err != nil {
		return nil, errors.New("Failed to get marketplace")
	}
	var mplace Marketplace
	json.Unmarshal(MarketplaceAsBytes, &mplace) //un stringify it aka JSON.parse()

	mplace.Tasks = append(mplace.Tasks, task) //append to marketplace
	fmt.Println("! appended task to marketplace")
	jsonAsBytes, _ = json.Marshal(mplace)
	err = stub.PutState(MarketplaceStr, jsonAsBytes) //rewrite marketplace
	if err != nil {
		return nil, err
	}
	fmt.Println("- end open trade")
	return nil, nil

}

// ============================================================================================================================
// update_task - upload submission into task
// ============================================================================================================================

func (t *SimpleChaincode) update_task(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	// var fail Marketplace

	fmt.Println("Number of args: ")
	fmt.Println(len(args))
	fmt.Println(args[0])
	fmt.Println(args[1])

	//   0            1
	// "uid", "user_who_updated"
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2: task uid and user who updated")
	}

	fmt.Println("- update task")

	// get all active tasks in marketplace
	marketplaceAsBytes, err := stub.GetState(MarketplaceStr)
	if err != nil {
		return nil, errors.New("Failed to get marketplace array")
	}
	var mplace []string
	json.Unmarshal(marketplaceAsBytes, &mplace) //un stringify it aka JSON.parse()

	fmt.Println(mplace)
	for i := range mplace { //iter through all the tasks
		fmt.Println("looking @ task name: " + mplace[i])

		// marbleAsBytes, err := stub.GetState(mplace[i]) //grab this task
		// if err != nil {
		// 	return fail, errors.New("Failed to get task")
		// }
		// res := Task{}
		// json.Unmarshal(marbleAsBytes, &res) //un stringify it aka JSON.parse()
		// fmt.Println("looking @ : " + res.Uid + "," + res.User + ", " + strconv.Itoa(res.Amount)) + "," + res.Title + ", " + res.Description + "," + res.Submissions);

		// //check for user && color && size
		// if (res.Uid) == args[0] {
		// 	fmt.Println("found the task: " + res.Uid)
		// 	res.
		// 	fmt.Println("! end find marble 4 trade")
		// 	return res, nil
		// }
	}

	fmt.Println("- end update task")
	return nil, nil

}
