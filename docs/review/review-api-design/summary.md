# レビュー結果サマリ

- **Status**: Passed
- **Reviewed At**: 2026-08-17
- **Target Files**:
  - [02_auth.md](file:///mnt/c/Users/ayumu_wkkd3w3/BoxForSomething/github/SyncTask/docs/design/api_design/02_auth.md)
  - [03_users.md](file:///mnt/c/Users/ayumu_wkkd3w3/BoxForSomething/github/SyncTask/docs/design/api_design/03_users.md)

本ブランチ (`review/api-design`) における認証 API 設計書 `docs/design/api_design/02_auth.md` およびユーザー管理 API 設計書 `docs/design/api_design/03_users.md` に対する詳細レビューを実施した結果、仕様上の不備・他ドキュメント（`01_overview.md`, `04_tasks.md`, `database_design.md`）との不整合・セキュリティ上の考慮漏れ・スキーマ定義の齟齬等は発見されず、すべての要件・整合性チェックを通過していることを確認しました。

## レビュー確認項目 (`02_auth.md`)

1. **アカウント登録フロー (`3.1.1` - `3.1.3`)**
   - **`request-otp`**: リクエスト評価順序（入力検証 400 `BAD_REQUEST` -> クールダウン 429 `OTP_RESEND_COOLDOWN` -> 遅延制御 200 OK）、アカウント列挙防止・Timing Attack 対策としての固定レスポンス遅延（1.0s ± 0.1s）およびダミーセッション発行仕様の整合性
   - **`verify-otp`**: リクエスト評価順序（入力検証 400 `BAD_REQUEST` -> セッション状態/期限検証 400/410 `GONE` -> OTP照合 400/422 `OTP_REISSUED_DUE_TO_FAILURES`）、本登録成功時の自動ログイン・セッションCookie（`sync_task_sid`）および CSRF トークンCookie（`XSRF-TOKEN`）の発行、既存セッション物理削除、レスポンススキーマ（`user` オブジェクト完全性）の完全性
   - **`resend-otp`**: セッション状態・最大有効期限（15分）検証、クールダウン検証、試行失敗回数リセット・有効期限延長・Timing Attack 対策遅延の仕様整合性

2. **ログイン・ログアウトフロー (`3.1.4` - `3.1.5`)**
   - **`login`**: リクエスト評価順序（入力検証 400 `BAD_REQUEST` -> IPレートリミット/アカウントロック 429 `RATE_LIMIT_EXCEEDED` -> 認証照合 401 `UNAUTHORIZED`）、旧セッション物理削除、Cookie 発行およびレスポンス `user` スキーマの整合性
   - **`logout`**: 冪等性確保（未ログイン/期限切れ時も 200 OK）、有効セッション存在時の CSRF 検証（403 `FORBIDDEN`）、Cookie 消去ヘッダー（`Max-Age=0`）の完全性

3. **パスワードリセットフロー (`3.1.6` - `3.1.9`)**
   - **`request-otp` / `verify-otp` / `resend-otp`**: アカウント列挙防止遅延、検証成功時の `verified` ステータス遷移および15分間延長仕様の明示性
   - **`reset`**: リクエスト評価順序（入力検証 400 `BAD_REQUEST` -> `verified` セッション状態/期限検証 403 `FORBIDDEN`/410 `GONE` -> 検証済み `OTP_SESSION` 経由ユーザー属性参照によるユーザー名・メールローカル部含有および現パスワード同一性検証 422 `SAME_AS_CURRENT_PASSWORD`/`INVALID_PASSWORD_CONTENT` -> 更新・全セッション破棄・Cookie消去 200 OK）の論理的一貫性

4. **メールアドレス変更フロー (`3.1.10` - `3.1.12`)**
   - **`request-otp`**: 認証・CSRF検証（401/403）、同一メールアドレス指定検証（422 `SAME_AS_CURRENT_EMAIL`）、Timing Attack 対策遅延（200 OK）
   - **`verify-otp` / `resend-otp`**: 認証・CSRF検証、他ユーザー所有 `otp_session_id` 指定の認可チェック（403 `FORBIDDEN`）、変更確定時の全ログインセッション破棄・Cookie消去・完了通知メール送信仕様の完全性

## レビュー確認項目 (`03_users.md`)

1. **`GET users/{user_id}` (3.2.1)**
   - 認証検証 (401 `UNAUTHORIZED`) および IDOR/BOLA 秘匿検証 (404 `NOT_FOUND`) の評価順序の妥当性
   - レスポンス JSON スキーマ (`user.id`, `user.username`, `user.email`, `user.created_at`, `user.updated_at`) の完全性

2. **`PUT users/{user_id}` (3.2.2)**
   - CSRF トークン検証 (`X-CSRF-Token`)、リクエストボディ (`username`) のトリム・文字数・使用可能文字バリデーション (400 `BAD_REQUEST`)、IDOR 検証 (404 `NOT_FOUND`)、同一ユーザー名チェック (422 `SAME_AS_CURRENT_USERNAME`) の評価順序の明示性
   - 既存エラー定義との整合性

3. **`DELETE users/{user_id}` (3.2.3)**
   - アカウント論理削除処理 (`IS_DELETED=true`, `DELETED_AT=NOW()`) およびメールアドレス退避フォーマット (`deleted_<USER_ID>_<EMAIL>`) の設計整合性
   - 所有タスクデータ・全ログインセッション (`LOGIN_SESSION`)・アクティブ OTP セッション (`OTP_SESSION`) の即時物理削除仕様
   - パスワード再認証検証の評価順序、失敗時の 1.0s ± 0.1s 固定遅延、5回連続失敗時のセッション強制破棄 (`SESSION_DESTROYED`) および Cookie 消去ヘッダー (`Set-Cookie: sync_task_sid=; Path=/; Max-Age=0`, `Set-Cookie: XSRF-TOKEN=; Path=/; Max-Age=0`) の整合性

4. **`PATCH users/{user_id}/password` (3.2.4)**
   - 新パスワード制約（8〜128文字、英大小数字記号3種以上、ユーザー名・メールローカル部（4文字以上）含有不可（Case-Insensitive））のバリデーション仕様 (400 `BAD_REQUEST`)
   - パスワード再認証失敗時の遅延制御・5回連続失敗時のセッション破棄・成功時の `REAUTH_FAILED_COUNT` リセット
   - 同一パスワード指定チェック (422 `SAME_AS_CURRENT_PASSWORD`) の評価順序の正当性（再認証成功後のチェック）
   - パスワード変更成功時の全ログインセッション・OTP セッション一括物理削除および Cookie 削除ヘッダーの完全性
