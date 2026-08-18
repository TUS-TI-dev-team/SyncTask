# MAIL_AUTH_LOGテーブルのEVENT_TYPE備考欄における全角カンマの混入

- **Status**: Open
- **Severity**: Minor
- **Created At**: 2026-08-18 22:15:00
- **Target Files**:
  - [database_design.md](docs/design/database_design.md)

## 1. 問題の概要
`docs/design/database_design.md` の「6.3 メール認証ログ (MAIL_AUTH_LOG)」テーブル定義において、`EVENT_TYPE` カラムの備考欄に全角カンマ（`，`）が混入しており、ドキュメントの書式に不揃いが生じています。

## 2. 詳細な指摘内容
`database_design.md` L169 に以下の記述があります：

```markdown
| 処理イベント種別 | `EVENT_TYPE` | `VARCHAR(30)` / `NOT NULL` | `ISSUED` (発行), `VERIFY_SUCCESS` (検証成功), `VERIFY_FAILED` (検証失敗), `RESEND_REQUESTED` (手動再送), `AUTO_RESEND` (5回失敗時自動処理)，`EXPIRED` (有効期限切れ) |
```

### 問題点：
- `AUTO_RESEND` の直後の区切り記号が全角カンマ `，` になっており、他の列挙値との区切り記号（半角カンマ + 半角スペース `, `）と不揃いになっています。

## 3. 推奨される修正案
全角カンマを半角カンマおよびスペース `, ` に修正してください：

```markdown
| 処理イベント種別 | `EVENT_TYPE` | `VARCHAR(30)` / `NOT NULL` | `ISSUED` (発行), `VERIFY_SUCCESS` (検証成功), `VERIFY_FAILED` (検証失敗), `RESEND_REQUESTED` (手動再送), `AUTO_RESEND` (5回失敗時自動処理), `EXPIRED` (有効期限切れ) |
```
