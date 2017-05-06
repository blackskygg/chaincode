package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type Config map[string][]Object

type Object map[string]interface{}

type Rule string

func FromFile(file string) (Config, error) {
	var config Config
	dat, err := ioutil.ReadFile(file)
	err = json.Unmarshal(dat, &config)
	return config, err
}

func processTable(stub shim.ChaincodeStubInterface, content map[string]interface{}) {
	var cols_def []*shim.ColumnDefinition
	tbl_name := content["name"].(string)

	for _, k := range content["columns"].([]interface{}) {
		is_key := false
		t := shim.ColumnDefinition_STRING
		name := k.(map[string]interface{})["name"].(string)

		if name == "id" {
			is_key = true
		}
		if k.(map[string]interface{})["type"].(string) == "int" {
			t = shim.ColumnDefinition_INT32
		}
		cols_def = append(cols_def, &shim.ColumnDefinition{
			Type: t, Name: name, Key: is_key})
	}
	stub.CreateTable(tbl_name, cols_def)

	row := shim.Row{}
	row.Columns = append(row.Columns, &shim.Column{&shim.Column_String_{tbl_name}})
	row.Columns = append(row.Columns, &shim.Column{&shim.Column_String_{content["rule"].(string)}})
	r, e := stub.InsertRow("table_rules", row)
	if e != nil {
		os.Exit(10)
	}
	fmt.Print(r)
	fmt.Print(e)
}

func processAction(stub shim.ChaincodeStubInterface, content map[string]interface{}) {

}

func processObjects(stub shim.ChaincodeStubInterface, content map[string]interface{}) {
}

func processUsers(stub shim.ChaincodeStubInterface, content map[string]interface{}) {
}

func createTables(stub shim.ChaincodeStubInterface) {
	var cols_def []*shim.ColumnDefinition
	cols_def = append(cols_def, &shim.ColumnDefinition{
		Type: shim.ColumnDefinition_STRING,
		Name: "name",
		Key:  true})
	cols_def = append(cols_def, &shim.ColumnDefinition{
		Type: shim.ColumnDefinition_STRING,
		Name: "rule"})
	if r := stub.CreateTable("table_rules", cols_def); r != nil {
		fmt.Print(r)
		os.Exit(9)
	}
}

// This should only be invoked once.
func (conf *Config) ApplyConfig(stub shim.ChaincodeStubInterface) {
	objs := (*conf)["objects"]

	createTables(stub)
	for _, v := range objs {
		if v["type"] == "table" {
			processTable(stub, v)
		} else {
			processAction(stub, v)
		}
	}
}

func init() {}
