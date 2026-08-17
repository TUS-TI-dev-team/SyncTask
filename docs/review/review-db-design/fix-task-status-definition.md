# TASKテーブルのSTATUS備考における「未着手」ステータスの欠落

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-17 12:45:00
- **Target Files**:
  - [database_design.md](docs/design/database_design.md)
  - [requirements.md](docs/req-def/requirements.md)
  - [screen_design.md](docs/design/screen_design.md)

## 1. 問題の概要
`TASK` テーブルの `STATUS` カラムの備考に `IN_PROGRESS, COMPLETED` のみが記載されており、初期ステータスおよび要件定義・画面設計で必須とされている「未着手」ステータスが欠落しています。

## 2. 詳細な指摘内容
- `docs/design/database_design.md` の47行目において、`STATUS` カラムの備考が `IN_PROGRESS, COMPLETED` と記載されています。
- しかし、`docs/req-def/requirements.md` の82行目および98行目では以下のように定義されています：
  > 82: 初期ステータスは「未着手」
  > 98: タスク状態更新（未着手、進行中、完了）
- また、`docs/design/screen_design.md` の39行目・41行目でも「未着手 / 進行中 / 完了」の3状態を扱うことが明記されています。
- 設計書の備考欄から「未着手」に該当する値（例: `NOT_STARTED` または `TODO`）が抜けているため、実装時の定義揺れを招く恐れがあります。

## 3. 推奨される修正案
`TASK` テーブルの `STATUS` カラムの備考に、取り得る3つのステータス値（例: `NOT_STARTED` [未着手・初期値], `IN_PROGRESS` [進行中], `COMPLETED` [完了]）を明記し、デフォルト値として `NOT_STARTED` を指定してください。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 13:22:30
- **Status**: Resolved

### 実施した修正内容
`TASK` テーブルの `STATUS` カラムの制約を `VARCHAR(20) / NOT NULL, DEFAULT 'NOT_STARTED'` に更新し、備考欄の記述を `NOT_STARTED (未着手・初期値), IN_PROGRESS (進行中), COMPLETED (完了)` に修正しました。

### 変更したファイル
- [database_design.md](docs/design/database_design.md)
