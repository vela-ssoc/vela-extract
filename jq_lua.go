package extract

import (
	"github.com/vela-ssoc/vela-kit/lua"
)

func (jq *JsonQuery) String() string                         { return "" }
func (jq *JsonQuery) Type() lua.LValueType                   { return lua.LTObject }
func (jq *JsonQuery) AssertFloat64() (float64, bool)         { return 0, false }
func (jq *JsonQuery) AssertString() (string, bool)           { return "", false }
func (jq *JsonQuery) AssertFunction() (*lua.LFunction, bool) { return jq.toLFunc(), true }
func (jq *JsonQuery) Peek() lua.LValue                       { return jq }

func (jq *JsonQuery) Call(L *lua.LState) int {
	if len(jq.code) == 0 {
		L.Push(emptyBox)
		return 1
	}

	n := L.GetTop()
	if n == 0 {
		L.Push(emptyBox)
		return 1
	}

	box := &Box{}
	for i := 1; i <= n; i++ {
		val := L.Get(i)
		switch val.Type() {
		case lua.LTNil, lua.LTFunction:
			continue
		default:
			jq.call(val.String(), box)
		}
	}

	L.Push(box)
	return 1
}

func (jq *JsonQuery) toLFunc() *lua.LFunction {
	return lua.NewFunction(jq.Call)
}

func (jq *JsonQuery) fileL(L *lua.LState) int {
	filepath := L.CheckString(1)

	L.Push(load(filepath, jq))
	return 1
}

func (jq *JsonQuery) requestL(L *lua.LState) int {
	url := L.CheckString(1)
	L.Push(request(url, jq))
	return 1
}

func (jq *JsonQuery) Index(L *lua.LState, key string) lua.LValue {
	switch key {
	case "file":
		return lua.NewFunction(jq.fileL)
	case "request":
		return lua.NewFunction(jq.requestL)
	}

	return lua.LNil
}

func jqL(L *lua.LState) int {
	jq := &JsonQuery{}
	helperL(L, jq.compile)
	L.Push(jq)
	return 1
}
