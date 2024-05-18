package save2md

import (
	"fmt"
	"leetcode2md/internal/pkg/request"
	"log/slog"
	"os"

	"github.com/spf13/viper"
)

func writeSolutionToFile(file *os.File, solution *request.SolutionArticle) {
	slog.Info("正在存储题解...\n")
	file.WriteString("## 题解\n\n")

	file.WriteString("### 我的题解\n\n")

	file.WriteString(fmt.Sprintf("### %s\n\n", solution.Title))

	file.WriteString("> [!quote] 原题解链接\n")
	file.WriteString(fmt.Sprintf("> [%s](%s)\n\n", solution.Title, viper.GetString("url")))

	file.WriteString(solution.Content)
	slog.Info("题解存储完成\n")
}
