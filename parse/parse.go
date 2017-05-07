package parse

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/blackskygg/chaincode/third/govaluate"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func makeParameterMap(stub shim.ChaincodeStubInterface, exp, id string) map[string]interface{} {
	result := make(map[string]interface{})
	re, _ := regexp.Compile(`attr(\[\w+\])+`)
	wre, _ := regexp.Compile(`\w+`)

	l := re.FindAllString(exp, -1)
	for _, w := range l {
		wl := wre.FindAllString(w, -1)
		table_name := wl[1]

		idx, err := strconv.Atoi(wl[2])

		row, err := stub.GetRow(table_name, []shim.Column{shim.Column{&shim.Column_String_{id}}})
		tbl, err := stub.GetTable(table_name)

		fmt.Print("good")
		fmt.Print(err)
		if tbl.ColumnDefinitions[idx].Type == shim.ColumnDefinition_STRING {
			result[w] = interface{}(row.Columns[idx].GetString_())
		} else if tbl.ColumnDefinitions[idx].Type == shim.ColumnDefinition_INT32 {
			result[w] = interface{}(row.Columns[idx].GetInt32())
		}

	}

	return result

}

func Eval(exp string, stub shim.ChaincodeStubInterface, id string) (bool, error) {
	expression, err := govaluate.NewEvaluableExpression(exp)
	var params map[string]interface{}
	if err != nil {
		return false, err
	}

	params = makeParameterMap(stub, exp, id)
	result, err := expression.Evaluate(params)
	return result.(bool), err
}
