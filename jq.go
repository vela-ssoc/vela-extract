package extract

import (
	"encoding/json"
	"github.com/bytedance/sonic"
	"github.com/itchyny/gojq"
	auxlib2 "github.com/vela-ssoc/vela-kit/auxlib"
	"io"
)

var emptyJsonQuery = &JsonQuery{}

type JsonQuery struct {
	code []*gojq.Code
}

func (jq *JsonQuery) compile(v string) error {
	q, err := gojq.Parse(v)
	if err != nil {
		return err
	}

	cde, err := gojq.Compile(q)
	if err != nil {
		return err
	}
	jq.code = append(jq.code, cde)
	return nil
}

func (jq *JsonQuery) call(chunk string, box *Box) {
	if len(chunk) == 0 || chunk == "nil" {
		return
	}
	var obj interface{}

	err := sonic.Unmarshal(auxlib2.S2B(chunk), &obj)
	if err != nil {
		return
	}

	for _, code := range jq.code {
		iter := code.Run(obj)
		for {
			item, ok := iter.Next()
			if !ok {
				return
			}

			if val := auxlib2.ToString(item); len(val) != 0 {
				box.append(val)
			}
		}

	}
}

func (jq *JsonQuery) scanner(reader io.Reader, box *Box) {
	var obj interface{}
	dec := json.NewDecoder(reader)
	err := dec.Decode(&obj)
	if err != nil {
		return
	}
	for _, code := range jq.code {
		iter := code.Run(obj)
		for {
			item, ok := iter.Next()
			if !ok {
				return
			}

			if val := auxlib2.ToString(item); len(val) != 0 {
				box.append(val)
			}
		}

	}
}
