# タスク作成ポップアップにおける単一/繰り返し作成UI切り替え、プレビュー要素、および締切日時クリア操作仕様の未確定・記載漏れ

- **Status**: Resolved
- **Severity**: Medium
- **Created At**: 2026-08-22 14:05:00
- **Target Files**:
  - [screen_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/screen_design.md#L34-L35)
  - [screen_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/screen_design.md#L46-L48)
  - [02_task_management.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/req-def/requirements/02_task_management.md#L48-L51)
  - [04_tasks.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/api_design/04_tasks.md#L124-L132)

## 1. 問題の概要
[screen_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/screen_design.md) の「タスク作成ポップアップ」および「タスク編集ポップアップ」において、以下のUI仕様が未確定または構成要素の定義と補足説明の間で食い違いがあります。
1. 「タスク作成ポップアップ」の構成要素テーブルに単一タスク用入力欄と繰り返し用入力欄が混在して羅列されており、トグル切り替え時の表示・非表示・非活性のUI切り替え制御が明記されていない。
2. 補足（47行目）で言及されている「該当件数のリアルタイムプレビュー」および上限超過/0件時のUI警告領域が、構成要素テーブル（34行目）に含まれていない。
3. タスク作成・編集ポップアップにおける「締切日時の解除（クリア）」操作のUI部品仕様が記載されていない（API上は `due_datetime: null` で解除可能）。

## 2. 詳細な指摘内容
1. **[screen_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/screen_design.md) 34行目（構成要素）**:
   - `タスク名入力欄<br>優先度選択プルダウン<br>締切り日時入力欄<br>コメント入力欄<br>「繰り返し作成」切り替えトグル/チェックボックス<br>期間指定（開始日・終了日）入力欄（カレンダーピッカー付き）<br>締め切り時刻(任意指定)<br>曜日選択チェックボックス（日〜土）<br>「決定」ボタン<br>「キャンセル」ボタン`
   - 単一作成時の「締切り日時入力欄（日付＋時刻）」と、繰り返し作成時の「期間指定（開始日・終了日）＋締め切り時刻(任意指定)」がどのように切り替わるのかが構成要素表から読み取れません。
2. **リアルタイムプレビュー領域の構成要素欠落**:
   - 47行目の補足には「該当件数のリアルタイムプレビューを提供」とありますが、34行目の構成要素には「生成件数プレビュー・警告表示領域」が含まれていません。要件定義（[02_task_management.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/req-def/requirements/02_task_management.md)）では0件時および100件超過時の事前警告やAll-or-Nothingの制御が定められており、UI上での明示的な件数表示（例: 「○件作成予定」）および100件超過時の「決定」ボタン非活性化または警告表示仕様が必要です。
3. **締切日時クリア操作の欠落**:
   - タスク編集ポップアップ（35行目）等において、設定済みの締切日時を解除（未設定 `null` に戻す）するためのUI操作（「締切をクリア」ボタン、またはピッカー内のクリアアイコン等）が明記されていません。

## 3. 推奨される修正案
1. **構成要素の明確化**:
   - 構成要素テーブル（34行目）を更新し、トグルによる動的切り替え要素（単一作成時：締切り日時入力欄／繰り返し作成時：期間指定ピッカー、締め切り時刻入力欄、曜日選択チェックボックス、生成件数プレビュー・エラー表示領域）を明確に区分して記述する。
   - 締切日時入力欄に「クリア（解除）ボタン」を含めることを明記する。
2. **タスク作成・編集補足の追記**:
   - 繰り返しトグルON時は、単一の「締切り日時入力欄」が非表示（または非活性）となり、「期間指定」「締切時刻」「曜日選択」「生成件数プレビュー（1〜100件、0件または100件超過時は警告メッセージを表示し「決定」ボタンをDisabled制御）」が表示される旨を補足に追記する。

---

## 修正完了報告

- **Resolved At**: 2026-08-22 14:10:00
- **Status**: Resolved

### 実施した修正内容
`docs/design/screen_design.md` の「タスク作成ポップアップ」および「タスク編集ポップアップ」において、単一/繰り返し作成の動的UI切り替え、生成件数プレビュー・警告・Disabled制御（1〜100件）、および締切日時クリアボタン仕様を明記しました。

### 変更したファイル
- [screen_design.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/screen_design.md)
