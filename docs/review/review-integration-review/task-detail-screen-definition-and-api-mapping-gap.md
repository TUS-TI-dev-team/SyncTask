# タスク詳細画面の画面設計書上の定義欠落とタスク詳細取得APIの利用対応の不整合

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-22 13:50:00
- **Target Files**:
  - [screen_design.md](docs/design/screen_design.md)
  - [04_tasks.md](docs/design/api_design/04_tasks.md)
  - [02_task_management.md](docs/req-def/requirements/02_task_management.md)

## 1. 問題の概要
API設計書（`04_tasks.md`）にはタスク個別取得用の `GET tasks/{task_id}` が定義され、要件定義書（`02_task_management.md`）にも「タスク詳細/編集画面への遷移」等の言及があるが、画面設計書（`screen_design.md`）のタスク関連画面一覧には「タスク詳細画面（または詳細ポップアップ）」が定義されておらず、「タスク編集ポップアップ」のみが定義されている。単体タスク詳細の閲覧UIの有無や、`GET tasks/{task_id}` がどの画面・ポップアップから呼び出されるかの対応関係が明記されていない。

## 2. 詳細な指摘内容
1. **API設計書の定義**:
   - `docs/design/api_design/04_tasks.md` 3.3.3 節に `GET tasks/{task_id}`（タスク詳細取得API）が定義されている。
2. **要件定義書の記述**:
   - `docs/req-def/requirements/02_task_management.md` 3.3.3 節において「クリックによるタスク詳細/編集画面への遷移」と言及されている。
3. **画面設計書の定義欠落**:
   - `docs/design/screen_design.md` のタスク関連画面定義（line 31-35）には、「タスク一覧画面」「タスク作成ポップアップ」「タスク編集ポップアップ」の3つしか記載されておらず、「タスク詳細画面」または「タスク詳細ポップアップ」という画面要素が存在しない。
   - タスク編集ポップアップが詳細閲覧の役割を兼ねているのか、あるいは一覧取得済みキャッシュを利用せずに編集ポップアップ展開時に `GET tasks/{task_id}` を呼び出すのかについての設計仕様が明示されていない。

## 3. 推奨される修正案
- `docs/design/screen_design.md` において以下のいずれかを明確化する：
  1. タスク詳細の単体閲覧専用画面は設けず、「タスク編集ポップアップ」が閲覧兼編集の役割を担うこと。
  2. 「タスク編集ポップアップ」を開く際（一覧クリック時、日付詳細ポップアップからの遷移時、または直接リンク遷移時）、`GET /api/tasks/{task_id}` を実行して最新のタスク情報を取得・表示すること。
- `docs/req-def/requirements/02_task_management.md` 3.3.3 節の「タスク詳細/編集画面への遷移」を「タスク編集ポップアップの表示」と表記を統一する。

---

## 修正完了報告

- **Resolved At**: 2026-08-22 13:55:00
- **Status**: Resolved

### 実施した修正内容
タスク単体の詳細閲覧専用画面は設けず「タスク編集ポップアップ」が閲覧兼編集の役割を担うこと、およびポップアップ展開時に `GET tasks/{task_id}` を呼び出して最新情報を取得・表示する仕様を明記しました。また要件定義書の表記を「タスク編集ポップアップの表示」へ統一しました。

### 変更したファイル
- [screen_design.md](file:///C:\Users\kazuh\Programming\repos\SyncTask\docs\design\screen_design.md)
- [04_tasks.md](file:///C:\Users\kazuh\Programming\repos\SyncTask\docs\design\api_design\04_tasks.md)
- [02_task_management.md](file:///C:\Users\kazuh\Programming\repos\SyncTask\docs\req-def\requirements\02_task_management.md)
