# `auth/resend-otp` および `auth/login` エラー仕様における `400 Bad Request` の定義欠落

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-17 16:45:00
- **Target Files**:
  - [02_auth.md](docs/design/api_design/02_auth.md)

## 1. 問題の概要
`02_auth.md` の以下の4つのエンドポイントのエラー仕様（`Errors`）において、リクエストパラメータ不足やJSONフォーマット不正時等に返却されるべき `400 Bad Request`（`code: "BAD_REQUEST"`）の定義が完全に欠落している。

1. 3.1.3 `POST auth/register/resend-otp` (L104-L107)
2. 3.1.4 `POST auth/login` (L142-L145)
3. 3.1.8 `POST auth/password-reset/resend-otp` (L253-L255)
4. 3.1.12 `POST auth/change-email/resend-otp` (L396-L400)

## 2. 詳細な指摘内容
- **3.1.4 `POST auth/login`**: `email` や `password` フィールドが未指定・空文字、または JSON パースエラー等の場合に返される `400 Bad Request` が記載されておらず、`401` と `429` のみが定義されている。
- **3.1.3 / 3.1.8 / 3.1.12 `resend-otp` 各種**: リクエストボディ内の `otp_session_id` が未指定・型不正等の場合に返される `400 Bad Request` が記載されていない。
- **共通設計方針との不整合**: `01_overview.md` 1.3「代表的なエラーコード一覧」において `400` / `BAD_REQUEST` は「リクエスト形式またはバリデーション不正」と定められており、パラメータ不足時のレスポンスコードが設計書上で未定義となっている。

## 3. 推奨される修正案
対象の4つのエンドポイントの `Errors` セクションに、それぞれ以下の項目を追記してください：

```markdown
- `400 Bad Request`: リクエストボディ不正・必須パラメータ欠落（code: `"BAD_REQUEST"`）
```

---

## 修正完了報告

- **Resolved At**: 2026-08-17 16:50:00
- **Status**: Resolved

### 実施した修正内容
`02_auth.md` の 3.1.3 (`register/resend-otp`), 3.1.4 (`login`), 3.1.8 (`password-reset/resend-otp`), 3.1.12 (`change-email/resend-otp`) の 4 つのエンドポイントの `Errors` セクションに `400 Bad Request`（`code: "BAD_REQUEST"`）を追記しました。

### 変更したファイル
- [02_auth.md](docs/design/api_design/02_auth.md)
