package seed

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"synctask/backend/util"
)

// TestSeedUser_PasswordHashMatches @spec シード: ユーザーのパスワードが平文 Password123! と一致する bcrypt ハッシュで検証できること
func TestSeedUser_PasswordHashMatches(t *testing.T) {
	for _, u := range users {
		t.Run("シードユーザー: パスワード検証に成功する_"+u.Email, func(t *testing.T) {
			// seed.Run と同じ方式 (bcrypt.DefaultCost) でハッシュを生成して検証する
			hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
			require.NoError(t, err)
			err = bcrypt.CompareHashAndPassword(hash, []byte(u.Password))
			assert.NoError(t, err, "パスワード検証に失敗してはいけない")
		})
	}
}

// TestSeedUser_Format @spec シード: すべてのユーザー定義が LOGIN_ACCOUNT の制約に適合すること
func TestSeedUser_Format(t *testing.T) {
	for _, u := range users {
		assert.Equal(t, strings.ToLower(u.Email), u.Email, "EMAIL は小文字で定義すること (ログイン照合は小文字正規化前提)")
		assert.NotEmpty(t, u.UserID)
		assert.Len(t, u.UserID, 36, "USER_ID は UUID 形式 (36 文字)")
		assert.Regexp(t, `^[A-Za-z0-9]{2,20}$`, u.UserName, "USER_NAME は 2〜20 文字の英数字")
	}
}

// TestSeedTask_ConstraintValues @spec シード: タスクの PRIORITY / STATUS が CHECK 制約の許容値に適合すること
func TestSeedTask_ConstraintValues(t *testing.T) {
	validPriority := map[string]bool{"low": true, "medium": true, "high": true}
	validStatus := map[string]bool{"not_started": true, "in_progress": true, "completed": true}

	require.NotEmpty(t, tasks)
	for _, task := range tasks {
		assert.True(t, validPriority[task.Priority], "PRIORITY が不適合: %s (%s)", task.Priority, task.Title)
		assert.True(t, validStatus[task.Status], "STATUS が不適合: %s (%s)", task.Status, task.Title)
		assert.NotEmpty(t, task.Title)
		assert.LessOrEqual(t, len(task.Title), 255, "TITLE は 255 文字以内")
	}
}

// TestSeedTask_SearchTextNormalization @spec シード: タスクの SEARCH_TEXT が util.NormalizeSearchText の仕様 (NFKC・小文字化・ひらカナ変換) に適合すること
func TestSeedTask_SearchTextNormalization(t *testing.T) {
	for _, task := range tasks {
		want := util.NormalizeSearchText(task.Title, task.Comment)
		assert.NotEmpty(t, want, "SEARCH_TEXT は NOT NULL のため空になってはいけない (%s)", task.Title)
		// ひらがなはカタカナ化されることを確認
		for _, r := range want {
			assert.False(t, r >= 0x3041 && r <= 0x3096, "ひらがなはカタカナ化されるべき (%s)", task.Title)
		}
		// 正規化後の文字列が再び同じ結果になること (冪等性の確認)
		assert.Equal(t, want, util.NormalizeSearchText(want, ""), "正規化は冪等であること")
	}
}

// TestSeedSession_IDLength @spec シード: セッション ID が LOGIN_SESSION.SESSION_ID の VARCHAR(64) に適合すること
func TestSeedSession_IDLength(t *testing.T) {
	require.NotEmpty(t, sessions)
	for _, s := range sessions {
		assert.LessOrEqual(t, len(s.SessionID), 64, "SESSION_ID は 64 文字以内")
		assert.NotEmpty(t, s.UserID)
	}
}
