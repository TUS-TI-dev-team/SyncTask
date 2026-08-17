# アカウント削除およびパスワード変更における5回再認証失敗時（SESSION_DESTROYED）の Cookie 削除ヘッダー仕様の漏れ

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-17 16:45:00
- **Target Files**:
  - [03_users.md](docs/design/api_design/03_users.md)

## 1. 問題の概要
`DELETE users/{user_id}` (3.2.3) および `PATCH users/{user_id}/password` (3.2.4) のエラー仕様において、パスワード再認証が5回連続で失敗した場合、レスポンスは `401 Unauthorized` (code: `"SESSION_DESTROYED"`) となり、文章で「セッション強制破棄・Cookie消去」と規定されています。
しかし、正常系レスポンス（200 OK）とは異なり、`401 SESSION_DESTROYED` のエラーレスポンス仕様に Cookie を削除・無効化するためのレスポンスヘッダー（`Set-Cookie: sync_task_sid=; Max-Age=0` および `Set-Cookie: XSRF-TOKEN=; Max-Age=0`）の明記が漏れています。

## 2. 詳細な指摘内容
1. `03_users.md` L90-L91 および L133-L134 の正常系レスポンス（200 OK）では以下のように Cookie 削除ヘッダーが明確に定義されています。
   - `Set-Cookie: sync_task_sid=; Max-Age=0`
   - `Set-Cookie: XSRF-TOKEN=; Max-Age=0`
2. しかし、5回連続失敗時のエラーレスポンス `401 Unauthorized` (L105, L148) では、補足テキストで `(セッション強制破棄・Cookie消去、code: "SESSION_DESTROYED"、遅延 1.0s ± 0.1s)` と記載されているのみで、`##### Errors` 内のヘッダー定義として `Set-Cookie` が明示されていません。
3. サーバー側で DB 上のセッションレコードを物理削除しても、レスポンスヘッダーに `Set-Cookie` による Max-Age=0 指定が含まれていない場合、クライアント（ブラウザ）側に無効化された Cookie が残留し続けるため、フロントエンド実装や状態管理において不整合や混乱を生じさせる原因となります。

## 3. 推奨される修正案
`03_users.md` の 3.2.3 (`DELETE users/{user_id}`) および 3.2.4 (`PATCH users/{user_id}/password`) の `##### Errors` セクションにおける `401 Unauthorized` (`code: "SESSION_DESTROYED"`) の説明に、レスポンスヘッダーとして `Set-Cookie` を明記してください。

```markdown
- `401 Unauthorized`: 
  - 未ログイン・セッション無効（code: `"UNAUTHORIZED"`）
  - パスワード再認証失敗 1〜4回目（code: `"REAUTH_FAILED"`、遅延 1.0s ± 0.1s）
  - パスワード再認証失敗 5回連続達成分（セッション強制破棄・Cookie消去、code: `"SESSION_DESTROYED"`、遅延 1.0s ± 0.1s）
    - **Set-Cookie**: `sync_task_sid=; Max-Age=0`
    - **Set-Cookie**: `XSRF-TOKEN=; Max-Age=0`
```

---

## 修正完了報告

- **Resolved At**: 2026-08-17 16:50:00
- **Status**: Resolved

### 実施した修正内容
`03_users.md` 3.2.3 および 3.2.4 の `Errors` における `401 SESSION_DESTROYED` の定義に、レスポンスヘッダー `Set-Cookie: sync_task_sid=; Max-Age=0` および `Set-Cookie: XSRF-TOKEN=; Max-Age=0` を明記追加しました。

### 変更したファイル
- [03_users.md](docs/design/api_design/03_users.md)
