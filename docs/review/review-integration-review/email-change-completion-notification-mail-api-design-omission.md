# メールアドレス変更確定APIにおける旧メールアドレス宛変更完了通知メール送信仕様の記載欠落

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-22 14:15:00
- **Target Files**:
  - [02_auth.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/api_design/02_auth.md#L563-L576)
  - [02_account_edit.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/02_account_edit.md#L141)
  - [01_account_and_auth.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/req-def/requirements/01_account_and_auth.md#L72)

## 1. 問題の概要
要件定義書（`01_account_and_auth.md` 2.2.2節）および処理設計書（`02_account_edit.md` 2.6節）において、メールアドレス変更確定（OTP照合成功）時に、不正変更検知・セキュリティ通知を目的として「更新前に退避した旧メールアドレス宛に変更完了通知メールを非同期送信する」仕様が規定されています。
しかし、API設計書（`02_auth.md` 3.1.11 `POST auth/change-email/verify-otp`）のリクエスト評価順序および処理説明において、この旧メールアドレス宛の変更完了通知メール送信に関する記述が欠落しています。

## 2. 詳細な指摘内容
1. **[01_account_and_auth.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/req-def/requirements/01_account_and_auth.md) 72行目**:
   - `変更確定と同時に、旧メールアドレス宛に変更完了通知メールを送信する。`
2. **[02_account_edit.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/02_account_edit.md) 141行目**:
   - `5. コミット後、更新前に退避した旧メールへ変更完了通知を非同期送信する。通知失敗で確定済み変更はロールバックせず、失敗を記録する。`
3. **[02_auth.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/api_design/02_auth.md) 574-576行目 (3.1.11節 リクエスト評価順序 ステップ5)**:
   - `5. メールアドレス更新確定・全ログインセッション物理削除 (200 OK): OTP照合成功時、ユーザーのメールアドレスを新アドレスへ更新し、セキュリティ要件として当該ユーザーのすべての既存ログインセッション（LOGIN_SESSION）および使用済み OTP_SESSION をDBから直ちに物理削除します。レスポンスヘッダーに Cookie 削除ヘッダー（Max-Age=0）を付与して返却し、フロントエンド側でログイン画面へリダイレクト（新メールアドレスでの再ログイン要求）させます。`
   - → 旧メールアドレス宛への非同期完了通知メール送信についての言及がありません。

## 3. 推奨される修正案
[02_auth.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/api_design/02_auth.md) 3.1.11節（`POST auth/change-email/verify-otp`）の評価順序ステップ5および補足説明に、以下の記述を追記してください。

```markdown
5. **メールアドレス更新確定・全ログインセッション物理削除 (`200 OK`)**:
   OTP照合成功時、ユーザーのメールアドレスを新アドレスへ更新し、セキュリティ要件として**当該ユーザーのすべての既存ログインセッション（`LOGIN_SESSION`）および使用済み `OTP_SESSION` をDBから直ちに物理削除**します。また、コミット後に更新前の旧メールアドレス宛へ変更完了通知メールを非同期送信します（通知送信失敗時も確定済みメールアドレス更新はロールバックせず、エラーログを記録）。レスポンスヘッダーに Cookie 削除ヘッダー（`Max-Age=0`）を付与して返却し、フロントエンド側でログイン画面へリダイレクト（新メールアドレスでの再ログイン要求）させます。
```

---

## 修正完了報告

- **Resolved At**: 2026-08-22 14:10:00
- **Status**: Resolved

### 実施した修正内容
`docs/design/api_design/02_auth.md` 3.1.11節（`POST auth/change-email/verify-otp`）のリクエスト評価順序ステップ5に、コミット後の旧メールアドレス宛変更完了通知メール非同期送信仕様（失敗時ログ記録、ロールバックなし）を追記しました。

### 変更したファイル
- [02_auth.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/api_design/02_auth.md)

