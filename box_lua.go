package extract

import (
	"github.com/vela-ssoc/vela-kit/pipe"
	"github.com/vela-ssoc/vela-kit/lua"
)

/*
local jq = vela.extract.jq(".a.b")

jq([[{"a":{"b":24}]]).pipe(function(v)
	print(v)
end)


*/

func (box *Box) String() string                         { return "" }
func (box *Box) Type() lua.LValueType                   { return lua.LTObject }
func (box *Box) AssertFloat64() (float64, bool)         { return 0, false }
func (box *Box) AssertString() (string, bool)           { return "", false }
func (box *Box) AssertFunction() (*lua.LFunction, bool) { return nil, false }
func (box *Box) Peek() lua.LValue                       { return box }

func (box *Box) toLFunc() *lua.LFunction {
	return lua.NewFunction(box.pipeL)
}

func (box *Box) pipeL(L *lua.LState) int {
	if len(box.value) == 0 {
		return 0
	}

	pip := pipe.NewByLua(L)

	for id, val := range box.value {
		pip.Call2(lua.S2L(val), id, L)
	}
	return 0
}

func (box *Box) Index(L *lua.LState, key string) lua.LValue {
	switch key {
	case "pipe":
		return lua.NewFunction(box.pipeL)
	case "err":
		if box.err == nil {
			return lua.LNil
		}
		return lua.S2L(box.err.Error())
	}
	return lua.LNil
}
