# TASKテーブルのTITLEおよびCOMMENTにおけるデータ型・文字数制約と要件定義書の乖離

- **Status**: Open
- **Severity**: Minor
- **Created At**: 2026-08-18 22:15:00
- **Target Files**:
  - [database_design.md](docs/design/database_design.md)
  - [requirements.md](docs/req-def/requirements.md)

## 1. 問題の概要
`docs/design/database_design.md` の「2. タスク管理 (TASK)」において、`TITLE`（タスク名）が `VARCHAR(255)`、`COMMENT`（コメント）が `TEXT` として定義されていますが、要件定義書（`requirements.md`）で規定された上限文字数（タイトル: 100文字、コメント: 1000文字）とデータ型のサイズ制限に乖離があります。

## 2. 詳細な指摘内容
1. **要件定義書の仕様 (`docs/req-def/requirements.md` L80-81, L97, L100)**:
   - タスク名: 1〜100文字必須（空白のみ・制御文字不可）
   - コメント: 0〜1000文字（任意）

2. **DB設計書の定義 (`docs/design/database_design.md` L47, L52)**:
   ```markdown
   | タスク名 | `TITLE` | `VARCHAR(255)` / `NOT NULL` | 1〜100文字（制御文字不可） |
   | コメント | `COMMENT` | `TEXT` | 補足メモ（0〜1000文字） |
   ```

### 問題点：
- `TITLE` が `VARCHAR(255)` となっているため、DB層では101〜255文字の文字列が保存可能となり、要件定義の上限（100文字）を超えるデータが混入する恐れがあります。
- `COMMENT` が `TEXT`（無制限）となっているため、DB層で1000文字を超える巨大なテキストが保存されるリスクを防止できません。

## 3. 推奨される修正案
テーブル定義のカラム型を要件定義書の上限文字数に合わせて修正するか、CHECK 制約を追加してDBレベルでのデータ整合性を保証してください：

```markdown
| 項目名 | カラム名 | データ型 / 制約 | 備考 |
| --- | --- | --- | --- |
| タスク名 | `TITLE` | `VARCHAR(100)` / `NOT NULL` | 1〜100文字（制御文字不可） |
| コメント | `COMMENT` | `VARCHAR(1000)` | 補足メモ（0〜1000文字） |
```
または `CHECK (length(COMMENT) <= 1000)` を制約として明記してください。
