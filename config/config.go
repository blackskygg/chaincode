package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type Config map[string][]Object

type Object map[string]interface{}

type Rule string

const dat = `
{
"objects" : [
	 {
		"type" : "table",
		"name" : "student",
		"columns" : [
			  {"name": "id", "type" : "string"},
			  {"name": "name", "type" : "string"},
			  {"name": "status", "type" : "string"}
			  ],
		"rule" : "usr[staff][department] == \"Academic\" && usr[staff][level] >= 3"
	 },
	 {
		"type" : "table",
		"name" : "pay",
		"columns" : [ {"name": "paid", "type" : "string"} ],
		"rule" : "usr[staff][department] == \"Financial\" && usr[staff][level] >= 7"
	 },
       	 {
		"type" : "table",
		"name" : "staff",
		"columns" : [
				{"name":"id", "type" : "string"},
				{"name":"name", "type" : "string"},
				{"name":"department", "type" : "string"},
				{"name":"level", "type" : "int"},
				{"name":"status", "type" : "string"}
			    ],

		"rule" : "usr[staff][department] == \"Personnel\" && usr[staff][level] >= 5 && usr[status] == \"normal\""
	 },
	 {
		"type" : "table",
		"name" : "netusr",
		"columns" : [{"name":"id", "type" : "string"},
				{"name":"balance", "type" : "int"}],
		"rule" : "usr[staff][department] == \"Network\""
	 },
	 {
		"type" : "action",
		"name" : "connect_network",
		"rule" : "target[student][status] == \"registered\" &&  target[netusr][balance] > 0"
	 },
	 {
		"type" : "action",
		"name" : "connect_network"
	 }
	 ],
"init_users": [
	{
		"id" : "PBoss",
		"[staff][department]" : "Personnel",
		"[staff][level]" : "7"
	},
	{
		"id" : "NBoss",
		"[staff][department]" : "Network",
		"[staff][level]" : "7"
	},
	{
		"id" : "ABoss",
		"[staff][department]" : "Academic",
		"[staff][level]" : "7"
	},
	{
		"id" : "FBoss",
		"[staff][department]" : "Financial",
		"[staff][level]" : "8"
	}

  ]
} 
`

func FromFile(file string) (Config, error) {
	var config Config
	//	dat, err := ioutil.ReadFile(file)
	err := json.Unmarshal([]byte(dat), &config)
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
