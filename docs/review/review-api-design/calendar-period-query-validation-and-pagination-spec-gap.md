# GET tasksにおけるカレンダー期間クエリ（start_date/end_date）のバリデーションおよびページネーション挙動の未定義

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-17 16:35:00
- **Target Files**:
  - [04_tasks.md](docs/design/api_design/04_tasks.md)

## 1. 問題の概要
`GET tasks` において、カレンダー表示用の期間指定クエリパラメータ（`start_date` および `end_date`）に関する片側省略時の挙動、日付逆転時のエラーハンドリング、最大取得期間の制限ルール、および期間取得時のページネーション挙動が未定義です。

## 2. 詳細な指摘内容
1. **期間パラメータのバリデーション仕様の不足**:
   - `start_date` のみが指定され `end_date` が省略された場合（またはその逆）に、サーバーがどのように処理するか（リクエスト拒否するのか、他方をデフォルト設定するのか）が明記されていません。
   - `start_date > end_date`（開始日が終了日より未来）が指定された場合のバリデーションエラー（`400 Bad Request` / `code: BAD_REQUEST`）の定義が不足しています。
   - `end_date - start_date` の最大許容期間（例: 最大42日間 / 6週間、または最大1年間）の制限が定義されていません。範囲無制限のクエリを許可した場合、極端に広い日付範囲（例: 100年分）の検索リクエストにより大量のDBレコードが走査され、DOS脆弱性やパフォーマンス劣化の原因となります。

2. **カレンダー期間取得におけるページネーションの不整合**:
   - `04_tasks.md` では `GET tasks` のデフォルト件数が `limit=20` と規定されていますが、カレンダー画面で `start_date` と `end_date` を指定して月間グリッド内の全タスクを取得しようとした際、1ヶ月内に20件以上のタスクが存在すると 21件目以降が取得されず、カレンダー上で表示が欠落します。
   - `start_date` / `end_date` 指定時はページネーションを無効化して期間内の全タスクを返却するのか、またはレスポンスの `pagination` オブジェクトがどのように変化するのかが明記されていません。

## 3. 推奨される修正案
1. `GET tasks` の `start_date` / `end_date` パラメータ説明欄およびエラー仕様に以下の制約を明記してください:
   - `start_date` と `end_date` はペアで指定必須とし、いずれか一方のみの指定や `start_date > end_date` の場合は `400 Bad Request`（`code: BAD_REQUEST`）を返却する。
   - 指定可能な最大期間幅を「最大42日間（6週間）」または「最大1年間」と制限し、超過した場合は `400 Bad Request` を返却する。
2. `start_date` および `end_date` を用いた期間クエリ実行時は、ページネーション limit を解除（または期間内全件返却）し、レスポンスの `pagination` メタデータの挙動について仕様書に注記を追記してください。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 16:42:00
- **Status**: Resolved

### 実施した修正内容
`GET tasks` における `start_date` と `end_date` のペア指定必須（片側指定や `start_date > end_date` 時は 400 Bad Request）、最大許容期間幅（42日間 / 6週間）、およびカレンダー期間取得時のページネーション limit 解除（全件返却）の挙動を追記・明確化しました。

### 変更したファイル
- [04_tasks.md](file:////mnt/c/Users/ayumu_wkkd3w3/BoxForSomething/github/SyncTask/docs/design/api_design/04_tasks.md)
