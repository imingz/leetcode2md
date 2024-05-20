package save2md

import (
	"fmt"
	"log/slog"
	"regexp"
	"strings"

	"github.com/spf13/viper"
)

type questionMd struct {
	Content     string
	Examples    []exampleMd
	Constraints []string
	Hints       []string
}

type exampleMd struct {
	Input   []string
	Output  []string
	Explain []string
}

func newQuestionMd(translatedContent string, hints []string) *questionMd {
	qm := &questionMd{
		Content:     "",
		Examples:    []exampleMd{},
		Constraints: []string{},
		Hints:       hints,
	}

	// [0]: Content, [1]: Examples, [2]: Constraints
	strings := strings.Split(translatedContent, "<p>&nbsp;</p>")
	qm.parseContent(strings[0])
	qm.parseExamples(strings[1])
	qm.parseConstraints(strings[2])

	return qm
}

// parseContent 解析题目内容
func (qm *questionMd) parseContent(content string) {
	qm.Content = replaceMath(replaceHtml(content))
}

// parseExamples 解析示例
func (qm *questionMd) parseExamples(str string) {
	str = replaceHtml(str)
	pre := regexp.MustCompile(`<pre>`)
	preMatches := pre.FindAllStringSubmatchIndex(str, -1)
	preEnd := regexp.MustCompile(`</pre>`)
	preEndMatches := preEnd.FindAllStringSubmatchIndex(str, -1)

	for i := range preMatches {
		subStr := str[preMatches[i][1]:preEndMatches[i][0]]
		var example exampleMd

		inputReg := regexp.MustCompile(`\*\*输入[:：\s\n]+\*\*`)
		outPutReg := regexp.MustCompile(`\*\*输出[:：\s\n]+\*\*`)
		explainReg := regexp.MustCompile(`\*\*解释[:：\s\n]+\*\*`)

		inputIndex := inputReg.FindAllStringSubmatchIndex(subStr, -1)[0]
		outPutIndex := outPutReg.FindAllStringSubmatchIndex(subStr, -1)[0]
		explainIndex := explainReg.FindAllStringSubmatchIndex(subStr, -1)

		example.Input = strings.Split(strings.Trim(subStr[inputIndex[1]:outPutIndex[0]], "\n"), "\n")
		if explainIndex != nil {
			example.Output = strings.Split(strings.Trim(subStr[outPutIndex[1]:explainIndex[0][0]], "\n"), "\n")
			explains := strings.Split(strings.Trim(subStr[explainIndex[0][1]:], "\n"), "\n")

			re_img := regexp.MustCompile(`<img .*?>`)

			for i, e := range explains {
				imageMatches := re_img.FindStringSubmatch(e)
				if imageMatches != nil {
					if len(imageMatches) > 1 {
						slog.Warn("未处理的情况，联系开发者。", "matches", imageMatches)
					}
					re_src := regexp.MustCompile(`src="(.*?)"`)
					src := re_src.FindStringSubmatch(imageMatches[0])[1]
					if viper.GetBool("save.image") {
						imgPath, err := saveImage(src)
						if err != nil {
							slog.Warn("图片保存失败，使用原链接", "err", err)
							explains[i] = fmt.Sprintf("![](%v)", src)
						} else {
							explains[i] = fmt.Sprintf("![](%v)", imgPath)
						}
					} else {
						explains[i] = fmt.Sprintf("![](%v)", src)
					}
				}
			}

			example.Explain = explains
		} else {
			example.Output = strings.Split(strings.Trim("subStr[outPutIndex[1]:]", "\n"), "\n")
		}
		qm.Examples = append(qm.Examples, example)
	}
}

// parseConstraints 解析约束
func (qm *questionMd) parseConstraints(str string) {
	re := regexp.MustCompile(`<li><code>(.*?)</code></li>`)
	matches := re.FindAllStringSubmatch(str, -1)

	for _, match := range matches {
		qm.Constraints = append(qm.Constraints, replaceMath(match[1]))
	}
}

// 替换数学语法
func replaceMath(str string) string {
	replace := []string{
		// `<sup>{{str}}</sup>` 为 `^{{str}}`
		"<sup>", "^", "</sup>", "",
		// `&lt;` --> `\lt`
		"&lt;", "\\lt",
		// 空格
		"&nbsp;", " ",
	}

	replacer := strings.NewReplacer(replace...)

	return replacer.Replace(str)
}

// 替换 html 语法
func replaceHtml(str string) (ans string) {
	defer func() {
		replaceAfter := []string{
			// 空格
			"&nbsp;", " ",
		}
		replacerAfter := strings.NewReplacer(replaceAfter...)
		ans = replacerAfter.Replace(ans)
	}()

	replace := []string{
		// 删除不可见字符
		"\u200B", "", "\u200C", "", "\u200D", "", "\uFEFF", "",
		// 删除 p
		"<p>", "", "</p>", "",
		// 删除 ul、ol，`<li>` --> `- `
		"<ul>\n", "", "</ul>\n", "", "	<li>", "- ", "</li>", "", "<ol>\n", "", "</ol>\n", "",
		// <code> </code>
		"$", "\\$", "<code>", "$", "</code>", "$",
		// `<strong>{{str}}</strong>` --> `**{{str}}**`
		"<strong>", "**", "</strong>", "**",
		// ` <em>{{str}}&nbsp;</em>` --> ` _{{str}}_ `
		" <em>", " _", "&nbsp;</em>", "_ ",
		"&quot;", "\"",
	}

	replacer := strings.NewReplacer(replace...)

	return replacer.Replace(str)
}
