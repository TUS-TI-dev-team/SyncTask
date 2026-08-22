# OTP検証5回連続失敗時の自動再送通知およびUIフィードバック仕様の記述漏れ

- **Status**: Resolved
- **Severity**: Major
- **Created At**: 2026-08-22 14:05:00
- **Target Files**:
  - [screen_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/screen_design.md#L7-L9)
  - [screen_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/screen_design.md#L62)
  - [screen_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/screen_design.md#L73-L78)
  - [02_auth.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/api_design/02_auth.md#L114)
  - [01_account_creation.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/01_account_creation.md#L58)

## 1. 問題の概要
API設計および処理設計では、OTP入力検証が5回連続で失敗した際、バックエンドが新OTPを自動生成・送信して `422 Unprocessable Entity`（code: `OTP_REISSUED_DUE_TO_FAILURES`）を返し、画面上に「新しいOTPを再送しました」と通知して入力欄をクリアし、60秒クールダウンを開始する仕様となっています。
しかし、[screen_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/screen_design.md) にはこの重要なUIフィードバックおよび状態更新（422受信時の挙動）が一切記載されていません。

## 2. 詳細な指摘内容
1. **[screen_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/screen_design.md) の現状記述**:
   - 7行目、9行目、62行目、76行目等において、「OTP発行後5分経過時のポップアップ通知」や「送信が5回連続で失敗してセッションが失効した場合の遷移」は記載されています。
   - 一方で、「OTP入力の誤り（検証照合）が5回連続した場合」の画面挙動についての記述が存在しません。
2. **[02_auth.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/api_design/02_auth.md) 114行目 および [01_account_creation.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/01_account_creation.md) 58行目**:
   - 不一致5回目は、失敗回数をリセットし、OTP自動再発行通知 `422 Unprocessable Entity`（code: `OTP_REISSUED_DUE_TO_FAILURES`、遅延 1.0s ± 0.1s）を返却する。
   - フロントエンドは「新しいOTPを再送しました」とポップアップ/トースト通知し、OTP入力欄をクリアし、この時点から「再送信」ボタンの60秒クールダウンを開始する。
3. **影響**:
   - 画面設計書に記載がないため、フロントエンド実装時に `422 OTP_REISSUED_DUE_TO_FAILURES` が単なる通常のエラー（400等）と同様に扱われ、ユーザーに自動再送された旨が伝わらず、古い無効化されたOTPを再入力し続けて混乱する恐れがあります。

## 3. 推奨される修正案
[screen_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/screen_design.md) の各OTP画面（7行目、9行目、62行目）の「役割・機能」欄、および「OTP・セキュリティ操作補足」（73〜78行目）に、以下の仕様を明記してください。

```markdown
- **OTP検証5回連続失敗時の自動再送挙動**: OTP入力の検証に5回連続で失敗した場合（APIより `422 OTP_REISSUED_DUE_TO_FAILURES` が返却された場合）、画面上に「入力試行回数の上限に達したため、新しい認証コードを再送信しました」という通知（ポップアップ/トースト）を表示し、OTP入力欄をクリアするとともに、「再送信」ボタンの60秒クールダウンカウントダウンを開始する。
```

---

## 修正完了報告

- **Resolved At**: 2026-08-22 14:10:00
- **Status**: Resolved

### 実施した修正内容
`docs/design/screen_design.md` の各OTP画面および「OTP・セキュリティ操作補足」に、5回連続検証失敗時の自動再送通知（`422 OTP_REISSUED_DUE_TO_FAILURES`）、入力欄クリア、および60秒クールダウンUI仕様を明記しました。

### 変更したファイル
- [screen_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/screen_design.md)
