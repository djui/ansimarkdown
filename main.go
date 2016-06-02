package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	md "gopkg.in/russross/blackfriday.v2"
)

func main() {
	log.SetFlags(0)
	flag.Parse()

	var in []byte
	var err error

	if flag.NArg() == 0 {
		in, err = ioutil.ReadAll(os.Stdin)
	} else {
		in, err = ioutil.ReadFile(flag.Arg(0))
	}

	if err != nil {
		log.Fatalln("Error:", err)
	}

	r := &ansiRenderer{}
	o := md.Options{
		Extensions: md.NoIntraEmphasis |
			md.Tables |
			md.FencedCode |
			md.Autolink |
			md.Strikethrough |
			md.SpaceHeaders |
			md.BackslashLineBreak |
			md.DefinitionLists,
	}

	out := r.Render(md.Parse(in, o))
	fmt.Println(string(out))
}
