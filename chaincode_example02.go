package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/blackskygg/chaincode/config"
	"github.com/blackskygg/chaincode/parse"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type SimpleChaincode struct {
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	conf, err := config.FromFile("init.conf")
	conf.ApplyConfig(stub)
	return nil, err
}

func (t *SimpleChaincode) checkPermission(table_name, role string) error {
	switch {
	case table_name == "student" && role == "ABoss":
	case table_name == "pay" && role == "FBoss":
	case table_name == "staff" && role == "PBoss":
	case table_name == "netusr" && role == "NBoss":
	default:
		return errors.New("Permission denied!")
	}

	return nil
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	role := string(args[0])
	switch function {
	case "add":
		table_name := string(args[1])
		if err := t.checkPermission(table_name, role); err != nil {
			return nil, err
		}

		var row shim.Row
		tbl, err := stub.GetTable(table_name)
		col_defs := tbl.GetColumnDefinitions()
		if err != nil {
			return nil, err
		}
		for i, v := range col_defs {
			if v.Type == shim.ColumnDefinition_STRING {
				value := string(args[2+i])
				row.Columns =
					append(row.Columns,
						&shim.Column{&shim.Column_String_{value}})
			} else if v.Type == shim.ColumnDefinition_INT32 {
				value, err := strconv.Atoi(args[2+i])
				if err != nil {
					return nil, err
				}
				row.Columns =
					append(row.Columns,
						&shim.Column{&shim.Column_Int32{int32(value)}})
			}
		}

		if _, err := stub.InsertRow(table_name, row); err != nil {
			return nil, err
		}

		return nil, nil

	case "del":
		table_name := string(args[1])
		key := string(args[2])
		if err := t.checkPermission(table_name, role); err != nil {
			return nil, err
		}

		err := stub.DeleteRow(table_name,
			[]shim.Column{shim.Column{&shim.Column_String_{key}}})
		if err != nil {
			return nil, err
		}

		return nil, nil

	case "update":
		table_name := string(args[1])
		key := string(args[2])
		if err := t.checkPermission(table_name, role); err != nil {
			return nil, err
		}
		idx, err := strconv.Atoi(args[2])
		if err != nil {
			return nil, err
		}

		rows, err := stub.GetRow(table_name, []shim.Column{shim.Column{&shim.Column_String_{key}}})
		if err != nil {
			return nil, err
		}

		tbl, err := stub.GetTable(table_name)
		if tbl.ColumnDefinitions[idx].Type == shim.ColumnDefinition_STRING {
			rows.Columns[idx] = &shim.Column{&shim.Column_String_{string(args[3])}}
		} else if tbl.ColumnDefinitions[idx].Type == shim.ColumnDefinition_INT32 {
			val, err := strconv.Atoi(args[3])
			if err != nil {
				return nil, err
			}
			rows.Columns[idx] = &shim.Column{&shim.Column_Int32{int32(val)}}
		}

		if _, err := stub.ReplaceRow(table_name, rows); err != nil {
			return nil, err
		}

		return nil, nil
	default:
		return nil, errors.New("")
	}

	return nil, nil
}

func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	id := string(args[0])
	expression := string(args[1])
	switch function {
	case "query":
		var result string
		var val bool
		var err error
		if val, err = parse.Eval(expression, stub, id); err != nil {
			return nil, err
		}

		if val {
			result = "true"
		} else {
			result = "false"
		}

		return []byte(result), nil
	default:
		return []byte{}, nil
	}

	return []byte{}, nil
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
