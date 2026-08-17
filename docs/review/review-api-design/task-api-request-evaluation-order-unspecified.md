# タスク管理系 API におけるリクエスト評価順序（リクエスト評価順序）の未規定

- **Status**: Resolved

---

## 修正完了報告

- **Resolved At**: 2026-08-17 17:03:00
- **Status**: Resolved

### 実施した修正内容
`04_tasks.md` に定義されている全 5 エンドポイント（`GET tasks`, `POST tasks`, `GET tasks/{task_id}`, `PATCH tasks/{task_id}`, `DELETE tasks/{task_id}`）に `##### リクエスト評価順序` セクションを追加し、認証/CSRF検証 → 入力バリデーション → 認可/IDOR検証 → ビジネスルール検証の統一順序を明記しました。

### 変更したファイル
- [04_tasks.md](docs/design/api_design/04_tasks.md)
- **Severity**: Medium
- **Created At**: 2026-08-17 17:01:00
- **Target Files**:
  - [04_tasks.md](docs/design/api_design/04_tasks.md)

## 1. 問題の概要
`04_tasks.md` に定義されている 5つのエンドポイント（`GET tasks`, `POST tasks`, `GET tasks/{task_id}`, `PATCH tasks/{task_id}`, `DELETE tasks/{task_id}`）において、他カテゴリの仕様書（`02_auth.md` や `03_users.md`）に存在する「リクエスト評価順序（`##### リクエスト評価順序`）」セクションが一切記述されていません。

## 2. 詳細な指摘内容
1. **エラー優先順位の不透明さ**:
   - 例えば、`PATCH tasks/{task_id}` や `DELETE tasks/{task_id}` において、未ログイン状態 (`401 Unauthorized`)、CSRFトークンヘッダー欠落 (`403 Forbidden`)、不正なリクエストボディ/バリデーションエラー (`400 Bad Request`)、存在しないまたは他ユーザー所有の `task_id` 指定 (`404 Not Found`) が同時に発生し得る状況において、どのエラー応答が最優先で返却されるべきかが未定義です。
2. **セキュリティ・IDOR検出プロービングリスク**:
   - 認可チェック・IDOR検証 (`404 Not Found`) より前にリクエストボディの構文・入力バリデーション (`400 Bad Request`) を実行するか、あるいは認証・CSRF検証を最優先するかによって、攻撃者がステータスコードの違いを利用してリソースの存在有無やバリデーションロジックを探索できるリスクが生じます。

## 3. 推奨される修正案
`04_tasks.md` 内の各エンドポイント（特に状態変更を伴う `POST tasks`, `PATCH tasks/{task_id}`, `DELETE tasks/{task_id}`）に `##### リクエスト評価順序` セクションを追加し、以下の統一順序を明記してください：
1. **認証・CSRF検証 (`401 Unauthorized` / `403 Forbidden`)**
2. **リクエスト構文・入力バリデーション (`400 Bad Request`)**
3. **認可チェック・IDOR/BOLA検証 (`404 Not Found`)**
4. **ビジネスルール検証 (`422 Unprocessable Entity`)**
