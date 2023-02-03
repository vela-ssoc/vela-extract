package extract

import (
	"regexp"
)

var ipv4 = &RegexQuery{
	regex: []*regexp.Regexp{
		regexp.MustCompile(`(?:(?:1[0-9][0-9]\.)|(?:2[0-4][0-9]\.)|(?:25[0-5]\.)|(?:[1-9][0-9]\.)|(?:[0-9]\.)){3}(?:(?:1[0-9][0-9])|(?:2[0-4][0-9])|(?:25[0-5])|(?:[1-9][0-9])|(?:[0-9]))`),
	},
}

var num = &RegexQuery{
	part: []int{2}, //partition 2
	regex: []*regexp.Regexp{
		regexp.MustCompile(`([0-9]+)`),
	},
}

var uRL = &RegexQuery{
	part: []int{1},
	regex: []*regexp.Regexp{
		regexp.MustCompile(`(https?://([\w-]+\.)+[\w-]+(/[\w-./?%&=]*)?)`),
	},
}
var phone = &RegexQuery{
	part: []int{1},
	regex: []*regexp.Regexp{
		regexp.MustCompile(`(1[3-9]\d{9})`),
	},
}
