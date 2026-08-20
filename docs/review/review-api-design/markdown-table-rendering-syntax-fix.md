# Markdownプレビューにおけるテーブル記法のレンダリング崩れ

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-17 13:17:00
- **Target Files**:
  - [api_design.md](docs/design/api_design.md)

## 1. 問題の概要
`docs/design/api_design.md` の L513 付近（`GET tasks` の Query Parameters）をはじめとする複数の箇所で、Markdown プレビュー（Antigravity / VS Code / 各種Markdownレンダラー）においてテーブルが正しく表形式としてレンダリングされず、平文テキストとして崩れて表示される問題が発生しています。

## 2. 詳細な指摘内容
1. **L512-527（`GET tasks` のクエリパラメータ表）の構文問題**:
   - リスト項目 `- **Query Parameters**:` の直後に空行を挟んでインデントなしのテーブル構文（`| パラメータ名 | ...`）が配置されています。
   - CommonMark / GFM 準拠のパーサーにおいて、箇条書きリストブロックの途中に突然テーブルが挿入されると、リスト継続判定が競合し、テーブルブロックとして認識されずにテキスト行として描画されてしまいます。
2. **他のエンドポイントにおけるテーブル定義形式の混在**:
   - `3.1.1 POST auth/register/request-otp`（L107）、`3.1.2 POST auth/register/verify-otp`（L138）、`3.2.2 PUT users/{user_id}`（L432）、`3.3.2 POST tasks`（L577）、`3.3.4 PUT tasks/{task_id}`（L663）等でも同様に、箇条書きリスト項目（`- **Request Body**:` 等）の直下にテーブルが置かれており、パーサー依存の描画崩れリスクがあります。

## 3. 推奨される修正案
1. `- **Query Parameters**:` や `- **Request Body**:` などの箇条書きリスト記法（`- `）を廃止し、小見出し（例: `##### Query Parameters`, `##### Request Body フィールド定義`）または独立した段落見出しへ統一してください。
2. テーブルブロックの前後には必ず空行を配置し、標準的な GFM テーブルとして確実にレンダリングされるドキュメント構造に修正してください。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 13:22:00
- **Status**: Resolved

### 実施した修正内容
- `docs/design/api_design.md` 内のすべてのエンドポイントにおいて、箇条書きリスト記法（`- **Query Parameters**:` や `- **Request Body**:`）を小見出し（`##### Query Parameters`、`##### Request Body フィールド定義` など）に統一しました。
- すべてのテーブルブロックの前後に確実に空行を配置し、CommonMark / GFM 準拠のパーサーおよび Antigravity プレビューで確実にテーブル形式としてレンダリングされる構造に修正しました。

### 変更したファイル
- [api_design.md](docs/design/api_design.md)
