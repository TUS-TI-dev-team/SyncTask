# LOGIN_FAILED_COUNTのリセット期間（5分）と要件（15分）の不整合

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-17 12:45:00
- **Target Files**:
  - [database_design.md](docs/design/database_design.md)
  - [requirements.md](docs/req-def/requirements.md)

## 1. 問題の概要
`LOGIN_ACCOUNT` テーブルの `LOGIN_FAILED_COUNT` の備考に「5分ごとにリセット」と記載されていますが、要件定義書に定められている「直近15分間に5回」「最後の失敗試行から15分経過でリセット」という定義と矛盾しています。

## 2. 詳細な指摘内容
- `docs/design/database_design.md` の21行目において、`LOGIN_FAILED_COUNT` の備考に `5分ごとにリセット (判定時動的リセット)` と記載されています。
- しかし、`docs/req-def/requirements.md` の194行目および196行目では以下のように定義されています：
  > 194: メールアドレス単位ロックアウト: 同一メールアドレスに対して、直近15分間に5回連続でログイン認証に失敗した場合、該当メールアドレスに対するログイン試行を30分間ロックアウトする
  > 196: 失敗カウントの保持期間: 最後の失敗試行から15分経過、またはログイン成功時に失敗カウントを0にリセットする
- そのため、5分でリセットしてしまうと、例えば10分間に5回失敗した場合にロックアウトされず、要件定義のセキュリティ基準（直近15分間の判定窓）を満たせなくなります。

## 3. 推奨される修正案
`LOGIN_ACCOUNT` テーブルの `LOGIN_FAILED_COUNT` の備考を「最後の失敗から15分経過でリセット (またはログイン成功時にリセット)」に修正してください。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 13:22:30
- **Status**: Resolved

### 実施した修正内容
`LOGIN_ACCOUNT` テーブルの `LOGIN_FAILED_COUNT` の備考記述を「直近15分間に5回失敗で30分間ロック / 最後の失敗から15分経過またはログイン成功時に0にリセット」に修正し、要件定義書のロックアウト判定仕様と整合させました。

### 変更したファイル
- [database_design.md](docs/design/database_design.md)
