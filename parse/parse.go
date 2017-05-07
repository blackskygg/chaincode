package parse

import (
	"regexp"
	"strconv"

	govaluate "github.com/blackskygg/chaincode/third/govaluate_modified"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func makeParameterMap(stub shim.ChaincodeStubInterface, exp, id string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	re, err := regexp.Compile(`attr(_\w+)+`)
	wre, err := regexp.Compile(`[a-z0-9]+`)
	if err != nil {
		return result, err
	}

	l := re.FindAllString(exp, -1)
	for _, w := range l {
		wl := wre.FindAllString(w, -1)
		table_name := wl[1]

		idx, err := strconv.Atoi(wl[2])
		if err != nil {
			return result, err
		}

		row, err := stub.GetRow(table_name, []shim.Column{shim.Column{&shim.Column_String_{id}}})
		if err != nil {
			return result, err
		}

		tbl, err := stub.GetTable(table_name)
		if err != nil {
			return result, err
		}

		if tbl.ColumnDefinitions[idx].Type == shim.ColumnDefinition_STRING {
			result[w] = interface{}(row.Columns[idx].GetString_())
		} else if tbl.ColumnDefinitions[idx].Type == shim.ColumnDefinition_INT32 {
			result[w] = interface{}(row.Columns[idx].GetInt32())
		}

	}

	return result, nil

}

func Eval(exp string, stub shim.ChaincodeStubInterface, id string) (bool, error) {
	expression, err := govaluate.NewEvaluableExpression(exp)
	var params map[string]interface{}
	if err != nil {
		return false, err
	}

	params, err = makeParameterMap(stub, exp, id)
	if err != nil {
		return false, err
	}

	result, err := expression.Evaluate(params)
	return result.(bool), err
}
