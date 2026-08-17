# パスワードリセット完了API（3.1.9）における無効・未存在 OTP セッション ID 指定時のエラー定義漏れ

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-17 17:07:00
- **Target Files**:
  - [02_auth.md](docs/design/api_design/02_auth.md)

## 1. 問題の概要
`02_auth.md` 3.1.9 `POST auth/password-reset/reset` の Errors セクションにおいて、リクエストで指定された `otp_session_id` がデータベース上に存在しない場合（未存在・形式不正時）のエラーレスポンス定義が漏れています。

## 2. 詳細な指摘内容
- `3.1.9` の Errors セクション（L337-L342）では以下のエラーコードが定義されています：
  - `400 Bad Request`: `新パスワード要件違反（文字数・文字種・ユーザー名/メール含有違反等、code: "BAD_REQUEST"）`
  - `403 Forbidden`: `未検証のOTPセッション（verified でない場合）でのリセット試行（code: "FORBIDDEN"）`
  - `410 Gone`: `OTP検証完了後の仮セッション有効期限切れ（検証成功後15分経過、code: "GONE"）`
  - `422 Unprocessable Entity`: `現在のパスワードと同一のパスワード（code: "SAME_AS_CURRENT_PASSWORD"）`
- 他の OTP 検証・再送 API（3.1.2, 3.1.3, 3.1.7, 3.1.8 等）では、無効なセッションIDや未存在セッションの指定に対して `400 Bad Request`（code: `"BAD_REQUEST"`）を返却する仕様が明記されています。
- 一方、`3.1.9` の `400 Bad Request` は「新パスワード要件違反」のみが対象となっており、存在しない `otp_session_id` や形式不正な `otp_session_id` が送信された場合に 400 エラーを返却する記述が含まれていません。

## 3. 推奨される修正案
`02_auth.md` 3.1.9 の Errors セクションにおける `400 Bad Request` の定義を以下のように修正・拡充してください：

```markdown
##### Errors
- `400 Bad Request`: 新パスワード要件違反（文字数・文字種・ユーザー名/メール含有違反等）または無効・存在しない `otp_session_id` 指定（code: `"BAD_REQUEST"`）
- `403 Forbidden`: 未検証のOTPセッション（`verified` でない場合）でのリセット試行（code: `"FORBIDDEN"`）
- `410 Gone`: OTP検証完了後の仮セッション有効期限切れ（検証成功後15分経過、code: `"GONE"`）
- `422 Unprocessable Entity`: 現在のパスワードと同一のパスワード（code: `"SAME_AS_CURRENT_PASSWORD"`）
```

---

## 修正完了報告

- **Resolved At**: 2026-08-17 17:10:00
- **Status**: Resolved

### 実施した修正内容
`02_auth.md` 3.1.9 `POST auth/password-reset/reset` の Errors セクションおよびリクエスト評価順序において、未存在・形式不正等の無効な `otp_session_id` が指定された場合に 400 Bad Request（code: `"BAD_REQUEST"`）を返却する旨を明確に規定しました。

### 変更したファイル
- [02_auth.md](docs/design/api_design/02_auth.md)
