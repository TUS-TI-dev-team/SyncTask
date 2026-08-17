# CSRFトークン取得方法・エンドポイントの欠落

- **Status**: Open
- **Severity**: Major
- **Created At**: 2026-08-17 15:01:00
- **Target Files**:
  - [api_design.md](docs/design/api_design.md)
  - [requirements.md](docs/req-def/requirements.md)

## 1. 問題の概要
API設計書のセキュリティ仕様（1.2）にて「状態を変更するすべてのHTTPメソッド（`POST`, `PUT`, `PATCH`, `DELETE`）において CSRFトークンの検証を必須とします。クライアントは `X-CSRF-Token` リクエストヘッダーに有効なCSRFトークンを付与して送信します。」と規定されているが、クライアントがCSRFトークンをどのように取得するかの仕様がAPI設計書に一切記載されていない。

## 2. 詳細な指摘内容
CSRFトークンの検証を必須としているにもかかわらず、以下の情報が欠落している:

1. **CSRFトークンの発行・取得方法**: 一般的な実装パターンとして以下がある:
   - 専用のCSRFトークン取得エンドポイント（例: `GET /api/csrf-token`）
   - ログイン成功レスポンスのボディまたはヘッダーにCSRFトークンを含める
   - `Set-Cookie` で別途CSRFトークン用Cookie（`XSRF-TOKEN` 等、`HttpOnly=false`）を発行し、クライアントがJSで読み取ってヘッダーに付与する（Double Submit Cookie パターン）
   - HTMLレンダリング時にメタタグで埋め込む

2. **CSRFトークンの有効期限**: セッション単位で固定か、リクエストごとにローテーションするか。

3. **未ログイン状態のPOSTリクエストとCSRF**: `POST auth/login`, `POST auth/register/request-otp` 等の認証不要POSTエンドポイントにもCSRF検証を適用するのか否か。要件定義書 L64 にて「ログアウト」にのみCSRFトークン検証必須が明記されているが、API設計書では `POST auth/change-email/*` 等の認証必須エンドポイントにも `X-CSRF-Token` ヘッダーが記載されている一方、`POST auth/login` には記載がない。この判断基準が不明確。

4. **エンドポイント一覧（セクション2）にCSRFトークン取得APIが未記載**: エンドポイント一覧テーブルにCSRFトークン取得のためのエントリが存在しない。

## 3. 推奨される修正案
以下のいずれかの方式でCSRFトークンの取得仕様を明記する:

**案A: 専用エンドポイント方式**
```
| **セキュリティ** | `GET` | `csrf-token` | CSRFトークンの取得 | 不要 |
```
レスポンスで `{ "csrf_token": "..." }` を返却し、クライアントはこの値を `X-CSRF-Token` ヘッダーに設定する。

**案B: Double Submit Cookie 方式**
セッション発行時（ログイン / 新規登録成功時）に `XSRF-TOKEN` Cookie（`HttpOnly=false`）を同時に設定し、クライアントはJavaScriptでこのCookieを読み取り `X-CSRF-Token` ヘッダーに設定する旨を共通仕様に明記する。

いずれの方式を採用するかを決定し、セクション1.2に取得方法とライフサイクルを明記し、エンドポイント一覧に取得用エントリ（必要な場合）を追加する。
