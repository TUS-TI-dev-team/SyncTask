# PATCH tasks/{task_id} における due_datetime の日付のみ指定（YYYY-MM-DD）サポートの記述抜け

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-17 16:45:00
- **Target Files**:
  - [04_tasks.md](docs/design/api_design/04_tasks.md)

## 1. 問題の概要
`PATCH tasks/{task_id}` の `due_datetime` 更新において、`POST tasks` と同様に日付のみ（`YYYY-MM-DD`）のフォーマットを許容しデフォルト時刻（`23:59:00+09:00`）を補完するのかについての仕様記述が欠落しています。

## 2. 詳細な指摘内容
新規タスク作成エンドポイント `POST tasks` のフィールド定義（`04_tasks.md` L96）では、`due_datetime` について以下の通り明記されています。
> ISO 8601 日時文字列（例: `2026-08-20T23:59:00+09:00`）、または日付のみ `YYYY-MM-DD`（時刻省略時は `23:59:00+09:00` を設定）。

しかし、タスク部分更新エンドポイント `PATCH tasks/{task_id}` のフィールド定義（`04_tasks.md` L197）では以下のようにのみ記載されています。
> ISO 8601 日時文字列（`null` 指定で締切解除）

フロントエンドのUIで日付ピッカーから締切日のみを変更する場合に `YYYY-MM-DD` 形式の文字列が送信される可能性がありますが、`PATCH` エンドポイントがこの形式を受け付けるのか、あるいはバリデーションエラー（`400 Bad Request`）とするのかが不透明です。また、バックエンド実装者によって補完時刻（`23:59:00+09:00`）の取り扱いが不一致となるリスクがあります。

## 3. 推奨される修正案
`PATCH tasks/{task_id}` の `due_datetime` リクエストボディ定義（`04_tasks.md` L197）を以下のように更新し、`POST tasks` と一貫性を持たせてください。

```markdown
| `due_datetime` | string / null | × | ISO 8601 日時文字列、または日付のみ `YYYY-MM-DD`（時刻省略時は `23:59:00+09:00` を設定。`null` 指定で締切解除） |
```

---

## 修正完了報告

- **Resolved At**: 2026-08-17 16:50:00
- **Status**: Resolved

### 実施した修正内容
`04_tasks.md` 3.3.4 (`PATCH tasks/{task_id}`) の `due_datetime` パラメータ定義を更新し、日付のみ `YYYY-MM-DD` 形式入力時にデフォルト時刻（`23:59:00+09:00`）を自動補完する仕様を明記しました。

### 変更したファイル
- [04_tasks.md](docs/design/api_design/04_tasks.md)
