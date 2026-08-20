# セッション伝送方式の曖昧さとCSRF対策対象エンドポイントの不整合

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-17 12:05:00
- **Target Files**:
  - [api_design.md](docs/design/api_design.md)
  - [requirements.md](docs/req-def/requirements.md)

## 1. 問題の概要
要件定義書（`requirements.md`）では HttpOnly/Secure Cookie（`sync_task_sid`）によるセッション管理と状態変更系リクエストに対するCSRF対策が定められていますが、`api_design.md` では各エンドポイントの入力欄に「セッションID」が直接記載されており、Cookie自動送信なのかリクエストボディ送信なのかが曖昧です。また、CSRF対策の適用対象が一部の `PUT` API のみに限定して記載されており不整合です。

## 2. 詳細な指摘内容
1. **セッション伝送方式の不統一**:
   - `docs/req-def/requirements.md` L268-273 では `HttpOnly`, `Secure`, `SameSite=Lax` Cookie によるセッション管理と明記されています。
   - しかし `docs/design/api_design.md` では、`auth/logout` のみ「入力: Cookie (session_id)」と書かれ、他の保護されたエンドポイント（`users/{user_id}`, `tasks/`, `tasks/{task_id}` 等）では「入力: セッションID、...」と記載されており、リクエストパラメータ/ヘッダーとして明示送信するのか Cookie で送信するのかが混同されています。
2. **CSRF対策の適用範囲の不整合**:
   - `docs/design/api_design.md` では `users/{user_id}` (PUT) と `tasks/{task_id}` (PUT) の2箇所にのみ「送り付け対策(CSRF)必須」と記載されています。
   - しかし、Cookie ベースの認証を行う場合、状態変更を伴うすべてのリクエスト（`POST tasks/`, `DELETE tasks/{task_id}`, `DELETE users/{user_id}`, `PATCH users/{user_id}/password`, `POST auth/logout` 等）において CSRF 対策（CSRFトークンのヘッダー送信等）が必須です（`requirements.md` L64 にもログアウト時 CSRF 検証必須と記載）。

## 3. 推奨される修正案
1. 共通仕様セクションに「ログイン状態が必要なAPIはすべて HTTP Cookie（`sync_task_sid`）を介してセッションを認証する」旨を明記し、個別の入力欄から「セッションID」の個別記述を整理してください。
2. CSRF対策のポリシー（例: 状態変更を伴う POST / PUT / PATCH / DELETE メソッドでは `X-CSRF-Token` ヘッダーによる検証を必須とする等）を共通仕様として明記し、全該当エンドポイントに統一して適用してください。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 12:40:00
- **Status**: Resolved

### 実施した修正内容
- `docs/design/api_design.md` の第1章「概要・共通仕様」において、Cookieベースのセッション管理方式（Cookie名: `sync_task_sid`）および、すべての状態変更系メソッド（POST, PUT, PATCH, DELETE）に対する `X-CSRF-Token` リクエストヘッダー検証ポリシーを明記しました。
- 各エンドポイント定義から冗長・曖昧な「セッションID」の個別表記を排除し、認証要否およびCSRFヘッダー要否を統一的に定義しました。

### 変更したファイル
- [api_design.md](docs/design/api_design.md)
