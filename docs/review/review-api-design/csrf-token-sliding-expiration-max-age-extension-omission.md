# Sliding Expiration 適用時における CSRFトークンCookie (XSRF-TOKEN) Max-Age 延長仕様の未記述

- **Status**: Resolved
- **Severity**: Major
- **Created At**: 2026-08-17 16:45:00
- **Target Files**:
  - [01_overview.md](docs/design/api_design/01_overview.md)

## 1. 問題の概要
`01_overview.md` 1.1 節においてログインセッションCookie（`sync_task_sid`）は有効期限 30 日間であり「APIアクセスごとに自動延長（Sliding Expiration）されます」と規定されている。一方、1.2 節の CSRF トークンCookie（`XSRF-TOKEN`）の発行仕様ではログイン成功時およびアカウント新規登録完了時に `Set-Cookie: XSRF-TOKEN=<csrf_token>; Secure; SameSite=Lax; Path=/; Max-Age=2592000` を付与する旨が記載されているものの、Sliding Expiration によるログインセッション自動延長時にも `XSRF-TOKEN` Cookie の有効期限（Max-Age）が同様に更新・延長されるかどうかの規定が存在しない。

## 2. 詳細な指摘内容
1. **ログインセッションと CSRF トークンCookie の有効期限ズレリスク**:
   - `sync_task_sid` Cookie は API アクセスごとに Sliding Expiration により有効期限が継続的に30日延長される。
   - しかし、`XSRF-TOKEN` Cookie がログイン時/会員登録時の単一発行（Max-Age=2592000秒）にとどまり、Sliding Expiration のタイミングで `Set-Cookie` による有効期限更新が行われない場合、初回ログインから30日経過した時点で `XSRF-TOKEN` Cookie のみがブラウザ上で有効期限切れとなり消去される。

2. **アクティブユーザーの操作不能障害（403 Forbidden）**:
   ユーザーが30日間以上にわたり毎日アクティブにサービスを利用し続けていても、`sync_task_sid` は維持される一方で `XSRF-TOKEN` Cookie が消失する。これにより、状態変更リクエスト（`POST`, `PUT`, `PATCH`, `DELETE`）送信時に `X-CSRF-Token` ヘッダーを付与できなくなり、突如すべての更新操作が `403 Forbidden` エラーで拒否されるデッドロック状態が発生する。

## 3. 推奨される修正案
`01_overview.md` の 1.1 節および 1.2 節に、セッション自動延長（Sliding Expiration）実行時における `XSRF-TOKEN` Cookie の有効期限延長ルールを明記してください。

```markdown
- **Sliding Expiration 時の CSRF トークンCookie 延長**:
  - APIアクセスに伴うログインセッション（`sync_task_sid`）の自動延長（Sliding Expiration）時には、レスポンスヘッダーにおいて `XSRF-TOKEN` Cookie も同様に `Set-Cookie: XSRF-TOKEN=<csrf_token>; Secure; SameSite=Lax; Path=/; Max-Age=2592000` を出力し、有効期限を最新のセッション有効期限と同期して更新延長します。
```

---

## 修正完了報告

- **Resolved At**: 2026-08-17 16:50:00
- **Status**: Resolved

### 実施した修正内容
`01_overview.md` 1.1 節に「Sliding Expiration 時の CSRF トークンCookie 延長」の規定を追加し、ログインセッション自動延長時に `Set-Cookie: XSRF-TOKEN=<csrf_token>; ...; Max-Age=2592000` を出力して `XSRF-TOKEN` の有効期限をセッションと同期延長する旨を明記しました。

### 変更したファイル
- [01_overview.md](docs/design/api_design/01_overview.md)
