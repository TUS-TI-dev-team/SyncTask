# 前提ドキュメント参照パスおよびタイムゾーン仕様の記載不備

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-17 17:15:00
- **Target Files**:
  - [job_design.md](docs/design/job_design.md)
  - [requirements.md](docs/req-def/requirements.md)

## 1. 問題の概要
1. `job_design.md` の前提ドキュメントとして挙げられている要件定義書のパスが `docs/requirements.md` となっており、実際のリポジトリ内パス `docs/req-def/requirements.md` と一致していません。
2. 要件定義書では日次 Cron ジョブの実行時刻が「JST（日本標準時）」基準で定義されていますが、`job_design.md` ではスケジューラおよび DB クエリのタイムゾーン（JST / UTC）に関する明記が不足しています。

## 2. 詳細な指摘内容
1. [job_design.md:L4](docs/design/job_design.md#L4) において `前提: docs/requirements.md` と記載されていますが、正しくは `docs/req-def/requirements.md` です。
2. 要件定義書（[requirements.md:L240-241](docs/req-def/requirements.md#L240-L241), [requirements.md:L275](docs/req-def/requirements.md#L275)）では `毎日 00:00 JST / Cron: 0 0 * * *` のように日本標準時が明示されていますが、`job_design.md` では Cron 式（`0 0 * * *`）が UTC なのか JST なのか、またクエリ内の `NOW()` の評価基準タイムゾーンについて言及されていません。クラウド環境等で実行環境のデフォルトが UTC の場合、意図しない時刻に実行されるリスクがあります。

## 3. 推奨される修正案
1. [job_design.md:L4](docs/design/job_design.md#L4) の前提パスを `docs/req-def/requirements.md` に修正してください。
2. 各 Cron スケジュールの基準タイムゾーン（JST 基準で動作させるか、または UTC 換算した Cron 式を環境変数で設定するか）および DB タイムゾーンの扱いについて明記してください。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 17:13:00
- **Status**: Resolved

### 実施した修正内容
- `job_design.md` 冒頭の前提ドキュメントパスを `docs/req-def/requirements.md` に修正しました。
- 第1章および第6章（環境変数一覧）に、スケジューラの基準タイムゾーンが日本標準時（JST / `Asia/Tokyo`）である旨および環境変数 `CRON_TIMEZONE=Asia/Tokyo` を明記しました。

### 変更したファイル
- [job_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/job_design.md)
