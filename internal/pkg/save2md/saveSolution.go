package save2md

import (
	"fmt"
	"leetcode2md/internal/pkg/request"
	"log/slog"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/viper"
)

func writeSolutionToFile(file *os.File, solution *request.SolutionArticle) {
	slog.Info("正在存储题解...\n")
	file.WriteString("## 题解\n\n")

	file.WriteString("### 我的题解\n\n")

	file.WriteString(fmt.Sprintf("### %s\n\n", solution.Title))

	file.WriteString("> [!quote] 原题解链接\n")
	file.WriteString(fmt.Sprintf("> [%s](%s)\n\n", solution.Title, viper.GetString("url")))

	if viper.GetBool("save.filterLang") {
		// 检查一个字符串是否在一个字符串列表中
		contains := func(list []string, s string) bool {
			for _, v := range list {
				if v == s {
					return true
				}
			}
			return false
		}

		// 目标语言列表
		targetLanguages := viper.GetStringSlice("save.languages")

		slog.Info("已开启筛选语言", "languages", targetLanguages)

		re := regexp.MustCompile("(?s)```(.*?)```.*?")
		codeBlocks := re.FindAllStringSubmatch(solution.Content, -1)

		// 遍历所有的代码块
		for _, block := range codeBlocks {
			// 检查语言标签
			languageTag := strings.Split(block[1], " ")[0]
			if !contains(targetLanguages, languageTag) {
				// 如果不在目标语言列表中，删除这个代码块
				solution.Content = strings.Replace(solution.Content, block[0], "", -1)
			}
		}

		// 递归删除多余的空行
		for strings.Contains(solution.Content, "\n\n\n") {
			solution.Content = strings.Replace(solution.Content, "\n\n\n", "\n\n", -1)
		}
	}

	file.WriteString(solution.Content)
	slog.Info("题解存储完成\n")
}
