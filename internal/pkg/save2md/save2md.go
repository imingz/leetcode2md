package save2md

import (
	"fmt"
	"leetcode2md/internal/pkg/request"
	"log/slog"
	"os"

	"github.com/spf13/viper"
)

func Save(titleSlug, solutionSlug string) {
	file := getFile(titleSlug)
	defer func() {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}()

	title := saveQuestion(file, titleSlug)
	defer func() {
		err := os.Rename(file.Name(), fmt.Sprintf("%s/%s.md", viper.GetString("dir.md"), title))
		if err != nil {
			panic(err)
		}
	}()

	if solutionSlug != "" {
		saveSolution(file, solutionSlug)
	}
}

func getFile(titleSlug string) *os.File {
	slog.Info("正在创建文件...\n")
	file, err := os.Create(fmt.Sprintf("%s/%s.md", viper.GetString("dir.md"), titleSlug))
	if err != nil {
		slog.Error("创建文件失败")
		panic(err)
	}

	return file
}

func saveQuestion(file *os.File, titleSlug string) string {
	var question *request.Question
	var questionErr error

	slog.Info("正在请求题目...\n")
	for i := 0; i < 3; i++ {
		question, questionErr = request.GetQuestion(titleSlug)
		if questionErr == nil {
			break
		}
		slog.Info(fmt.Sprintf("请求题目失败，正在进行第 %d 次重试...\n", i+1))
	}

	if questionErr != nil {
		panic(questionErr)
	}
	slog.Info("题目请求完成\n")

	title := fmt.Sprintf("%s. %s", question.QuestionFrontendID, question.TranslatedTitle)

	writeQuestionToFile(file, title, question)

	return title
}

func saveSolution(file *os.File, solutionSlug string) {
	var solution *request.SolutionArticle
	var solutionErr error

	slog.Info("正在请求解题思路...\n")
	for i := 0; i < 3; i++ {
		solution, solutionErr = request.GetSolution(solutionSlug)
		if solutionErr == nil {
			break
		}
		slog.Info(fmt.Sprintf("请求解题思路失败，正在进行第 %d 次重试...\n", i+1))
	}

	if solutionErr != nil {
		panic(solutionErr)
	}
	slog.Info("解题思路请求完成\n")

	writeSolutionToFile(file, solution)
}
