# 新規登録およびパスワードリセットにおけるユーザー名・メールアドレスローカル部包含チェックの大文字小文字非区別規定漏れ

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-17 17:22:00
- **Target Files**:
  - [02_auth.md](docs/design/api_design/02_auth.md)

## 1. 問題の概要
`02_auth.md` の 3.1.1 節 (`POST auth/register/request-otp`) および 3.1.9 節 (`POST auth/password-reset/reset`) のパスワード制約において、「ユーザー名・メールのローカル部（4文字以上の場合）を含まないこと」と規定されていますが、大文字・小文字を区別せずに比較する（Case-Insensitive）旨が明記されていません。

## 2. 詳細な指摘内容
1. `username` は「大文字小文字可」、`email` は「小文字正規化」と定義されていますが、ユーザー入力時のパスワード `password` / `new_password` 内に大文字・小文字表記が異なるユーザー名またはメールローカル部（例: ユーザー名 `JohnDoe` に対して `johndoe123!` や `JOHNDOE123!`）が含まれていた場合、完全一致（Case-Sensitive）チェックではバリデーションを通過してしまう恐れがあります。
2. パスワード変更 API (`PATCH users/{user_id}/password`) では既に「大文字小文字を区別せず比較」が明記されていますが、新規登録 (`3.1.1`) およびパスワードリセット (`3.1.9`) ではこの規定が抜け落ちており、仕様の揺れが生じています。

## 3. 推奨される修正案
`02_auth.md` の 3.1.1 節のリクエストパラメータテーブル (`password`) および 3.1.9 節のリクエストパラメータテーブル (`new_password`) の説明をそれぞれ以下のように修正してください。
- 3.1.1: `8〜128文字、英大文字/英小文字/数字/記号のうち3種以上を含む。ユーザー名・メールのローカル部（4文字以上の場合、大文字小文字を区別せず比較）を含まないこと`
- 3.1.9: `8〜128文字、英大文字/英小文字/数字/記号のうち3種以上を含む。検証済み OTP_SESSION 経由で取得したユーザーのユーザー名・メールのローカル部（4文字以上の場合、大文字小文字を区別せず比較）を含まないこと。現在のパスワードと同一の場合は 422 エラー`

---

## 修正完了報告

- **Resolved At**: 2026-08-17 17:27:30
- **Status**: Resolved

### 実施した修正内容
`02_auth.md` の 3.1.1 節 (`POST auth/register/request-otp`) および 3.1.9 節 (`POST auth/password-reset/reset`) のパラメータテーブル (`password`, `new_password`) の説明欄に、「4文字以上の場合、大文字小文字を区別せず比較」の規定を明記し、全パスワード変更/設定エンドポイント間でバリデーションルールの統一を図りました。

### 変更したファイル
- [02_auth.md](file:///mnt/c/Users/ayumu_wkkd3w3/BoxForSomething/github/SyncTask/docs/design/api_design/02_auth.md)
