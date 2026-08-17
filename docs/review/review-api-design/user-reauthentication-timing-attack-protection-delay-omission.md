# アカウント削除およびパスワード変更の再認証失敗時におけるタイミング攻撃対策（応答遅延 1.0s ± 0.1s）の記述漏れ

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-17 16:35:00
- **Target Files**:
  - [03_users.md](docs/design/api_design/03_users.md)
  - [01_overview.md](docs/design/api_design/01_overview.md)

## 1. 問題の概要
`01_overview.md`（1.2 節）では、Timing Attack 対策およびブルートフォース保護のため、ログイン失敗・OTP検証失敗・パスワード照合を伴う認証失敗時に一律 `1.0s ± 0.1s` のレスポンス遅延を適用することが定められている。

しかし、`03_users.md` の `DELETE users/{user_id}` (3.2.3) および `PATCH users/{user_id}/password` (3.2.4) におけるパスワード再認証失敗時のエラー処理（`401 Unauthorized` / `REAUTH_FAILED`）において、応答遅延（`1.0s ± 0.1s`）の適用ルールが明記されていない。

## 2. 詳細な指摘内容
- `01_overview.md` L36-L37 (1.2 遅延制御):
  `遅延制御 (Timing Attack 対策): ログイン失敗、OTP検証失敗、アカウント存在有無のダミー処理時は、一律 1.0s ± 0.1s のレスポンス遅延を適用します。`
- `03_users.md` 3.2.3 / 3.2.4:
  再認証失敗時（パスワード誤り）のエラー記述において、応答遅延に関する注記が存在しない。

パスワードのハッシュ化（bcrypt / Argon2等）の計算時間差やパスワード比較処理の高速失敗を悪用した Timing Attack（タイミング攻撃）を防止するため、ログイン認証時と同様にアカウント削除時およびパスワード変更時の再認証失敗応答に対しても一律 `1.0s ± 0.1s` のレスポンス遅延を強制する必要がある。

## 3. 推奨される修正案
`03_users.md` の 3.2.3 節 (`DELETE users/{user_id}`) および 3.2.4 節 (`PATCH users/{user_id}/password`) の Errors セクション、または機能概要に、再認証失敗時の一律レスポンス遅延注記を追加してください。

```markdown
※ パスワード再認証失敗時（1〜4回目および5回連続達成分）は、Timing Attack 対策として一律 `1.0s ± 0.1s` のレスポンス遅延を適用してエラー（`401 Unauthorized`）を返却します。
```

---

## 修正完了報告

- **Resolved At**: 2026-08-17 16:42:00
- **Status**: Resolved

### 実施した修正内容
`DELETE users/{user_id}` (3.2.3) および `PATCH users/{user_id}/password` (3.2.4) の再認証失敗時（1〜4回目および5回連続達成分）において、Timing Attack 対策として一律 `1.0s ± 0.1s` のレスポンス遅延を適用する仕様・注記を追加しました。

### 変更したファイル
- [03_users.md](file:////mnt/c/Users/ayumu_wkkd3w3/BoxForSomething/github/SyncTask/docs/design/api_design/03_users.md)
