# GET tasks における page および limit の下限値（0以下・負数）および型不整合バリデーション仕様の不足

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-17 17:07:00
- **Target Files**:
  - [04_tasks.md](docs/design/api_design/04_tasks.md)

## 1. 問題の概要
`GET tasks` (3.3.1 節) において、ページネーション用クエリパラメータ `page` および `limit` に 0 以下の数値や不正な値が指定された場合のバリデーションおよびエラー返却ルールが明確に規定されていません。

## 2. 詳細な指摘内容
1. **下限値（`page < 1`, `limit < 1`）に対する検証ルールの欠落**:
   - `04_tasks.md` L28 の「リクエスト評価順序」では `limit 最大100件超過` に対する 400 Bad Request 検証が挙げられていますが、`page=0` や `page=-1`、`limit=0` や `limit=-10` のように 0 以下または負数のパラメータ値が渡された場合の挙動（400エラー返却か、デフォルト値 `page=1`, `limit=20` への自動補正か）が規定されていません。
2. **整数以外の型指定時の挙動不透明さ**:
   - `page=abc` や `limit=1.5` 等の非整数値が指定された場合のレスポンスが明記されていません。

## 3. 推奨される修正案
`04_tasks.md` 3.3.1 節の「Query Parameters」および「リクエスト評価順序」セクションを更新し、以下のバリデーション規則を明記してください：
- `page`: 1以上の整数（1未満の数値または非数値指定時は 400 `BAD_REQUEST`）
- `limit`: 1以上100以下の整数（1未満または100超の数値、あるいは非数値指定時は 400 `BAD_REQUEST`）

---

## 修正完了報告

- **Resolved At**: 2026-08-17 17:10:00
- **Status**: Resolved

### 実施した修正内容
`04_tasks.md` 3.3.1 節 (`GET tasks`) の Query Parameters およびリクエスト評価順序において、`page` (1以上の整数必須) と `limit` (1以上100以下の整数必須) の下限値・上限値および非数値・非整数入力に対する 400 Bad Request（code: `"BAD_REQUEST"`）バリデーションを明記しました。

### 変更したファイル
- [04_tasks.md](docs/design/api_design/04_tasks.md)
