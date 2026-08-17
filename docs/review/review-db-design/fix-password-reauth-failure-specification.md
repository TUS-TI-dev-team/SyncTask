# パスワード再認証失敗カウント（CHPASS_FAILED_COUNT）の仕様とセッション破棄要件の不整合

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-17 13:16:00
- **Target Files**:
  - [database_design.md](docs/design/database_design.md)
  - [requirements.md](docs/req-def/requirements.md)
  - [screen_design.md](docs/design/screen_design.md)

## 1. 問題の概要
`LOGIN_ACCOUNT` テーブルの `CHPASS_FAILED_COUNT` の備考に「5分ごとにリセット」と記載されていますが、要件定義書に定められている「パスワード変更時およびアカウント削除時の5回連続失敗によるログインセッション破棄（物理削除）」という要件と整合していません。また、アカウント削除時の再認証失敗の扱いが曖昧です。

## 2. 詳細な指摘内容
- `docs/design/database_design.md` の24〜25行目において、以下のように定義されています：
  - `CHPASS_FAILED_COUNT`: `INT` / `DEFAULT 0`、備考: `5分ごとにリセット (判定時動的リセット)`
  - `CHPASS_LAST_FAILED_AT`: `TIMESTAMPTZ`、備考: `最終失敗タイムスタンプ`
- 一方で、`docs/req-def/requirements.md` の51行目、58行目、212〜213行目および `docs/design/screen_design.md` の50〜51行目では以下のように定義されています：
  > 51: アカウントの削除: 5回連続でパスワード再認証に失敗した場合はログインセッションを破棄［物理削除］してログイン画面にリダイレクトする
  > 58: パスワード変更: 5回連続で既存パスワード認証に失敗した場合はログインセッションを破棄してログイン画面にリダイレクト
  > 212-213: パスワード変更時およびアカウント削除時の「現在のパスワード」認証: 5回連続でパスワードの認証に失敗したら、セッションを破棄（物理削除）してログイン画面に戻す
- 「5分ごとにリセット」という時間ベースのリセットルールは要件定義に存在せず、要件では「連続5回失敗」時に該当端末/セッションを物理削除して強制ログアウトさせることが求められています。
- また、カラム名が `CHPASS_FAILED_COUNT`（パスワード変更失敗回数）となっているため、同一のセキュリティ仕様が適用される「アカウント削除時のパスワード再認証失敗」が含まれるのかがDB定義上不明確です。

## 3. 推奨される修正案
- カラムの用途・命名を「パスワード変更・アカウント削除時のパスワード再認証失敗回数」（例: `REAUTH_FAILED_COUNT` または `CHPASS_FAILED_COUNT` の説明更新）として整理してください。
- 備考の記述を「5分ごとにリセット」から、要件定義に沿って「5回連続失敗時にログインセッションを物理削除してログイン画面へリダイレクト（成功時またはセッション破棄時に0クリア）」に修正してください。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 13:22:30
- **Status**: Resolved

### 実施した修正内容
`LOGIN_ACCOUNT` テーブルの `CHPASS_FAILED_COUNT` / `CHPASS_LAST_FAILED_AT` を `REAUTH_FAILED_COUNT` / `REAUTH_LAST_FAILED_AT` にリネームし、パスワード変更およびアカウント削除時の再認証共通の失敗管理カラムとして位置づけました。備考を「パスワード変更・アカウント削除時の再認証失敗回数。5回連続失敗でログインセッション物理削除（成功時またはセッション破棄時に0リセット）」に修正しました。

### 変更したファイル
- [database_design.md](docs/design/database_design.md)
