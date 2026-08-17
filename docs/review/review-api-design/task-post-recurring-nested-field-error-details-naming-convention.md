# POST tasks における recurring_rule ネストフィールドバリデーションエラー時の error.details[].field 命名規約の未定義

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-17 17:07:00
- **Target Files**:
  - [04_tasks.md](docs/design/api_design/04_tasks.md)

## 1. 問題の概要
`POST tasks` (3.3.2 節) の繰り返し一括作成処理において、`recurring_rule` オブジェクト内の特定フィールド（例: `start_date`, `end_date`, `days_of_week`, `due_time`）の入力不備によりエラーが発生した際、エラーレスポンス `error.details[].field` に格納されるキー名の命名規約がドキュメント上曖昧です。

## 2. 詳細な指摘内容
1. **ネストフィールドにおける `field` 表現の揺れ**:
   - `04_tasks.md` L148-149 の Errors 例では、生成件数オーバーや該当曜日なしに対して `"field": "recurring_rule"` と記載されています。
   - 一方で、`recurring_rule.start_date` の日付形式違反や `start_date > end_date` 違反が発生した場合、`error.details[].field` がドット記法（例: `"recurring_rule.start_date"`）となるのか、あるいは末尾フィールド名（例: `"start_date"`）または親オブジェクト名（例: `"recurring_rule"`）となるのかが明確化されていません。
2. **フロントエンドのエラーハイライト連動不具合リスク**:
   - フォームの個別の入力フィールド（開始日入力欄、曜日選択UI等）にバリデーションエラーメッセージを紐付ける際、バックエンドからの `field` キー名が統一されていないとフロントエンド側のバインド処理が正常に動作しない恐れがあります。

## 3. 推奨される修正案
`04_tasks.md` 3.3.2 節の Errors または共通エラー仕様注記において、`recurring_rule` 内のネストされた個別フィールドでバリデーションエラーが発生した場合の `field` 名のフォーマット規約（例: `"recurring_rule.start_date"`, `"recurring_rule.days_of_week"` 等のドット記法）を明記してください。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 17:10:00
- **Status**: Resolved

### 実施した修正内容
`04_tasks.md` 3.3.2 節 (`POST tasks`) の Errors 注記において、`recurring_rule` オブジェクト内のネストされた個別フィールド（`start_date`, `days_of_week` 等）でバリデーションエラーが発生した場合の `error.details[].field` キー名はドット記法（例: `"recurring_rule.start_date"`）となる命名規約を明確に規定しました。

### 変更したファイル
- [04_tasks.md](docs/design/api_design/04_tasks.md)
