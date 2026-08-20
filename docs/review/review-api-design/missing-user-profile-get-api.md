# ログインユーザーのプロフィール情報取得APIの欠落

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-17 12:05:00
- **Target Files**:
  - [api_design.md](docs/design/api_design.md)
  - [screen_design.md](docs/design/screen_design.md)

## 1. 問題の概要
画面設計書で定義されている「プロフィール表示画面」および「ヘッダー」において、現在ログイン中のユーザー情報（ユーザー名、メールアドレス）を表示するための取得API（`GET`）がAPI一覧に定義されていません。

## 2. 詳細な指摘内容
- `docs/design/screen_design.md` L47 では「プロフィール表示画面: アカウント情報（ユーザー名、メールアドレス）を表示」、L57 ではヘッダー表示が定義されています。
- しかし `docs/design/api_design.md` にはユーザー関連APIとして `users/{user_id}` の `PUT`（更新）と `DELETE`（削除）、`PATCH`（パスワード変更）しか定義されておらず、自身のプロフィール情報を取得する `GET` エンドポイントが存在しません。
- また、パスパラメータ `{user_id}` を指定する形式の場合、クライアントが事前に自身の `user_id` を知る手段がないか、他人の `user_id` を指定した際の認可制御（BOLA/IDOR対策）が必要になります。

## 3. 推奨される修正案
1. ログインユーザー自身のプロフィール情報を取得するエンドポイント `GET users/me`（または `GET users/{user_id}`）を追加してください。
2. 返却するレスポンス情報として `user_id`, `username`, `email` などのフィールドを定義してください。
3. セキュリティおよびRESTful設計の観点から、自身の情報操作については `users/me` へのルーティングを検討してください。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 12:40:00
- **Status**: Resolved

### 実施した修正内容
- `docs/design/api_design.md` に `GET users/{user_id}` エンドポイントを追加定義しました。
- ユーザーID、ユーザー名、メールアドレス、作成日時等を返却するレスポンススキーマを定義し、セッションと一致しない `{user_id}` や他ユーザーの指定時は `404 Not Found` を返却する認可制御仕様を明記しました。

### 変更したファイル
- [api_design.md](docs/design/api_design.md)
