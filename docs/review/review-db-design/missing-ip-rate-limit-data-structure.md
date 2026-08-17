# IPアドレス単位レートリミット（パスワードスプレー対策）を実現するデータ構造の欠如

- **Status**: Resolved
- **Severity**: Major
- **Created At**: 2026-08-17 14:08:00
- **Target Files**:
  - [database_design.md](docs/design/database_design.md)
  - [requirements.md](docs/req-def/requirements.md)

## 1. 問題の概要
要件定義書にて「同一IPアドレスから直近5分間に累計30回以上のログイン認証失敗が発生した場合、該当IPアドレスからのログインAPIリクエストを一律15分間一時遮断（HTTP 429 Too Many Requests）する」と定義されていますが、このIPアドレス単位のレートリミットを実現するためのデータ構造（テーブルまたはカラム）がDB設計に一切存在しません。

## 2. 詳細な指摘内容
- `docs/req-def/requirements.md` の195行目:
  > IPアドレス単位レートリミット（パスワードスプレー対策）: 同一IPアドレスから直近5分間に累計30回以上のログイン認証失敗が発生した場合、該当IPアドレスからのログインAPIリクエストを一律15分間一時遮断（HTTP 429 Too Many Requests、遅延1s±0.1s）する
- 現在の `LOGIN_ACCOUNT` テーブルは `LOGIN_FAILED_COUNT` / `LOGIN_LAST_FAILED_AT` / `LOGIN_LOCK_UNTIL` によって **メールアドレス単位**（=アカウント単位）のロックアウトのみを管理しています。
- しかし、パスワードスプレー攻撃（異なるメールアドレスに対して同一IPから大量のログイン試行を行う攻撃）に対する防御には、**IPアドレスを主キーとする失敗回数の追跡**が必要です。
- `LOGIN_LOG` テーブルにIPアドレスとログイン成否は記録されていますが、毎回のログインリクエストでログテーブルを5分間のウィンドウで集計クエリを実行する方式では性能上の問題があります。また、この方式を採用する場合でも、設計書にその旨の方針が明記されていません。

## 3. 推奨される修正案
以下のいずれかの方式を設計書に明記してください：

### A. 専用テーブルの追加（推奨）
IPアドレス単位のログイン失敗追跡用テーブル（例: `LOGIN_IP_RATE_LIMIT`）を新設する：

| カラム名 | データ型 / 制約 | 備考 |
| --- | --- | --- |
| `IP_ADDRESS` | `VARCHAR(45)` / `PRIMARY KEY` | 対象IPアドレス |
| `FAILED_COUNT` | `INT` / `DEFAULT 0` | 直近5分間の失敗回数（リクエスト時に `LAST_FAILED_AT` から5分超過で0リセット） |
| `LAST_FAILED_AT` | `TIMESTAMPTZ` | 最終失敗日時 |
| `BLOCKED_UNTIL` | `TIMESTAMPTZ` | 30回到達時に `NOW() + 15分` を設定。この時刻まで該当IPからのログインを遮断 |

### B. アプリケーション層（インメモリ）での管理
Redis等のインメモリストアや、アプリケーションレイヤーのレートリミッター（ミドルウェア）にて管理する方式とする場合、その旨をDB設計書または別途の設計書に明記する。

### C. LOGIN_LOGからの動的集計
`LOGIN_LOG` テーブルを使用して動的に `COUNT(*) WHERE IP_ADDRESS = ? AND IS_SUCCESS = FALSE AND CREATED_AT >= NOW() - INTERVAL '5 minutes'` を集計する方式とする場合、パフォーマンスへの影響とインデックス設計をDB設計書に明記する。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 14:26:00
- **Status**: Resolved

### 実施した修正内容
- `docs/design/database_design.md` に新セクション「5. ログインレートリミット管理 (`LOGIN_IP_RATE_LIMIT`)」を追加しました。
- `IP_ADDRESS` (PK), `FAILED_COUNT`, `LAST_FAILED_AT`, `BLOCKED_UNTIL`, `UPDATED_AT` を持つ専用テーブルを定義し、要件定義書に定められた「同一IPから直近5分間に30回以上失敗で15分間遮断（HTTP 429）」をDB層で高速かつ確実に管理できるようにしました。また、日次パージ方針も記載しました。

### 変更したファイル
- [database_design.md](docs/design/database_design.md)
