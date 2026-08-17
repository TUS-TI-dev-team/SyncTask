# 認証・アカウント登録系 API 群（02_auth.md）におけるレスポンスボディ定義テーブルの全域的欠落

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-17 17:22:00
- **Target Files**:
  - [02_auth.md](docs/design/api_design/02_auth.md)

## 1. 問題の概要
`02_auth.md` の全 12 エンドポイント（3.1.1 〜 3.1.12）において、リクエストボディパラメータには詳細な定義テーブルが用意されているのに対し、レスポンスボディについては JSON のレスポンス例 (`Response (200 OK)` や `Response (201 Created)`) が記載されているのみで、レスポンスボディの各フィールドに関する型、必須／Null許容性、文字数・数値範囲制約、説明を定義するテーブルが一切存在しません。

## 2. 詳細な指摘内容
1. レスポンス例の JSON スニペットのみでは、返却されるデータ型（文字列・数値・オブジェクト）、Null 許容性、フォーマット（例: `expires_in_seconds` が整数型かつ非 Null であること、`otp_session_id` のプレフィックスおよび最大長等）の正式なデータ構造仕様が曖昧になります。
2. 特に以下のエンドポイントで返却されるレスポンスフィールドの定義が不足しています：
   - `3.1.1`, `3.1.6`, `3.1.10` (`request-otp`): `otp_session_id`, `masked_email`, `expires_in_seconds` の型および制約
   - `3.1.2`, `3.1.4` (`verify-otp`, `login`): `user` およびその配下 `id`, `username`, `email` の型および非 Null 保証
   - `3.1.3`, `3.1.8`, `3.1.12` (`resend-otp`): `message`, `masked_email`, `expires_in_seconds` の型および制約
   - `3.1.7`, `3.1.9`, `3.1.11` (`verify-otp`, `reset`, `verify-otp`): `message` の型および非 Null 保証
3. スキーマ定義テーブルが存在しないため、フロントエンドの型定義（TypeScript 型/インターフェース）作成やバックエンドのシリアライザ実装時に仕様の解釈差が生じるリスクがあります。

## 3. 推奨される修正案
`02_auth.md` の各エンドポイント（3.1.1 〜 3.1.12）の `Response (200 OK / 201 Created)` セクションに、レスポンスボディフィールド定義テーブルを追加してください。
例 (`3.1.1 POST auth/register/request-otp`):

| フィールド | 型 | 必須 | 説明 |
| :--- | :--- | :---: | :--- |
| `otp_session_id` | string | ○ | 生成されたOTPセッションID（例: `otp_sess_a1b2c3d4e5`） |
| `masked_email` | string | ○ | マスク処理された送信先メールアドレス |
| `expires_in_seconds` | integer | ○ | OTPの有効期限（秒、デフォルト: 300） |

---

## 修正完了報告

- **Resolved At**: 2026-08-17 17:27:30
- **Status**: Resolved

### 実施した修正内容
`02_auth.md` の全 12 エンドポイント（3.1.1 〜 3.1.12）の `Response (200 OK / 201 Created)` セクションへ、レスポンスボディの各フィールド構造・型・必須性・説明を正確に規定するレスポンス定義テーブルを追加しました。

### 変更したファイル
- [02_auth.md](file:///mnt/c/Users/ayumu_wkkd3w3/BoxForSomething/github/SyncTask/docs/design/api_design/02_auth.md)
