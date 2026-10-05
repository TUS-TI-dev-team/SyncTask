// Package seed は開発時の動作テスト用シードデータを提供します。
//
// @spec シード: LOGIN_ACCOUNT と TASK 関連の開発用データを冪等に投入できること
package seed

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"synctask/backend/util"
)

// SeedUser はシード投入するユーザー (LOGIN_ACCOUNT) の定義です。
type SeedUser struct {
	UserID   string
	UserName string
	Email    string
	Password string // 平文パスワード (投入時に bcrypt ハッシュ化)
}

// SeedTask はシード投入するタスク (TASK) の定義です。
type SeedTask struct {
	TaskID    string
	UserID    string // SeedUser.UserID
	Title     string
	Comment   string
	Priority  string // low / medium / high
	Status    string // not_started / in_progress / completed
	IsPinned  bool
	DueOffset *time.Duration // 現在時刻からの相対期限 (nil は期限なし)
}

// SeedSession はシード投入するログインセッション (LOGIN_SESSION) の定義です。
type SeedSession struct {
	SessionID string
	UserID    string // SeedUser.UserID
	UserAgent string
}

// users はシード対象ユーザーの一覧です。
// パスワードはすべて開発用の固定値 "Password123!" です。
var users = []SeedUser{
	{
		UserID:   "00000000-0000-4000-8000-000000000001",
		UserName: "alice",
		Email:    "alice@example.com",
		Password: "Password123!",
	},
	{
		UserID:   "00000000-0000-4000-8000-000000000002",
		UserName: "bob",
		Email:    "bob@example.com",
		Password: "Password123!",
	},
	{
		UserID:   "00000000-0000-4000-8000-000000000003",
		UserName: "carol",
		Email:    "carol@example.com",
		Password: "Password123!",
	},
}

var (
	tenDays  = 10 * 24 * time.Hour
	oneDay   = 24 * time.Hour
	threeDay = 3 * 24 * time.Hour
)

// tasks はシード対象タスクの一覧です。
// SEARCH_TEXT は投入時に util.NormalizeSearchText で生成します。
var tasks = []SeedTask{
	// --- alice ---
	{
		TaskID:    "00000000-0000-4000-8000-000000000101",
		UserID:    "00000000-0000-4000-8000-000000000001",
		Title:     "四半期レポートの作成",
		Comment:   "上司への提出用。売上データをまとめる",
		Priority:  "high",
		Status:    "in_progress",
		IsPinned:  true,
		DueOffset: &threeDay,
	},
	{
		TaskID:   "00000000-0000-4000-8000-000000000102",
		UserID:   "00000000-0000-4000-8000-000000000001",
		Title:    "チームMTGの議事録まとめ",
		Priority: "medium",
		Status:   "not_started",
	},
	{
		TaskID:   "00000000-0000-4000-8000-000000000103",
		UserID:   "00000000-0000-4000-8000-000000000001",
		Title:    "Blog article draft",
		Comment:  "Tech blog for next release",
		Priority: "low",
		Status:   "not_started",
	},
	{
		TaskID:   "00000000-0000-4000-8000-000000000104",
		UserID:   "00000000-0000-4000-8000-000000000001",
		Title:    "経費精算",
		Priority: "medium",
		Status:   "completed",
	},
	{
		TaskID:    "00000000-0000-4000-8000-000000000105",
		UserID:    "00000000-0000-4000-8000-000000000001",
		Title:     "本番リリース手順書レビュー",
		Comment:   "インフラチームと共有する",
		Priority:  "high",
		Status:    "not_started",
		DueOffset: &oneDay,
	},
	// --- bob ---
	{
		TaskID:    "00000000-0000-4000-8000-000000000201",
		UserID:    "00000000-0000-4000-8000-000000000002",
		Title:     "新機能の設計レビュー",
		Priority:  "high",
		Status:    "in_progress",
		IsPinned:  true,
		DueOffset: &tenDays,
	},
	{
		TaskID:   "00000000-0000-4000-8000-000000000202",
		UserID:   "00000000-0000-4000-8000-000000000002",
		Title:    "バグチケットのトリアージ",
		Priority: "medium",
		Status:   "not_started",
	},
	{
		TaskID:   "00000000-0000-4000-8000-000000000203",
		UserID:   "00000000-0000-4000-8000-000000000002",
		Title:    "Weekly sync notes",
		Priority: "low",
		Status:   "completed",
	},
	// --- carol ---
	{
		TaskID:    "00000000-0000-4000-8000-000000000301",
		UserID:    "00000000-0000-4000-8000-000000000003",
		Title:     "請求書の確認", // 期限切れタスクのテスト用
		Priority:  "medium",
		Status:    "not_started",
		DueOffset: func() *time.Duration { d := -2 * 24 * time.Hour; return &d }(),
	},
}

