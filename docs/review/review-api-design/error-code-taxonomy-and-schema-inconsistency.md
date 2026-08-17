# 共通エラーレスポンスコード定義と個別エラーコード体系の不整合

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-17 16:25:00
- **Target Files**:
  - [01_overview.md](docs/design/api_design/01_overview.md)

## 1. 問題の概要
`01_overview.md` の 1.3「代表的なエラーコード一覧」テーブルにおいて、エラーコード（`code`）列に HTTP ステータス名由来の汎用文字列（`UNPROCESSABLE_ENTITY`, `RATE_LIMIT_EXCEEDED` 等）が記載されているが、各エンドポイントの詳細設計（`02_auth.md`, `03_users.md`）では `SAME_AS_CURRENT_USERNAME`, `SAME_AS_CURRENT_PASSWORD`, `OTP_REISSUED_DUE_TO_FAILURES`, `REAUTH_FAILED`, `SESSION_DESTROYED` などの個別具体的な文字列が `code` 値として定義されており、レスポンス JSON の `error.code` フィールドの役割および体系が曖昧になっている。

## 2. 詳細な指摘内容
1. **`error.code` の抽象度と表現ルールの混在**:
   1.3 のテーブル（L56-66）では、`422` の `code` として `UNPROCESSABLE_ENTITY` が掲げられ、説明欄の括弧書きに `SAME_AS_CURRENT_PASSWORD` 等の具体的なコードが記載されている。
   一方、`03_users.md` L60・L122 や `02_auth.md` L74 では、レスポンス JSON 内の `"code"` フィールドの値そのものとして `"SAME_AS_CURRENT_USERNAME"`, `"SAME_AS_CURRENT_PASSWORD"`, `"OTP_REISSUED_DUE_TO_FAILURES"` が直接指定されている。
   クライアント実装者が「`error.code` に入る値は HTTP ステータス文字列（`UNPROCESSABLE_ENTITY`）なのか、固有のエラー種別コード（`SAME_AS_CURRENT_USERNAME`）なのか」を混同する原因となる。

2. **追加されたエラーコードの未掲載**:
   `03_users.md` にて追加された再認証失敗時のエラーコード（`REAUTH_FAILED`, `SESSION_DESTROYED`）や、`429` 発生時の具体的なコード（`RATE_LIMIT_EXCEEDED` / `OTP_RESEND_COOLDOWN` 等）が `01_overview.md` 1.3 の代表的エラーコード一覧に網羅されていない。

## 3. 推奨される修正案
`01_overview.md` 1.3 セクションに以下を明確に追記・修正してください:
1. JSON 内の `error.code` フィールドの設計ルール（「汎用エラーコードと固有エラーコードのどちらを格納するか」の基準）を明記する。
2. 代表的なエラーコード一覧テーブルを更新し、各 HTTP ステータスに対してレスポンスされる具体エラーコードの一覧（Taxonomy）を整理して記載する。

例:
```markdown
#### エラーコード (`code`) 設計方針
`error.code` には、フロントエンドが画面表示や制御分岐を判定できるよう、大分類コード（`BAD_REQUEST`, `UNAUTHORIZED` 等）または具体的なビジネスルール違反コード（`SAME_AS_CURRENT_PASSWORD`, `OTP_REISSUED_DUE_TO_FAILURES` 等）を格納します。
```

---

## 修正完了報告

- **Resolved At**: 2026-08-17 16:42:00
- **Status**: Resolved

### 実施した修正内容
`01_overview.md` 1.3 節のエラーコード設計方針をアップデートし、`error.code` に格納される大分類コードおよび個別固有コードの役割を明確化するとともに、「代表的なエラーコード一覧」に各HTTPステータスに対応する具体的エラーコードの一覧を整理して反映しました。

### 変更したファイル
- [01_overview.md](file:////mnt/c/Users/ayumu_wkkd3w3/BoxForSomething/github/SyncTask/docs/design/api_design/01_overview.md)
