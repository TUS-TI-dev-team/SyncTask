# パスワードリセット完了画面の画面名不整合および仮セッションタイマー表示仕様の曖昧さ

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-22 14:05:00
- **Target Files**:
  - [screen_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/screen_design.md#L9-L10)
  - [06_password_reset.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/06_password_reset.md#L12)
  - [01_account_and_auth.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/req-def/requirements/01_account_and_auth.md#L126-L127)

## 1. 問題の概要
パスワードリセットフローにおける最終画面の画面名称について、[screen_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/screen_design.md) では「新パスワード入力画面」と命名されている一方、処理設計書（[06_password_reset.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/06_password_reset.md)）では「パスワードリセット画面」と記載されており、ドキュメント間で画面名に揺れがあります。
また、本画面には「15分間のOTP仮セッション有効期限」が存在しますが、画面上の構成要素にセッションタイマーUIを含めるか否かの仕様が曖昧です。

## 2. 詳細な指摘内容
1. **画面名の揺れ**:
   - [screen_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/screen_design.md) 10行目: 「新パスワード入力画面」
   - [06_password_reset.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/06_password_reset.md) 12行目、148行目: 「パスワードリセット画面」
2. **仮セッションタイマーUI仕様の曖昧さ**:
   - [screen_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/screen_design.md) 10行目において「15分のOTPセッション期限が切れた場合は自動的にパスワードリセット/メールアドレス入力画面へリダイレクトされる」とリダイレクト要件が記載されています。
   - しかし、同10行目の「構成要素」には OTP入力画面（7行目・9行目）にあるような「15分のセッションタイマー」の記載が含まれていません。
   - 画面上に残り時間のカウントダウンタイマーを表示するのか、あるいはバックグラウンドのタイマー監視/リクエスト時の410ハンドリングのみで制御するのかが未確定です。

## 3. 推奨される修正案
1. ドキュメント間で画面名称を統一する（例: 「パスワードリセット/新パスワード入力画面」または「パスワードリセット画面」に統一）。
2. [screen_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/screen_design.md) の構成要素に「15分の仮セッション有効期限タイマー表示（または有効期限カウントダウン）」を含めるか、あるいはバックグラウンドタイマー監視によるリダイレクト制御とするかを明記してください。

---

## 修正完了報告

- **Resolved At**: 2026-08-22 14:10:00
- **Status**: Resolved

### 実施した修正内容
- 画面名称を「パスワードリセット/新パスワード入力画面」に統一しました（`docs/design/screen_design.md`、`docs/design/process_design/06_password_reset.md`）。
- `docs/design/screen_design.md` の構成要素に「15分の仮セッション有効期限タイマー（カウントダウン）」を追記しました。

### 変更したファイル
- [screen_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/screen_design.md)
- [06_password_reset.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/06_password_reset.md)
