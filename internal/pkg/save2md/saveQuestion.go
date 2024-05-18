package save2md

import (
	"fmt"
	"leetcode2md/internal/pkg/request"
	"log/slog"
	"os"
)

func writeQuestionToFile(file *os.File, title string, question *request.Question) {
	slog.Info("正在存储题目...\n")
	writeYaml(file, title, question)

	writeQuoteLink(file, title, question)

	writeQuestionContent(file, question)

	slog.Info("题目存储完成\n")
}

// 写 yaml
func writeYaml(file *os.File, title string, question *request.Question) {
	file.WriteString("---\n")
	file.WriteString(fmt.Sprintf("title: %s\n", title))
	file.WriteString("tags:\n  - 算法\n  - LeetCode\n")
	for _, tag := range question.TopicTags {
		file.WriteString(fmt.Sprintf("  - %s\n", tag.TranslatedName))
	}
	file.WriteString(fmt.Sprintf("LeetCode_Level: %s\n", question.Difficulty))

	file.WriteString("---\n\n")
}

// 写标题
func writeQuoteLink(file *os.File, title string, question *request.Question) {
	file.WriteString(fmt.Sprintf("# %s\n\n", title))
	// 写题目地址
	file.WriteString("> [!quote] 题目地址\n")
	file.WriteString(fmt.Sprintf("> [%s. %s - 力扣（LeetCode）](https://leetcode-cn.com/problems/%s/description/)\n\n", question.QuestionFrontendID, question.TranslatedTitle, question.TitleSlug))
}

// 存储题目内容
func writeQuestionContent(file *os.File, question *request.Question) {
	file.WriteString("## 题目\n\n")

	md := newQuestionMd(question.TranslatedContent, question.Hints)

	// 内容
	file.WriteString(md.Content)

	// 示例
	file.WriteString("### 示例\n\n")
	for i, example := range md.Examples {
		file.WriteString(fmt.Sprintf("> [!example]+ 示例 %d\n", i+1))
		file.WriteString(fmt.Sprintf("> **输入**：`%s`\n", example.Input[0]))
		for _, str := range example.Input[1:] {
			file.WriteString(fmt.Sprintf("> %s\n", str))
		}
		file.WriteString(fmt.Sprintf("> **输出**：`%s`\n", example.Output[0]))
		for _, str := range example.Output[1:] {
			file.WriteString(fmt.Sprintf("> %s\n", str))
		}
		if example.Explain != nil {
			file.WriteString(fmt.Sprintf("> **解释**：%s\n", example.Explain[0]))
			for _, str := range example.Explain[1:] {
				file.WriteString(fmt.Sprintf("> %s\n", str))
			}
		}
		file.WriteString("\n")
	}

	// 约束
	file.WriteString("### 约束\n\n")
	file.WriteString("> [!caution]+ 约束\n")
	for _, constraint := range md.Constraints {
		file.WriteString(fmt.Sprintf("> - $%s$\n", constraint))
	}
	file.WriteString("\n")

	// 提示
	if len(md.Hints) != 0 {
		file.WriteString("### 提示\n\n")
		for i, hint := range md.Hints {
			file.WriteString(fmt.Sprintf("> [!hint]- 提示 %d\n", i+1))
			file.WriteString(fmt.Sprintf("> %s\n\n", hint))
		}
	}
}
