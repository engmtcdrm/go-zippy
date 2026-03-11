package main

import (
	"log/slog"

	"github.com/engmtcdrm/go-eggy"
	"github.com/engmtcdrm/go-lager"
	pp "github.com/engmtcdrm/go-prettyprint"
	"github.com/engmtcdrm/go-zippy/examples/internal"
)

func main() {
	opts := &lager.HandlerOptions{Level: lager.LevelDebug}
	sout := lager.NewStreamHandler(lager.StreamStdout, opts)
	logger := slog.New(sout)
	slog.SetDefault(logger)

	ex := eggy.NewExamplePrompt(internal.AllExamples).
		Title(pp.Yellow("Examples of Zippy"))
	ex.Show()
}
