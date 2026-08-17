# POST tasks の毎週繰り返し一括作成における recurring_rule.start_date への過去日指定許容の仕様明示漏れ

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-17 17:22:00
- **Target Files**:
  - [04_tasks.md](docs/design/api_design/04_tasks.md)
  - [requirements.md](docs/req-def/requirements.md)

## 1. 問題の概要
要件定義書 (`requirements.md` L87) には「開始日には過去日も指定可能」と明確に規定されていますが、`04_tasks.md` L119 の `recurring_rule.start_date` の定義・補足説明には過去日が指定可能である旨が記載されておらず、実装者が誤って過去日入力を 400 Bad Request として拒否する不具合を生じさせるリスクがあります。

## 2. 詳細な指摘内容
1. **要件定義書 (`requirements.md` L87) の規定**:
   - `requirements.md` L87 の期間・生成件数バリデーション項目に「開始日には過去日も指定可能」と明記されています。
2. **API設計書 (`04_tasks.md` L119) の記述**:
   - しかし `04_tasks.md` L119 のフィールド定義では `開始日（YYYY-MM-DD）。start_date <= end_date` としか記載されておらず、`start_date` に過去日付（システム日付より前の日付）が指定された場合の挙動が明確ではありません。
3. **影響**:
   - バックエンド実装時に、実装者が一般的な日付バリデーションとして「過去日の開始日は不可」と判断し、`start_date < 今日` の場合に 400 Bad Request を返してしまう実装のブレが発生するリスクがあります。

## 3. 推奨される修正案
`04_tasks.md` 3.3.2 節の `recurring_rule.start_date` フィールド定義および補足説明において、以下のように仕様を明確化してください：

`recurring_rule.start_date` の説明欄に「開始日（`YYYY-MM-DD`）。`start_date <= end_date`。なお、業務要件に基づき `start_date` には過去日付の指定も可能であり、`start_date <= end_date` かつ期間が1年以内（生成件数1〜100件）の条件を満たしていれば、過去日に該当するタスクも正常に一括生成されます」と追記してください。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 17:27:35
- **Status**: Resolved

### 実施した修正内容
`04_tasks.md` 3.3.2 節 (`POST tasks`) の `recurring_rule.start_date` のパラメータ定義において、要件定義書 (`requirements.md` L87) に基づき「`start_date` には過去日付の指定も可能であり、`start_date <= end_date` かつ期間が1年以内（生成件数1〜100件）を満たしていれば過去日に該当するタスクも正常に一括生成される」旨を明確に追記しました。

### 変更したファイル
- [04_tasks.md](file:///mnt/c/Users/ayumu_wkkd3w3/BoxForSomething/github/SyncTask/docs/design/api_design/04_tasks.md)
