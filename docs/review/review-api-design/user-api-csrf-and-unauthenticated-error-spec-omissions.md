# ユーザー管理API群におけるCSRF検証失敗（403）および未ログイン（401）エラー仕様の記述漏れ

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-17 16:25:00
- **Target Files**:
  - [03_users.md](docs/design/api_design/03_users.md)

## 1. 問題の概要
`01_overview.md`（1.2/1.3 節）では、状態を変更するすべてのHTTPリクエスト（`PUT`, `DELETE`, `PATCH`）において CSRF トークン（`X-CSRF-Token`）の検証を必須とし、不正時は `403 Forbidden` (`FORBIDDEN`) を返却すること、および未認証リクエストには `401 Unauthorized` (`UNAUTHORIZED`) を返却することが定義されている。

しかし、`03_users.md` の各エンドポイント（`PUT users/{user_id}`, `DELETE users/{user_id}`, `PATCH users/{user_id}/password`）のエラー仕様セクションにおいて、`403 Forbidden` および `401 Unauthorized` の記述が漏れている。

## 2. 詳細な指摘内容
- `PUT users/{user_id}` (3.2.2 L58-L62):
  - エラー一覧に `401 Unauthorized` (未ログイン) および `403 Forbidden` (CSRFトークン不正) が記載されていない。
- `DELETE users/{user_id}` (3.2.3 L89-L92):
  - エラー一覧に `403 Forbidden` (CSRFトークン不正) が記載されていない。また、未認証時の `401 Unauthorized` と再認証失敗時の `401` が区別されていない。
- `PATCH users/{user_id}/password` (3.2.4 L120-L124):
  - エラー一覧に `403 Forbidden` (CSRFトークン不正) および未ログイン時の `401 Unauthorized` が記載されていない。

## 3. 推奨される修正案
`03_users.md` 内の全エンドポイントの `##### Errors` セクションに、共通仕様に基づき `401 Unauthorized` (`UNAUTHORIZED`) および `403 Forbidden` (`FORBIDDEN`) を明記してください。

```markdown
##### Errors
- `401 Unauthorized`: 未ログインまたはセッションが無効・期限切れ（code: `"UNAUTHORIZED"`）
- `403 Forbidden`: CSRFトークンヘッダーの欠落または不一致（code: `"FORBIDDEN"`）
```

---

## 修正完了報告

- **Resolved At**: 2026-08-17 16:42:00
- **Status**: Resolved

### 実施した修正内容
`03_users.md` 内の全エンドポイント（`GET`, `PUT`, `DELETE`, `PATCH`）の `##### Errors` セクションに、未ログイン時の `401 Unauthorized` (`UNAUTHORIZED`) および CSRFトークン不正時の `403 Forbidden` (`FORBIDDEN`) の定義を追加しました。

### 変更したファイル
- [03_users.md](file:////mnt/c/Users/ayumu_wkkd3w3/BoxForSomething/github/SyncTask/docs/design/api_design/03_users.md)
