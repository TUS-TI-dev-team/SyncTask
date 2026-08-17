# POST tasks の毎週繰り返し作成における recurring_rule.due_time のフォーマット検証・JST正規化・エラーキー名の仕様不足

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-17 17:22:00
- **Target Files**:
  - [04_tasks.md](docs/design/api_design/04_tasks.md)

## 1. 問題の概要
`POST tasks` (3.3.2 節) の `recurring_rule.due_time` において、`HH:mm` フォーマット（省略時 `23:59`）とされていますが、時刻形式違反（例: `"25:00"`, `"9:0"`, 不正文字列等）が指定された場合のバリデーション制約・エラーレスポンス構造（`error.details[].field: "recurring_rule.due_time"`）、および一括生成される各タスクの `due_datetime`（`YYYY-MM-DDTHH:mm:00+09:00`）へのJST正規化に関する詳細仕様が不足しています。

## 2. 詳細な指摘内容
1. **`recurring_rule.due_time` の入力フォーマット制約の曖昧さ**:
   - `04_tasks.md` L122 の `recurring_rule.due_time` フィールド定義には `締切時刻 HH:mm（省略時は 23:59）` とのみ記載されており、24時間表記（`00:00`〜`23:59`）の範囲外の値（例: `"25:00"`, `"12:60"`）や秒数を含む形式（`"18:00:00"`）が送信された場合のバリデーションルールが明記されていません。
2. **エラーレスポンスにおけるフィールド名のドット表記**:
   - L153 にはネストされたフィールドのエラー表記としてドット記法（`"recurring_rule.start_date"` 等）が追記されていますが、`due_time` の形式エラー発生時についての具体的なキー名が明示されていません。
3. **生成タスクの `due_datetime` へのJST正規化仕様の欠落**:
   - 繰り返し一括作成時に、各日付と `due_time`（省略時 `23:59`）が結合され、タイムゾーンオフセット `+09:00`（JST）が付与されて ISO 8601 形式 `YYYY-MM-DDTHH:mm:00+09:00` として登録・返却される仕様が明確に記述されていません。

## 3. 推奨される修正案
`04_tasks.md` 3.3.2 節 (`POST tasks`) の `recurring_rule.due_time` フィールド定義および補足説明・Errors セクションを以下のように更新してください：

1. **フィールド定義の修正 (`recurring_rule.due_time`)**:
   - 24時間表記の `HH:mm`（`00:00`〜`23:59`）形式必須。省略時は `23:59` を適用。
2. **リクエスト評価順序・Errors の更新**:
   - `due_time` が `HH:mm` 形式でない場合、または `00:00`〜`23:59` の範囲外の場合は 400 Bad Request（`code: "BAD_REQUEST"`, `error.details: [{ "field": "recurring_rule.due_time", "message": "締切時刻の形式が不正です（HH:mm形式で指定してください）" }]`）を返却する。
3. **生成されるタスクの属性に関する注記**:
   - 繰り返し一括生成される各タスクの `due_datetime` には、該当日に `due_time`（省略時は `23:59`）と JST オフセット `+09:00` が結合された ISO 8601 日時文字列（例: `"2026-08-22T18:00:00+09:00"`）が自動的に設定されて登録・返却される旨を明記する。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 17:27:35
- **Status**: Resolved

### 実施した修正内容
`04_tasks.md` 3.3.2 節 (`POST tasks`) において以下の通り修正を行いました：
1. `recurring_rule.due_time` のフィールド定義を「24時間表記の HH:mm（00:00〜23:59、省略時は 23:59）。HH:mm 形式でない場合、または 00:00〜23:59 の範囲外の数値や秒数を含む形式が指定された場合は 400 Bad Request」と明確化。
2. 繰り返し一括生成されるタスクの `due_datetime` に該当日の日付と `due_time`（省略時 23:59）に JST オフセット `+09:00` が結合された ISO 8601 日時文字列が自動設定される旨を注記。
3. Errors セクションの 400 Bad Request 内に `recurring_rule.due_time` 形式不正時の `error.details` 例を追記。

### 変更したファイル
- [04_tasks.md](file:///mnt/c/Users/ayumu_wkkd3w3/BoxForSomething/github/SyncTask/docs/design/api_design/04_tasks.md)
