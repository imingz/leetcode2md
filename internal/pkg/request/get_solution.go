package request

import (
	"encoding/json"
	"fmt"
	"strings"
)

var (
	ErrorSolutionNotFound = fmt.Errorf("未找到解题思路，请检查 url")
)

func UnmarshalSolutionResponse(data []byte) (SolutionResponse, error) {
	var r SolutionResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *SolutionResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type SolutionResponse struct {
	Data SolutionArticleData `json:"data"`
}

type SolutionArticleData struct {
	SolutionArticle *SolutionArticle `json:"solutionArticle"`
}

type SolutionArticle struct {
	Content string `json:"content"`
	Title   string `json:"title"`
}

func generateSolutionPayload(slug string) string {
	query := `
query solutionArticle($slug: String) {
	solutionArticle(slug: $slug) {
		title
		content
	}
}`
	variables := fmt.Sprintf("{\"slug\":\"%s\"}", slug)

	payloadString := fmt.Sprintf("{\"query\":\"%s\",\"variables\":%s}", query, variables)

	// 替换换行符
	replaceNewLine := func(s string) string {
		return strings.ReplaceAll(s, "\n", "\\r\\n")
	}

	// 移除多余的空格
	removeExtraSpaces := func(s string) string {
		return strings.Join(strings.Fields(s), " ")
	}

	return removeExtraSpaces(replaceNewLine(payloadString))
}

func GetSolution(titleSlug string) (*SolutionArticle, error) {
	payload := strings.NewReader(generateSolutionPayload(titleSlug))

	body, err := request(payload)
	if err != nil {
		return nil, err
	}

	// 解析为 json
	r, err := UnmarshalSolutionResponse(body)
	if err != nil {
		return nil, err
	}

	if r.Data.SolutionArticle == nil {
		return nil, ErrorSolutionNotFound
	}

	return r.Data.SolutionArticle, nil
}
