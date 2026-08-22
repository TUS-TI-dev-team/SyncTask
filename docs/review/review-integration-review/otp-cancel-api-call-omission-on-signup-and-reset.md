# アカウント作成およびパスワードリセットOTP画面の戻る・キャンセル時における破棄API呼び出し記述の欠落

- **Status**: Resolved
- **Severity**: Major
- **Created At**: 2026-08-22 14:05:00
- **Target Files**:
  - [screen_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/screen_design.md#L7-L9)
  - [01_overview.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/api_design/01_overview.md#L21)
  - [02_auth.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/api_design/02_auth.md#L59-L91)

## 1. 問題の概要
[screen_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/screen_design.md) において、「アカウント関連/OTP入力画面（メールアドレス変更用）」の戻る・キャンセル操作には `POST auth/otp-session/cancel` を呼び出してサーバー側セッションを即時無効化する旨が明記されていますが、「アカウント作成/OTP入力画面」および「パスワードリセット/OTP入力画面」の戻る・キャンセル操作には当該API呼び出しの記述が欠落しています。

## 2. 詳細な指摘内容
1. **[screen_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/screen_design.md) 62行目（メール変更OTP画面）**:
   - `戻る・キャンセル：プロフィール編集画面（クライアント側 otp_session_id を破棄し、POST auth/otp-session/cancel を呼び出してサーバー側セッションを即時無効化）` と記載されています。
2. **[screen_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/screen_design.md) 7行目・9行目（アカウント作成・PWリセットOTP画面）**:
   - 7行目: `戻る・キャンセル：アカウント作成/情報入力画面`
   - 9行目: `戻る・キャンセル：パスワードリセット/メールアドレス入力画面`
   - セッション破棄APIの呼び出しについて触れられていません。
3. **[01_overview.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/api_design/01_overview.md) 21行目 および [02_auth.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/api_design/02_auth.md) 60行目**:
   - 「OTP入力画面からの離脱や「戻る」ボタン押下時は、クライアント側の `otp_session_id` を破棄するとともに、セッション破棄API（`POST auth/otp-session/cancel`）を呼び出してサーバー側 `OTP_SESSION` を即座に物理削除（無効化）します。」と定義されており、認証不要で会員登録・PWリセット・メール変更の全種別で利用可能とされています。
4. **影響**:
   - 画面設計の記述の差異により、新規登録やパスワードリセット時にユーザーが「戻る」を押下した際、フロントエンドが `POST auth/otp-session/cancel` を呼び出さずに画面遷移してしまう実装リスクが生じます。この場合、サーバー側に15分間セッションが残存し、同一メールアドレスでの再操作時にダミーセッション扱いや排他ロックによる不整合が発生する懸念があります。

## 3. 推奨される修正案
[screen_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/screen_design.md) の「アカウント作成/OTP入力画面」（7行目）および「パスワードリセット/OTP入力画面」（9行目）の遷移先欄、ならびに「OTP・セキュリティ操作補足」に、戻る・キャンセル時の `POST auth/otp-session/cancel` 呼び出し仕様を追記・統一してください。

```markdown
| アカウント作成/OTP入力画面 | ... | 決定（認証成功）：ホーム画面<br>戻る・キャンセル：アカウント作成/情報入力画面（クライアント側 otp_session_id を破棄し、POST auth/otp-session/cancel を呼び出してサーバー側セッションを即時無効化） | ... |
| パスワードリセット/OTP入力画面 | ... | 決定（認証成功）：新パスワード入力画面<br>戻る・キャンセル：パスワードリセット/メールアドレス入力画面（クライアント側 otp_session_id を破棄し、POST auth/otp-session/cancel を呼び出してサーバー側セッションを即時無効化） | ... |
```

---

## 修正完了報告

- **Resolved At**: 2026-08-22 14:10:00
- **Status**: Resolved

### 実施した修正内容
`docs/design/screen_design.md` の「アカウント作成/OTP入力画面」および「パスワードリセット/OTP入力画面」の画面遷移仕様に、戻る・キャンセル時の `POST auth/otp-session/cancel` 呼び出しおよびクライアント側セッション破棄仕様を明記・統一しました。

### 変更したファイル
- [screen_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/screen_design.md)
