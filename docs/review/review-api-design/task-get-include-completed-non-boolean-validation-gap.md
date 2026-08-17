# GET tasks における include_completed クエリパラメータの非 boolean 指定時バリデーション仕様の欠落

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-17 17:45:00
- **Target Files**:
  - [04_tasks.md](docs/design/api_design/04_tasks.md)

## 1. 問題の概要
`04_tasks.md` の `GET tasks` において、`page`, `limit`, `priority`, `status`, `view_type`, `sort_by`, `keyword`, `start_date`, `end_date`, `due_date` などの各パラメータでは「型違反・定義外の値指定時は 400 Bad Request（code: `"BAD_REQUEST"`）」と明記されているのに対し、`include_completed` パラメータのみ非 boolean 文字列（例: `include_completed=invalid` や `include_completed=123`）が指定された場合のエラーハンドリング定義が漏れています。

## 2. 詳細な指摘内容
1. **バリデーション挙動の一貫性欠落**:
   - `04_tasks.md` L15（Query Parameters 表）:
     `include_completed` | boolean | × | 通常: `false`, カレンダー: `true` | 完了タスクを含めるか（`true` / `false`）...
   - 他のパラメータ（`page`, `limit`, `priority`, `status` 等）では、定義外の値や型違反の場合に 400 Bad Request を返却する旨が注記されていますが、`include_completed` の表および「リクエスト評価順序」（L28）では非 boolean 値に対する挙動が規定されていません。
   - 不正な文字列が指定された場合、サーバーがエラーとせず `false` または `true` へ自動キャスト（ファジー処理）してしまうと、他の厳格バリデーションパラメータと一貫性が崩れます。

## 3. 推奨される修正案
`04_tasks.md` L15 の `include_completed` パラメータ説明および L28 の「リクエスト評価順序」に、非 boolean 値指定時に 400 Bad Request を返却する旨を明確に記述してください。

```markdown
| `include_completed`| boolean | × | 通常: `false`<br>カレンダー: `true` | 完了タスクを含めるか（`true` / `false`）。`true` または `false` 以外の文字列・数値が指定された場合は 400 Bad Request（code: `"BAD_REQUEST"`）を返却。通常一覧時はデフォルト `false`、`start_date` / `end_date` 指定（カレンダー表示）時はデフォルト `true`。なお `status` が明示指定された場合は本パラメータは無視されます |
```

---

## 修正完了報告

- **Resolved At**: 2026-08-17 17:50:00
- **Status**: Resolved

### 実施した修正内容
`04_tasks.md` 3.3.1 節 (`GET tasks`) の `include_completed` クエリパラメータ説明、リクエスト評価順序、および Errors セクションに、`true` または `false` 以外の非 boolean 文字列・数値が指定された場合は 400 Bad Request (code: `"BAD_REQUEST"`) を返却する旨を明記しました。

### 変更したファイル
- [04_tasks.md](docs/design/api_design/04_tasks.md)
