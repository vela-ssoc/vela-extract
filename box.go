package extract

var emptyBox = &Box{}

type Box struct {
	err   error
	value []string
}

func (box *Box) append(v ...string) {
	box.value = append(box.value, v...)
}
