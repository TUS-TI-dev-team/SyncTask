# 概要書におけるパスワードリセット完了時のセッション破棄・Cookie消去記載漏れ

- **Status**: Open
- **Severity**: Medium
- **Created At**: 2026-08-17 17:55:00
- **Target Files**:
  - [01_overview.md](docs/design/api_design/01_overview.md)

## 1. 問題の概要
`01_overview.md` の「1.1 セッション管理 & 認証方式」における「セッション破棄・Cookie消去仕様」の適用トリガー一覧に、パスワードリセット完了処理（`auth/password-reset/reset`）が記載漏れしている。

## 2. 詳細な指摘内容
`02_auth.md` 3.1.9 节 (`POST auth/password-reset/reset`) の仕様では、パスワードリセット完了成功時に「当該OTPセッションおよび該当ユーザーのすべての既存ログインセッションをDBから直ちに物理削除し、Cookieを消去して再ログインを要求する」と規定されており、レスポンスヘッダーとして以下の Cookie 削除ヘッダーが出力されます。

- `Set-Cookie: sync_task_sid=; HttpOnly; Secure; SameSite=Lax; Path=/; Max-Age=0`
- `Set-Cookie: XSRF-TOKEN=; Secure; SameSite=Lax; Path=/; Max-Age=0`

しかし、全体共通仕様を規定する `01_overview.md` の 1.1 節「セッション破棄・Cookie消去仕様」では、対象トリガーとして `ログアウト（auth/logout）、アカウント削除（users/{user_id}）、メールアドレス変更完了（auth/change-email/verify-otp）、パスワード変更（users/{user_id}/password）、および再認証連続失敗による強制破棄（SESSION_DESTROYED）` のみが列挙されており、`パスワードリセット完了（auth/password-reset/reset）` が欠落しています。これにより、共通仕様書と詳細設計書の間でセッション破棄対象イベントの不整合が生じています。

## 3. 推奨される修正案
`01_overview.md` 1.1 節の「セッション破棄・Cookie消去仕様」の本文に `パスワードリセット完了（auth/password-reset/reset）` を追加し、記述を以下のように補正してください。

```markdown
- **セッション破棄・Cookie消去仕様**:
  - ログアウト（`auth/logout`）、アカウント削除（`users/{user_id}`）、メールアドレス変更完了（`auth/change-email/verify-otp`）、パスワード変更（`users/{user_id}/password`）、パスワードリセット完了（`auth/password-reset/reset`）、および再認証連続失敗による強制破棄（`SESSION_DESTROYED`）の発生時は、サーバー側で DB 上のセッションレコードを物理削除すると同時に、レスポンスヘッダーで以下の Cookie 削除ヘッダーを出力してクライアント側の Cookie を直ちに無効化・消去します。
    - `Set-Cookie: sync_task_sid=; HttpOnly; Secure; SameSite=Lax; Path=/; Max-Age=0`
    - `Set-Cookie: XSRF-TOKEN=; Secure; SameSite=Lax; Path=/; Max-Age=0`
```
