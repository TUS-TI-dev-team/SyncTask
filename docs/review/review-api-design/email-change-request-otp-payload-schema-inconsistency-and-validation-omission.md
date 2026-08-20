# メールアドレス変更リクエストAPI (`POST auth/change-email/request-otp`) におけるレスポンス構造不整合・バリデーション表およびエラー定義の欠落

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-17 16:35:00
- **Target Files**:
  - [02_auth.md](docs/design/api_design/02_auth.md)

## 1. 問題の概要
`POST auth/change-email/request-otp`（3.1.10）において、レスポンスJSON構造が他機能の `request-otp` / `resend-otp`（`3.1.1`, `3.1.6`, `3.1.12`）と不統一であり `masked_email` フィールドが欠落しています。また、`new_email` のバリデーション制約表、および `##### Errors` セクションが存在しません。

## 2. 詳細な指摘内容
1. **レスポンス JSON 構造の不統一（`masked_email` の欠落）**:
   - `3.1.1 POST auth/register/request-otp`, `3.1.6 POST auth/password-reset/request-otp`, および同機能の再送API `3.1.12 POST auth/change-email/resend-otp` のレスポンスには、ユーザーに送信先を確認させるための `masked_email`（例: `"new_**********@example.com"`）が含まれています。
   - しかし `3.1.10 POST auth/change-email/request-otp`（L268-L274）のレスポンスには `otp_session_id` と `expires_in_seconds` のみが定義されており、`masked_email` が欠落しています。
2. **バリデーション制約表の欠落**:
   - `3.1.1` や `3.1.4` ではリクエストパラメータの制約表が定義されていますが、`3.1.10` には `new_email` に対するバリデーション表が存在せず、有効メールアドレス形式、空白トリム、小文字正規化、現在登録中のメールアドレスと同一の場合のエラー（`422 Unprocessable Entity`）などの仕様が明示されていません。
3. **エラー仕様 (`##### Errors`) セクションの欠落**:
   - `3.1.10` に `##### Errors` セクションが存在せず、`400 Bad Request`（入力形式違反）、`401 Unauthorized`（未ログイン）、`403 Forbidden`（CSRFトークン不正）、`422 Unprocessable Entity`（現在のメールアドレスと同一）のエラー定義が不透明です。

## 3. 推奨される修正案
1. `3.1.10` の Response (200 OK) に `masked_email` を追加し、レスポンス構造を統一してください:
   ```json
   {
     "otp_session_id": "otp_sess_chg_998877",
     "masked_email": "new_**********@example.com",
     "expires_in_seconds": 300
   }
   ```
2. `3.1.10` に `new_email` のバリデーションテーブルを追加してください:
   | フィールド | 型 | 必須 | 制約・バリデーション |
   | :--- | :--- | :---: | :--- |
   | `new_email` | string | ○ | 有効なメールアドレス形式、前後の空白トリム、小文字正規化。現在のメールアドレスと同一の場合は 422 エラー |
3. `3.1.10` に `##### Errors` セクションを追加し、`400 Bad Request`, `401 Unauthorized`, `403 Forbidden`, `422 Unprocessable Entity` を明記してください。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 16:42:00
- **Status**: Resolved

### 実施した修正内容
`POST auth/change-email/request-otp` (3.1.10) のレスポンスに `masked_email` を追加して他機能と構造を統一し、`new_email` のバリデーションテーブル（形式・トリム・小文字化・422チェック）、CSRFヘッダー指定、アカウント列挙防止注記、および `##### Errors` セクションを追加しました。

### 変更したファイル
- [02_auth.md](file:////mnt/c/Users/ayumu_wkkd3w3/BoxForSomething/github/SyncTask/docs/design/api_design/02_auth.md)
