# OTP_SESSIONのパージ条件と全体最大有効期限の矛盾および排他制御用部分一意インデックスの不足

- **Status**: Resolved
- **Severity**: Major
- **Created At**: 2026-08-17 15:20:00
- **Target Files**:
  - [database_design.md](docs/design/database_design.md)
  - [job_design.md](docs/design/job_design.md)
  - [requirements.md](docs/req-def/requirements.md)

## 1. 問題の概要
1. 定期クリーンアップ（Cronバッチ）におけるパージ条件（`WHERE EXPIRES_AT < NOW()`）が、要件定義書で規定されている「再送を含む手続き全体の最大有効期限（初回発行から15分間）」と矛盾しており、単発OTPの有効期限（5分）が切れた直後にセッションレコードが物理削除され、15分以内の手動再送機能が破綻するリスクがあります。
2. 同一メールアドレスに対する手続き中の排他制御（他ユーザーからの重複登録・変更要求拒否）をDBレベルで確実に担保するための部分一意インデックス（Partial Unique Index）の推奨設計が定義されていません。

## 2. 詳細な指摘内容
- `docs/req-def/requirements.md`（216, 228, 264行目）では、以下の通り定義されています：
  - 「有効期限: 発行から5分間（ただし再送信を行っても、1回の手続きにおけるセッションおよび排他ロックの全体最大有効期限は初回発行から15分間とする）」
  - 「手動再送依頼: 直前の送信から60秒間のクールダウン期間を設定し、再送信回数の上限は設けないが、初回発行から15分の全体有効期限が経過した時点で手続きは無効となる」
- しかし、`docs/design/job_design.md` の第2-1節および第3節では以下のように定義されています：
  ```sql
  DELETE FROM OTP_SESSION
  WHERE EXPIRES_AT < NOW();
  ```
  - `EXPIRES_AT` は単発OTPの発行から5分で到来するため、ユーザーが5分以内にOTPを入力せず、5分〜15分の間に「OTPを再送信する」を押そうとした場合、Cron（15分周期）によりセッション自体が物理削除されてしまい、手続きが継続不能（セッションロスト）となる重大な不具合が発生します。
  - パージすべき対象は「全体最大有効期限が経過したレコード（`MAX_EXPIRES_AT < NOW()`）」または「明示的に無効化・失効済みのレコード（`STATUS IN ('expired', 'locked', 'completed') AND EXPIRES_AT < NOW()`）」である必要があります。
- さらに、`docs/req-def/requirements.md`（33, 49, 201行目）では「OTPセッション有効期間中は、新規登録・変更・パスワードリセット手続き中のメールアドレスに対する重複登録・変更リクエストを排他維持する」と規定されていますが、`docs/design/database_design.md` の第7.2節のインデックス定義には通常の非一意インデックス（`idx_otp_session_pending_email ON OTP_SESSION (PENDING_EMAIL, STATUS, EXPIRES_AT)`）しかなく、同時並行リクエスト時に同一メールアドレスで複数の有効セッションが生成されるレースコンディションをDB制約レベルで防止する部分一意インデックス（例: `CREATE UNIQUE INDEX idx_otp_session_active_email ON OTP_SESSION (PENDING_EMAIL) WHERE STATUS IN ('active', 'verified');`）が定義されていません。

## 3. 推奨される修正案
1. **パージ条件の修正**:
   - `docs/design/database_design.md`（Note欄および7.2節）および `docs/design/job_design.md`（2-1節、3節）におけるクリーンアップSQLの削除条件を以下のように修正・整合させてください：
     ```sql
     DELETE FROM OTP_SESSION
     WHERE MAX_EXPIRES_AT < NOW()
        OR (STATUS IN ('expired', 'locked', 'completed') AND EXPIRES_AT < NOW());
     ```
2. **排他制御用インデックスの追加・見直し**:
   - `docs/design/database_design.md` の第7.2節に、有効な認証手続き中のメールアドレスに対するレースコンディション防止用として、部分一意インデックス（Partial Unique Index）の定義を追加・推奨してください：
     ```sql
     -- 同一メールアドレスに対する有効なOTPセッションの重複発行防止（排他制御）
     CREATE UNIQUE INDEX uq_otp_session_active_pending_email ON OTP_SESSION (PENDING_EMAIL) WHERE STATUS IN ('active', 'verified');
     ```

---

## 修正完了報告

- **Resolved At**: 2026-08-17 15:22:00
- **Status**: Resolved

### 実施した修正内容
1. `docs/design/database_design.md`（Note欄および7.2節）および `docs/design/job_design.md`（第1節、第2-1節、第3節）において、`CLEANUP_OTP_SESSIONS` の物理削除条件を `MAX_EXPIRES_AT < NOW() OR (STATUS IN ('expired', 'locked', 'completed') AND EXPIRES_AT < NOW())` に更新し、5分〜15分の間の手動再送機能がCronパージによって壊れないよう整合させました。
2. `docs/design/database_design.md` の第7.2節に、認証手続き中のメールアドレスに対する同時リクエストレースコンディション防止用の部分一意インデックス（`uq_otp_session_active_pending_email`）を追加し、パージ用インデックス（`idx_otp_session_purge`）のキー構成も最適化しました。

### 変更したファイル
- [database_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/database_design.md)
- [job_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/job_design.md)
