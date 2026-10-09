# 認証ミドルウェア自作ステップ(学習用ドキュメント)

このドキュメントは、**自分の手**で SyncTask バックエンドにセッション認証ミドルウェアを実装するための学習ロードマップです。
答えを書くのではなく、「各ステップで何をやるか」「何を調べるか」「詰まったらどこを見るか」を示します。コード例は最小限の骨組みのみです。

各ステップの詳細は以下のファイルを参照してください。

| Step | 内容 | ファイル |
| --- | --- | --- |
| Step 1 | Ginミドルウェアの仕組み理解(ロガー自作) | [step-1.md](step-1.md) |
| Step 2 | Cookieを読む | [step-2.md](step-2.md) |
| Step 3 | セッション検証の層を作る | [step-3.md](step-3.md) |
| Step 4 | 検証結果を後続に渡す | [step-4.md](step-4.md) |
| Step 5 | 401応答の設計とルート保護 | [step-5.md](step-5.md) |
| Step 6 | テストを書く | [step-6.md](step-6.md) |
| Step 7 | まとめと次の一歩 | [step-7.md](step-7.md) |

## 前提知識

| 項目 | 内容 |
| --- | --- |
| 技術レベル想定 | Go は A Tour of Go 完了レベル。Web開発は Hono で2ヶ月程度の経験(記憶は薄れ気味) |
| 既存実装(活用するもの) | `backend/handler/login.go` がログイン成功時に **`sync_task_sid`**(HttpOnly セッションCookie)と **`XSRF-TOKEN`**(CSRF Cookie)を発行済み |
| 既存テーブル | `LOGIN_SESSION`(`SESSION_ID`, `USER_ID`, `EXPIRES_AT`, `USER_AGENT`, `CREATED_AT`)。`backend/db/migrations/000001_init.up.sql` 参照 |
| このドキュメントのスコープ | **セッション認証のみ**(CSRF検証は次の学習テーマとして切り出し) |

## 全体工程図

```
[Step 1] Ginミドルウェアの仕組み理解(ロガー自作)
    ↓
[Step 2] Cookieを読む(c.Cookie)
    ↓
[Step 3] セッション検証の層を作る(repository / service)
    ↓
[Step 4] 検証結果を後続に渡す(c.Set / c.Get)
    ↓
[Step 5] 401応答の設計とルート保護(AppError + tasksグループ適用)
    ↓
[Step 6] テストを書く(httptest + fake service)
    ↓
[Step 7] まとめと次の一歩
```

各ステップの目安は **1日** です。Step 2〜5 は順番に依存しますが、Step 6 は Step 3 以降の各ステップで少しずつ先取りして構いません。

## 実装方針(先に決めておくこと)

以下は本ドキュメント策定時に決定済みの方針です。迷ったらこれに従ってください。

- **配置**: 新規パッケージ `backend/middleware/`(`middleware/auth.go`)
- **レイヤー経路**: `repository`(新規 LoginSessionRepository)→ `service`(新規 SessionService)→ `middleware`。既存の `router.go` のDIパターン(`router.go` の 57〜65 行目付近)に倣う
- **保護対象**: `POST/GET/PATCH /api/tasks` のみ。`/api/auth/*` と `/health-check` は公開のまま
- **401応答**: 既存の `model.AppError` / `model.ErrorResponse` パターンで統一(code: `UNAUTHORIZED`)
- **テスト**: `httptest` + service の fake(interface 差し替え)。DBや testcontainers は使わない

## 進め方のヒント

- 各ステップの「調べるキーワード」で公式ドキュメント(pkg.go.dev / gin-gonic.com / go.dev)を**必ず一次情報として**確認する。Qiitaやブログ記事はあくまで補助
- Hono の経験と対応づけると理解が早い: `app.use()` ≒ `r.Use()`、`c.set()` ≒ `c.Set()`、`c.env` ≒ 依存性注入(関数引数)
- 詰まったら、まず既存コード(`service/login.go` や `router_test.go`)に同じパターンがないか探す。このリポジトリはパターンの見本が揃っている
- 各ステップの末尾に「詰まったら(トグルして開く)」セクションがある。**すぐ開かず、まず自分で10分は調べてみてから**段階的に開くこと

## 参考リンクの検証

外部リンクは **2026-10-09 時点**で全件の到達(HTTP 200)を確認済みです。Gin公式ドキュメントは 2026年に URL 構造が変更され、日本語版は `gin-gonic.com/ja/docs/...` 配下になっています。
