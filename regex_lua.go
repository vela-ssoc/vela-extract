package extract

import "github.com/vela-ssoc/vela-kit/lua"

func (re *RegexQuery) String() string                         { return "" }
func (re *RegexQuery) Type() lua.LValueType                   { return lua.LTObject }
func (re *RegexQuery) AssertFloat64() (float64, bool)         { return 0, false }
func (re *RegexQuery) AssertString() (string, bool)           { return "", false }
func (re *RegexQuery) AssertFunction() (*lua.LFunction, bool) { return re.toLFunc(), true }
func (re *RegexQuery) Peek() lua.LValue                       { return re }

func (re *RegexQuery) toLFunc() *lua.LFunction {
	return lua.NewFunction(re.Call)
}

func (re *RegexQuery) Call(L *lua.LState) int {
	if len(re.regex) == 0 {
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
			re.call(val, box)
		}
	}

	L.Push(box)
	return 1

}

func (re *RegexQuery) partL(L *lua.LState) int {
	n := L.GetTop()
	if n == 0 {
		L.Push(re)
		return 1
	}

	for i := 1; i <= n; i++ {
		p := L.IsInt(i)
		if p <= 0 {
			continue
		}
		re.part = append(re.part, p)
	}

	L.Push(re)
	return 1
}

func (re *RegexQuery) debugL(L *lua.LState) int {
	re.Debug(L, L.Get(1))
	return 0
}

func (re *RegexQuery) fileL(L *lua.LState) int {
	filepath := L.CheckString(1)
	L.Push(load(filepath, re))
	return 1
}

func (re *RegexQuery) requestL(L *lua.LState) int {
	url := L.CheckString(1)
	L.Push(request(url, re))
	return 1
}

func (re *RegexQuery) Index(L *lua.LState, key string) lua.LValue {
	switch key {
	case "part":
		return lua.NewFunction(re.partL)
	case "debug":
		return lua.NewFunction(re.debugL)
	case "file":
		return lua.NewFunction(re.fileL)
	case "request":
		return lua.NewFunction(re.requestL)
	}

	return lua.LNil
}

func regexL(L *lua.LState) int {
	re := &RegexQuery{}
	helperL(L, re.compile)
	L.Push(re)
	return 1
}
