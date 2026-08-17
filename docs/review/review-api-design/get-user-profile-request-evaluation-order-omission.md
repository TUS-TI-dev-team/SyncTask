# `GET users/{user_id}` API における「リクエスト評価順序」セクションの欠落

- **Status**: Resolved

---

## 修正完了報告

- **Resolved At**: 2026-08-17 17:03:00
- **Status**: Resolved

### 実施した修正内容
`03_users.md` 3.2.1 節 (`GET users/{user_id}`) に `##### リクエスト評価順序` セクションを追加し、認証検証 (401) を認可チェック/IDOR検証 (404) より先に評価する順序を明記しました。

### 変更したファイル
- [03_users.md](docs/design/api_design/03_users.md)
- **Severity**: Minor
- **Created At**: 2026-08-17 17:01:00
- **Target Files**:
  - [03_users.md](docs/design/api_design/03_users.md)

## 1. 問題の概要
`03_users.md` 内のエンドポイント群（`PUT users/{user_id}`, `DELETE users/{user_id}`, `PATCH users/{user_id}/password`）にはすべて `##### リクエスト評価順序` セクションが用意され、各検証処理の優先順位が明確化されているが、`GET users/{user_id}` (3.2.1 節) のみ「リクエスト評価順序」セクションが存在しない。

## 2. 詳細な指摘内容
1. **評価順序の未定義**:
   - `GET users/{user_id}` では `401 Unauthorized`（未ログイン・セッション無効）と `404 Not Found`（他ユーザーの `user_id` 指定による認可エラー/IDOR/BOLA）が定義されているが、どちらの検証を先に行うかの順序が明記されていない。
   - 他のユーザー管理 API との統一性を保ち、認証検証（401）を認可チェック（404）より先に評価する挙動を明確にするため、「リクエスト評価順序」セクションを新設することが望ましい。

## 3. 推奨される修正案
`03_users.md` の 3.2.1 節 (`GET users/{user_id}`) に以下の `##### リクエスト評価順序` セクションを追加してください。

**追加案の例**:
```markdown
##### リクエスト評価順序
1. **認証検証 (`401 Unauthorized`)**:
   ログインセッションの有効性を確認（未ログインまたはセッション無効時は 401 `UNAUTHORIZED`）。
2. **認可チェック・IDOR/BOLA検証 (`404 Not Found`)**:
   パスパラメータ `user_id` とセッションユーザーIDの一致を検証（不一致または存在しない場合は 404 `NOT_FOUND`）。
```
