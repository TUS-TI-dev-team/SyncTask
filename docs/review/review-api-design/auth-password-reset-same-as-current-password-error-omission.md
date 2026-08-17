# パスワードリセット完了API (`POST auth/password-reset/reset`) における同一パスワード指定エラー (`422 SAME_AS_CURRENT_PASSWORD`) の定義漏れ

- **Status**: Resolved

---

## 修正完了報告

- **Resolved At**: 2026-08-17 17:03:00
- **Status**: Resolved

### 実施した修正内容
`02_auth.md` 3.1.9 節 (`POST auth/password-reset/reset`) の `Errors` に `422 SAME_AS_CURRENT_PASSWORD` を追加し、現在のパスワードと同一のパスワードが指定された場合に同一パスワード変更を禁止するエラーを返却する仕様を明記しました。また `01_overview.md` のエラーコード一覧表の説明欄も「パスワード変更/リセット時の同一パスワード指定」に更新しました。

### 変更したファイル
- [02_auth.md](docs/design/api_design/02_auth.md)
- [01_overview.md](docs/design/api_design/01_overview.md)
- **Severity**: Medium
- **Created At**: 2026-08-17 17:01:00
- **Target Files**:
  - [02_auth.md](docs/design/api_design/02_auth.md)
  - [01_overview.md](docs/design/api_design/01_overview.md)

## 1. 問題の概要
`01_overview.md` の 1.3 節「代表的なエラーコード一覧」では、`422 SAME_AS_CURRENT_PASSWORD` が定義されており、`03_users.md` の 3.2.4 (`PATCH users/{user_id}/password`) では現在のパスワードと同一のパスワードを設定しようとした際に `422 SAME_AS_CURRENT_PASSWORD` を返却することが明確に定義されています。しかし、`02_auth.md` の 3.1.9 `POST auth/password-reset/reset` のエラー仕様 (`##### Errors`) では `422 SAME_AS_CURRENT_PASSWORD` が完全に欠落しています。

## 2. 詳細な指摘内容
1. **パスワード再設定ポリシーの不整合**:
   - ユーザーがパスワードリセット手続きを実行し、新パスワードとして現在のパスワードと全く同じパスワードを入力した場合の挙動が 3.1.9 で定義されていません。
   - ログイン中のパスワード変更 API (`PATCH users/{user_id}/password`) ではセキュリティポリシーとして `422 SAME_AS_CURRENT_PASSWORD` を返却して同一パスワードへの変更を拒否しているにもかかわらず、パスワードリセット API (`POST auth/password-reset/reset`) の Errors 一覧（L312-L316）には `400 Bad Request`, `403 Forbidden`, `410 Gone` のみが定義されており、422 エラーの扱いが不明です。
2. **仕様上の曖昧さ**:
   - リセット時も現在のパスワードと同一のパスワードを禁止する仕様である場合、422 エラーが漏れていることになります。
   - 逆にリセット時は現在のパスワードとの比較を行わず変更を許可する仕様である場合も、そのセキュリティ上の判断方針が文書化されていません。

## 3. 推奨される修正案
`02_auth.md` の 3.1.9 `POST auth/password-reset/reset` の本文および Errors セクションに、パスワードリセット時における現在のパスワードと同一パスワード設定時の挙動を明記してください。

同パスワード再設定を禁止する場合は、以下を Errors に追加してください：
- `- 422 Unprocessable Entity: 現在のパスワードと同一のパスワード（code: "SAME_AS_CURRENT_PASSWORD"）`
