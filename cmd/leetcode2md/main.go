package main

import (
	"fmt"
	"leetcode2md/internal/leetcode2md"
	"log/slog"
	"os"
	"runtime"

	_ "go.uber.org/automaxprocs"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	if runtime.GOOS == "windows" {
		defer func() {
			fmt.Println("按回车键退出...")
			fmt.Scanln()
		}()
	}

	command := leetcode2md.NewLeetcode2mdCommand()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
