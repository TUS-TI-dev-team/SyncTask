# セッションCookieのセキュリティ属性要件およびCSRFトークン検証仕様の不足

- **Status**: Resolved
- **Severity**: Major
- **Created At**: 2026-08-13 22:15:00
- **Target Files**:
  - [requirements.md](docs/req-def/requirements.md)

## 1. 問題の概要
要件定義書（L197, L214-L215）において、「セッション管理（トークンではなくセッション）」および「CSRF対策 → CSRFトークンの導入」が記載されていますが、セッションIDを保持するCookieのセキュリティ属性（HttpOnly, Secure, SameSite）およびCSRFトークンの生成・配布・検証方式の具体的な仕様が不足しています。

## 2. 詳細な指摘内容
1. **セッションCookieのセキュリティ属性未定義 (L214-L215)**:
   - ログインセッションを Cookie ベースのセッション ID で管理する際、XSS によるセッション強奪（セッションハイジャック）を防ぐための `HttpOnly` 属性、中間者攻撃を防ぐための `Secure` 属性（HTTPS通信限定）、および クロスサイトリクエスト偽造を多重防護するための `SameSite` 属性（`Lax` または `Strict`）が規定されていません。
2. **CSRFトークン検証仕様の具体的記述の欠落 (L197)**:
   - 単に「CSRFトークンの導入」と書かれているのみで、トークンの発行タイミング（ログイン成功時やセッション確立時にCookie等で発行）、送信方式（Double-Submit CookieパターンやカスタムHTTPヘッダー `X-CSRF-Token` 経由での送信）、および検証対象のリクエスト（状態を変更する `POST`, `PUT`, `PATCH`, `DELETE` リクエスト）が明記されていません。
3. **認証前エンドポイント（ログイン・OTP検証）でのCSRF・オリジン検証**:
   - ログイン API や OTP 発行・検証 API などの認証前リクエストにおける CSRF/Origin ヘッダーチェックの要件が不明確です。

## 3. 推奨される修正案
`docs/req-def/requirements.md` の「セキュリティ -> セッション管理」および「脆弱性対策」要件を更新・追記してください：

```markdown
- セッションCookie要件
  - ログインセッションIDを保持するCookieには以下のセキュリティ属性を必須で付与する:
    - `HttpOnly`: JavaScript（`document.cookie`）からの読み取りを禁止しXSS対策を実施
    - `Secure`: HTTPS通信時のみCookieを送信（開発環境を除き常時有効）
    - `SameSite=Lax`: サードパーティコンテキストでの無駄なCookie送信を制限
  - セッションCookieの名称は標準的な名称（例: `sync_task_sid`）とし、有効期限はSliding Expirationに従い更新する
- CSRF対策仕様
  - ログイン後の状態変更リクエスト（`POST`, `PUT`, `PATCH`, `DELETE`）に対し、CSRFトークン検証を実施する
  - CSRFトークンはログイン成功時およびセッション作成時にサーバー側で暗号学的に安全なランダム文字列として発行し、レスポンスヘッダーまたは専用Cookie経由でクライアントに受け渡す
  - クライアント（フロントエンド）は状態変更リクエストの送信時にカスタムHTTPヘッダー（`X-CSRF-Token`）にトークンを付与して送信し、バックエンドはセッションに紐づくトークンと一致することを検証する
```

---

## 修正完了報告

- **Resolved At**: 2026-08-13 22:20:00
- **Status**: Resolved

### 実施した修正内容
要件定義書（`docs/req-def/requirements.md`）にログインセッションCookieのセキュリティ属性（`HttpOnly`、`Secure`、`SameSite=Lax`）を必須付与する旨を追加しました。また、CSRF対策についてはCSRFトークン導入方針を明記し、具体的な検証方式や伝送形式等の詳細仕様は詳細設計フェーズにて定義する旨を整理・反映しました。

### 変更したファイル
- [requirements.md](docs/req-def/requirements.md)
