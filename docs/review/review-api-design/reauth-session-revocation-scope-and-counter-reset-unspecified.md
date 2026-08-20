# パスワード再認証失敗 5 回達成分におけるセッション破棄対象範囲およびログアウト時の失敗カウンターリセット規定の不足

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-17 16:54:00
- **Target Files**:
  - [03_users.md](docs/design/api_design/03_users.md)
  - [database_design.md](docs/design/database_design.md)

## 1. 問題の概要
`03_users.md` の `DELETE users/{user_id}` (3.2.3) および `PATCH users/{user_id}/password` (3.2.4) のエラー仕様において、パスワード再認証が 5 回連続で失敗した場合の `401 SESSION_DESTROYED` エラーについて「セッション強制破棄」と記載されている。
しかし、この強制破棄の対象が「リクエストを送信した操作中の当該ログインセッションのみ」なのか「該当ユーザーのアカウントに紐づくすべてのログインセッション」なのかが明記されていない。

また、`database_design.md` (L24) では `LOGIN_ACCOUNT` テーブルの `REAUTH_FAILED_COUNT` が「成功時またはセッション破棄時に0リセット」と規定されているが、ユーザーが自主的にログアウト (`auth/logout`) した場合やセッション有効期限切れ時に `REAUTH_FAILED_COUNT` をリセットする運用方針が `03_users.md` に明記されていない。

## 2. 詳細な指摘内容
1. **セッション破棄範囲の曖昧さ**:
   - `DELETE users/{user_id}` の成功時 (200 OK) や `PATCH users/{user_id}/password` の成功時 (200 OK) では、概要欄にて「所有タスクデータおよび全セッションは物理削除」や「全セッションを一括物理削除」と全セッション破棄であることが明示されている。
   - 一方、再認証 5 回連続失敗時の `401 SESSION_DESTROYED` レスポンス（L113, L168）では「セッション強制破棄・Cookie消去」としか記載されておらず、当該単一セッションの無効化なのか全セッションの一括無効化なのかが開発者によって解釈の揺れが生じる。

2. **`REAUTH_FAILED_COUNT` のリセットタイミングの不透明さ**:
   - DB 設計上、`REAUTH_FAILED_COUNT` はアカウント単位 (`LOGIN_ACCOUNT`) で管理されているため、単一セッションが破棄された場合でも別セッションや次回ログイン時にカウンター値が残留・波及する可能性がある。
   - 再認証成功時、5 回連続失敗に伴うセッション破棄時、およびログアウト操作時に `REAUTH_FAILED_COUNT` を確実かつ明示的に 0 へリセットする挙動を `03_users.md` にも統一的に記載することが望ましい。

## 3. 推奨される修正案
`03_users.md` の 3.2.3 および 3.2.4 の再認証失敗仕様・Errors セクションに以下の補足規定を明記してください。

```markdown
※ パスワード再認証失敗が5回連続に達した場合（`SESSION_DESTROYED`）、セキュリティ保護のため操作中の該当ログインセッションを直ちに物理削除し、Cookieを消去します。また、再認証失敗カウンター（`REAUTH_FAILED_COUNT`）は、再認証成功時、5回失敗によるセッション強制破棄時、およびログアウト時に 0 にリセットされます。
```

## 修正完了報告

- **Resolved At**: 2026-08-17 16:56:00
- **Status**: Resolved

### 実施した修正内容
`03_users.md` の 3.2.3 および 3.2.4 にて、パスワード再認証 5 回連続失敗時の `SESSION_DESTROYED` による破棄対象が「操作中の該当ログインセッションのみ」であること、および `REAUTH_FAILED_COUNT` カウンターが再認証成功時・5回失敗時・ログアウト（`auth/logout`）時に 0 リセットされることを明確に追記しました。

### 変更したファイル
- [03_users.md](docs/design/api_design/03_users.md)
