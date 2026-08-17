# PATCH tasks/{task_id} における非 Null 許容フィールドへの null 指定時のバリデーション挙動未定義

- **Status**: Resolved

---

## 修正完了報告

- **Resolved At**: 2026-08-17 17:03:00
- **Status**: Resolved

### 実施した修正内容
`04_tasks.md` 3.3.4 節 (`PATCH tasks/{task_id}`) の注記およびリクエスト評価順序において、非 Null 許容フィールド（`title`, `priority`, `status`, `is_pinned`）に対して明示的に `null` が指定された場合は `400 Bad Request`（code: `"BAD_REQUEST"`）を返却する旨を明記しました。

### 変更したファイル
- [04_tasks.md](docs/design/api_design/04_tasks.md)
- **Severity**: Minor
- **Created At**: 2026-08-17 17:01:00
- **Target Files**:
  - [04_tasks.md](docs/design/api_design/04_tasks.md)

## 1. 問題の概要
`PATCH tasks/{task_id}`（3.3.4）のフィールド定義表において、`comment` および `due_datetime` は `string / null` 型として明示的に `null` 指定によるクリア挙動が定義されていますが、非 Null 許容フィールドである `title`, `priority`, `status`, `is_pinned` に対して明示的に `null`（例: `{"title": null}`, `{"status": null}`）が送信された場合のバリデーション挙動が未定義です。

## 2. 詳細な指摘内容
`04_tasks.md` L199-206 のフィールド定義において、`title` は `string`（1〜100文字）、`priority` は `string` (`high`, `medium`, `low`)、`status` は `string` (`not_started`, `in_progress`, `completed`)、`is_pinned` は `boolean` と規定されています。

クライアントがこれらのキーに対して `null` を指定して送信した場合に：
- 無視されて更新されないのか
- `400 Bad Request`（code: `"BAD_REQUEST"`）として明確に拒否されるのか
が不透明です。DBスキーマ上 `TITLE`, `PRIORITY`, `STATUS`, `IS_PINNED` はいずれも `NOT NULL` 制約が設定されているため、`null` が渡された場合はバックエンドで `400 Bad Request` エラーを即座に返却する必要があります。

## 3. 推奨される修正案
`PATCH tasks/{task_id}` のフィールド定義表または注記に、「非 Null 許容フィールド（`title`, `priority`, `status`, `is_pinned`）に対して `null` が指定された場合は 400 Bad Request（code: `"BAD_REQUEST"`）を返却する」旨を明確に記載してください。
