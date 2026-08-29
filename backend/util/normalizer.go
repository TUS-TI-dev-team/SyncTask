package util

import (
	"strings"

	"golang.org/x/text/unicode/norm"
)

// NormalizeSearchText はタイトルとコメントを結合し、
// 小文字化、NFKC正規化、ひらがな→カタカナ変換を行った検索用文字列を生成します。
func NormalizeSearchText(title, comment string) string {
	combined := strings.TrimSpace(title)
	commentTrimmed := strings.TrimSpace(comment)

	if combined == "" && commentTrimmed == "" {
		return ""
	}

	if combined != "" && commentTrimmed != "" {
		combined = combined + " " + commentTrimmed
	} else if commentTrimmed != "" {
		combined = commentTrimmed
	}

	// 1. NFKC正規化（全角英数記号の半角化、半角カナの全角カタカナ化、濁点・半濁点の合成）
	normalized := norm.NFKC.String(combined)

	// 2. 英字を小文字化
	lowered := strings.ToLower(normalized)

	// 3. ひらがなを全角カタカナに変換（\u3041〜\u3096 -> \u30a1〜\u30f6）
	var sb strings.Builder
	sb.Grow(len(lowered))
	for _, r := range lowered {
		if r >= 0x3041 && r <= 0x3096 {
			sb.WriteRune(r + 0x60)
		} else {
			sb.WriteRune(r)
		}
	}

	return sb.String()
}
