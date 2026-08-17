# 共通仕様における標準レスポンスヘッダーおよびキャッシュ制御方針の記述欠落

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-17 16:25:00
- **Target Files**:
  - [01_overview.md](docs/design/api_design/01_overview.md)

## 1. 問題の概要
`01_overview.md` の 1. 概要・共通仕様において、ベースURL、JSON形式、UTF-8エンコーディング、ISO 8601日時フォーマットが定義されているが、すべてのAPI応答における標準レスポンスヘッダー（`Content-Type: application/json; charset=utf-8`）および個人情報・タスク情報を保護するためのキャッシュ制御ヘッダー（`Cache-Control: no-store` 等）の規定が欠落している。

## 2. 詳細な指摘内容
1. **`Content-Type` ヘッダーの非明記**:
   通信形式が JSON であることは記載されているが、HTTP レスポンスヘッダーとして `Content-Type: application/json; charset=utf-8` を付与することが共通ルールとして明確に定義されていない。

2. **個人情報・認証情報のキャッシュ防止（`Cache-Control`）**:
   ユーザー情報（`users/{user_id}`）やタスク情報（`tasks`）は個人の機密情報を含むため、ブラウザやプロキシサーバーによる共有キャッシュを防止するために `Cache-Control: no-store` や `Pragma: no-cache` などのレスポンスヘッダーを適用する方針をセキュリティ / 共通仕様に明記する必要がある。

## 3. 推奨される修正案
`01_overview.md` の 1. 概要・共通仕様に、以下の共通レスポンスヘッダー仕様を追加してください:

```markdown
- **共通レスポンスヘッダー**:
  - `Content-Type: application/json; charset=utf-8`
  - **キャッシュ制御**: 個人情報・認証情報・タスクデータの漏洩を防ぐため、全API応答に `Cache-Control: no-store, no-cache, must-revalidate` および `Pragma: no-cache` を付与します。
```

---

## 修正完了報告

- **Resolved At**: 2026-08-17 16:42:00
- **Status**: Resolved

### 実施した修正内容
`01_overview.md` 1. 共通仕様に標準レスポンスヘッダー `Content-Type: application/json; charset=utf-8` およびキャッシュ制御ヘッダー `Cache-Control: no-store, no-cache, must-revalidate` / `Pragma: no-cache` の適用規則を追加しました。

### 変更したファイル
- [01_overview.md](file:////mnt/c/Users/ayumu_wkkd3w3/BoxForSomething/github/SyncTask/docs/design/api_design/01_overview.md)
