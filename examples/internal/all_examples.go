package internal

import (
	"github.com/engmtcdrm/go-eggy"
)

const (
	tempDirName = "go-zippy-example-*"
	zipFileName = "go-zippy-example.zip"
)

var AllExamples = []eggy.Example{
	{Name: "Contents Example", Fn: ExampleContents},
	{Name: "Extract Example", Fn: ExtractExample},
}
