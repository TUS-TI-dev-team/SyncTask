# `PATCH users/{user_id}/password` におけるビジネスルール検証 (422) と再認証検証 (401) の評価順序不備によるセキュリティ・カウンター回避リスク

- **Status**: Resolved
- **Severity**: Major
- **Created At**: 2026-08-17 16:54:00
- **Target Files**:
  - [03_users.md](docs/design/api_design/03_users.md)

## 1. 問題の概要
`03_users.md` の `PATCH users/{user_id}/password` (3.2.4) の「リクエスト評価順序」において、ステップ 3「ビジネスルール検証 (`422 SAME_AS_CURRENT_PASSWORD`)」が、ステップ 4「パスワード再認証検証 (`401 REAUTH_FAILED` / `SESSION_DESTROYED`)」の前に実行される定義になっている。

リクエストボディ内の `new_password` と `current_password` の文字列比較を、DB 内の実際の現在のパスワードハッシュ照合（ステップ 4）より前に実施するため、誤った `current_password` を送信した場合であっても `new_password == current_password` であれば 422 エラーが返却され、パスワード照合（401）、失敗カウンター加算、および Timing Attack 対策の応答遅延（1.0s ± 0.1s）がバイパスされるリスクが存在する。

## 2. 詳細な指摘内容
1. `3.2.4` L134-L137:
   - ステップ 3: `new_password` が `current_password` と同一か検証（同一の場合は即座に 422 `SAME_AS_CURRENT_PASSWORD` を返却、遅延・失敗カウンター加算なし）
   - ステップ 4: `current_password` のハッシュ照合を実施（不一致時は一律 `1.0s ± 0.1s` の遅延を適用し、失敗回数をカウントアップ。5回達成分はセッション強制破棄）

2. 攻撃者または第三者が誤った `current_password`（例: `"WrongPassword123!"`）を指定し、`new_password` に同一の `"WrongPassword123!"` を設定してリクエストを送信した場合：
   - ステップ 3 により `new_password == current_password` が判定され、ステップ 4 のパスワードハッシュ照合を実行することなく即座に `422 SAME_AS_CURRENT_PASSWORD` が返却される。
   - 本来実行されるべき再認証失敗カウンター（`REAUTH_FAILED_COUNT`）のカウントアップや、5回失敗時のセッション強制破棄（`SESSION_DESTROYED`）、および 1.0s レスポンス遅延がすべてスキップされてしまう。

3. また、現在のパスワードの正しい照合（認証完了）が完了していない段階で「現在のパスワードと同一である」という情報を 422 エラーとして返却することは、レスポンスオラクルとして機能する恐れがある。

## 3. 推奨される修正案
`3.2.4` の「リクエスト評価順序」において、パスワード再認証検証（DB ハッシュ照合による `current_password` の正しい検証）をビジネスルール検証（`SAME_AS_CURRENT_PASSWORD`）の前に実行するように評価順序を修正してください。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 16:56:00
- **Status**: Resolved

### 実施した修正内容
`03_users.md` の 3.2.4 (`PATCH users/{user_id}/password`) の「リクエスト評価順序」を修正し、パスワード再認証検証（401）をビジネスルール検証（422 `SAME_AS_CURRENT_PASSWORD`）よりも前に評価・実行する順序に変更しました。

### 変更したファイル
- [03_users.md](docs/design/api_design/03_users.md)

**修正後の評価順序案 (`PATCH users/{user_id}/password`)**:
```markdown
##### リクエスト評価順序
1. **リクエスト構文・入力バリデーション (`400 Bad Request`)**:
   リクエストボディの JSON 形式、必須フィールドの有無、および `new_password` の文字数・文字種・ユーザー名/メールローカル部含有検証を最優先で実施。不備がある場合は即座に 400 エラーを返却（遅延・失敗カウンター加算なし）。
2. **認可チェック・IDOR/BOLA検証 (`404 Not Found`)**:
   パスパラメータ `user_id` とセッションユーザーIDの一致を検証。不一致時は即座に 404 エラーを返却（遅延・失敗カウンター加算なし）。
3. **パスワード再認証検証 (`401 REAUTH_FAILED` / `SESSION_DESTROYED`)**:
   `current_password` のハッシュ照合を実施。不一致時は一律 `1.0s ± 0.1s` の遅延を適用し、失敗回数をカウントアップ。5回達成分はセッション強制破棄および Cookie 削除ヘッダーを付与。
4. **ビジネスルール検証 (`422 Unprocessable Entity`)**:
   再認証成功後、`new_password` が照合済みの `current_password` と同一か検証（同一の場合は 422 `SAME_AS_CURRENT_PASSWORD` を返却）。
```
