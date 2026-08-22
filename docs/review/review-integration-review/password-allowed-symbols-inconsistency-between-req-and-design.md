# パスワード使用可能記号種における要件定義書と設計書の不整合

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-22 13:50:00
- **Target Files**:
  - [01_account_and_auth.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/req-def/requirements/01_account_and_auth.md)
  - [01_overview.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/api_design/01_overview.md)
  - [02_auth.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/api_design/02_auth.md)
  - [03_users.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/api_design/03_users.md)
  - [01_account_creation.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/01_account_creation.md)
  - [06_password_reset.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/06_password_reset.md)
  - [07_password_change.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/process_design/07_password_change.md)

## 1. 問題の概要
パスワードの文字種要件（記号）について、要件定義書（`01_account_and_auth.md`）に記載されている記号一覧（28種類）と、API共通仕様書（`01_overview.md`）および各設計書に記載されている「ASCII印字可能半角記号全32種類」との間で記号の定義に食い違いがあります。
特に `@`（アットマーク）や `` ` ``（バッククォート）が要件定義書側のリストから漏れており、許容される記号の範囲がドキュメント間で一致していません。

## 2. 詳細な指摘内容

1. **要件定義書の記述**:
   - `docs/req-def/requirements/01_account_and_auth.md` 1.3節:
     > 記号 (`!"#$%&'()=-~^\|\{}[]+;:*_/?.>,<`)
     - 上記文字列に含まれる記号は28種類です（`@` や `` ` `` 等が含まれていません）。

2. **API設計書・処理設計書の記述**:
   - `docs/design/api_design/01_overview.md` 1.6節:
     > 許可記号（ASCII印字可能半角記号全32種類: ``! " # $ % & ' ( ) * + , - . / : ; < = > ? @ [ \ ] ^ _ ` { | } ~`` / 正規表現文字クラス: `[\x21-\x2f\x3a-\x40\x5b-\x60\x7b-\x7e]`）
   - `02_auth.md`、`03_users.md`、`process_design/` 各ファイルでも「全32種」と明記されています。

## 3. 推奨される修正案

- `docs/req-def/requirements/01_account_and_auth.md` 1.3節のパスワード要件の記号定義を、API設計書および処理設計書で統一されている「ASCII印字可能半角記号全32種類（``! " # $ % & ' ( ) * + , - . / : ; < = > ? @ [ \ ] ^ _ ` { | } ~``）」に合わせて修正・同期する。

---

## 修正完了報告

- **Resolved At**: 2026-08-22 13:55:00
- **Status**: Resolved

### 実施した修正内容
要件定義書（`01_account_and_auth.md`）のパスワード使用可能記号定義を、API設計書および処理設計書で定義されている「ASCII印字可能半角記号全32種類」に同期・更新しました。

### 変更したファイル
- [01_account_and_auth.md](file:///C:\Users\kazuh\Programming\repos\SyncTask\docs\req-def\requirements\01_account_and_auth.md)
