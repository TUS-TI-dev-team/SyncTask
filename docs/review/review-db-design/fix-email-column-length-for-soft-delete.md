# 論理削除時のEMAIL退避フォーマットによる桁あふれ・更新エラーの防止

- **Status**: Resolved
- **Severity**: Major
- **Created At**: 2026-08-17 14:00:00
- **Target Files**:
  - [database_design.md](docs/design/database_design.md)
  - [requirements.md](docs/req-def/requirements.md)

## 1. 問題の概要
`LOGIN_ACCOUNT` テーブルの `EMAIL` カラムが `VARCHAR(255)` で定義されていますが、アカウント退会（論理削除）時にメールアドレスの一意性制約の衝突を回避するため `deleted_<USER_ID>_<EMAIL>` 形式へ更新する仕様となっています。
この退避文字列は最大で約300文字に達するため、メールアドレスの長さによっては `VARCHAR(255)` のカラム上限を超過し、退会処理時にDBエラー（文字列切り捨て・サイズ超過エラー）が発生して退会できなくなります。

## 2. 詳細な指摘内容
- `database_design.md` の17行目および32行目:
  - `EMAIL`: `VARCHAR(255)` / `UNIQUE, NOT NULL`
  - 「論理削除実行時に `EMAIL` カラムの値を退避フォーマット（例: `deleted_<USER_ID>_<EMAIL>`）に更新し、有効なアカウント間でのみ一意性を維持します。」
- 文字数計算:
  - プレフィックス `deleted_` (8文字) + `USER_ID` (UUID 36文字) + アンダースコア `_` (1文字) + `EMAIL` (規格上最大254〜255文字)
  - 合計: 8 + 36 + 1 + 255 = **300文字**
- メールアドレスが211文字以上のアカウントが退会しようとした場合、300文字の更新文字列が `VARCHAR(255)` に収まらず、データベースの更新処理が失敗（`Value too long for type character varying(255)` 等）します。

## 3. 推奨される修正案
以下のいずれか（または組み合わせ）の対応を実施することを推奨します：

1. **カラム長の拡張（推奨）**:
   - `LOGIN_ACCOUNT.EMAIL` のデータ型を `VARCHAR(320)`（または `VARCHAR(300)` 以上）に拡張する。
2. **退避フォーマットの見直し**:
   - 退避フォーマットを `deleted_<USER_ID>` （メールアドレス部分を含めない）とするか、またはメールアドレスをハッシュ化（例: `deleted_<USER_ID>_<SHA256(EMAIL)[:16]>`）して255文字以内に確実に収める。
3. **部分インデックス（Partial Unique Index）の採用検討**:
   - 退避更新を行わずに、`IS_DELETED = FALSE` のレコードのみを対象とする部分一意インデックス（`CREATE UNIQUE INDEX idx_login_account_email ON LOGIN_ACCOUNT(EMAIL) WHERE IS_DELETED = FALSE;`）を採用する方式を検討・明記する。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 14:26:00
- **Status**: Resolved

### 実施した修正内容
- `docs/design/database_design.md` において、`LOGIN_ACCOUNT` テーブルの `EMAIL` カラムのデータ型を `VARCHAR(320)` に拡張しました。
- これにより、論理削除時の退避フォーマット `deleted_<USER_ID>_<EMAIL>`（最大約300文字）を確実に格納可能とし、退会処理時の桁あふれエラーを防止しました。

### 変更したファイル
- [database_design.md](docs/design/database_design.md)
