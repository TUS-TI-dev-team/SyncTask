---
name: herdr-review-loop
description: "herdr CLIを使用してReviewerエージェントとFixerエージェントを別Paneで自動制御・オーケストレーションし、要件定義レビュー（docs/req-def/requirements.md等）のレビュー実行、未解消指摘（Status: Open）の抽出、人間への対処方針ヒアリング、修正の実行、および再レビュー検証ループを半自動化・一言指示で実行するためのスキルです。"
---

# herdr-review-loop Skill

このスキルは、`herdr` CLI（マルチエージェントオーケストレーションランタイム）を活用し、**親エージェント（Orchestrator）が「Reviewerエージェント」と「Fixerエージェント」を別Paneで全自動制御**します。
人間は「各指摘に対する対処方針の指示（意志決定）」のみを行い、レビュー実行・指摘抽出・修正指示・再レビュー検証ループの煩雑な操作をすべて自動化します。

---

## 🎯 役割とメンタルモデル

```text
 +--------------------------------------------------------------------------------+
 | [Orchestrator (親エージェント)]                                                  |
 |  - herdr CLI で Reviewer / Fixer を別 Pane に分割・起動                        |
 |  - 指摘ファイル(Status: Open)の全件スキャン                                        |
 |  - 人間へ「〇〇の指摘に対する方針を入力してください」と提示 (★唯一の人間介入ポイント)     |
 +-----------------------------------+--------------------------------------------+
                                     |
           +-------------------------+-------------------------+
           |                                                   | herdr agent prompt / wait
           v                                                   v
 +-----------------------------------+   +----------------------------------------+
 | [Reviewer Pane]                   |   | [Fixer Pane]                           |
 | - /review /review-changes を実行  |   | - /apply-review-fixes を実行          |
 | - 指摘ファイルを生成               |   | - 人間の指示に従って仕様書・コードを修正|
 +-----------------------------------+   +----------------------------------------+
```

---

## 📋 実行フロー

### 1. ブランチ名の取得と対象ディレクトリの確定
1. `git branch --show-current` でブランチ名を取得し、スラッシュ等をハイフンに置換して `<sanitized-branch>` を確定します（例: `review-requirement-definition`）。
2. 指摘ファイル保存先ディレクトリ `docs/review/<sanitized-branch>/` を確認・特定します。
3. レビュー対象ファイル（指定がない場合は `docs/req-def/requirements.md` またはブランチ差分）を確定します。

---

### 2. herdr 環境の確認および Helper Agent の起動

1. `herdr agent list` を実行し、現在生存中のエージェントを確認します。
2. `reviewer` または `fixer` という名前のエージェントが存在しない場合は、以下の手順で別 Pane を作成・起動します：

```bash
# 1. Reviewer 用 Pane の分割と起動
SPLIT_REV=$(herdr pane split --current --direction right --no-focus)
REV_PANE=$(printf '%s\n' "$SPLIT_REV" | jq -r '.result.pane.pane_id')
herdr agent start reviewer --kind antigravity --pane "$REV_PANE"

# 2. Fixer 用 Pane の分割と起動
SPLIT_FIX=$(herdr pane split "$REV_PANE" --direction down --no-focus)
FIX_PANE=$(printf '%s\n' "$SPLIT_FIX" | jq -r '.result.pane.pane_id')
herdr agent start fixer --kind antigravity --pane "$FIX_PANE"
```

---

### 3. レビューの自動トリガーと完了待機

Orchestrator から `reviewer` エージェントにプロンプトを送信し、完了を待機します：

```bash
herdr agent prompt reviewer "/review /review-changes <target_file>" --wait --until idle --timeout 300000
```

`--wait --until idle` により、Reviewer エージェントが要件定義の査読・指摘ファイル生成を完了して `idle` 状態になるまで Orchestrator は自動待機します。

---

### 4. 未解消指摘（Status: Open）の検出と人間へのヒアリング

1. `docs/review/<sanitized-branch>/` 配下の全 `.md` ファイルをスキャンし、先頭付近に `- **Status**: Open` が含まれるファイルを抽出します。
2. **`Status: Open` が 0 件の場合**:
   - 「すべてのレビュー指摘が解消されました」と人間に報告し、ループを終了します。
3. **`Status: Open` が 1件以上存在する場合**:
   - 各 Open な指摘ファイルについて、概要と推奨修正案を読み取り、人間に方針を提示します。

#### 人間への質問フォーマット例（★人間の唯一の介入ポイント）
> **【レビュー指摘】 `username-change-specification-deficiencies.md`**  
> - **概要**: ユーザー名変更時の重複チェック処理および表示遅延時の排他制御仕様が未定義です。  
> - **推奨案**: データベースでの一意制約および楽観的ロックの仕様を追記する。  
>  
> **👉 この指摘に対する対処方針を指示してください:**  
> (例: 「推奨案通りデータベース一意制約と楽観的ロックを追記して」「A案で修正」「現状仕様で割り切るためドキュメントのみ注記」等)

---

### 5. Fixer エージェントへの修正指示と完了待機

人間から入力された対処方針を受け取ったら、Orchestrator が `fixer` エージェントへ修正指示をアトミック送信します：

```bash
PROMPT="/apply-review-fixes 指摘ファイル: docs/review/<sanitized-branch>/<file_name>.md
対処方針: <人間が指定した方針>
上記の方針に従って仕様書およびコードの修正を行い、テストで検証したうえでステータスを Resolved に更新してください。"

herdr agent prompt fixer "$PROMPT" --wait --until idle --timeout 600000
```

`fixer` エージェントが修正および検証を完了し `idle` になるまで待機します。

---

### 6. 再検証ループ（Iteration）

1. 指定された指摘の修正が完了したら、`reviewer` エージェントへ再検証を依頼します：
   ```bash
   herdr agent prompt reviewer "/review /review-changes docs/review/<sanitized-branch>/ の修正結果を再検証し、指摘状態を更新してください" --wait --until idle --timeout 300000
   ```
2. 再び Step 4 に戻り、`Status: Open` が完全に 0件になるまで繰り返します。

---

## 🛠 トラブルシューティング

- **`agent_prompt_stalled` エラーが発生した場合**:
  エージェントがプロンプトを受信しても 5秒以内に状態変化が観測されなかったことを示します。`herdr agent list` や `herdr agent explain <name>` で状態を確認し、必要に応じて `herdr agent send-keys <name> enter` を送信して復帰させてください。
- **エージェントが応答停止した場合**:
  `herdr agent read <name> --source recent-unwrapped` で最新ログを確認し、ダイアログ待ち等の場合は `herdr agent send-keys <name> esc` や `enter` でレスキューします。
