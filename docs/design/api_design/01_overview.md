# API Design (API設計書)

## 1. 概要・共通仕様

- **ベースURL**: `https://<domain>/api/`
- **通信形式**: JSON (HTTP REST API)
- **文字コード**: UTF-8
- **共通レスポンスヘッダー**:
  - `Content-Type: application/json; charset=utf-8`
  - **キャッシュ制御**: 個人情報・認証情報・タスクデータの漏洩を防ぐため、全API応答に `Cache-Control: no-store, no-cache, must-revalidate` および `Pragma: no-cache` を付与します。
- **タイムゾーン / 日時フォーマット**:
  - 日時文字列: ISO 8601 拡張形式 / 日本標準時（例: `2026-08-17T12:00:00+09:00` または `YYYY-MM-DDTHH:mm:ss+09:00`）。※リクエスト時にタイムゾーンオフセットを含まない ISO 8601 文字列が指定された場合はデフォルトで JST (`+09:00`) の日時として解釈し、UTC (`Z`) や他タイムゾーンの日時文字列が指定された場合は JST (`+09:00`) に変換・正規化して登録・返却します。
  - 日付文字列: `YYYY-MM-DD`（カレンダー・締切日絞り込み時など）

### 1.1 セッション管理 & 認証方式
- **ログインセッション**:
  - トークンをリクエスト本文で送受信するのではなく、`HttpOnly`, `Secure`, `SameSite=Lax`, `Path=/`, `Max-Age=2592000` 属性が付与されたセッションCookie（名称: `sync_task_sid`、例: `Set-Cookie: sync_task_sid=<session_id>; HttpOnly; Secure; SameSite=Lax; Path=/; Max-Age=2592000`）によって管理します。
  - 認証が必要なAPIリクエストでは、ブラウザにより自動送信される Cookie からセッションを検証します。
  - セッション有効期限は 43,200分（1ヶ月）であり、APIアクセスごとに自動延長（Sliding Expiration）されます。
  - **Sliding Expiration 時の CSRF トークンCookie 延長**:
    APIアクセスに伴うログインセッション（`sync_task_sid`）の自動延長（Sliding Expiration）時には、レスポンスヘッダーにおいて `XSRF-TOKEN` Cookie も同様に `Set-Cookie: XSRF-TOKEN=<csrf_token>; Secure; SameSite=Lax; Path=/; Max-Age=2592000` を出力し、有効期限を最新のセッション有効期限と同期して更新延長します。
- **セッション破棄・Cookie消去仕様**:
  - ログアウト（`auth/logout`）、アカウント削除（`users/{user_id}`）、メールアドレス変更完了（`auth/change-email/verify-otp`）、パスワード変更（`users/{user_id}/password`）、および再認証連続失敗による強制破棄（`SESSION_DESTROYED`）の発生時は、サーバー側で DB 上のセッションレコードを物理削除すると同時に、レスポンスヘッダーで以下の Cookie 削除ヘッダーを出力してクライアント側の Cookie を直ちに無効化・消去します。
    - `Set-Cookie: sync_task_sid=; HttpOnly; Secure; SameSite=Lax; Path=/; Max-Age=0`
    - `Set-Cookie: XSRF-TOKEN=; Secure; SameSite=Lax; Path=/; Max-Age=0`
- **OTPセッション**:
  - アカウント新規作成、パスワードリセット、メールアドレス変更の手続き中は、手続きごとの `otp_session_id` をリクエストボディおよびレスポンスボディで送受信します。
  - OTP有効期限は発行から5分（手続き全体の最大有効期限は15分）です。

### 1.2 セキュリティ & CSRF・アカウント列挙対策
- **CSRF対策**:
  - Cookieベースの認証を行うため、**認証を必要とするすべての状態変更リクエスト（`POST`, `PUT`, `PATCH`, `DELETE`）**において CSRFトークン（`X-CSRF-Token` ヘッダー）の検証を必須とします（未認証のログイン・会員登録・パスワードリセット等のリクエストを除く）。
  - CSRFトークンは **Double Submit Cookie 方式** にて管理します。ログイン成功（`auth/login`）およびアカウント新規登録完了（`auth/register/verify-otp`）時に、レスポンスヘッダーで `Set-Cookie: XSRF-TOKEN=<csrf_token>; Secure; SameSite=Lax; Path=/; Max-Age=2592000`（JavaScriptから読み取り可能な `HttpOnly=false`）を発行します。
  - クライアントは JavaScript で Cookie から CSRF トークンを取得し、状態変更リクエストの `X-CSRF-Token` ヘッダーに付与して送信します。
