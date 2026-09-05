# 信頼済みプロキシ未設定時における X-Forwarded-For 経由の IP レート制限バイパス

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-09-05 20:42:00
- **Target Files**:
  - [router.go](backend/router/router.go)
  - [config.go](backend/config/config.go)
  - [.env.example](backend/.env.example)
  - [auth-login.md](docs/plans/backend/auth-login.md)

## 1. 問題の概要
`TRUSTED_PROXIES` 環境変数が未設定（または空文字）の状態でサーバーが起動された場合、`r.SetTrustedProxies(nil)` が呼び出され、Gin の仕様によりすべてのプロキシヘッダーが信頼される設定となります。これにより、直接アクセス環境において攻撃者が偽装した `X-Forwarded-For` ヘッダーを送信するだけで任意の IP アドレスを名乗ることが可能となり、IP 単位のレート制限（30 回失敗で 15 分遮断）をクライアント側から容易に回避・バイパスできるセキュリティリスクが存在します。

## 2. 詳細な指摘内容
1. `backend/config/config.go` の `getCSVEnv("TRUSTED_PROXIES")` は、環境変数が未設定または空文字の際に `nil` を返します。
2. `backend/router/router.go` (L32-34) では以下のようにそのまま Gin に渡しています：
   ```go
   if err := r.SetTrustedProxies(options.TrustedProxies); err != nil {
       panic("invalid trusted proxy configuration: " + err.Error())
   }
   ```
3. Gin では `SetTrustedProxies(nil)` が渡されると、内部的に警告を出力しつつすべての転送元を信頼する挙動（すべてのヘッダーを受け入れる）となります。
4. `backend/.env.example` の説明（L10-11）には「leave empty for direct access（直接アクセスの場合は空にせよ）」と記載されていますが、空にすると「直接アクセスだからヘッダーを無視する」のではなく、真逆に「任意の `X-Forwarded-For` ヘッダーをクライアント IP として採用する」状態になります。
5. `backend/handler/login.go` の `LoginHandler` では `c.ClientIP()` をそのままレートリミット用の識別 IP として利用しているため、攻撃者がリクエスト毎に異なる IP を `X-Forwarded-For` に付与することで、同一端末から無制限にログイン試行（ブルートフォース攻撃）を実行可能になります。
6. `docs/plans/backend/auth-login.md` の初期テストケース候補（L95:「準正常系: 信頼済みプロキシ設定時のみ転送元IPを採用し、任意の転送ヘッダーを無条件に信用しない」）が実装およびテストコードで検証されていません。

## 3. 推奨される修正案
1. `TRUSTED_PROXIES` が未設定または空の場合は、転送ヘッダーを無条件に信用しないよう、明示的に空のスライス `[]string{}` を渡す、あるいは直接アクセス時に安全なデフォルト設定を適用してください。
   ```go
   // router.go または config.go での対応例
   proxies := options.TrustedProxies
   if len(proxies) == 0 {
       proxies = []string{} // どのプロキシも信頼しない（RemoteAddr を優先採用）
   }
   if err := r.SetTrustedProxies(proxies); err != nil {
       panic("invalid trusted proxy configuration: " + err.Error())
   }
   ```
2. `backend/router/router_test.go` または `backend/handler/login_test.go` において、未信頼のプロキシから送られた `X-Forwarded-For` ヘッダーが無視され、`RemoteAddr` がクライアント IP として採用されることを確認する単体テストを追加してください。
3. `backend/.env.example` のコメントを実際の挙動に合わせて更新してください。

---

## 修正完了報告

- **Resolved At**: 2026-09-05 21:16:00
- **Status**: Resolved

### 実施した修正内容
1. `backend/router/router.go` において、`options.TrustedProxies` が空または `nil` の場合に空スライス `[]string{}` を渡すよう正規化しました。これにより、Gin の `SetTrustedProxies(nil)` による全プロキシ信頼化（偽装ヘッダー採用）を防ぎ、直接アクセス環境において全プロキシを不信頼（`RemoteAddr` のみ採用）とする安全なデフォルト挙動を実現しました。
2. `backend/.env.example` の `TRUSTED_PROXIES` のコメントを実際の挙動に合わせ、空設定時は全プロキシ不信頼（安全な直接アクセス）となる旨を明記しました。
3. `backend/router/router_test.go` に `TestSetupRouter_TrustedProxies_DirectAccess` および `TestSetupRouter_TrustedProxies_WithTrustedProxy` を新設し、以下を検証しました：
   - デフォルト（未設定）時：偽装 `X-Forwarded-For` ヘッダーが付与されていても無視され、`RemoteAddr` がクライアント IP として評価されること。
   - 信頼済みプロキシ設定時：信頼済みプロキシ経由のリクエストでは `X-Forwarded-For` が採用され、未信頼プロキシ経由では無視されて `RemoteAddr` が採用されること。

### 変更したファイル
- [router.go](backend/router/router.go)
- [router_test.go](backend/router/router_test.go)
- [.env.example](backend/.env.example)

