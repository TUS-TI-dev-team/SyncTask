# LOGIN_SESSIONテーブルの有効期限（1週間）と要件（1ヶ月）の不整合

- **Status**: Resolved
- **Severity**: Major
- **Created At**: 2026-08-17 12:45:00
- **Target Files**:
  - [database_design.md](docs/design/database_design.md)
  - [requirements.md](docs/req-def/requirements.md)

## 1. 問題の概要
`LOGIN_SESSION` テーブルの `EXPIRES_AT` カラム備考に「1週間(10080分)」と記載されていますが、要件定義書に定められている「1ヶ月（43200分） / Sliding Expiration方式」と矛盾しています。

## 2. 詳細な指摘内容
- `docs/design/database_design.md` の63行目において、`EXPIRES_AT` の備考が `1週間(10080分) / CRONで日次削除` と定義されています。
- しかし、`docs/req-def/requirements.md` の191行目、267行目、274行目では以下のように明記されています：
  > 191: ログイン認証のセッション期限は1ヶ月とする（APIアクセスごとに自動延長される Sliding Expiration 方式）
  > 267: 有効期限 43200分(1ヶ月)
  > 274: ユーザーがログイン中にリクエスト（API呼び出し）を行うたびに、セッションの有効期限をアクセス時刻から43200分（1ヶ月）後に自動更新（Sliding Expiration）する。最後のアクセスから43200分経過したセッションを無効（期限切れ）とする
- 設計書側のセッション有効期間の記述が短く（1週間）、要件（1ヶ月）と整合していません。

## 3. 推奨される修正案
`LOGIN_SESSION` テーブルの `EXPIRES_AT` の備考記述を、要件定義に合わせて「1ヶ月 (43200分) / APIアクセス時にSliding Expirationで更新 / 期限切れは日次Cronで物理削除」のように修正してください。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 13:22:30
- **Status**: Resolved

### 実施した修正内容
`LOGIN_SESSION` テーブルの `EXPIRES_AT` の備考記述を「1ヶ月 (43200分) / APIアクセス時にSliding Expirationで自動延長 / 期限切れは日次Cron（00:00 JST）で物理削除」に修正し、要件定義書（Sliding Expiration、43200分）と完全に整合させました。

### 変更したファイル
- [database_design.md](docs/design/database_design.md)
