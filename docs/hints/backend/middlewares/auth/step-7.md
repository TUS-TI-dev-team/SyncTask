# Step 7: まとめと次の一歩

[← Step 6](step-6.md) | [README に戻る](README.md)

## やること(概要)

実装の振り返りと、今回あえてスコープから外した項目の接続先を整理して、学習を次に繋げます。

## やること(詳細)

1. 自分の書いた `middleware/auth.go` を上から読み直し、各判断(期限判定の場所、401の応答内容、コンテキストキー名など)に理由コメントが書かれているか確認する
2. 振り返りメモをこのファイルの下に追記する(下にテンプレあり)
3. 次の学習テーマの候補を把握する:
   - **CSRF検証**: 発行済みの `XSRF-TOKEN` をダブルサブミットCookie方式で検証するミドルウェア
   - **ログアウト**: `LOGIN_SESSION` の削除+Cookie失効
   - **user_id によるタスクの所有者検証**: Step 4 で注入した user_id を handler/service に渡し、「他人のタスクを見れない」ようにする

## 参考リンク

- OWASP Session Management Cheat Sheet(セッション設計のチェックリスト。自分の実装と突き合わせる) — https://cheatsheetseries.owasp.org/cheatsheets/Session_Management_Cheat_Sheet.html
- OWASP Authentication Cheat Sheet — https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html

## 詰まったら(トグルして開く)

<details>
<summary>📱 Step 7 で詰まったら開く</summary>

<details>
<summary>レベル1: まずどこを見る？</summary>

- 自分の `middleware/auth.go` を上から読み直す。外部ドキュメントはもう見なくてよい段階
- OWASP の2つのチートシートを「チェックリスト」として使い、自分の実装に該当項目(セッション無効化、有効期限、Cookie属性、失敗応答の粒度)があるかを照合する

</details>

<details>
<summary>レベル2: 詰まりやすいポイント</summary>

- 振り返りで書くのは「事実 + 解決方法」のセット。「詰まった」だけで終わると次に活きない。どう解決したか、どのリンク/コードが助けになったかまで書く
- 設計判断の理由コメントが抜けている箇所の典型的なのは「期限判定の場所」と「401メッセージの粒度」。Step 3 と Step 5 でメモした内容がコードコメントに反映されているか戻る
- スコープ外に置いた3テーマ(CSRF/ログアウト/所有者検証)は、**今回作った部品のどこに接続するか**を1文で言えるようにする。言えなければ理解がまだ足りないサイン

</details>

<details>
<summary>レベル3: 方向づけ</summary>

- 次テーマは CSRF(`XSRF-TOKEN` のダブルサブミット)が自然な流れ。ログイン時に Cookie がどう発行されたかを `backend/handler/login.go` で見返してから取りかかると理解が早い
- 振り返りテンプレートをこのファイル末尾にコピーして埋める。埋められない項目があれば、それは次の学習テーマ候補のメモになる

</details>
</details>

---

## 振り返りメモ(テンプレート)

```markdown
## 振り返り(Step 7 完了時に追記)

- うまくいったこと:
- 詰まったことと解決方法:
- 自分が下した設計判断とその理由:
  - 期限判定の場所(SQL側 or service側):
  - 401メッセージの粒度:
  - コンテキストキーの設計:
- 次に学びたいこと:
```
