package extract

import "github.com/vela-ssoc/vela-kit/lua"

func (x *XpathQuery) String() string                         { return "" }
func (x *XpathQuery) Type() lua.LValueType                   { return lua.LTObject }
func (x *XpathQuery) AssertFloat64() (float64, bool)         { return 0, false }
func (x *XpathQuery) AssertString() (string, bool)           { return "", false }
func (x *XpathQuery) AssertFunction() (*lua.LFunction, bool) { return x.toLFunc(), true }
func (x *XpathQuery) Peek() lua.LValue                       { return x }

func (x *XpathQuery) toLFunc() *lua.LFunction {
	return lua.NewFunction(x.Call)
}

func (x *XpathQuery) Call(L *lua.LState) int {
	if len(x.xpath) == 0 {
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
			x.call(val.String(), box)
		}
	}

	L.Push(box)
	return 1

}

func (x *XpathQuery) attrL(L *lua.LState) int {
	v := L.CheckString(1)
	if len(v) == 0 {
		L.Push(x)
		return 1
	}

	x.attribute = v
	L.Push(x)
	return 1
}

func (x *XpathQuery) requestL(L *lua.LState) int {
	url := L.CheckString(1)
	L.Push(request(url, x))
	return 1
}

func (x *XpathQuery) Index(L *lua.LState, key string) lua.LValue {
	switch key {
	case "attr":
		return lua.NewFunction(x.attrL)
	case "request":
		return lua.NewFunction(x.requestL)
	}

	return lua.LNil
}

func xpathL(L *lua.LState) int {
	x := &XpathQuery{}
	helperL(L, x.compile)
	L.Push(x)
	return 1
}
