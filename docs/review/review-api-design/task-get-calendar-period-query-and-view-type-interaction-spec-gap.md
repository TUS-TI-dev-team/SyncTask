# GET tasks におけるカレンダー期間取得（start_date / end_date）と view_type・due_datetime null タスクの相互作用仕様の未定義

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-17 16:53:00
- **Target Files**:
  - [04_tasks.md](docs/design/api_design/04_tasks.md)

## 1. 問題の概要
`GET tasks` API において、カレンダー期間指定用クエリパラメータ `start_date` / `end_date` と他の絞り込み条件（`view_type`, `status` 等）を同時に指定した場合の挙動、および締切日時未設定（`due_datetime: null`）のタスクの除外仕様についての明確な記述が欠落しています。

## 2. 詳細な指摘内容
`04_tasks.md` L20-L21 には、`start_date` と `end_date` について以下の説明があります：
> `start_date`: カレンダー表示用: グリッド取得開始日（`YYYY-MM-DD`）。`end_date` とペアで指定必須。指定時はページネーション limit を解除し期間内の全タスクを返却。最大許容期間幅は 42日間（6週間）  
> `end_date`: カレンダー表示用: グリッド取得終了日（`YYYY-MM-DD`）。`start_date` とペアで指定必須（`start_date <= end_date`）

しかし、以下のケースでの挙動が明確化されていません：
1. **締切未設定（`null`）タスクの扱い**: カレンダー表示はグリッド上の各日付（`YYYY-MM-DD`）にタスクを割り当てるための API であるため、`due_datetime` が `null` のタスクは当然検索結果（`start_date 00:00:00+09:00 <= due_datetime <= end_date 23:59:59+09:00`）から除外されますが、その点が明記されていません。
2. **`view_type` や `status` との併用**: `start_date` / `end_date` を指定しながら `view_type=high_priority` や `view_type=pinned` を指定した場合に、期間内かつ該当ビューのタスクのみがフィルタリング返却されるのか、あるいは併用不可として `400 Bad Request` となるのかが未定義です。

## 3. 推奨される修正案
`04_tasks.md` 3.3.1 のレスポンス補足注記（L50 付近）に、以下の内容を追加明記してください。

```markdown
※ `start_date` / `end_date` を指定したカレンダー期間取得時は、`due_datetime` が設定されているタスクのうち `start_date 00:00:00+09:00 <= due_datetime <= end_date 23:59:59+09:00` の範囲に該当するタスクのみが抽出返却されます（`due_datetime` が `null` のタスクは除外されます）。なお、`start_date` / `end_date` 指定時に `view_type`, `priority`, `status`, `keyword` 等の絞り込みパラメータが併用された場合は、指定された期間内で該当する絞り込み条件を満たすタスクのみが一括返却されます。
```

## 修正完了報告

- **Resolved At**: 2026-08-17 16:56:00
- **Status**: Resolved

### 実施した修正内容
`04_tasks.md` の 3.3.1 (`GET tasks`) の補足注記に、`start_date` / `end_date` 指定のカレンダー期間取得時において `due_datetime` が `null` のタスクは対象外となること、および `view_type`, `priority`, `status`, `keyword` 等のパラメータと併用して AND 条件で絞り込めることを明記しました。

### 変更したファイル
- [04_tasks.md](docs/design/api_design/04_tasks.md)
