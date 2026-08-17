# タスク詳細取得APIの欠落およびステータス・ピン留め更新APIの設計不備

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-17 12:05:00
- **Target Files**:
  - [api_design.md](docs/design/api_design.md)
  - [screen_design.md](docs/design/screen_design.md)

## 1. 問題の概要
タスク編集画面の初期表示やURL直接指定アクセス時に必要な個別タスク詳細取得API（`GET tasks/{task_id}`）が欠落しています。また、画面上でのステータス直接変更やピン留めトグル操作に対して、全体更新（`PUT`）しか定義されておらず、部分更新（PATCH）または更新対象フィールドの仕様が不明確です。

## 2. 詳細な指摘内容
1. **個別タスク取得API（`GET tasks/{task_id}`）の欠落**:
   - `docs/design/screen_design.md` L26 のタスク編集画面やL40の日付詳細ポップアップからのタスク詳細遷移において、対象タスクの最新状態を単一取得するエンドポイントが必要です。
   - `docs/design/api_design.md` には `GET tasks/`（一覧）しか定義されておらず、単一リソース取得エンドポイントがありません。
2. **ステータス・ピン留めの部分更新設計**:
   - `screen_design.md` L30-31, L41 では、一覧画面やカレンダー画面の日付セル/ポップアップ上で「ボタン1つで状態変更」「ピン留めトグル」を行うUIが規定されています。
   - `api_design.md` L29 では `PUT tasks/{task_id}` （変更内容、CSRFトークン）と大まかに書かれており、全フィールド送信が必要な `PUT` なのか、部分更新（`PATCH`）を許容するのか、あるいはステータス/ピン留め専用更新エンドポイント（例: `PATCH tasks/{task_id}` で特定フィールドのみ更新）とするのかが曖昧です。

## 3. 推奨される修正案
1. 個別タスク取得API `GET tasks/{task_id}` を追加してください。
2. タスク更新において、全項目更新（`PUT tasks/{task_id}`）と部分更新（`PATCH tasks/{task_id}`、または `PATCH` への統一）を明確に定義し、ステータス（`status`）やピン留め（`is_pinned`）のみの更新ペイロード仕様を明記してください。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 12:40:00
- **Status**: Resolved

### 実施した修正内容
- `docs/design/api_design.md` に単一タスク詳細取得用 `GET tasks/{task_id}` エンドポイントを追加定義しました。
- `PUT tasks/{task_id}` において、全項目一括更新に加え、リクエストボディに含まれるフィールドのみを更新する部分更新兼用仕様（ステータス `status` やピン留め `is_pinned` の単体更新含む）を明記しました。

### 変更したファイル
- [api_design.md](docs/design/api_design.md)
