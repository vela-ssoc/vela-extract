package extract

import (
	"io"
	"net/http"
	"os"
)

type extractor interface {
	scanner(io.Reader, *Box)
}

func load(filename string, extract extractor) *Box {
	box := &Box{}
	file, err := os.Open(filename)
	if err != nil {
		box.err = err
		return box
	}
	defer file.Close()

	extract.scanner(file, box)
	return box
}

func request(url string, extract extractor) *Box {
	box := &Box{}
	r, err := http.Get(url)
	if err != nil {
		box.err = err
		return box
	}
	defer r.Body.Close()
	extract.scanner(r.Body, box)
	return box
}
