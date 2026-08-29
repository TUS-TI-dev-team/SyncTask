# POST /tasks エンドポイント開発計画書

本ドキュメントでは、SyncTask バックエンドにおける新規タスク作成 API（`POST /api/tasks`）の実装計画を定義します。
単一タスク作成および毎週繰り返しタスクの即時一括作成（最大100件）、日本語検索用正規化文字列（`SEARCH_TEXT`）の自動生成、ならびに `backend/TESTING_GUIDE.md` に準拠した単体テスト作成・検証手順を具体的に定めます。

---

## 1. 概要・要件定義

### 1.1 エンドポイント仕様概要
- **パス / メソッド**: `POST /api/tasks`（ベースURL `/api/` 配下）
- **認証**: 必須（Cookie セッション）
- **CSRF検証**: 必須（`X-CSRF-Token` ヘッダー）
- **機能概要**:
  1. **単一タスク作成**: タイトル、コメント、優先度、締切日時、ピン留めフラグを指定して 1 件のタスクを登録。
  2. **毎週繰り返し一括作成 (`is_recurring: true`)**: 開始日〜終了日、指定曜日、締切時刻を指定し、該当するタスクを 1〜100 件の範囲で即時一括生成して登録。
  3. **検索用テキスト自動生成**: タイトルおよびコメントから小文字化・NFKC正規化・ひらがな→カタカナ変換を行った正規化文字列を生成し、`TASK.SEARCH_TEXT` カラムに保存。

### 1.2 リクエスト仕様
#### 単一作成リクエスト例
```json
{
  "title": "課題レポート提出",
  "comment": "第5章の要約を含むこと",
  "priority": "high",
  "due_datetime": "2026-08-20T23:59:00+09:00",
  "is_pinned": false
}
```

#### 毎週繰り返し一括作成リクエスト例
```json
{
  "title": "週次ゼミ発表準備",
  "comment": "進捗スライド作成",
  "priority": "medium",
  "is_pinned": false,
  "is_recurring": true,
  "recurring_rule": {
    "start_date": "2026-08-22",
    "end_date": "2026-10-31",
    "days_of_week": ["saturday"],
    "due_time": "18:00"
  }
}
```

#### バリデーション制約
- `title`: 前後空白トリム後 1〜100 文字必須。改行・タブ等の制御文字は不可（400 `BAD_REQUEST`）。
- `comment`: 0〜1000 文字（トリム後）。改行は `\n` に正規化。未入力時は `""`。
- `priority`: `'high'`, `'medium'`, `'low'`（デフォルト: `'medium'`）。明示的な `null` や不正値は不可（400 `BAD_REQUEST`）。
- `due_datetime`: ISO 8601 文字列または `YYYY-MM-DD`（時刻省略時は `23:59:00+09:00`）。タイムゾーンなしは JST（`+09:00`）解釈、他タイムゾーンは JST 変換。省略または `null` 指定時は `null`。（※`is_recurring: true` 時は無視）
- `is_pinned`: boolean（デフォルト: `false`）。明示的な `null` や非 boolean は不可（400 `BAD_REQUEST`）。
- `is_recurring`: boolean（デフォルト: `false`）。
- `recurring_rule` (`is_recurring: true` 時のみ必須):
  - `start_date`: `YYYY-MM-DD`（過去日指定可、`start_date <= end_date`）
  - `end_date`: `YYYY-MM-DD`（最大 1 年間 / 52 週以内）
  - `days_of_week`: 1 つ以上の曜日必須（`["monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"]`）。小文字正規化・重複排除。
  - `due_time`: `HH:mm`（`00:00`〜`23:59`、省略時は `23:59`）。不正形式・秒数指定は不可（400 `BAD_REQUEST`）。
  - 生成件数が 0 件または 101 件以上の場合は 400 `BAD_REQUEST`。

