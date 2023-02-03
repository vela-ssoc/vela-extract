package extract

import (
	"encoding/json"
	"github.com/vela-ssoc/vela-kit/audit"
	"github.com/vela-ssoc/vela-kit/auxlib"
	"github.com/vela-ssoc/vela-kit/lua"
	"io"
	"regexp"
	"strings"
)

type RegexQuery struct {
	part  []int
	regex []*regexp.Regexp
}

func (re *RegexQuery) compile(v string) error {
	r, e := regexp.Compile(v)
	if e != nil {
		return e
	}
	re.regex = append(re.regex, r)
	return nil
}

func (re *RegexQuery) partition(sub [][]string, box *Box) {
	if len(sub) == 0 {
		return
	}

	if len(re.part) == 0 {
		for _, item := range sub {
			box.append(strings.Join(item, ""))
		}
		return
	}

	for _, item := range sub {
		var s2 []string
		for _, p := range re.part {
			if len(item) <= p-1 {
				continue
			}
			s2 = append(s2, item[p-1])
		}

		if len(s2) == 0 {
			continue
		}
		box.append(strings.Join(s2, ""))
	}
}

func (re *RegexQuery) scanner(reader io.Reader, box *Box) {
	chunk, err := io.ReadAll(reader)
	if err != nil {
		box.err = err
		return
	}

	for _, r := range re.regex {
		sub := r.FindAllStringSubmatch(auxlib.B2S(chunk), -1)
		re.partition(sub, box)
	}
}

func (re *RegexQuery) call(lv lua.LValue, box *Box) {
	chunk := lv.String()
	if len(chunk) == 0 || chunk == "nil" {
		return
	}

	for _, r := range re.regex {
		sub := r.FindAllStringSubmatch(chunk, -1)
		re.partition(sub, box)
	}

}

func (re *RegexQuery) Debug(L *lua.LState, lv lua.LValue) {
	chunk := lv.String()
	if len(chunk) == 0 || chunk == "nil" {
		return
	}

	for _, r := range re.regex {
		sub := r.FindAllStringSubmatch(chunk, -1)
		if len(sub) == 0 {
			xEnv.Errorf("not found sub")
			continue
		}

		info, _ := json.MarshalIndent(sub, "", "    ")

		ev := audit.Debug("%s:\nextract regex:\n%s", chunk, auxlib.B2S(info)).From(L.CodeVM())
		ev.Put()
	}
}
