# セッション・CSRF Cookie 破棄ヘッダーにおけるセキュリティ属性（HttpOnly, Secure, SameSite）の指定漏れ

- **Status**: Resolved
- **Severity**: Major
- **Created At**: 2026-08-17 17:43:00
- **Target Files**:
  - [01_overview.md](docs/design/api_design/01_overview.md)

## 1. 問題の概要
`01_overview.md` 1.1 節の「セッション破棄・Cookie消去仕様」において、ログアウトやアカウント削除等の各種セッション破棄時にレスポンス出力する Cookie 削除ヘッダーとして `Set-Cookie: sync_task_sid=; Path=/; Max-Age=0` および `Set-Cookie: XSRF-TOKEN=; Path=/; Max-Age=0` が定義されているが、Cookie 発行・更新時に付与されていた `Secure`, `SameSite=Lax` 属性（および `sync_task_sid` の `HttpOnly` 属性）が欠落している。

## 2. 詳細な指摘内容
1. **発行時・自動延長時の Cookie 仕様（`01_overview.md` L17, L21, L33）**:
   - `sync_task_sid`: `Set-Cookie: sync_task_sid=<session_id>; HttpOnly; Secure; SameSite=Lax; Path=/; Max-Age=2592000`
   - `XSRF-TOKEN`: `Set-Cookie: XSRF-TOKEN=<csrf_token>; Secure; SameSite=Lax; Path=/; Max-Age=2592000`

2. **セッション破棄時の削除ヘッダー記述（`01_overview.md` L24-L25）**:
   - `Set-Cookie: sync_task_sid=; Path=/; Max-Age=0`
   - `Set-Cookie: XSRF-TOKEN=; Path=/; Max-Age=0`

3. **セキュリティおよびブラウザ互換性上の影響**:
   RFC 6265bis および現代の主要ブラウザ（Chrome, Safari, Firefox 等）のセキュリティ仕様では、`Secure` や `SameSite=Lax`, `HttpOnly` フラグが付与されて保存された Cookie に対して、属性が一致しない削除ヘッダー（`Set-Cookie: sync_task_sid=; Path=/; Max-Age=0`）を送信した場合、ブラウザの Cookie ストレージ上で既存のセキュア Cookie の破棄が適切に行われず、クライアント側に無効なセッション Cookie が残存するリスクが存在する。

## 3. 推奨される修正案
`01_overview.md` 1.1 節の「セッション破棄・Cookie消去仕様」における Cookie 削除ヘッダーの定義に `HttpOnly`, `Secure`, `SameSite=Lax` 属性を明記し、以下のように修正してください。

```markdown
- **セッション破棄・Cookie消去仕様**:
  - ログアウト（`auth/logout`）、アカウント削除（`users/{user_id}`）、メールアドレス変更完了（`auth/change-email/verify-otp`）、パスワード変更（`users/{user_id}/password`）、および再認証連続失敗による強制破棄（`SESSION_DESTROYED`）の発生時は、サーバー側で DB 上のセッションレコードを物理削除すると同時に、レスポンスヘッダーで以下の Cookie 削除ヘッダーを出力してクライアント側の Cookie を直ちに無効化・消去します。
    - `Set-Cookie: sync_task_sid=; HttpOnly; Secure; SameSite=Lax; Path=/; Max-Age=0`
    - `Set-Cookie: XSRF-TOKEN=; Secure; SameSite=Lax; Path=/; Max-Age=0`
```

---

## 修正完了報告

- **Resolved At**: 2026-08-17 17:50:00
- **Status**: Resolved

### 実施した修正内容
`01_overview.md` 1.1 節の「セッション破棄・Cookie消去仕様」における Cookie 削除ヘッダーに `HttpOnly`, `Secure`, `SameSite=Lax` 属性を追加し、発行時と属性が完全に一致したヘッダー（`Set-Cookie: sync_task_sid=; HttpOnly; Secure; SameSite=Lax; Path=/; Max-Age=0` および `Set-Cookie: XSRF-TOKEN=; Secure; SameSite=Lax; Path=/; Max-Age=0`）が出力されるよう修正しました。

### 変更したファイル
- [01_overview.md](docs/design/api_design/01_overview.md)
