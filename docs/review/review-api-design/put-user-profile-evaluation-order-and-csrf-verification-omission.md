# `PUT users/{user_id}` API におけるリクエスト評価順序の未定義および CSRF トークン検証ステップの記載漏れ

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-17 16:54:00
- **Target Files**:
  - [03_users.md](docs/design/api_design/03_users.md)
  - [01_overview.md](docs/design/api_design/01_overview.md)

## 1. 問題の概要
`03_users.md` の `DELETE users/{user_id}` (3.2.3) および `PATCH users/{user_id}/password` (3.2.4) には「リクエスト評価順序」セクションが追加・修正されているが、`PUT users/{user_id}` (3.2.2) には評価順序セクションが定義されていない。

また、`PUT`, `DELETE`, `PATCH` のいずれのエンドポイントにおいても、Errors セクションに定義されている `403 Forbidden`（CSRFトークン検証失敗）の実行タイミングが「リクエスト評価順序」のステップ内に明記されていない。

## 2. 詳細な指摘内容
1. **`PUT users/{user_id}` の評価順序欠落**:
   - `3.2.2` では `400 Bad Request`（入力バリデーション違反）、`401 Unauthorized`（未ログイン）、`403 Forbidden`（CSRFトークン不正）、`404 Not Found`（IDOR/BOLA認可エラー）、`422 Unprocessable Entity`（`SAME_AS_CURRENT_USERNAME`）が定義されているが、優先順位が不明である。
   - 例えば、他ユーザーの `user_id` を指定した不当なリクエスト（IDOR/BOLA）において、入力エラー（400）や同一ユーザー名（422）が先に評価されると、認可チェック（404）の前に内部バリデーション状態が漏洩するリスクが生じる。

2. **CSRF トークン検証（403）の位置づけの不透明さ**:
   - `01_overview.md` (1.2 節) では状態変更リクエストにおける CSRF トークン検証が必須と規定されているが、`3.2.3` や `3.2.4` の「リクエスト評価順序」ステップ（入力検証 → 認可チェック → ビジネスルール → パスワード再認証）内に `403 Forbidden` の検証ステップが含まれていない。
   - CSRF トークン検証を未ログイン認証（401）や認可チェック（404）、パスワードハッシュ照合等のどの段階で実行すべきかが曖昧であり、実装上の不整合や不要な処理の実行につながる懸念がある。

## 3. 推奨される修正案
1. `03_users.md` の `3.2.2` (`PUT users/{user_id}`) に `##### リクエスト評価順序` セクションを追加してください。
2. `3.2.2`, `3.2.3`, `3.2.4` の各エンドポイントの「リクエスト評価順序」に CSRF トークン検証（`403 Forbidden`）のステップを明記してください。

**修正後の評価順序例 (`PUT users/{user_id}`)**:
```markdown
##### リクエスト評価順序
1. **認証・CSRF検証 (`401 Unauthorized` / `403 Forbidden`)**:
   ログインセッションの有効性を確認（未ログイン時は 401 `UNAUTHORIZED`）、および `X-CSRF-Token` ヘッダーを検証（欠落・不一致時は 403 `FORBIDDEN`）。
2. **リクエスト構文・入力バリデーション (`400 Bad Request`)**:
   リクエストボディの `username` に対し、トリム後の文字数（2〜20文字）・使用可能文字（英数字）を検証（不備時は 400 `BAD_REQUEST`）。
3. **認可チェック・IDOR/BOLA検証 (`404 Not Found`)**:
   パスパラメータ `user_id` とセッションユーザーIDの一致を検証（不一致または存在しない場合は 404 `NOT_FOUND`）。
4. **ビジネスルール検証 (`422 Unprocessable Entity`)**:
   トリム後の `username` が現在のユーザー名と同一か検証（同一の場合は 422 `SAME_AS_CURRENT_USERNAME`）。
```
