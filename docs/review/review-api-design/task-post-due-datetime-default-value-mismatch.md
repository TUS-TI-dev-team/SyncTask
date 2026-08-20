# POST tasks における due_datetime のデフォルト値および未設定（NULL）扱いの不整合

- **Status**: Resolved
- **Severity**: Major
- **Created At**: 2026-08-17 16:30:00
- **Target Files**:
  - [04_tasks.md](docs/design/api_design/04_tasks.md)
  - [requirements.md](docs/req-def/requirements.md)
  - [database_design.md](docs/design/database_design.md)

## 1. 問題の概要
`04_tasks.md` の `POST tasks` (L95) において、`due_datetime` の説明に「省略時は当日 23:59:00+09:00」と定義されていますが、これは要件定義書 (`requirements.md`) および DB設計書 (`database_design.md`) の仕様と矛盾しており、締切日時が存在しないタスク（未設定 / NULL）を新規作成できない不具合の原因となります。

## 2. 詳細な指摘内容
1. **締切日時未設定（NULL）タスク作成の不可**:
   - `requirements.md` (L81, L96, L106, L116, L123, L132) および `database_design.md` (L49) では、タスクの締切日時は「任意設定」であり、締切日時が「未設定（NULL）」のタスクが存在することを前提としたソート順ルール（締切あり昇順 ➔ 締切なし/同一作成日時降順）や絞り込みロジックが定義されています。
   - しかし `04_tasks.md` L95 の「単一作成時用（省略時は当日 23:59:00+09:00）」という定義に従うと、単一作成時に `due_datetime` を省略した場合にすべて「当日23:59:00」が強制設定されてしまい、作成時に締切なし（NULL）のタスクを作成できなくなります。

2. **「時刻未指定時のデフォルト」と「日付自体の省略」の混同**:
   - `requirements.md` L81 の「時刻部分［hh:mm］が未指定［省略］の場合はデフォルト締切時刻として 23:59 JST を適用する」という規定は、ユーザーが締切日（日付）を指定しつつ時刻を省略した場合の挙動を指しています。
   - `POST tasks` の仕様で `due_datetime` フィールドごと省略（または `null` 指定）された場合は「締切日時未設定（NULL）」として登録されるべきです。

## 3. 推奨される修正案
1. `04_tasks.md` L95 の `due_datetime` フィールド定義を以下のように修正してください:
   - **型**: `string / null`
   - **制約・バリデーション**: ISO 8601 日時文字列（例: `2026-08-20T23:59:00+09:00`）、または日付のみ `YYYY-MM-DD`（時刻省略時は `23:59:00+09:00` を付与）。省略時または `null` 指定時は「締切日時未設定（NULL）」として作成。
2. リクエストボディ例および単一タスク作成レスポンスに、締切日時未設定 (`due_datetime: null`) のパターンについての補足を明記してください。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 16:42:00
- **Status**: Resolved

### 実施した修正内容
`POST tasks` における `due_datetime` の型を `string / null` に更新し、省略時または `null` 指定時は「締切日時未設定（NULL）」として作成し、日付のみ `YYYY-MM-DD` 指定時のみ時刻 `23:59:00+09:00` を補完するよう仕様を適正化しました。

### 変更したファイル
- [04_tasks.md](file:////mnt/c/Users/ayumu_wkkd3w3/BoxForSomething/github/SyncTask/docs/design/api_design/04_tasks.md)
