package request

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
)

// questionResponseType
type questionResponseType struct {
	Data questionData `json:"data"`
}

type questionData struct {
	Question *Question `json:"question"`
}

type Question struct {
	TitleSlug          string     `json:"titleSlug"`
	Difficulty         string     `json:"difficulty"`
	QuestionFrontendID string     `json:"questionFrontendId"`
	TopicTags          []TopicTag `json:"topicTags"`
	TranslatedContent  string     `json:"translatedContent"`
	TranslatedTitle    string     `json:"translatedTitle"`
	Hints              []string   `json:"hints"`
}

type TopicTag struct {
	TranslatedName string `json:"translatedName"`
}

func unmarshalResponseType(data []byte) (questionResponseType, error) {
	var r questionResponseType
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *questionResponseType) marshal() ([]byte, error) {
	return json.Marshal(r)
}

func generateQuestionPayload(titleSlug string) string {
	query := `
		query question($titleSlug: String!) {
			question(titleSlug: $titleSlug) {
				titleSlug
				difficulty
				questionFrontendId
				translatedTitle
				topicTags {
					translatedName
				}
				hints
				translatedContent
			}
		}`
	variables := fmt.Sprintf("{\"titleSlug\":\"%s\"}", titleSlug)

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

func GetQuestion(titleSlug string) (*Question, error) {
	payload := strings.NewReader(generateQuestionPayload(titleSlug))

	body, err := request(payload)
	if err != nil {
		return nil, err
	}

	// 解析为 json
	r, err := unmarshalResponseType(body)
	if err != nil {
		return nil, err
	}

	if r.Data.Question == nil {
		slog.Error("未找到题目，请检查 url", "titleSlug", titleSlug)
		return nil, fmt.Errorf("未找到题目，请检查 url")
	}

	return r.Data.Question, nil
}
