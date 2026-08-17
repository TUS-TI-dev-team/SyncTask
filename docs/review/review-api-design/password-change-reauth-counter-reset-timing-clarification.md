# `PATCH users/{user_id}/password` における再認証成功時の `REAUTH_FAILED_COUNT` リセットタイミングの明確化

- **Status**: Resolved

---

## 修正完了報告

- **Resolved At**: 2026-08-17 17:03:00
- **Status**: Resolved

### 実施した修正内容
`03_users.md` 3.2.4 節 (`PATCH users/{user_id}/password`) の「リクエスト評価順序」ステップ4（パスワード再認証検証）の文末に、ハッシュ照合成功時に即座に `REAUTH_FAILED_COUNT` を 0 にリセットしてステップ5へ進む旨を明記しました。

### 変更したファイル
- [03_users.md](docs/design/api_design/03_users.md)
- **Severity**: Minor
- **Created At**: 2026-08-17 17:01:00
- **Target Files**:
  - [03_users.md](docs/design/api_design/03_users.md)
  - [database_design.md](docs/design/database_design.md)

## 1. 問題の概要
`03_users.md` の 3.2.4 節 (`PATCH users/{user_id}/password`) の「リクエスト評価順序」において、ステップ 4「パスワード再認証検証 (401)」に成功した後、ステップ 5「ビジネスルール検証 (422 `SAME_AS_CURRENT_PASSWORD`)」が評価される。しかし、ステップ 4 でパスワード一致が確認された時点で DB 上の再認証失敗カウンター（`REAUTH_FAILED_COUNT`）が 0 にリセットされるタイミングが明確に記述されていない。

## 2. 詳細な指摘内容
1. **カウンターリセットのタイミング曖昧さ**:
   - 注記（L175）には「再認証失敗カウンター（`REAUTH_FAILED_COUNT`）は、再認証成功時 ... 0 にリセットされます」と記載されている。
   - ステップ 4（再認証検証）で現在のパスワードが一致した直後に `REAUTH_FAILED_COUNT` を 0 にリセットすることが明記されていないと、後続のステップ 5（ビジネスルール検証 422）でエラーが返却された場合に、過去の再認証失敗カウントがリセットされずに残存してしまう実装上の不備が生じるリスクがある。
   - 正しいパスワードが入力された時点（ステップ 4 成功時）で `REAUTH_FAILED_COUNT` は即座に 0 にリセットされるべきである。

## 3. 推奨される修正案
`03_users.md` 3.2.4 節の「リクエスト評価順序」ステップ 4 の記述に、ハッシュ照合成功時に `REAUTH_FAILED_COUNT` を 0 へ即座にリセットする旨を明記してください。

**修正案の例 (`03_users.md` 3.2.4 リクエスト評価順序)**:
```markdown
4. **パスワード再認証検証 (`401 REAUTH_FAILED` / `SESSION_DESTROYED`)**:
   `current_password` のハッシュ照合を実施。不一致時は一律 `1.0s ± 0.1s` の遅延を適用し、失敗回数をカウントアップ。5回達成分は操作中の該当ログインセッションのみを強制破棄および Cookie 削除ヘッダーを付与。**照合成功時は即座に `REAUTH_FAILED_COUNT` を 0 にリセットしてステップ5に進む。**
```
