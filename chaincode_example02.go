package main

import (
	"errors"
	"fmt"

	"github.com/blackskygg/chaincode/attributes"
	"github.com/blackskygg/chaincode/config"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type SimpleChaincode struct {
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	conf, err := config.FromFile("init.conf")
	conf.ApplyConfig(stub)
	return nil, err
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var err error
	//	table_name := args[0]

	switch function {
	case "add":
		uuid := args[0]
		if _, err := attributes.FromJson([]byte(args[1])); err != nil {
			return nil, errors.New("bad json")
		}
		info := args[1]
		if err = stub.PutState(uuid, []byte(info)); err != nil {
			return nil, errors.New("failed to enroll")
		}
	case "del":
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

func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var info []byte
	var uuid string
	var err error

	switch function {
	case "query":
		var cols []shim.Column
		var rows shim.Row
		cols = append(cols, shim.Column{&shim.Column_String_{"student"}})

		rows, err := stub.GetRow("table_rules", []shim.Column{shim.Column{&shim.Column_String_{"student"}}})

		return []byte(rows.String()), err
	case "cert":
		return stub.GetCallerCertificate()
	default:
		return []byte{}, nil
	}

	uuid = args[0]

	if info, err = stub.GetState(uuid); info == nil {
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
