# `01_overview.md` 1.1 節におけるログインセッション Cookie (`sync_task_sid`) の `Path` および `Max-Age` 属性定義の欠落

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-17 17:15:00
- **Target Files**:
  - [01_overview.md](docs/design/api_design/01_overview.md)

## 1. 問題の概要
`01_overview.md` の 1.1 節「セッション管理 & 認証方式」において、CSRF トークン Cookie（`XSRF-TOKEN`）には `Path=/; Max-Age=2592000` の属性指定が具体的に明記されているのに対し、ログインセッション Cookie（`sync_task_sid`）には `HttpOnly`, `Secure`, `SameSite=Lax` のみが記載されており、`Path=/` および秒単位の `Max-Age=2592000`（43,200分 / 30日）の `Set-Cookie` 属性定義が欠落している。

## 2. 詳細な指摘内容
1. **Cookie 送信スコープの不完全な定義（L17）**:
   L17 にて `HttpOnly`, `Secure`, `SameSite=Lax` 属性が付与された Cookie と記載されているが、`Path=/` 属性が定義されていないため、Cookie のスコープに関する仕様が不完全である。
2. **ヘッダーフォーマット記述の不統一（L17, L21, L29）**:
   CSRF トークン Cookie については L21 や L29 で `Set-Cookie: XSRF-TOKEN=<csrf_token>; Secure; SameSite=Lax; Path=/; Max-Age=2592000` と具体的な `Set-Cookie` レスポンスヘッダー仕様が示されているのに対し、`sync_task_sid` の項には具現化されたヘッダーフォーマットが示されていない。`02_auth.md`（L178）では `Set-Cookie: sync_task_sid=<session_id>; HttpOnly; Secure; SameSite=Lax; Path=/; Max-Age=2592000` と定義されており、概要書側で属性の記載漏れが生じている。

## 3. 推奨される修正案
`01_overview.md` 1.1 節 L17-L19 の記述を更新し、`Path=/` および `Max-Age=2592000` を含めた `Set-Cookie` ヘッダーフォーマットを明記してください。

```markdown
- **ログインセッション**:
  - トークンをリクエスト本文で送受信するのではなく、`HttpOnly`, `Secure`, `SameSite=Lax`, `Path=/`, `Max-Age=2592000` 属性が付与されたセッションCookie（名称: `sync_task_sid`、例: `Set-Cookie: sync_task_sid=<session_id>; HttpOnly; Secure; SameSite=Lax; Path=/; Max-Age=2592000`）によって管理します。
```

---

## 修正完了報告

- **Resolved At**: 2026-08-17 17:18:00
- **Status**: Resolved

### 実施した修正内容
`01_overview.md` 1.1 節「セッション管理 & 認証方式」のログインセッション Cookie (`sync_task_sid`) の記述を更新し、`Path=/` および `Max-Age=2592000` 属性と `Set-Cookie: sync_task_sid=<session_id>; HttpOnly; Secure; SameSite=Lax; Path=/; Max-Age=2592000` の具体例を追記しました。

### 変更したファイル
- [01_overview.md](docs/design/api_design/01_overview.md)

