# POST tasksにおける毎週繰り返し一括作成の days_of_week バリデーションおよび生成レスポンスソート順の未定義

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-17 16:35:00
- **Target Files**:
  - [04_tasks.md](docs/design/api_design/04_tasks.md)

## 1. 問題の概要
`POST tasks` において毎週繰り返し一括作成（`is_recurring: true`）を行う際の `recurring_rule.days_of_week` 配列に対するバリデーション仕様（空配列・重複要素・大文字小文字表記）、およびレスポンス (201 Created) 内で返却される一括作成タスク配列 `tasks` のソート順が定義されていません。

## 2. 詳細な指摘内容
1. **`days_of_week` の境界値・表記バリデーションの不足**:
   - `recurring_rule.days_of_week` に空配列 `[]` が渡された場合のエラー判定（`400 Bad Request` / `code: BAD_REQUEST`）が明記されていません。
   - 配列内に同一の曜日が重複して含まれていた場合（例: `["monday", "monday"]`）の扱い（重複を排除して処理するか、バリデーションエラーとするか）が未定義です。
   - 曜日の文字列比較において、大文字・小文字表記（例: `["Monday", "MONDAY"]`）の許容性および小文字正規化の規約が記載されていません。

2. **レスポンス `tasks` 配列のソート順の未規定**:
   - 繰り返し一括生成成功時に `201 Created` レスポンスで返却される `tasks` 配列（最大100件）の並び順（生成日時の古い順 / `due_datetime` の昇順）が指定されておらず、クライアント側で受領後の表示順が不安定になる懸念があります。

## 3. 推奨される修正案
1. `POST tasks` の `recurring_rule.days_of_week` 仕様に以下を追記してください:
   - 1つ以上の要素が必須であり、空配列 `[]` または無効な曜日文字列が含まれる場合は `400 Bad Request` を返却する。
   - 受信した曜日文字列は小文字へ正規化され、重複要素は自動的にユニーク化（デデュプリケーション）されて処理される。
2. `POST tasks` のレスポンス仕様注記に「返却される `tasks` 配列は、`due_datetime` の昇順（時系列順）でソートされて返却されます」と明記してください。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 16:42:00
- **Status**: Resolved

### 実施した修正内容
`POST tasks` の `recurring_rule.days_of_week` に対し、空配列禁止・小文字正規化・デデュプリケーション（重複排除）・無効文字列時 400 エラーのバリデーション規則を追加し、一括生成レスポンスの `tasks` 配列が `due_datetime` 昇順でソートされる仕様を明確化しました。

### 変更したファイル
- [04_tasks.md](file:////mnt/c/Users/ayumu_wkkd3w3/BoxForSomething/github/SyncTask/docs/design/api_design/04_tasks.md)
