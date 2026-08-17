# アカウント削除およびパスワード変更レスポンスにおける CSRF トークン Cookie (`XSRF-TOKEN`) 削除指定の欠落

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-17 16:35:00
- **Target Files**:
  - [03_users.md](docs/design/api_design/03_users.md)
  - [02_auth.md](docs/design/api_design/02_auth.md)

## 1. 問題の概要
`01_overview.md`（1.2 節）および `02_auth.md`（3.1.5 節 `auth/logout`）では、ログインセッションの破棄に伴い、セッション Cookie（`sync_task_sid`）と CSRF トークン Cookie（`XSRF-TOKEN`）の両方を `Set-Cookie: <name>=; Max-Age=0` で無効化・破棄することが定義されている。

しかし、`03_users.md` における `DELETE users/{user_id}` (3.2.3) および `PATCH users/{user_id}/password` (3.2.4) の成功レスポンス（`200 OK`）ヘッダー指定において、`sync_task_sid=; Max-Age=0` のみが記載されており、`XSRF-TOKEN=; Max-Age=0` の削除指示が欠落している。

## 2. 詳細な指摘内容
- **`DELETE users/{user_id}` (3.2.3 L81)**:
  `- **Set-Cookie**: sync_task_sid=; Max-Age=0`
- **`PATCH users/{user_id}/password` (3.2.4 L112)**:
  `- **Set-Cookie**: sync_task_sid=; Max-Age=0`（再ログイン要求）

アカウントが削除された場合、およびパスワード変更により全セッションが即座に無効化・物理削除された場合、認証セッション Cookie だけでなく、クライアント（ブラウザ）に保持されている CSRF トークン Cookie (`XSRF-TOKEN`) も同時にクリアされる必要がある。
これが欠落していると、ログアウト処理 (`auth/logout`) と Cookie 無効化挙動に不整合が生じ、無効化されたセッションに対応する無効な `XSRF-TOKEN` がブラウザ上に残留する不具合やセキュリティ状態の不一致の原因となる。

## 3. 推奨される修正案
`03_users.md` の 3.2.3 節 (`DELETE users/{user_id}`) および 3.2.4 節 (`PATCH users/{user_id}/password`) の `Response (200 OK)` セクションに `XSRF-TOKEN` Cookie の削除指定を追記してください。

```markdown
##### Response (200 OK)
- **Set-Cookie**: `sync_task_sid=; Max-Age=0`
- **Set-Cookie**: `XSRF-TOKEN=; Max-Age=0`
```

---

## 修正完了報告

- **Resolved At**: 2026-08-17 16:42:00
- **Status**: Resolved

### 実施した修正内容
`DELETE users/{user_id}` (3.2.3) および `PATCH users/{user_id}/password` (3.2.4) の Response 200 OK ヘッダーに `Set-Cookie: XSRF-TOKEN=; Max-Age=0` を追加し、セッション破棄時の CSRF トークン Cookie 消去指定の漏れを補正しました。

### 変更したファイル
- [03_users.md](file:////mnt/c/Users/ayumu_wkkd3w3/BoxForSomething/github/SyncTask/docs/design/api_design/03_users.md)
