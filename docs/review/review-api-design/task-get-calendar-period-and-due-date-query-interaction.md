# GET tasks におけるカレンダー期間取得（start_date/end_date）と締切日絞り込み（due_date）併用時の挙動・相互作用の未定義

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-17 17:07:00
- **Target Files**:
  - [04_tasks.md](docs/design/api_design/04_tasks.md)

## 1. 問題の概要
`GET tasks` (3.3.1 節) において、カレンダー期間取得パラメータ（`start_date` / `end_date`）と単一日締切絞り込みパラメータ（`due_date`）が同時に指定された場合のフィルタロジックおよび優先関係が未定義です。

## 2. 詳細な指摘内容
1. **併用時の条件結合ロジックの不透明さ**:
   - `04_tasks.md` L61 の注記には、「`start_date` / `end_date` 指定時に `view_type`, `priority`, `status`, `keyword` 等の絞り込みパラメータが併用された場合は、指定された期間内で該当する絞り込み条件を満たすタスクのみが一括返却されます」と記載されていますが、`due_date` パラメータが明示的に除外されています。
   - 例えば、`GET /api/tasks?start_date=2026-08-01&end_date=2026-08-31&due_date=2026-08-15` というクエリが送信された場合：
     - ① カレンダー期間条件 (`start_date <= due_datetime <= end_date`) と `due_date` 条件 (`due_datetime <= due_date 23:59:59`) の AND 条件として評価されるのか
     - ② `due_date` が無視されてカレンダー期間取得が優先されるのか
     - ③ パラメータの非互換として `400 Bad Request` を返却するのか
     が規定されていません。
2. **`due_date` の特殊挙動とのコンフリクト**:
   - `due_date` パラメータは「過去の未完了・期限超過タスクを含み検索する」という独自の検索範囲仕様を持っているため、カレンダー期間グリッド（期間内のタスクのみ抽出）と仕様上の衝突が発生します。

## 3. 推奨される修正案
`04_tasks.md` 3.3.1 節に以下のいずれかの評価規則を明記してください：
- 案A: `start_date` / `end_date` 指定時は `due_date` パラメータの指定を不可とし、同時に指定された場合は 400 `BAD_REQUEST` を返却する。
- 案B: `start_date` / `end_date` 指定時に `due_date` が指定された場合は、カレンダー期間のフィルタ範囲（`start_date 00:00:00 <= due_datetime <= end_date 23:59:59`）かつ `due_datetime <= due_date 23:59:59` の AND 条件で抽出される（過去の期限超過タスクの自動包含挙動はカレンダー表示範囲内に限定される）旨を明記する。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 17:10:00
- **Status**: Resolved

### 実施した修正内容
ヒアリングに基づき「案A」を採用し、`04_tasks.md` 3.3.1 節 (`GET tasks`) の Query Parameters、リクエスト評価順序、注記および Errors において、カレンダー期間指定パラメータ（`start_date` / `end_date`）と締切日絞り込みパラメータ（`due_date`）の同時指定は不可とし、同時に指定された場合は 400 Bad Request（code: `"BAD_REQUEST"`）を返却する仕様を明記しました。

### 変更したファイル
- [04_tasks.md](docs/design/api_design/04_tasks.md)
