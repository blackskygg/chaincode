package parse

import (
	"fmt"

	"github.com/Knetic/govaluate"
)

func Eval(exp string, params map[string]interface{}) (bool, error) {
	fmt.Println(exp)
	expression, err := govaluate.NewEvaluableExpression(exp)
	if err != nil {
		return false, err
	}

	result, err := expression.Evaluate(params)
	return result.(bool), err
}

func main() {
	fmt.Println(Eval("0 < 10", nil))
}
