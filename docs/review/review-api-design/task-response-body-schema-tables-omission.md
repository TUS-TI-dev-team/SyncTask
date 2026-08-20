# タスク管理系 API におけるレスポンスボディ フィールド定義テーブルの全般的な欠落

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-17 17:33:00
- **Target Files**:
  - [04_tasks.md](docs/design/api_design/04_tasks.md)

## 1. 問題の概要
`04_tasks.md` に定義されている全 5 エンドポイント（`GET tasks`, `POST tasks`, `GET tasks/{task_id}`, `PATCH tasks/{task_id}`, `DELETE tasks/{task_id}`）において、レスポンスボディに含まれるオブジェクト（`task`, `tasks[]`, `items[]`, `pagination`）の各フィールド名、データ型、必須/Null許容性、および説明を明記した「レスポンスボディ フィールド定義テーブル」が記述されていません。

## 2. 詳細な指摘内容
1. **他仕様書（02_auth.md, 03_users.md）とのフォーマット不一致**:
   - `02_auth.md` や `03_users.md` では、レスポンスボディに含まれる全フィールドについて「フィールド名 | 型 | 必須 | 説明」の形式でテーブル化され、型や Null 許容性（`string / null` 等）が明示されています。
   - 一方、`04_tasks.md` では JSON のレスポンスサンプルと注記のみが記載されており、`task` リソースモデルの各属性（`id`, `user_id`, `title`, `comment`, `priority`, `status`, `due_datetime`, `is_pinned`, `created_at`, `updated_at`）や `pagination` モデルの型定義テーブルが存在しません。

2. **スキーマ決定・自動コード生成・テストにおける曖昧性**:
   - レスポンスサンプルのみに依存すると、OpenAPI や TypeScript 型定義、バックエンドDTOスキーマを作成する際に、どのフィールドが必須（省略不可）でどのフィールドが Null 許容（`null` の可能性があるか）が不透明となり、フロントエンド・バックエンド間の実装ミスの原因となります。

## 3. 推奨される修正案
`04_tasks.md` の各エンドポイント（`GET tasks`, `POST tasks`, `GET tasks/{task_id}`, `PATCH tasks/{task_id}`, `DELETE tasks/{task_id}`）に `##### Response Body フィールド定義` テーブルを追加し、`task` リソースモデルおよび `pagination` モデルの属性（型・必須/Null許容・説明）を明示的に定義してください。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 17:38:45
- **Status**: Resolved

### 実施した修正内容
`docs/design/api_design/04_tasks.md` の全5つのエンドポイント（`GET tasks`, `POST tasks`, `GET tasks/{task_id}`, `PATCH tasks/{task_id}`, `DELETE tasks/{task_id}`）に `##### Response Body フィールド定義` テーブルを新たに追加し、レスポンスオブジェクト（`task`, `tasks[]`, `items[]`, `pagination`, `message`）の全属性のデータ型・必須/Null許容性・詳細な説明を明確に定義しました。

### 変更したファイル
- [04_tasks.md](docs/design/api_design/04_tasks.md)

