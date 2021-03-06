package attributes

import (
	json "encoding/json"
	"errors"
	"fmt"
)

type Attribute map[string]string

func FromJson(j []byte) (Attribute, error) {
	var a Attribute
	err := json.Unmarshal(j, &a)
	return a, err
}

func (a Attribute) ToJson() ([]byte, error) {
	return json.Marshal(a)
}

func (a Attribute) UpdateKey(key string, val string) Attribute {
	a[key] = val
	return a
}

func (a Attribute) DelKey(key string) (string, error) {
	val, ok := a[key]
	if ok != true {
		return "", errors.New("no such key")
	} else {
		delete(a, key)
		return val, nil
	}
}

func init() {}
func test() {
	b := []byte("{\"a\":\"aaa\", \"b\":\"bbb\"}")
	a, _ := FromJson(b)
	fmt.Print(a)

	a = Attribute{}
	a.UpdateKey("k", "kkk")
	a.UpdateKey("l", "lll")
	a.UpdateKey("m", "mmm")
	a.DelKey("o")
	a.DelKey("m")

	j, _ := a.ToJson()
	fmt.Print(string(j))
}
