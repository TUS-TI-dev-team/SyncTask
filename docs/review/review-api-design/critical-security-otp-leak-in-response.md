# メールアドレス変更リクエストAPIにおけるOTPレスポンス返却の脆弱性

- **Status**: Resolved
- **Severity**: Major
- **Created At**: 2026-08-17 12:05:00
- **Target Files**:
  - [api_design.md](docs/design/api_design.md)

## 1. 問題の概要
`auth/change-email/request-otp` の出力（Response）仕様に平文の `OTP` が含まれており、メール受信環境を持たない第三者であってもレスポンスを見るだけで認証を突破できてしまう重大なセキュリティ脆弱性（CWE-319 / CWE-640）が存在します。

## 2. 詳細な指摘内容
`docs/design/api_design.md` の 22行目に以下の記載があります：

```markdown
| `auth/change-email/request-otp` | `POST` | メールアドレス変更用OTP作成・送信 | プロフィール編集画面 | メールアドレス変更後「決定」時 | セッションID、新メールアドレス | 検証ステータス、OTP | セッションID検証必須 |
```

- 出力に `OTP`（ワンタイムパスワード）が返却される仕様となっています。
- OTP を HTTP レスポンスボディに含めてクライアントに返してしまうと、メールアドレスの所有権確認（メール疎通確認）という OTP の本来のセキュリティ目的が完全に無効化されます。
- また、他の OTP 発行エンドポイント（`auth/register/request-otp` や `auth/password-reset/request-otp`）では OTP 自体は返却されておらず、レスポンス仕様としても一貫性を欠いています。

## 3. 推奨される修正案
1. `auth/change-email/request-otp` の出力から `OTP` を削除し、処理結果ステータス（および必要に応じて OTP セッション識別子など）のみを返却する仕様に修正してください。
2. 実際の OTP はバックエンドから新メールアドレス宛のメール送信でのみ通知されるようにしてください。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 12:40:00
- **Status**: Resolved

### 実施した修正内容
- `docs/design/api_design.md` の `POST auth/change-email/request-otp` レスポンス仕様から平文 `OTP` を完全に削除しました。
- レスポンスには `otp_session_id`, `masked_email`, `expires_in_seconds` のみを返し、OTP 自体はメール送信でのみ通知される安全な仕様へ修正しました。

### 変更したファイル
- [api_design.md](docs/design/api_design.md)
