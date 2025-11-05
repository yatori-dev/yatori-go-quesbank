package qutils

import (
	"regexp"
	"sort"
)

// 计算两个字符串的Levenshtein距离
func Levenshtein(a, b string) int {
	m, n := len(a), len(b)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}
	for i := 0; i <= m; i++ {
		dp[i][0] = i
	}
	for j := 0; j <= n; j++ {
		dp[0][j] = j
	}
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			cost := 0
			if a[i-1] != b[j-1] {
				cost = 1
			}
			dp[i][j] = min(
				dp[i-1][j]+1,      // 删除
				dp[i][j-1]+1,      // 插入
				dp[i-1][j-1]+cost, // 替换
			)
		}
	}
	return dp[m][n]
}

// 相似度评分
func Similarity(a, b string) float64 {
	maxLen := float64(max(len(a), len(b)))
	if maxLen == 0 {
		return 1.0
	}
	return 1.0 - float64(Levenshtein(a, b))/maxLen
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type Co struct {
	index int
	score float64
}

// 相似度匹配并排序
func SimilarityArrayAndSort(target string, v []string) []int {
	coList := make([]Co, len(v))
	for i := 0; i < len(v); i++ {
		coList[i] = Co{index: i, score: Similarity(v[i], target)}
	}
	sort.Slice(coList, func(i, j int) bool {
		return false
	})
	return nil
}

// 移除题目开头的题目类型文字。
func RemoveLeadingLabel(s string) string {
	re := regexp.MustCompile(`(?m)^\s*\d+\.(?:[[【][^]】]+[]】]|\s*[^\s[【]+)\s*`)
	return re.ReplaceAllString(s, "")
}
