# POST tasks および GET tasks における priority / status / view_type の不正列挙値指定時のバリデーション仕様欠落

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-17 17:33:00
- **Target Files**:
  - [04_tasks.md](docs/design/api_design/04_tasks.md)

## 1. 問題の概要
`POST tasks` (3.3.2 節) の `priority` フィールドへの `null` または定義外の列挙値（例: `"URGENT"`）の送信、および `GET tasks` (3.3.1 節) の `priority`, `status`, `view_type`, `sort_by` クエリパラメータに不正な値が指定された場合のエラー挙動（400 Bad Request / `code: "BAD_REQUEST"`）に関する明記が不足しています。

## 2. 詳細な指摘内容
1. **`POST tasks` における `priority` の `null` / 不正値検証**:
   - `POST tasks` では `is_pinned` について「明示的に `null` または非 boolean 型が指定された場合は 400 Bad Request (`code: "BAD_REQUEST"`)」と明確に規定されていますが、同じリクエストボディ内の `priority` フィールドに対する `null` や未知の列挙値（例: `"URGENT"`）送信時のエラーハンドリングルールが具体的に記載されていません。

2. **`GET tasks` のクエリパラメータにおける Enum 値検証**:
   - `GET tasks` のリクエスト評価順序 2（入力バリデーション）では `page`, `limit`, `keyword` 長さ, 日付フォーマットの違反は 400 `BAD_REQUEST` となる旨が挙げられていますが、`priority` (`high`, `medium`, `low`), `status` (`not_started`, `in_progress`, `completed`), `view_type` (`high_priority`, `near_deadline`, `pinned`), `sort_by` の各 Enum パラメータに定義外の文字列が指定された場合のバリデーションルール記述が漏れています。

## 3. 推奨される修正案
1. `POST tasks` のフィールド定義およびリクエスト評価順序 2 に、「`priority` フィールドに `null` または定義外の列挙値（`high`, `medium`, `low` 以外）が指定された場合は 400 Bad Request (`code: "BAD_REQUEST"`) を返却する」旨を追記してください。
2. `GET tasks` のリクエスト評価順序 2 に、「`priority`, `status`, `view_type`, `sort_by` パラメータに定義外の無効な値指定時は 400 Bad Request (`code: "BAD_REQUEST"`) を返却する」旨を明確化してください。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 17:38:45
- **Status**: Resolved

### 実施した修正内容
`docs/design/api_design/04_tasks.md` において以下を修正・追記しました：
1. `POST tasks` (3.3.2 節) のリクエスト評価順序 2、フィールド定義、および Errors (400 Bad Request) に、`priority` フィールドへの `null` や定義外列挙値指定時に 400 `BAD_REQUEST` を返却する旨を明記。
2. `GET tasks` (3.3.1 節) の Query Parameters テーブル、リクエスト評価順序 2、および Errors (400 Bad Request) に、`priority`, `status`, `view_type`, `sort_by` の各 Enum クエリパラメータへ無効な値が指定された場合に 400 `BAD_REQUEST` を返却する旨を明記。

### 変更したファイル
- [04_tasks.md](docs/design/api_design/04_tasks.md)