### 1.3 レスポンス仕様 (201 Created)
```json
{
  "created_count": 1,
  "tasks": [
    {
      "id": "7c9e6679-7425-40de-944b-e07fc1f90ae7",
      "user_id": "550e8400-e29b-41d4-a716-446655440000",
      "title": "課題レポート提出",
      "comment": "第5章の要約を含むこと",
      "priority": "high",
      "status": "not_started",
      "due_datetime": "2026-08-20T23:59:00+09:00",
      "is_pinned": false,
      "created_at": "2026-08-17T12:00:00+09:00",
      "updated_at": "2026-08-17T12:00:00+09:00"
    }
  ]
}
```
※繰り返し作成時は `due_datetime` 昇順でソートされて返却。

### 1.4 エラーレスポンス仕様 (400 Bad Request 等)
```json
{
  "error": {
    "code": "BAD_REQUEST",
    "message": "入力内容に不備があります。",
    "details": [
      {
        "field": "title",
        "message": "タイトルは必須です。"
      }
    ]
  }
}
```

---

## 2. アーキテクチャ・設計決定事項

### 2.1 パッケージ構成（クリーン多層構造）
```
backend/
├── config/
├── db/
├── handler/
│   ├── task.go          # HTTPリクエスト解釈・レスポンス構築・認証情報取得
│   └── task_test.go     # ハンドラー層単体テスト
├── service/
│   ├── task.go          # バリデーション、繰り返しタスク生成、正規化呼び出し
│   └── task_test.go     # サービス層単体テスト
├── repository/
│   ├── task.go          # DBアクセス層（sql.DB / pgx / トランザクション制御）
│   └── task_test.go     # リポジトリ層単体テスト（sqlmock利用）
├── model/
│   ├── task.go          # タスクエンティティ、リクエスト/レスポンス構造体
│   ├── error.go         # 共通エラーレスポンス構造体
│   └── task_test.go     # モデルバリデーション単体テスト
├── util/
│   ├── normalizer.go    # 日本語同一視（小文字化+NFKC+ひらがな→カタカナ変換）
│   └── normalizer_test.go # 正規化ユーティリティ単体テスト
└── router/
    ├── router.go        # ルート登録（/api/tasks）
    └── router_test.go
```

### 2.2 認証・認可方針
- 本タスクでは認証ミドルウェアの本格実装は行わず、Handler が Gin Context（`c.GetString("userID")`）から `userID` を取得する設計とする。
- `userID` が空の場合は `401 UNAUTHORIZED` を返却する。
- 単体テストでは、Gin Context にモック `userID`（UUID文字列）をセットしてテストを実行する。

### 2.3 日本語正規化（`SEARCH_TEXT`）生成ロジック
- **処理ステップ**:
  1. タイトルとコメントを結合（例: `title + " " + comment`）。
  2. `strings.ToLower` で英字を小文字化。
  3. `golang.org/x/text/unicode/norm.NFKC.String()` で全角英数・記号・半角カナを NFKC 正規化。
  4. ひらがな文字コード（`\u3041`〜`\u3096`）をカタカナ文字コード（`\u30a1`〜`\u30f6`）に変換（差分 `+0x60`）。

### 2.4 繰り返しタスク生成アルゴリズム
1. `start_date` と `end_date` の期間整合性を検証（`start_date <= end_date`、期間 <= 366日）。
2. `due_time`（`HH:mm`）のフォーマットと数値範囲（00:00〜23:59）を検証。
3. `days_of_week`（配列）を正規化し、曜日セットを作成。
4. `start_date` から 1 日ずつインクリメントしながらループし、該当日が指定曜日に含まれる場合に `due_datetime`（`YYYY-MM-DDT{due_time}:00+09:00`）を生成。
5. 生成件数が 0 件の場合はエラー（`error.details: [{ field: "recurring_rule", message: "指定された期間内に該当する曜日が存在しません" }]`）。
6. 生成件数が 100 件を超える場合はエラー（`error.details: [{ field: "recurring_rule", message: "生成件数が上限（100件）を超えています" }]`）。
7. 生成された各タスクに個別の UUID、共通の `title`, `comment`, `priority`, `is_pinned`, `status: "not_started"`, `search_text`, `created_at`, `updated_at` を設定。

