package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeSearchText(t *testing.T) {
	t.Run("正常系: 半角英大文字・小文字が小文字に統一されること", func(t *testing.T) {
		got := NormalizeSearchText("Hello WORLD", "Go Program")
		assert.Equal(t, "hello world go program", got)
	})

	t.Run("正常系: 全角英数字が半角小文字にNFKC正規化されること", func(t *testing.T) {
		got := NormalizeSearchText("ＴＥＳＴ１２３", "Ａｂｃ４５６")
		assert.Equal(t, "test123 abc456", got)
	})

	t.Run("正常系: 半角カタカナが全角カタカナに変換されること", func(t *testing.T) {
		got := NormalizeSearchText("ﾃｽﾄ", "ﾀｽｸ")
		assert.Equal(t, "テスト タスク", got)
	})

	t.Run("正常系: ひらがなが全角カタカナに変換されること", func(t *testing.T) {
		got := NormalizeSearchText("あいうえお", "やゆよ わをん")
		assert.Equal(t, "アイウエオ ヤユヨ ワヲン", got)
	})

	t.Run("正常系: 濁音・半濁音（が、ぱ、ｶﾞ、ﾊﾟ）が正しく合成・変換されること", func(t *testing.T) {
		got := NormalizeSearchText("が ぱ", "ｶﾞ ﾊﾟ")
		assert.Equal(t, "ガ パ ガ パ", got)
	})

	t.Run("正常系: タイトルとコメントが結合されて正規化されること", func(t *testing.T) {
		got := NormalizeSearchText("タスク作成", "詳細内容")
		assert.Equal(t, "タスク作成 詳細内容", got)
	})

	t.Run("境界値: 空文字の場合は空文字が返却されること", func(t *testing.T) {
		got := NormalizeSearchText("", "")
		assert.Equal(t, "", got)
	})
}
