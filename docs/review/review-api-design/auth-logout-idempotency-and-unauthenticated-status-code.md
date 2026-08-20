# `POST auth/logout` における未ログイン・セッション無効時の 401 エラー返却による冪等性要件違反

- **Status**: Resolved
- **Severity**: Major
- **Created At**: 2026-08-17 16:45:00
- **Target Files**:
  - [02_auth.md](docs/design/api_design/02_auth.md)

## 1. 問題の概要
`02_auth.md` の 3.1.5 `POST auth/logout` において、未ログインまたはセッション無効な状態でのログアウト要求に対して `- 401 Unauthorized: 未ログイン (code: "UNAUTHORIZED")` を返却する仕様となっている。
しかし、要件定義書（`requirements.md` L67）では「セッションが既に無効化または期限切れの状態でログアウトが実行された場合もエラーとせず、セッションCookieを消去してログイン画面へリダイレクトする（冪等性の確保）」と定められており、401 エラーの返却は要件定義書の冪等性要件に違反している。

## 2. 詳細な指摘内容
`docs/design/api_design/02_auth.md` L158-L161:
```markdown
##### Errors
- 401 Unauthorized: 未ログイン（code: "UNAUTHORIZED"）
- 403 Forbidden: CSRFトークン不正（code: "FORBIDDEN"）
```

- **要件定義書との乖離**: `requirements.md` の L67 にて、セッションが無効または期限切れの場合でもエラーにせず `200 OK` を返してセッションCookieを消去する「冪等性（Idempotency）」の担保が明記されている。
- **クライアントへの影響**: 401 Unauthorized を返却してしまうと、フロントエンド側で不要なエラーハンドリングやトースト表示が発生し、ユーザー体験を損ねる恐れがある。また、セッションCookieの消去処理がクライアント側で確実に完了しないリスクが生じる。

## 3. 推奨される修正案
`POST auth/logout` のレスポンスおよびエラー仕様を以下のように修正してください：

1. 未ログインまたはセッション期限切れの状態で `POST auth/logout` が呼び出された場合でも、`200 OK` レスポンスを返却し、`Set-Cookie: sync_task_sid=; Max-Age=0` および `Set-Cookie: XSRF-TOKEN=; Max-Age=0` を付与してCookieを確実に破棄する仕様に変更する。
2. エラーレスポンスからは `401 Unauthorized` を削除し、CSRFトークン検証失敗時の `403 Forbidden` のみを残す。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 16:50:00
- **Status**: Resolved

### 実施した修正内容
`02_auth.md` 3.1.5 (`POST auth/logout`) の仕様を改定し、未ログイン/セッション期限切れ時もエラーにせず `200 OK` と Cookie 削除ヘッダーを返却して冪等性を確保するよう修正しました。`Errors` から `401 Unauthorized` を削除しました。

### 変更したファイル
- [02_auth.md](docs/design/api_design/02_auth.md)
