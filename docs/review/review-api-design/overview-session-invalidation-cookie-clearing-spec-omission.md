# 概要書におけるセッション破棄・Cookie消去共通仕様の記載漏れ

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-17 17:33:00
- **Target Files**:
  - [01_overview.md](docs/design/api_design/01_overview.md)

## 1. 問題の概要
`01_overview.md` の「1.1 セッション管理 & 認証方式」および「1.2 セキュリティ & CSRF・アカウント列挙対策」において、ログインセッションCookie（`sync_task_sid`）およびCSRFトークンCookie（`XSRF-TOKEN`）の発行・自動延長（Sliding Expiration）に関する仕様は記載されているが、ログアウト・アカウント削除・メールアドレス変更完了・パスワード変更および再認証連続失敗時の「セッション破棄に伴うCookie消去レスポンスヘッダー仕様」（`Max-Age=0` による無効化）の共通定義が漏れている。

## 2. 詳細な指摘内容
`02_auth.md`（ログアウト `auth/logout`）、`03_users.md`（アカウント削除 `DELETE users/{user_id}`、パスワード変更 `PATCH users/{user_id}/password`）、`02_auth.md`（メールアドレス変更 `POST auth/change-email/verify-otp`）、および再認証5回連続失敗時（`SESSION_DESTROYED`）の仕様では、サーバーがログインセッションおよびCSRFトークンを無効化する際に、以下のレスポンスヘッダーを出力することが規定されています。

- `Set-Cookie: sync_task_sid=; Path=/; Max-Age=0`
- `Set-Cookie: XSRF-TOKEN=; Path=/; Max-Age=0`

しかし、全体共通仕様を司る `01_overview.md` の 1.1 セッション管理節には、セッション発行時の属性（`Max-Age=2592000`）と自動延長時の属性（Sliding Expiration）のみが記述されており、セッション破棄時に両Cookieをどのように削除・無効化するかの共通ルールが明記されていません。

## 3. 推奨される修正案
`01_overview.md` の「1.1 セッション管理 & 認証方式」に、「セッション破棄・Cookie無効化」に関する段落または小節を追加し、以下のように明記してください。

```markdown
- **セッション破棄・Cookie消去仕様**:
  - ログアウト（`auth/logout`）、アカウント削除（`users/{user_id}`）、メールアドレス変更完了（`auth/change-email/verify-otp`）、パスワード変更（`users/{user_id}/password`）、および再認証連続失敗による強制破棄（`SESSION_DESTROYED`）の発生時は、サーバー側で DB 上のセッションレコードを物理削除すると同時に、レスポンスヘッダーで以下の Cookie 削除ヘッダーを出力してクライアント側の Cookie を直ちに無効化・消去します。
    - `Set-Cookie: sync_task_sid=; Path=/; Max-Age=0`
    - `Set-Cookie: XSRF-TOKEN=; Path=/; Max-Age=0`
```

---

## 修正完了報告

- **Resolved At**: 2026-08-17 17:38:45
- **Status**: Resolved

### 実施した修正内容
`docs/design/api_design/01_overview.md` の「1.1 セッション管理 & 認証方式」に、「セッション破棄・Cookie消去仕様」を追記し、ログアウト・アカウント削除・メールアドレス変更完了・パスワード変更・再認証連続失敗による強制破棄時の Cookie 無効化ヘッダー（`Max-Age=0`）の共通定義を明記しました。

### 変更したファイル
- [01_overview.md](docs/design/api_design/01_overview.md)

