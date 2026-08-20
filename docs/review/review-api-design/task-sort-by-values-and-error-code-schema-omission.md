# GET tasks の sort_by 許容値定義の不足および各エンドポイントのエラーコード指定の欠落

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-17 16:30:00
- **Target Files**:
  - [04_tasks.md](docs/design/api_design/04_tasks.md)
  - [01_overview.md](docs/design/api_design/01_overview.md)

## 1. 問題の概要
`04_tasks.md` において、`GET tasks` の `sort_by` パラメータの `default` 以外の許容値（選択肢）が未定義である点、および 3.3.1〜3.3.5 の全エンドポイントのエラー仕様において `01_overview.md` で規定されている構造化エラーコード文字列 (`code`) が具体的に記載されていない点についての不備があります。

## 2. 詳細な指摘内容
1. **`sort_by` パラメータの許容値・列挙値の不足**:
   - `04_tasks.md` L22 では `sort_by` のデフォルト値として `default`（ピン留め優先 ➔ 締切昇順 ➔ 作成日時降順）が挙げられていますが、他に使用可能なソート指定子（例: `due_date_asc`, `created_at_desc`, `priority_desc` 等）が存在するのか、あるいは `default` のみが唯一サポートされている値なのかが記載されていません。
   - クライアント側でソート条件の切り替えパラメータを送信する際のスキーマが不明確です。

2. **エラー仕様におけるエラーコード (`code`) 文字列の欠落**:
   - `01_overview.md` Section 1.3 では、すべてのエラーレスポンスが JSON 構造 `{ "error": { "code": "...", "message": "..." } }` で返却され、`BAD_REQUEST`, `UNAUTHORIZED`, `FORBIDDEN`, `NOT_FOUND`, `UNPROCESSABLE_ENTITY` などのエラーコード文字列が規定されています。
   - しかし `04_tasks.md` の全エンドポイント（3.3.1 Errors, 3.3.2 Errors, 3.3.3 Errors, 3.3.4 Errors, 3.3.5 Errors）の Errors セクションでは、`- 400 Bad Request: クエリパラメータ不正` や `- 404 Not Found: 認可エラー` のような定型文の簡略箇条書きのみとなっており、具体的にどの `code` 文字列が返却されるのか（例: `BAD_REQUEST`, `UNAUTHORIZED`, `NOT_FOUND`, `FORBIDDEN`）が明記されていません。
   - また、繰り返しタスク作成時の件数超過（0件または101件以上）や期間不整合発生時のビジネスルールエラーコード（`UNPROCESSABLE_ENTITY` または `BAD_REQUEST`）と詳細メッセージの関連付けが欠落しています。

## 3. 推奨される修正案
1. `GET tasks` の `sort_by` パラメータ仕様に、指定可能な列挙値（Enum値）の一覧を明確に定義してください（`default` のみがサポート対象である場合はその旨を明記）。
2. 3.3.1〜3.3.5 のすべての `Errors` セクションにおいて、`01_overview.md` に準拠したレスポンス形式および具体的な `code` 値（例: `"code": "BAD_REQUEST"`, `"code": "NOT_FOUND"`）を併記した定義に更新してください。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 16:42:00
- **Status**: Resolved

### 実施した修正内容
`GET tasks` の `sort_by` パラメータに対し指定可能な Enum 値の一覧（`default`, `due_date_asc`, `due_date_desc`, `created_at_desc`, `priority_desc`）を明記し、3.3.1〜3.3.5 の全 Errors セクションに具体的な `code` 値（`"BAD_REQUEST"`, `"UNAUTHORIZED"`, `"FORBIDDEN"`, `"NOT_FOUND"`, `"UNPROCESSABLE_ENTITY"`）を追記しました。

### 変更したファイル
- [04_tasks.md](file:////mnt/c/Users/ayumu_wkkd3w3/BoxForSomething/github/SyncTask/docs/design/api_design/04_tasks.md)
