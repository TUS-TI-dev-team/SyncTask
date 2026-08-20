# POST auth/password-reset/reset におけるリクエスト評価順序の不備とユーザー参照コンテキスト不足

- **Status**: Resolved
- **Severity**: Major
- **Created At**: 2026-08-17 17:22:00
- **Target Files**:
  - [02_auth.md](docs/design/api_design/02_auth.md)

## 1. 問題の概要
`02_auth.md` の 3.1.9 節 (`POST auth/password-reset/reset`) の「リクエスト評価順序」において、ステップ 1「リクエスト構文・入力バリデーション (400 Bad Request)」の時点で `new_password` の「ユーザー名・メールのローカル部含有制約」を検証すると規定されています。しかし、リクエストボディには `otp_session_id` と `new_password` のみが含まれており、ステップ 2「OTPセッション状態・期限検証」で DB 上の `OTP_SESSION` および紐づく `USER` レコードを確認する前にユーザーの `username` や `email` を参照することは不可能です。

## 2. 詳細な指摘内容
1. 3.1.9 節のリクエストボディには `otp_session_id` と `new_password` のみが存在します（`username` や `email` は含まれません）。
2. にもかかわらず、ステップ 1（L369）において「`new_password` の文字数・文字種・ユーザー名/メール含有制約を検証します」と記述されており、DB 照合前かつ `otp_session_id` の存在・状態（`verified`）判定前にユーザーの属性情報を参照しようとする構造的な矛盾が発生しています。
3. 万が一 `otp_session_id` が未存在・不正・ダミーセッション・期限切れである場合、ステップ 1 の時点で参照すべきユーザーが存在しないため、バリデーション処理がエラーまたは予期せぬ挙動を引き起こすリスクがあります。

## 3. 推奨される修正案
`02_auth.md` 3.1.9 節の「リクエスト評価順序」およびパラメータ制約記述を以下のように修正してください：
1. **ステップ 1（400 Bad Request）**: 単体の構文チェック（`new_password` の 8〜128文字、文字種要件、`otp_session_id` の形式チェック）のみを行う旨を明記。
2. **ステップ 2（403 Forbidden / 410 Gone）**: `otp_session_id` の存在・ステータス（`verified`）・有効期限を検証。
3. **ステップ 3（422 Unprocessable Entity）**: 検証成功した `OTP_SESSION` から対象ユーザーの `username` および `email` を参照し、`new_password` 内への含有チェック（大文字小文字非区別）および現在のパスワードとの同一性チェックを実施。
4. パラメータテーブル `new_password` の説明に「検証済み `OTP_SESSION` 経由で DB から取得した対象ユーザーのユーザー名・メールアドレスのローカル部を含まないこと」と追記。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 17:27:30
- **Status**: Resolved

### 実施した修正内容
`02_auth.md` 3.1.9 節 (`POST auth/password-reset/reset`) において以下の通り評価順序およびパラメータ記述を更新・明確化しました：
1. ステップ 1（400 Bad Request）を `new_password` の文字数・文字種要件および `otp_session_id` の形式チェック等の単体構文検証に限定。
2. ステップ 2（403 Forbidden / 410 Gone）で `otp_session_id` の存在・ステータス（`verified`）・有効期限を検証。
3. ステップ 3（422 Unprocessable Entity）にて検証成功した `OTP_SESSION` 経由で取得した対象ユーザーの `username` および `email` ローカル部（大文字小文字非区別）の `new_password` への含有チェックおよび現パスワード同一性チェックを実施。
4. パラメータテーブルおよび Errors セクションの記述を更新。

### 変更したファイル
- [02_auth.md](file:///mnt/c/Users/ayumu_wkkd3w3/BoxForSomething/github/SyncTask/docs/design/api_design/02_auth.md)
