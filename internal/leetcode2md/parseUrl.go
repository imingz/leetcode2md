package leetcode2md

import (
	"log/slog"
	"regexp"
)

// 从 url 中解析 titleSlug 和 solutionSlug
func parseUrl(url string) (titleSlug, solutionSlug string) {
	slog.Info("正在解析 url...\n")
	re := regexp.MustCompile(`^https://leetcode.cn/problems/([^/]+)/solutions/[^/]+/([^/]+)/.*$`)
	matches := re.FindStringSubmatch(url)
	if len(matches) < 3 {
		slog.Warn("解析 url 失败，尝试解析 titleSlug...")
		titleSlug = parseUrlOnlyTitle(url)
		return
	}

	slog.Info("解析 url 成功", "titleSlug", matches[1], "solutionSlug", matches[2])
	return matches[1], matches[2]
}

// 从 url 中解析 titleSlug
func parseUrlOnlyTitle(url string) (titleSlug string) {
	slog.Info("正在解析 titleSlug...\n")
	re := regexp.MustCompile(`^https://leetcode.cn/problems/([^/]+)/.*$`)
	matches := re.FindStringSubmatch(url)
	if len(matches) < 2 {
		slog.Error("解析 titleSlug 失败 ", "url", url)
		panic("解析 url 失败")
	}

	slog.Info("解析 titleSlug 成功", "titleSlug", matches[1])
	return matches[1]
}
