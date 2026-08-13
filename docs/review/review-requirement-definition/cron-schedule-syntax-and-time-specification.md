# ログインセッション削除Cronスケジュールの表記不備・曖昧さ

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-13 20:55:00
- **Target Files**:
  - [requirements.md](docs/req-def/requirements.md)

## 1. 問題の概要
セッション管理要件におけるログインセッションのCron定期削除処理の指定時刻表記が「毎日00/00」となっており、標準的なCron構文および時刻表記として不適切・曖昧です。

## 2. 詳細な指摘内容
- **[requirements.md](file:///c:/Users/kazuh/Programming/repos/SyncTask/docs/req-def/requirements.md#L204)**:
  ```markdown
  - ログインセッション
    - 有効期限　43200分(1ヶ月)
    - セッションの内無効になったものをCRON, 毎日00/00に削除する
  ```

「毎日00/00」という表記は標準的なCron式（`0 0 * * *`）や時刻指定（`00:00 JST`）と異なり、開発・運用実装時に指定時刻の誤認や実装不備を引き起こす恐れがあります。

## 3. 推奨される修正案
Cron実行のスケジュールを標準的な時刻指定表記に修正してください。

```markdown
- ログインセッション
  - 有効期限　43200分(1ヶ月)
  - 有効期限切れ（無効化済み）のセッションは、日次Cronジョブ（毎日 00:00 JST / Cron: `0 0 * * *`）にてDBから物理削除する
```

---

## 修正完了報告

- **Resolved At**: 2026-08-13 21:00:00
- **Status**: Resolved

### 実施した修正内容
`docs/req-def/requirements.md` の「非機能要件 > セッション管理 > ログインセッション」の削除スケジュール表記を修正しました。
- 「セッションの内無効になったものをCRON, 毎日00/00に削除する」から「有効期限切れ（無効化済み）のセッションは、日次Cronジョブ（毎日 00:00 JST / Cron: `0 0 * * *`）にてDBから物理削除する」へ明確化しました。

### 変更したファイル
- [requirements.md](docs/req-def/requirements.md)

