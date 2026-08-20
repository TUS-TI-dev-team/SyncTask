# パスワード変更API (`PATCH users/{user_id}/password`) におけるリクエストボディのバリデーション仕様定義の欠落

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-17 16:25:00
- **Target Files**:
  - [03_users.md](docs/design/api_design/03_users.md)
  - [requirements.md](docs/req-def/requirements.md)

## 1. 問題の概要
`PUT users/{user_id}` (3.2.2) にはリクエストボディの型・必須チェック・制約条件を明記したテーブルが用意されているが、`PATCH users/{user_id}/password` (3.2.4) にはリクエストボディのフィールド定義テーブルが存在しない。

要件定義書 (`requirements.md` L28, L252) および API基本仕様 (`02_auth.md` L21) で定められている新パスワードのバリデーション要件が `03_users.md` 上に記載されていないため、実装時およびテスト時の入力検証ルールが不透明になっている。

## 2. 詳細な指摘内容
`03_users.md` 3.2.4 節（L95-L124）において、JSONリクエストボディ例 (`current_password`, `new_password`) は掲載されているが、パラメータテーブルが存在せず、以下のルールが明記されていない:
1. `current_password`: 必須、文字列。
2. `new_password`: 必須、8〜128文字、英大文字・英小文字・数字・記号の4種類のうち3種類以上含む、ユーザー名およびメールアドレスのローカル部（4文字以上の場合）を含まない、前後の空白除去（トリム）。

## 3. 推奨される修正案
`03_users.md` の 3.2.4 `PATCH users/{user_id}/password` 節に以下のフィールド定義テーブルを追加してください。

```markdown
##### Request Body
```json
{
  "current_password": "Password123!",
  "new_password": "NewSecurePassword456!"
}
```

| フィールド | 型 | 必須 | 制約・バリデーション |
| :--- | :--- | :---: | :--- |
| `current_password` | string | ○ | 現在のパスワード。不一致時は 401 (REAUTH_FAILED) エラー |
| `new_password` | string | ○ | 8〜128文字。英大文字/英小文字/数字/記号のうち3種以上を含む。ユーザー名・メールのローカル部（4文字以上の場合）を含まないこと。現在のパスワードと同一の場合は 422 エラー |
```

---

## 修正完了報告

- **Resolved At**: 2026-08-17 16:42:00
- **Status**: Resolved

### 実施した修正内容
`PATCH users/{user_id}/password` (3.2.4) 節に `current_password` および `new_password` の型・必須チェック・文字数・文字種・ユーザー名/メール含有禁止・同一パスワード指定禁止（422）などの制約条件を明記した Request Body フィールド定義テーブルを追加しました。

### 変更したファイル
- [03_users.md](file:////mnt/c/Users/ayumu_wkkd3w3/BoxForSomething/github/SyncTask/docs/design/api_design/03_users.md)
