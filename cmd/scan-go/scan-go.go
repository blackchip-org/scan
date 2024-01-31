package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/blackchip-org/scan"
	"github.com/blackchip-org/scan/scango"
)

var jsonOutput bool

func main() {
	log.SetFlags(0)
	flag.BoolVar(&jsonOutput, "json", false, "output tokens as json")
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
	}

	f, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	s := scan.NewScanner(flag.Arg(0), f)
	ctx := scango.NewContext()
	r := scan.NewRunner(s, ctx.RuleSet)
	toks := r.All()
	if jsonOutput {
		for i, t := range toks {
			if t.Val == t.Lit {
				toks[i].Lit = ""
			}
		}
		data, err := json.MarshalIndent(toks, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(data))
	} else {
		fmt.Println(scan.FormatTokenTable(toks))
	}
}
