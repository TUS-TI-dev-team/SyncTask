# MAIL_AUTH_LOGテーブル定義における改行混入によるMarkdownテーブル構文の崩れ

- **Status**: Open
- **Severity**: Minor
- **Created At**: 2026-08-18 22:05:00
- **Target Files**:
  - [database_design.md](docs/design/database_design.md)

## 1. 問題の概要
`docs/design/database_design.md` の「6.3 メール認証ログ (MAIL_AUTH_LOG)」テーブル定義において、カラム名 `ACCESS_AT` の途中に不正な改行コードが混入しており、Markdown テーブルの構文が崩れて表示およびパーサー処理に支障をきたしています。

## 2. 詳細な指摘内容
`database_design.md` の L172-173 に以下の記述があります：

```markdown
| ダミー処理区分 | `IS_DUMMY` | `BOOLEAN` / `NOT NULL` | `TRUE` (ダミー表示) / `FALSE` (実処理) |
| アクセス日時 | `ACCESS_A
T` | `TIMESTAMPTZ` / `NOT NULL` | ログ記録日時（インデックス対象） |
```

`ACCESS_AT` の文字列の `ACCESS_A` と `T` の間で改行が発生しているため、Markdown テーブルの行が意図せず分割され、テーブルとして正しくレンダリングされません。

## 3. 推奨される修正案
改行を除去し、1行でカラム名が `ACCESS_AT` となるよう修正してください：

```markdown
| アクセス日時 | `ACCESS_AT` | `TIMESTAMPTZ` / `NOT NULL` | ログ記録日時（インデックス対象） |
```
