# パスワード変更画面およびアカウント削除画面における再認証失敗時（1〜4回目）のUIエラー表示仕様の曖昧さ

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-22 14:05:00
- **Target Files**:
  - [screen_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/screen_design.md#L63-L64)
  - [03_users.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/api_design/03_users.md#L150-L153)
  - [03_users.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/api_design/03_users.md#L207-L210)
  - [03_account_delete.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/03_account_delete.md#L60-L65)

## 1. 問題の概要
[screen_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/screen_design.md) の「パスワード変更画面」および「アカウント削除/パスワード再認証画面」において、パスワード再認証に「5回連続で失敗した場合のセッション破棄・ログイン画面リダイレクト」は記載されていますが、1〜4回目の失敗時（APIから `401 REAUTH_FAILED` が返却された場合）におけるUIフィードバック（エラー表示文言、残試行回数の表示有無、入力欄のクリア等）の仕様が明記されていません。

## 2. 詳細な指摘内容
1. **[screen_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/screen_design.md) 63行目・64行目の現状記述**:
   - 「5回連続で既存パスワード認証に失敗した場合もログインセッションを破棄してログイン画面にリダイレクト」と5回目の強制破棄のみが強調されています。
2. **APIおよび処理設計での仕様**:
   - [03_users.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/api_design/03_users.md) および [03_account_delete.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/03_account_delete.md) では、1〜4回目の失敗時は `401 REAUTH_FAILED`（遅延 1.0s ± 0.1s）が返却され、5回目の失敗時は `401 SESSION_DESTROYED`（遅延 1.0s ± 0.1s、Cookie消去）が返却されます。
3. **曖昧点**:
   - 1〜4回目の失敗時に、画面のインラインエラーまたはアラートバナーにどのようなメッセージを表示するのか（「パスワードが正しくありません」等の汎用エラーのみか、あるいはセキュリティ上残試行回数を表示しないのか）が画面設計書上で定義されていません。
   - 失敗時にパスワード入力欄をクリアしてフォーカスを当てるかどうかのUI挙動が未定義です。

## 3. 推奨される修正案
[screen_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/screen_design.md) の63行目・64行目および「共通UI・エラー表示補足」（79〜81行目）に、再認証失敗時（1〜4回目の `REAUTH_FAILED` 受信時）の画面挙動を具体的に追記してください。

```markdown
- **パスワード再認証失敗時の画面制御**:
  - 1〜4回目の認証失敗（`401 REAUTH_FAILED` 受信時）: 現在のパスワード入力欄の下部に「パスワードが正しくありません」とインラインエラーを表示し、入力欄の値をクリアして再入力を促す（※アカウント列挙・推測攻撃防止のため残試行回数は画面表示しない）。
  - 5回目の連続認証失敗（`401 SESSION_DESTROYED` 受信時）: アラート通知とともに直ちに全セッション破棄・Cookie消去を行い、ログイン画面へ強制リダイレクトする。
```

## 修正完了報告

- **Resolved At**: 2026-08-22 14:10:00
- **Status**: Resolved

### 実施した修正内容
- `docs/design/screen_design.md` のパスワード変更画面・アカウント削除画面および共通補足欄に、1〜4回目再認証失敗時（`401 REAUTH_FAILED`）のインラインエラー表示「パスワードが正しくありません」、入力欄クリア、および残試行回数非表示仕様を明記しました。
- `docs/design/process_design/03_account_delete.md` にも同様の画面挙動を具体化しました。

### 変更したファイル
- [screen_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/screen_design.md)
- [03_account_delete.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/03_account_delete.md)
