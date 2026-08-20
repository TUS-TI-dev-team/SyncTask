# 推奨インデックス設計（Section 7）におけるLOGIN_ACCOUNTテーブルのインデックス定義欠落

- **Status**: Open
- **Severity**: Medium
- **Created At**: 2026-08-18 22:15:00
- **Target Files**:
  - [database_design.md](docs/design/database_design.md)

## 1. 問題の概要
`docs/design/database_design.md` の「7. 推奨インデックス設計 (INDEXES)」において、`TASK`、`LOGIN_SESSION`、`OTP_SESSION`、`LOGIN_IP_RATE_LIMIT`、ログテーブル群のインデックス定義が記載されていますが、認証処理の最重要テーブルである `LOGIN_ACCOUNT` に対するインデックス設計が完全に抜け落ちています。

## 2. 詳細な指摘内容
ログイン認証（`EMAIL` によるアカウント検索）、メール重複チェック、および論理削除済みアカウントの除外検索は、本システムで最も高頻度に実行されるクリティカルなDB操作です。

しかし、`database_design.md` の「7. 推奨インデックス設計 (INDEXES)」には以下のセクションしか存在しません：
- 7.1 タスク管理 (`TASK`)
- 7.2 セッション管理 (`LOGIN_SESSION`, `OTP_SESSION`)
- 7.3 レートリミット管理 (`LOGIN_IP_RATE_LIMIT`)
- 7.4 ログテーブル (`LOGIN_LOG`, `ACCESS_LOG`, `MAIL_AUTH_LOG`)

### 問題点：
- `LOGIN_ACCOUNT.EMAIL` に `UNIQUE` 制約を付与した場合の一意インデックス、または論理削除アカウントを除外して高速検索・一意性を担保するための部分一意インデックス（Partial Unique Index）の設計方針が Section 7 に明記されていません。
- DDL 生成時やマイグレーション実装時にインデックス作成が漏れ、ログイン検索時のフルテーブルスキャン（パフォーマンス劣化）を招く恐れがあります。

## 3. 推奨される修正案
「7. 推奨インデックス設計」に `LOGIN_ACCOUNT` のセクションを追加し、インデックス作成 SQL を明記してください：

```sql
### 7.X アカウント管理 (`LOGIN_ACCOUNT`)
```sql
-- メールアドレス検索および一意性保証用（有効アカウント対象）
CREATE INDEX idx_login_account_email ON LOGIN_ACCOUNT (EMAIL);

-- ※ 論理削除を考慮した部分一意インデックスを採用する場合の例:
-- CREATE UNIQUE INDEX uq_login_account_active_email ON LOGIN_ACCOUNT (EMAIL) WHERE IS_DELETED = FALSE;
```
```
