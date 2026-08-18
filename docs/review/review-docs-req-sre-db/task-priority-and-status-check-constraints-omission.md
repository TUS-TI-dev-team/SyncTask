# TASKテーブルのPRIORITYおよびSTATUSにおけるCHECK制約定義の欠落

- **Status**: Open
- **Severity**: Minor
- **Created At**: 2026-08-18 22:05:00
- **Target Files**:
  - [database_design.md](docs/design/database_design.md)

## 1. 問題の概要
`docs/design/database_design.md` の「2. タスク管理 (TASK)」において、`PRIORITY` および `STATUS` カラムの許容値が固定の列挙値（ENUMドメイン）であるにもかかわらず、データ型が `VARCHAR(20)` のみで CHECK 制約が定義されておらず、不正値の混入を防ぐDBレベルの整合性担保が明記されていません。

## 2. 詳細な指摘内容
`database_design.md` の L48-50 に以下の定義があります：

| 項目名 | カラム名 | データ型 / 制約 | 備考 |
| --- | --- | --- | --- |
| 優先度 | `PRIORITY` | `VARCHAR(20)` / `NOT NULL, DEFAULT 'MEDIUM'` | `LOW` (低), `MEDIUM` (中・初期値), `HIGH` (高) |
| タスクステータス | `STATUS` | `VARCHAR(20)` / `NOT NULL, DEFAULT 'NOT_STARTED'` | `NOT_STARTED` (未着手・初期値), `IN_PROGRESS` (進行中), `COMPLETED` (完了) |

### 問題点：
- アプリケーション側のバリデーション漏れや直接SQL操作・マイグレーション時のヒューマンエラーにより、`'URGENT'` や `'IN-PROGRESS'`, `'pending'` といった未定義文字列が混入するリスクがあります。
- `VARCHAR(20)` のまま運用する場合、DB制約として `CHECK (PRIORITY IN ('LOW', 'MEDIUM', 'HIGH'))` および `CHECK (STATUS IN ('NOT_STARTED', 'IN_PROGRESS', 'COMPLETED'))` を明記しておくことで、データ完全性が保証されます。

## 3. 推奨される修正案
テーブル定義の「データ型 / 制約」欄に CHECK 制約の定義を追記してください：

```markdown
| 優先度 | `PRIORITY` | `VARCHAR(20)` / `NOT NULL, DEFAULT 'MEDIUM', CHECK (PRIORITY IN ('LOW', 'MEDIUM', 'HIGH'))` | `LOW` (低), `MEDIUM` (中・初期値), `HIGH` (高) |
| 締切日時 | `DUE_DATE` | `TIMESTAMPTZ` | 任意設定（未指定時は該当日 23:59 JST を適用） |
| タスクステータス | `STATUS` | `VARCHAR(20)` / `NOT NULL, DEFAULT 'NOT_STARTED', CHECK (STATUS IN ('NOT_STARTED', 'IN_PROGRESS', 'COMPLETED'))` | `NOT_STARTED` (未着手・初期値), `IN_PROGRESS` (進行中), `COMPLETED` (完了) |
```
