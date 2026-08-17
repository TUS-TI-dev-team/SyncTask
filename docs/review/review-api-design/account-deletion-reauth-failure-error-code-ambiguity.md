# アカウント削除APIにおける再認証失敗エラーレスポンスの曖昧さ

- **Status**: Open
- **Severity**: Minor
- **Created At**: 2026-08-17 15:01:00
- **Target Files**:
  - [api_design.md](docs/design/api_design.md)

## 1. 問題の概要
`DELETE users/{user_id}`（3.2.3）のエラー仕様（L548）にて「`400 Bad Request`: パスワード再認証失敗（5回連続失敗時はセッション強制破棄・401 Unauthorized）」と記載されているが、1つの括弧内に2つの異なるHTTPステータスコード（400と401）が混在しており、エラーレスポンスの条件分岐が不明確である。

同様の問題は `PATCH users/{user_id}/password`（3.2.4, L579）にも存在する。

## 2. 詳細な指摘内容
- **API設計書 L548**: `400 Bad Request`: パスワード再認証失敗（5回連続失敗時はセッション強制破棄・401 Unauthorized）
- **API設計書 L579**: `400 Bad Request`: 新パスワード要件違反、または現在のパスワード不一致（5回連続失敗でセッション破棄・401）

以下の点が不明確:
1. パスワード不一致（1〜4回目）は `400 Bad Request` を返すのか、`401 Unauthorized` を返すのか。パスワードの認証失敗は `401` が適切ではないか。
2. 5回目の連続失敗時はセッション破棄と共に `401 Unauthorized` を返すとのことだが、これは別のエラーコードを使用するのか（例: `"code": "SESSION_DESTROYED"` 等）。
3. L579 では「新パスワード要件違反」と「現在のパスワード不一致」が同じ `400` に混在しており、フロントエンドがエラー原因を区別できない（バリデーションエラーなのか認証失敗なのか）。

## 3. 推奨される修正案
エラー仕様を明確に分離して記載する:

```
##### Errors
- `400 Bad Request`: 新パスワード要件違反（文字数・文字種・包含禁止）
- `401 Unauthorized`: 現在のパスワード不一致（code: `REAUTH_FAILED`。5回連続失敗時は
  セッション強制破棄・Cookie消去。code: `SESSION_DESTROYED`）
- `404 Not Found`: 認可エラー
- `422 Unprocessable Entity`: 新パスワードが現在のパスワードと同一
```
