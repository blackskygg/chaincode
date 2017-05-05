/*
Copyright IBM Corp. 2016 All Rights Reserved.

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

//WARNING - this chaincode's ID is hard-coded in chaincode_example04 to illustrate one way of
//calling chaincode from a chaincode. If this example is modified, chaincode_example04.go has
//to be modified as well with the new ID of chaincode_example02.
//chaincode_example05 show's how chaincode ID can be passed in as a parameter instead of
//hard-coding.

import (
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	return nil, nil
}

// Transaction makes payment of X units from A to B
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var err error

	switch function {
	case "enroll":
		uuid := args[0]
		info := args[1]
		if err = stub.PutState(uuid, []byte(info)); err != nil {
			return nil, errors.New("failed to enroll")
		}
	case "unenroll":
		uuid := args[0]
		if err = stub.DelState(uuid); err != nil {
			return nil, errors.New("failed to unenroll")
		}
	case "update":
		uuid := args[0]
		info := args[1]
		if _, err = stub.GetState(uuid); err != nil {
			return nil, errors.New("Not a valid student")
		}

		if err = stub.PutState(uuid, []byte(info)); err != nil {
			return nil, err
		}

	default:
		return nil, err
	}

	return nil, err
}

// Query callback representing the query of a chaincode
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var info []byte
	var uuid string
	var err error

	uuid = args[0]

	if info, err = stub.GetState(uuid); err != nil {
		jsonResp := "{\"Error\":\"Not a valid student\"}"
		return nil, errors.New(jsonResp)
	}

	return info, err
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
