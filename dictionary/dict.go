package dictionary

import (
	"github.com/fiorix/go-diameter/diam/dict"
	"log"
	"os"
)

func Load() *dict.Parser {
	var err error
	gopath := os.Getenv("GOPATH")

	parser, err := dict.NewParser(gopath + "/src/server/dictionary/base.xml")

	var run = func(filename string) {
		if err != nil {
			log.Println(err.Error())
			return
		}
		parser.LoadFile(filename)
	}

	run(gopath + "/src/server/dictionary/creditcontrol.xml")
	run(gopath + "/src/server/dictionary/tgpp_ro_rf.xml")

	return parser
}