- **アカウント列挙防止 (User Enumeration 対策)**:
  - 新規登録（`auth/register/request-otp`）、パスワードリセット（`auth/password-reset/request-otp`）、メールアドレス変更（`auth/change-email/request-otp`）および各種 OTP 再送（`resend-otp`）において、指定されたメールアドレスの登録有無、他ユーザーとの重複、または**他ユーザーの有効なOTPセッション期間中（手続き中）**の指定にかかわらず、**正常成功時・ダミー発行時を一貫して区別せずレスポンス遅延（1.0s ± 0.1s）を適用した上で `200 OK` を返却**します。これにより、エラーコード・レスポンス構造・応答時間のあらゆる差異からメールアドレスの登録状況が推測されることを完全に防止します。
  - ダミーOTPセッションに対する後続の検証（`verify-otp`）や再送（`resend-otp`）に対しても、実セッションと全く同一のエラーコード（400, 410, 422, 429）および応答遅延（1.0s ± 0.1s）を適用します。
  - ログイン失敗時は、メールアドレス不一致・パスワード不一致・論理削除済みアカウントのいずれも一律で `401 Unauthorized`（code: `UNAUTHORIZED`、遅延 1.0s ± 0.1s）を返却します。
- **認可制御 (IDOR / BOLA 対策)**:
  - ユーザー情報（`users/{user_id}`）およびタスク情報（`tasks/{task_id}`）へのアクセス・変更・削除時は、セッション内のログインユーザーIDとリソースの所有ユーザーIDの一致を厳格に検証します。
  - 他ユーザー所有のリソースまたは存在しないリソースへのアクセスに対しては、リソースの存在有無を秘匿するため一律 `404 Not Found` を返却します。
- **遅延制御 (Timing Attack 対策)**:
  - ログイン失敗（連続ログイン失敗によるアカウントロックアウト・IPレートリミット超過時 429 を含む）、OTP発行・再送処理（正常成功時およびアカウント存在有無のダミー処理時を含む一括）、OTP検証失敗、パスワード再認証失敗時は、一律 `1.0s ± 0.1s` のレスポンス遅延を適用します。

### 1.3 共通エラーレスポンス構造
すべてのエラー応答は以下のJSONフォーマットで返却されます（HTTP ステータスコード: `4xx` または `5xx`）。

```json
{
  "error": {
    "code": "BAD_REQUEST",
    "message": "入力内容に不備があります。",
    "details": [
      {
        "field": "email",
        "message": "メールアドレスの形式が正しくありません。"
      }
    ]
  }
}
```

#### エラーレスポンス スキーマ定義

| フィールド | 型 | 必須 | 説明 |
| :--- | :--- | :---: | :--- |
| `error.code` | string | ○ | エラーコード（例: `BAD_REQUEST`, `UNAUTHORIZED`, `SAME_AS_CURRENT_PASSWORD` 等） |
| `error.message` | string | ○ | ユーザー向け汎用エラーメッセージ |
| `error.details` | array | ○ | フィールド単位のバリデーション詳細情報リスト。対象フィールドが存在しないエラー応答の場合は空配列 `[]` を返却（`null` やキー省略は不可） |
| `error.details[].field` | string | ○ | エラー対象のフィールド名またはクエリパラメータ名。リクエストボディがネストされたオブジェクト構造を持つ場合は `recurring_rule.due_time` のようにドット記法（`親オブジェクト.子フィールド`）で指定し、`GET` リクエスト等のクエリパラメータバリデーション違反時は対象のクエリパラメータ名（例: `priority`, `page`, `start_date`）を指定 |
| `error.details[].message` | string | ○ | フィールド固有のエラーメッセージ |

#### エラーコード (`code`) 設計方針
`error.code` には、フロントエンドが画面表示や制御分岐（再試行・画面遷移・フィールドエラー強調等）を適切に判定できるよう、HTTPステータス分類に対応する大分類コード（`BAD_REQUEST`, `UNAUTHORIZED`, `FORBIDDEN` 等）または具体的なビジネスルール違反コード（`SAME_AS_CURRENT_PASSWORD`, `SAME_AS_CURRENT_USERNAME`, `SAME_AS_CURRENT_EMAIL`, `REAUTH_FAILED`, `SESSION_DESTROYED`, `OTP_REISSUED_DUE_TO_FAILURES` 等）を格納します。

#### 代表的なエラーコード一覧