// sessions はシード対象ログインセッションの一覧です。
// 有効なセッション Cookie を用いた API 動作テストに使用できます。
var sessions = []SeedSession{
	{
		SessionID: "seed-session-alice-00000000000000000000000000000001",
		UserID:    "00000000-0000-4000-8000-000000000001",
		UserAgent: "seed-data (alice)",
	},
	{
		SessionID: "seed-session-bob-000000000000000000000000000000001",
		UserID:    "00000000-0000-4000-8000-000000000002",
		UserAgent: "seed-data (bob)",
	},
}

// Run はシードデータをデータベースに投入します。
//
//   - force=false の場合: 既存行と重複するものはスキップし、不足分のみ追加します (冪等)。
//   - force=true の場合: TASK / LOGIN_SESSION / LOGIN_ACCOUNT を全件削除してから再投入します。
//
// マイグレーションは呼び出し側で適用済みであることを前提とします。
func Run(ctx context.Context, database *sql.DB, force bool) error {
	tx, err := database.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if force {
		for _, stmt := range []string{
			`DELETE FROM TASK`,
			`DELETE FROM LOGIN_SESSION`,
			`DELETE FROM LOGIN_ACCOUNT`,
		} {
			if _, err := tx.ExecContext(ctx, stmt); err != nil {
				return fmt.Errorf("failed to clear tables (force): %w", err)
			}
		}
	}

	now := time.Now()

	for i, u := range users {
		hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password for %s: %w", u.Email, err)
		}
		createdAt := now.Add(-time.Duration(len(users)-i) * 24 * time.Hour)
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO LOGIN_ACCOUNT (USER_ID, USER_NAME, EMAIL, PASSWORD_HASH, CREATED_AT, UPDATED_AT)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (USER_ID) DO NOTHING
		`, u.UserID, u.UserName, u.Email, string(hash), createdAt, createdAt); err != nil {
			return fmt.Errorf("failed to insert LOGIN_ACCOUNT %s: %w", u.Email, err)
		}
	}

	for i, t := range tasks {
		due := any(nil)
		if t.DueOffset != nil {
			due = now.Add(*t.DueOffset)
		}
		createdAt := now.Add(-time.Duration(len(tasks)-i) * time.Hour)
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO TASK (TASK_ID, USER_ID, TITLE, PRIORITY, DUE_DATETIME, STATUS, IS_PINNED, COMMENT, SEARCH_TEXT, CREATED_AT, UPDATED_AT)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
			ON CONFLICT (TASK_ID) DO NOTHING
		`, t.TaskID, t.UserID, t.Title, t.Priority, due, t.Status, t.IsPinned, t.Comment, util.NormalizeSearchText(t.Title, t.Comment), createdAt, createdAt); err != nil {
			return fmt.Errorf("failed to insert TASK %s: %w", t.TaskID, err)
		}
	}

	expiresAt := now.Add(30 * 24 * time.Hour)
	for _, s := range sessions {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO LOGIN_SESSION (SESSION_ID, USER_ID, EXPIRES_AT, USER_AGENT, CREATED_AT)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (SESSION_ID) DO NOTHING
		`, s.SessionID, s.UserID, expiresAt, s.UserAgent, now); err != nil {
			return fmt.Errorf("failed to insert LOGIN_SESSION %s: %w", s.SessionID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit seed transaction: %w", err)
	}
	return nil
}
