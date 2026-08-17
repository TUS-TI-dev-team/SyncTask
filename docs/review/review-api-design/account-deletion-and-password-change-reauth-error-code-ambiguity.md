# アカウント削除およびパスワード変更APIにおける再認証失敗・入力エラーの400/401ステータスコード混在とコード未定義

- **Status**: Resolved
- **Severity**: Major
- **Created At**: 2026-08-17 16:25:00
- **Target Files**:
  - [03_users.md](docs/design/api_design/03_users.md)

## 1. 問題の概要
`DELETE users/{user_id}`（3.2.3）および `PATCH users/{user_id}/password`（3.2.4）のエラー仕様において、`400 Bad Request` の項目に「パスワード再認証失敗（5回連続失敗時はセッション強制破棄・401 Unauthorized）」や「新パスワード要件違反、または現在のパスワード不一致」と記述されており、400と401のレスポンス定義が同一項目内に混在している。

このため、リクエストパラメータ自体の不足（400 Bad Request）とパスワード不一致（401 Unauthorized）の区別が曖昧であり、さらに再認証失敗時やセッション破棄時のエラーコード（`REAUTH_FAILED`, `SESSION_DESTROYED`）が明記されていない。

## 2. 詳細な指摘内容
- **`03_users.md` L90 (`DELETE users/{user_id}`)**:
  ` - 400 Bad Request: パスワード再認証失敗（5回連続失敗時はセッション強制破棄・401 Unauthorized）`
- **`03_users.md` L121 (`PATCH users/{user_id}/password`)**:
  ` - 400 Bad Request: 新パスワード要件違反、または現在のパスワード不一致（5回連続失敗でセッション破棄・401）`

### 問題点
1. **HTTPステータスコードの誤用と混在**: パスワードの認証失敗は `401 Unauthorized` を返却すべきであり、`400 Bad Request` 内に `(401)` を括弧書きでネストして記述するのはエラー応答の仕様として不正である。
2. **フロントエンドのエラー識別不能**: `PATCH users/{user_id}/password` において「新パスワードのバリデーションエラー」と「現在のパスワード不一致」が同じ `400` として括られており、クライアント側でどちらのエラーが発生したのか判別できない。
3. **エラーコードの未定義**: パスワード再認証失敗時（1〜4回目）の `code: "REAUTH_FAILED"` や、5回連続失敗によるセッション強制破棄時の `code: "SESSION_DESTROYED"`、入力形式不正時の `code: "BAD_REQUEST"` の記述が欠落している。

## 3. 推奨される修正案
`DELETE users/{user_id}` および `PATCH users/{user_id}/password` のエラー定義を、以下のように明示的なステータスコードとエラーコード毎に分離して記述してください。

### `DELETE users/{user_id}` 修正案
```markdown
##### Errors
- `400 Bad Request`: リクエストボディ不正・パスワード未入力（code: `"BAD_REQUEST"`）
- `401 Unauthorized`: 
  - 未ログイン・セッション無効（code: `"UNAUTHORIZED"`）
  - パスワード再認証失敗 1〜4回目（code: `"REAUTH_FAILED"`）
  - パスワード再認証失敗 5回連続達成分（セッション強制破棄・Cookie消去、code: `"SESSION_DESTROYED"`）
- `403 Forbidden`: CSRFトークン不正（code: `"FORBIDDEN"`）
- `404 Not Found`: 認可エラー（他ユーザーID指定または存在しないユーザー）
```

### `PATCH users/{user_id}/password` 修正案
```markdown
##### Errors
- `400 Bad Request`: 新パスワード要件違反（文字数・文字種・ユーザー名/メール含有違反等、code: `"BAD_REQUEST"`）
- `401 Unauthorized`: 
  - 未ログイン・セッション無効（code: `"UNAUTHORIZED"`）
  - 現在のパスワード不一致 1〜4回目（code: `"REAUTH_FAILED"`）
  - 現在のパスワード不一致 5回連続達成分（セッション強制破棄・Cookie消去、code: `"SESSION_DESTROYED"`）
- `403 Forbidden`: CSRFトークン不正（code: `"FORBIDDEN"`）
- `404 Not Found`: 認可エラー
- `422 Unprocessable Entity`: 新パスワードが現在のパスワードと同一（code: `"SAME_AS_CURRENT_PASSWORD"`）
```

---

## 修正完了報告

- **Resolved At**: 2026-08-17 16:42:00
- **Status**: Resolved

### 実施した修正内容
`DELETE users/{user_id}` および `PATCH users/{user_id}/password` のエラー定義において、400 Bad Request と 401 Unauthorized を明確に分離し、再認証失敗時（1〜4回目）の `code: "REAUTH_FAILED"`、5回連続失敗時の `code: "SESSION_DESTROYED"`、新パスワード同一時の `code: "SAME_AS_CURRENT_PASSWORD"`、リクエスト不正時の `code: "BAD_REQUEST"` などの具体コードを併記しました。

### 変更したファイル
- [03_users.md](file:////mnt/c/Users/ayumu_wkkd3w3/BoxForSomething/github/SyncTask/docs/design/api_design/03_users.md)
