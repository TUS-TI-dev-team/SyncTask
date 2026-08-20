# プロフィール更新API (`PUT users/{user_id}`) におけるユーザー名バリデーション仕様の記述不足（空白トリムルールの欠落）

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-17 16:35:00
- **Target Files**:
  - [03_users.md](docs/design/api_design/03_users.md)
  - [02_auth.md](docs/design/api_design/02_auth.md)
  - [requirements.md](docs/req-def/requirements.md)

## 1. 問題の概要
`PUT users/{user_id}` (3.2.2) のリクエストパラメータテーブルにおいて、`username` フィールドの制約が「2〜20文字、英数字。現在のユーザー名と同一の場合は 422 エラー」と記載されており、新規アカウント登録時 (`02_auth.md` L19) や要件定義書 (`requirements.md` L27) に明記されている「前後の空白自動トリム (trimming)」のルールが欠落している。

## 2. 詳細な指摘内容
- `02_auth.md` (3.1.1 `POST auth/register/request-otp` L19):
  `username`: `2〜20文字、英数字（大文字小文字可）、前後の空白トリム`
- `requirements.md` (L27):
  `ユーザー名: 2〜20文字の英数字（大文字小文字可）。前後の空白は自動トリム。`
- `03_users.md` (3.2.2 L45):
  `username`: `2〜20文字、英数字。現在のユーザー名と同一の場合は 422 エラー`

`03_users.md` に「前後の空白トリム」が記載されていないため、クライアントが `" example "` のように前後空白を含むユーザー名を送信した際に、以下の懸念が生じる:
1. トリム前の文字列長でバリデーションが行われるのか、トリム後に 2〜20 文字チェックを行うのかが不透明。
2. トリム前の文字列と現在のユーザー名を比較した場合、トリム処理漏れにより `SAME_AS_CURRENT_USERNAME` (422) チェックを通り抜けて同名（実質同一）として更新されてしまうリスク。

## 3. 推奨される修正案
`03_users.md` 3.2.2 節のパラメータテーブルおよびエラー説明を以下のように修正してください。

| フィールド | 型 | 必須 | 制約・バリデーション |
| :--- | :--- | :---: | :--- |
| `username` | string | ○ | 2〜20文字、英数字（大文字小文字可）。前後の空白は自動トリム。トリム後のユーザー名が現在のユーザー名と同一の場合は 422 エラー |

---

## 修正完了報告

- **Resolved At**: 2026-08-17 16:42:00
- **Status**: Resolved

### 実施した修正内容
`PUT users/{user_id}` (3.2.2) の `username` バリデーションテーブルに「前後の空白は自動トリム。トリム後のユーザー名が現在のユーザー名と同一の場合は 422 エラー」の規定を追記しました。

### 変更したファイル
- [03_users.md](file:////mnt/c/Users/ayumu_wkkd3w3/BoxForSomething/github/SyncTask/docs/design/api_design/03_users.md)
