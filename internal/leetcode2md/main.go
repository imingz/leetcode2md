package leetcode2md

import (
	"fmt"
	"leetcode2md/internal/pkg/save2md"
	"leetcode2md/pkg/version"
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configFile string

func NewLeetcode2mdCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "leetcode2md",
		Short:   "转化 LeetCode 题目和题解为 Markdown 格式",
		Long:    "转化 LeetCode 题目和题解为 Markdown 格式",
		Version: "",
		// 命令出错时，不打印帮助信息。不需要打印帮助信息，设置为 true 可以保持命令出错时一眼就能看到错误信息
		// SilenceUsage: true,
		// 指定调用 cmd.Execute() 时，执行的 Run 函数，函数执行失败会返回错误信息
		RunE: func(cmd *cobra.Command, args []string) error {
			version.PrintAndExitIfRequested()
			return run()
		},
		// 这里设置命令运行时，不需要指定命令行参数
		Args: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}

			return nil
		},
	}

	// 以下设置，使得 initConfig 函数在每个命令运行时都会被调用以读取配置
	cobra.OnInitialize(initConfig)

	// 在这里您将定义标志和配置设置。
	// 添加 --version 标志
	version.AddFlags(cmd.PersistentFlags())

	return cmd
}

// run 函数是实际的业务代码入口函数.
func run() error {
	// 打印所有的配置项及其值
	slog.Debug("All settings", "settings", viper.AllSettings())

	save2md.Save(parseUrl(viper.GetString("url")))

	return nil
}