---

## 3. 開発手順・ステップ

### Step 1: テストデータ・単体テストプログラムの作成 (TDD)
`backend/TESTING_GUIDE.md` の規約（ファイル命名、日本語 `t.Run`、`require`/`assert` の使い分け、`@spec` 連携）に従い、以下の単体テストを作成します。

#### 1. `backend/util/normalizer_test.go`
- `正常系: 半角英大文字・小文字が小文字に統一されること`
- `正常系: 全角英数字が半角小文字にNFKC正規化されること`
- `正常系: 半角カタカナが全角カタカナに変換されること`
- `正常系: ひらがなが全角カタカナに変換されること`
- `正常系: 濁音・半濁音（が、ぱ、ｶﾞ、ﾊﾟ）が正しく合成・変換されること`
- `正常系: タイトルとコメントが結合されて正規化されること`
- `境界値: 空文字の場合は空文字が返却されること`

#### 2. `backend/model/task_test.go`
- `正常系: 有効な単一タスクリクエストがバリデーションを通過すること`
- `正常系: 有効な繰り返しタスクリクエストがバリデーションを通過すること`
- `異常系: タイトルが空文字または空白文字のみの場合にエラーとなること`
- `異常系: タイトルに改行やタブ等の制御文字が含まれる場合にエラーとなること`
- `境界値: タイトルが100文字の場合に通過し、101文字でエラーとなること`
- `境界値: コメントが1000文字の場合に通過し、1001文字でエラーとなること`
- `異常系: priority に無効な値が指定された場合にエラーとなること`
- `異常系: is_recurring=true で recurring_rule が nil の場合にエラーとなること`
- `異常系: start_date > end_date の場合にエラーとなること`
- `異常系: days_of_week が空配列または不正な曜日の場合にエラーとなること`
- `異常系: due_time のフォーマットが不正（秒付きや範囲外）な場合にエラーとなること`

#### 3. `backend/service/task_test.go`
- `正常系: 単一タスクの作成処理が成功し、生成されたタスクが返却されること`
- `正常系: 繰り返しタスクの作成処理で期間内の該当日タスクが昇順で正しく生成されること`
- `境界値: 繰り返し作成で生成件数がちょうど1件の場合に成功すること`
- `境界値: 繰り返し作成で生成件数がちょうど100件の場合に成功すること`
- `異常系: 繰り返し作成で生成件数が0件の場合に指定のエラーメッセージを返すこと`
- `異常系: 繰り返し作成で生成件数が101件以上の場合に指定のエラーメッセージを返すこと`
- `準正常系: due_datetime に日付のみ（YYYY-MM-DD）が指定された場合、23:59:00+09:00 が補完されること`
- `準正常系: UTCや他タイムゾーンの due_datetime が JST に変換・正規化されること`

#### 4. `backend/repository/task_test.go` (`sqlmock` 利用)
- `正常系: 単一タスクの INSERT が正常に実行されること`
- `正常系: 繰り返しタスクの複数件 INSERT（トランザクション）が正常に実行されること`
- `異常系: DBエラー発生時にロールバックされエラーが返却されること`

#### 5. `backend/handler/task_test.go`
- `正常系: 正しいリクエストで 201 Created とレスポンスJSONが返却されること`
- `正常系: 繰り返し作成で 201 Created と created_count, tasks 配列が返却されること`
- `異常系: 未ログイン（Context に userID なし）の場合に 401 UNAUTHORIZED を返すこと`
- `異常系: リクエストバリデーション違反時に 400 BAD_REQUEST と詳細 details を返すこと`
- `異常系: 不正な JSON ボディの場合に 400 BAD_REQUEST を返すこと`

