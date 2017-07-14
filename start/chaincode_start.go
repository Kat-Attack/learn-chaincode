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

var taskIndexStr = "_taskindex"            //name for the key/value that will store a list of all known tasks
var MarketplaceStr = "_marketplace"        // name for the key/value that will store all open tasks
var CompletedTasksStr = "_completedTasks " // name for the key/value that will store all completed tasks

type Task struct {
	Uid         string   `json:"id"`
	User        string   `json:"email"` // users are defined by their emails
	Amount      int      `json:"amount"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Submissions []string `json:"submissions"`
	CompletedBy string   `json:"completed_by"` // either null or a user's email
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
	} else if function == "add_submission" {
		return t.add_submission(stub, args)
	} else if function == "single_task_add_submission" {
		return t.single_task_add_submission(stub, args)
	} else if function == "delete_submission" {
		return t.delete_submission(stub, args)
	} else if function == "single_task_delete_submission" {
		return t.single_task_delete_submission(stub, args)
	}
	// else if function == "deposit" {
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
	task.Submissions = []string{args[5]}
	task.CompletedBy = args[6]

	fmt.Println("below is task: ")
	fmt.Println(task)

	////// store task with Uid as key for easy search /////
	taskAsBytes, _ := json.Marshal(task)
	err = stub.PutState(args[0], taskAsBytes)
	if err != nil {
		return nil, err
	}
	//////////////////////////////////////////////////////

	//get the open trade struct
	MarketplaceAsBytes, err := stub.GetState(MarketplaceStr)
	if err != nil {
		return nil, errors.New("Failed to get marketplace")
	}
	var mplace Marketplace
	json.Unmarshal(MarketplaceAsBytes, &mplace) //un stringify it aka JSON.parse()

	//// append task into marketplace ///////////////////////////
	mplace.Tasks = append(mplace.Tasks, task)
	fmt.Println("! appended task to marketplace")
	jsonAsBytes, _ := json.Marshal(mplace)
	err = stub.PutState(MarketplaceStr, jsonAsBytes) //rewrite marketplace
	if err != nil {
		return nil, err
	}
	//////////////////////////////////////////////////////////////

	fmt.Println("- end add task")
	return nil, nil

}

// ============================================================================================================================
// single_task_add_submission - update submission on task, only for single task key, not marketplace
// ============================================================================================================================
func (t *SimpleChaincode) single_task_add_submission(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	//   0       1
	// "uid",  "bob@email.com"
	if len(args) < 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	fmt.Println("- start single_task_add_submission")
	fmt.Println(args[0] + " - " + args[1])

	// get task from blockchain
	tasksAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return nil, errors.New("Failed to get task")
	}

	res := Task{}
	json.Unmarshal(tasksAsBytes, &res)                 //un stringify it aka JSON.parse()
	res.Submissions = append(res.Submissions, args[1]) // append submission

	fmt.Println("! appended submission to task")
	fmt.Println(res.Submissions)
	fmt.Println(res)

	// update task and push back into blockchain
	jsonAsBytes, _ := json.Marshal(res)
	err = stub.PutState(args[0], jsonAsBytes) //rewrite the task with id as key
	if err != nil {
		return nil, err
	}

	fmt.Println("- end single_task_add_submission")
	return nil, nil
}

// ============================================================================================================================
// single_task_delete_submission - delete submission on task, only for single task key, not marketplace
// ============================================================================================================================
func (t *SimpleChaincode) single_task_delete_submission(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	//   0       1
	// "uid",  "bob@email.com"
	if len(args) < 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	fmt.Println("- start single_task_delete_submission")
	fmt.Println(args[0] + " - " + args[1])

	// get task from blockchain
	tasksAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return nil, errors.New("Failed to get task")
	}

	res := Task{}
	json.Unmarshal(tasksAsBytes, &res) //un stringify it aka JSON.parse()

	for i, v := range res.Submissions { // remove submission from task
		if v == args[1] {
			res.Submissions = append(res.Submissions[:i], res.Submissions[i+1:]...)
			break
		}
	}

	fmt.Println("! deleted submission from task")
	fmt.Println(res.Submissions)
	fmt.Println(res)

	// update task and push back into blockchain
	jsonAsBytes, _ := json.Marshal(res)
	err = stub.PutState(args[0], jsonAsBytes) //rewrite the task with id as key
	if err != nil {
		return nil, err
	}

	fmt.Println("- end single_task_delete_submissionn")
	return nil, nil
}

// ============================================================================================================================
// add_submission - upload submission into task
// ============================================================================================================================

func (t *SimpleChaincode) add_submission(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	fmt.Println("Number of args: ")
	fmt.Println(len(args))
	fmt.Println(args[0])
	fmt.Println(args[1])

	//   0            1
	// "uid", "user_who_updated"
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2: task uid and user who updated")
	}

	fmt.Println("- add submission")

	// get all active tasks in marketplace
	MarketplaceAsBytes, err := stub.GetState(MarketplaceStr)
	if err != nil {
		return nil, errors.New("Failed to get marketplace array")
	}
	var mplace Marketplace
	json.Unmarshal(MarketplaceAsBytes, &mplace) //un stringify it aka JSON.parse()

	fmt.Print("Marketplace array: ")
	fmt.Println(mplace)

	//////// update submission in marketplace as well //////
	for i := range mplace.Tasks { //iter through all the tasks
		fmt.Print("looking @ task name: ")
		fmt.Println(mplace.Tasks[i])

		if mplace.Tasks[i].Uid == args[0] { // found the trade to update
			fmt.Println("Found trade to add submission")

			t.single_task_add_submission(stub, []string{args[0], args[1]})

			mplace.Tasks[i].Submissions = append(mplace.Tasks[i].Submissions, args[1]) // add submission to marketplace array
			fmt.Println("! appended submission to task in marketplace")
			fmt.Println(mplace.Tasks[i].Submissions)
			fmt.Println(mplace.Tasks[i])

			jsonAsBytes, _ := json.Marshal(mplace)
			err = stub.PutState(MarketplaceStr, jsonAsBytes) //rewrite the marketplace with new submission
			if err != nil {
				return nil, err
			}
			break
		}
	}
	fmt.Println("- end add submission")
	return nil, nil

}

// ============================================================================================================================
// delete_submission - delete submission from task
// ============================================================================================================================
func (t *SimpleChaincode) delete_submission(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Number of args: ")
	fmt.Println(len(args))
	fmt.Println(args[0])
	fmt.Println(args[1])

	//   0            1
	// "uid", "user_who_updated"
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2: task uid and user who updated")
	}

	fmt.Println("- delete submission")

	// get all active tasks in marketplace
	MarketplaceAsBytes, err := stub.GetState(MarketplaceStr)
	if err != nil {
		return nil, errors.New("Failed to get marketplace array")
	}
	var mplace Marketplace
	json.Unmarshal(MarketplaceAsBytes, &mplace) //un stringify it aka JSON.parse()

	fmt.Print("Marketplace array: ")
	fmt.Println(mplace)

	//////// update submission in marketplace as well //////
	for i := range mplace.Tasks { //iter through all the tasks
		fmt.Print("looking @ task name: ")
		fmt.Println(mplace.Tasks[i])

		if mplace.Tasks[i].Uid == args[0] { // found the trade to update
			fmt.Println("Found trade to delete submission")

			t.single_task_delete_submission(stub, []string{args[0], args[1]})

			for i, v := range mplace.Tasks[i].Submissions { // delete submission from marketplace array
				fmt.Print(i)
				fmt.Println(v)
				if v == args[1] {
					fmt.Println("found v")
					mplace.Tasks[i].Submissions = append(mplace.Tasks[i].Submissions[:i], mplace.Tasks[i].Submissions[i+1:]...)
					break
				}
			}

			fmt.Println("! deleted submission from task in marketplace")
			fmt.Println(mplace.Tasks[i].Submissions)
			fmt.Println(mplace.Tasks[i])

			jsonAsBytes, _ := json.Marshal(mplace)
			err = stub.PutState(MarketplaceStr, jsonAsBytes) //rewrite the marketplace with new submission
			if err != nil {
				return nil, err
			}
			break
		}
	}
	fmt.Println("- end delete submission")
	return nil, nil

}
