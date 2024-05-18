package main

import (
	"lc2md/internal/leetcode2md"
	"log/slog"
	"os"

	_ "go.uber.org/automaxprocs"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	command := leetcode2md.NewLeetcode2mdCommand()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
