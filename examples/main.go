package main

import (
	"github.com/engmtcdrm/go-eggy"
	pp "github.com/engmtcdrm/go-prettyprint"
	"github.com/engmtcdrm/go-zippy/examples/internal"
)

func main() {
	ex := eggy.NewExamplePrompt(internal.AllExamples).
		Title(pp.Yellow("Examples of Zippy"))
	ex.Show()
}
