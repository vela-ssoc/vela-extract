package extract

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/antchfx/xmlquery"
	"github.com/vela-ssoc/vela-kit/auxlib"
	"io"
	"strings"
)

type XpathQuery struct {
	attribute string
	xpath     []string
}

func (x *XpathQuery) compile(v string) error {
	if len(v) == 0 {
		return fmt.Errorf("invalid xpath string")
	}
	x.xpath = append(x.xpath, v)
	return nil
}

func (x *XpathQuery) xml(v string, box *Box) {
	doc, err := xmlquery.Parse(strings.NewReader(v))
	if err != nil {
		xEnv.Errorf("xml parse fail %v", err)
		return
	}

	for _, item := range x.xpath {
		nodes, er := xmlquery.QueryAll(doc, item)
		if er != nil {
			continue
		}

		for _, node := range nodes {
			var value string
			if x.attribute != "" {
				value = node.SelectAttr(x.attribute)
			} else {
				value = node.InnerText()
			}
			box.append(value)
		}
	}
}

func (x *XpathQuery) html(v string, box *Box) {
	doc, err := htmlquery.Parse(strings.NewReader(v))
	if err != nil {
		xEnv.Errorf("xml parse fail %v", err)
		return
	}

	for _, item := range x.xpath {
		nodes, er := htmlquery.QueryAll(doc, item)
		if er != nil {
			continue
		}

		for _, node := range nodes {
			var value string
			if x.attribute != "" {
				value = htmlquery.SelectAttr(node, x.attribute)
			} else {
				value = htmlquery.InnerText(node)
			}
			box.append(value)
		}
	}
}

func (x *XpathQuery) call(chunk string, box *Box) {
	if len(chunk) == 0 || chunk == "nil" {
		return
	}

	if strings.HasPrefix(chunk, "<?xml") {
		x.xml(chunk, box)
		return
	}

	x.html(chunk, box)
}

func (x *XpathQuery) scanner(reader io.Reader, box *Box) {
	chunk, err := io.ReadAll(reader)
	if err != nil {
		box.err = err
		return
	}

	x.call(auxlib.B2S(chunk), box)
}
