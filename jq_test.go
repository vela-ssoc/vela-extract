package extract

import (
	"github.com/bytedance/sonic"
	"github.com/itchyny/gojq"
	"testing"
)

func TestJq(t *testing.T) {
	query, err := gojq.Parse(`.a.b`)
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	cv, err := gojq.Compile(query)
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	var obj interface{}
	sonic.Unmarshal([]byte(`{"a":{"b":24}}`), &obj)
	iter := cv.Run(obj)
	for {
		item, ok := iter.Next()
		if !ok {
			return
		}
		t.Errorf("%v", item)
	}
}
