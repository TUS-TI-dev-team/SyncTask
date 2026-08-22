# プロフィール編集画面における確認ダイアログキャンセル時およびメール送信失敗時の画面遷移・留まり挙動の曖昧さ

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-22 14:05:00
- **Target Files**:
  - [screen_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/screen_design.md#L61)
  - [02_account_edit.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/02_account_edit.md#L53-L96)

## 1. 問題の概要
[screen_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/screen_design.md) の「プロフィール編集画面」において、表の「遷移先」欄の定義と「役割・機能」の説明文との間で記述の粒度に食い違いがあり、確認ダイアログでのキャンセル時やメール送信失敗時の画面遷移（留まり）挙動が未確定に見える状態となっています。

## 2. 詳細な指摘内容
1. **[screen_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/screen_design.md) 61行目（遷移先欄）**:
   - `メール変更無し：プロフィール表示画面<br>メール変更あり：アカウント関連/OTP入力画面`
   - 「キャンセル」ボタン押下時（プロフィール表示画面へ遷移）や、確認ダイアログ表示後に「キャンセル」を選択した場合（編集画面に留まる）の遷移先定義が遷移先列から欠落しています。
2. **メール送信失敗時（503）の挙動の違い**:
   - [02_account_edit.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/02_account_edit.md) 53〜96行目では：
     - 「ユーザー名＋メール両方変更時」: ユーザー名は既に更新されているため、メール送信失敗時はトースト通知を表示して「プロフィール表示画面」へ遷移する。
     - 「メールのみ変更時」: メール送信失敗時はエラートーストを表示して「プロフィール編集画面」に留まり再入力を促す。
   - この2つのケースにおける遷移先の違いが画面設計書の「役割・機能」および「遷移先」欄で区別されていません。

## 3. 推奨される修正案
[screen_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/screen_design.md) 61行目の「遷移先」欄および「役割・機能」欄を更新し、以下の内容を明確に記述してください。

```markdown
| プロフィール編集画面 | ... | 変更確定（メール変更なし）：プロフィール表示画面<br>変更確定（メール変更あり）：アカウント関連/OTP入力画面<br>キャンセル・確認ダイアログでキャンセル：プロフィール編集画面に留まる（編集キャンセルボタン押下時はプロフィール表示画面へ戻る）<br>メール送信失敗（両方変更時）：トースト通知の上でプロフィール表示画面へ遷移<br>メール送信失敗（メールのみ変更時）：エラー通知の上でプロフィール編集画面に留まる | ... |
```

---

## 修正完了報告

- **Resolved At**: 2026-08-22 14:10:00
- **Status**: Resolved

### 実施した修正内容
`docs/design/screen_design.md` 61行目の「プロフィール編集画面」におけるキャンセル時、ダイアログキャンセル時、およびメール送信失敗時（両方変更時 vs メールのみ変更時）の画面遷移・留まり挙動を明記しました。

### 変更したファイル
- [screen_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/screen_design.md)
