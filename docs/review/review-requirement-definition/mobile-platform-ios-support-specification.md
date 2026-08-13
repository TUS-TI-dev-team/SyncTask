# 対応プラットフォームにおけるモバイルOS（iOS）の記載欠落・定義不一致

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-13 20:55:00
- **Target Files**:
  - [requirements.md](docs/req-def/requirements.md)

## 1. 問題の概要
システム概要では「クロスプラットフォーム対応タスク管理Webアプリ」と定義されていますが、非機能要件の「対応プラットフォーム」のスマホ項目において iOS（Safari / Chrome）の記載がなく、Android のみが対象となっています。

## 2. 詳細な指摘内容
- **[requirements.md](file:///c:/Users/kazuh/Programming/repos/SyncTask/docs/req-def/requirements.md#L7)**:
  ```markdown
  - 個人向けのクロスプラットフォーム対応タスク管理Webアプリを作成し、端末間でタスク管理をできるようにすること
  ```
- **[requirements.md](file:///c:/Users/kazuh/Programming/repos/SyncTask/docs/req-def/requirements.md#L115-L120)**:
  ```markdown
  - スマホ
  	- OS
  		- Android
  	- ブラウザ
  		- Google Chrome
  		- Mozilla Firefox
  ```

「端末間でタスク管理をできるようにする」「クロスプラットフォーム対応」を目的としているにもかかわらず、スマホ対象ブラウザに iOS (Safari, Chrome) が記載されておらず、仕様上の意図的な対象外なのか記載漏れなのかが不明確です。

## 3. 推奨される修正案
iOS端末（iOS Safari, Mobile Chrome）を検証対象・サポート対象に含めるか、あるいは学習目的のスコープ定義として明記してください。

```markdown
- スマホ
	- OS
		- Android
		- iOS
	- ブラウザ
		- Google Chrome
		- Mozilla Firefox
		- Safari (iOS)
```

---

## 修正完了報告

- **Resolved At**: 2026-08-13 21:00:00
- **Status**: Resolved

### 実施した修正内容
`docs/req-def/requirements.md` の「非機能要件 > 対応プラットフォーム > スマホ」のOS項目を更新しました。
- iOSについて非対応（動作保証外）である旨を明記し、スコープ（Androidのみサポート）を明確化しました。

### 変更したファイル
- [requirements.md](docs/req-def/requirements.md)