| HTTP Status | エラーコード (`code`) | 説明 |
| :--- | :--- | :--- |
| 400 | `BAD_REQUEST` | リクエスト形式またはバリデーション不正、OTP不一致 |
| 401 | `UNAUTHORIZED` | 未ログイン、セッション無効・期限切れ、ログイン認証失敗 |
| 401 | `REAUTH_FAILED` | アカウント削除/パスワード変更時の再認証失敗（1〜4回目、遅延 1.0s ± 0.1s） |
| 401 | `SESSION_DESTROYED` | アカウント削除/パスワード変更時の再認証失敗 5回連続達成分（セッション強制破棄、遅延 1.0s ± 0.1s） |
| 403 | `FORBIDDEN` | CSRFトークン不正または権限不足、未検証OTPセッションでの更新試行、異ユーザー所有OTPセッション指定 |
| 404 | `NOT_FOUND` | 指定されたリソース（または他者所有リソース）が存在しない |
| 410 | `GONE` | OTPセッションの有効期限切れ（全体最大15分経過含む） |
| 422 | `UNPROCESSABLE_ENTITY` | 汎用ビジネスルール違反 |
| 422 | `INVALID_PASSWORD_CONTENT` | パスワード変更/リセット時のユーザー名・メールアドレスローカル部含有違反 |
| 422 | `SAME_AS_CURRENT_PASSWORD` | パスワード変更/リセット時の同一パスワード指定 |
| 422 | `SAME_AS_CURRENT_USERNAME` | ユーザー名変更時の同一ユーザー名指定 |
| 422 | `SAME_AS_CURRENT_EMAIL` | メールアドレス変更時の同一メールアドレス指定 |
| 422 | `OTP_REISSUED_DUE_TO_FAILURES` | OTP検証 5回連続失敗による自動再送・カウンターリセット |
| 429 | `RATE_LIMIT_EXCEEDED` | 連続ログイン試行失敗（アカウントロック）またはIPレートリミット超過 |
| 429 | `OTP_RESEND_COOLDOWN` | OTP再送クールダウン期間（60秒）中 |
| 500 | `INTERNAL_SERVER_ERROR` | サーバー内部エラー |

---

## 2. エンドポイント一覧

| カテゴリ | メソッド | URI | 役割・機能 | 認証要否 | CSRFヘッダー |
| :--- | :--- | :--- | :--- | :---: | :---: |
| **認証 (Auth)** | `POST` | `auth/register/request-otp` | 新規登録情報のバリデーション・OTP発行・メール送信 | 不要 | 不要 |
| | `POST` | `auth/register/verify-otp` | 新規登録OTP検証・アカウント本登録・セッション発行 | 不要 | 不要 |
| | `POST` | `auth/register/resend-otp` | 新規登録OTPの再送信 | 不要 | 不要 |
| | `POST` | `auth/login` | メールアドレス・パスワードによるログイン認証 | 不要 | 不要 |
| | `POST` | `auth/logout` | ログインセッションの破棄・ログアウト | 必須 | 必須 |
| | `POST` | `auth/password-reset/request-otp` | パスワードリセット用OTP発行・メール送信 | 不要 | 不要 |
| | `POST` | `auth/password-reset/verify-otp` | パスワードリセット用OTP検証 | 不要 | 不要 |
| | `POST` | `auth/password-reset/resend-otp` | パスワードリセット用OTPの再送信 | 不要 | 不要 |
| | `POST` | `auth/password-reset/reset` | 新パスワードの設定完了処理 | 不要 | 不要 |
| | `POST` | `auth/change-email/request-otp` | メールアドレス変更用OTP作成・送信 | 必須 | 必須 |
| | `POST` | `auth/change-email/verify-otp` | メールアドレス変更用OTP検証・変更確定 | 必須 | 必須 |
| | `POST` | `auth/change-email/resend-otp` | メールアドレス変更用OTPの再送信 | 必須 | 必須 |
| **ユーザー (Users)** | `GET` | `users/{user_id}` | ログインユーザーのプロフィール情報取得 | 必須 | 不要 |
| | `PUT` | `users/{user_id}` | プロフィール情報（ユーザー名等）の更新 | 必須 | 必須 |
| | `DELETE` | `users/{user_id}` | アカウント論理削除 | 必須 | 必須 |
| | `PATCH` | `users/{user_id}/password` | ログイン状態でのパスワード変更 | 必須 | 必須 |
| **タスク (Tasks)** | `GET` | `tasks` | タスク一覧取得（検索・絞り込み・カレンダー期間取得・ページネーション） | 必須 | 不要 |
| | `POST` | `tasks` | 新規タスク作成（単一作成 / 毎週繰り返し一括作成） | 必須 | 必須 |
| | `GET` | `tasks/{task_id}` | 単一タスクの詳細取得 | 必須 | 不要 |
| | `PATCH` | `tasks/{task_id}` | タスク情報の部分更新（特定フィールドのみの更新） | 必須 | 必須 |
| | `DELETE` | `tasks/{task_id}` | タスクの物理削除 | 必須 | 必須 |