---

### Step 2: プログラムの実装

#### 1. 共通エラー・タスクモデルの実装
- `backend/model/error.go`: `AppError`, `ErrorResponse`, `ErrorDetail` 定義
- `backend/model/task.go`: `Task`, `CreateTaskRequest`, `RecurringRule`, `CreateTaskResponse` 定義とバリデーションメソッド

#### 2. 日本語正規化ユーティリティの実装
- `backend/util/normalizer.go`: `NormalizeSearchText(title, comment string) string` 実装

#### 3. リポジトリ層の実装
- `backend/repository/task.go`:
  - `TaskRepository` インターフェース定義
  - `CreateTask(ctx context.Context, task *model.Task) error`
  - `CreateTasks(ctx context.Context, tasks []*model.Task) error` (トランザクション一括作成)

#### 4. サービス層の実装
- `backend/service/task.go`:
  - `TaskService` インターフェースおよび構造体
  - `CreateTask(ctx context.Context, userID string, req *model.CreateTaskRequest) (*model.CreateTaskResponse, error)`
  - 繰り返しタスク展開・JST日時計算・件数チェックロジック

#### 5. ハンドラー層の実装
- `backend/handler/task.go`:
  - `CreateTaskHandler(service service.TaskService) gin.HandlerFunc`
  - Context から `userID` 取得、リクエストバインド、サービス呼び出し、エラーハンドリング、201 レスポンス返却

#### 6. ルーターへの登録
- `backend/router/router.go`:
  - `/api/tasks` に対する `POST` ルートの追加

---

### Step 3: テスト実行・検証
`backend/TESTING_GUIDE.md` に従い、単体テストを実行・検証します。

```bash
# バックエンドディレクトリで実行
cd backend

# 全単体テストの実行
go test ./...

# 詳細出力（Verbose）での確認
go test -v ./...

# カバレッジの計測と確認
go test -cover ./...
```

**成功基準**:
1. すべてのパッケージ（`util`, `model`, `service`, `repository`, `handler`, `router`）の単体テストが PASS すること。
2. 仕様書に記載された境界値・異常系が網羅されていること。

---

### Step 4: プログラムの修正・リファクタリング（テスト失敗時のサイクル）
1. テスト失敗が発生した場合は、エラーログ・アサーション差分を確認。
2. 実装コード（またはテスト期待値）の不整合を修正。
3. 再度 `go test -v ./...` を実行し、全テスト通過を確認。
4. GoDoc コメント（`@spec` 記法）を追記・整備。

---

## 4. 変更対象ファイル一覧

| 操作 | ファイルパス | 説明 |
|---|---|---|
| 作成 | `backend/model/error.go` | 共通エラーレスポンス構造体 |
| 作成 | `backend/model/task.go` | タスクエンティティ・リクエスト/レスポンスモデル |
| 作成 | `backend/model/task_test.go` | モデルバリデーション単体テスト |
| 作成 | `backend/util/normalizer.go` | 日本語正規化ユーティリティ（小文字+NFKC+カナ統一） |
| 作成 | `backend/util/normalizer_test.go` | 日本語正規化単体テスト |
| 作成 | `backend/repository/task.go` | タスク永続化リポジトリ |
| 作成 | `backend/repository/task_test.go` | リポジトリ層単体テスト（sqlmock） |
| 作成 | `backend/service/task.go` | タスク作成ビジネスロジック・繰り返し展開 |
| 作成 | `backend/service/task_test.go` | サービス層単体テスト |
| 作成 | `backend/handler/task.go` | POST /api/tasks ハンドラー |
| 作成 | `backend/handler/task_test.go` | ハンドラー層単体テスト |
| 変更 | `backend/router/router.go` | ルート登録（/api/tasks） |
| 変更 | `backend/router/router_test.go` | ルーター単体テスト更新 |
