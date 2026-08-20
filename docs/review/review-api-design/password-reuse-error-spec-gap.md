# パスワード変更・リセットにおける同一パスワード再利用禁止エラーの仕様欠落

- **Status**: Resolved
- **Severity**: Minor
- **Created At**: 2026-08-17 13:17:00
- **Target Files**:
  - [api_design.md](docs/design/api_design.md)
  - [requirements.md](docs/req-def/requirements.md)

## 1. 問題の概要
`docs/design/api_design.md` の共通仕様（1.3 表 L57）ではエラーコード `422 UNPROCESSABLE_ENTITY` の説明として「パスワード再利用」が挙げられていますが、各パスワード変更・リセットエンドポイント（`3.1.9 POST auth/password-reset/reset` および `3.2.4 PATCH users/{user_id}/password`）のエラー一覧およびスキーマ定義に「現在と同じパスワードの再利用禁止（422）」に関する記述が欠落しています。

## 2. 詳細な指摘内容
1. **`PATCH users/{user_id}/password` のエラー定義不足**:
   - `docs/design/api_design.md` L500-502 では `400 Bad Request`（新パスワード要件違反、現在のパスワード不一致）と `404 Not Found` のみが記載されており、新パスワードが現在のパスワードと同一である場合のハンドリングが記載されていません。
2. **`POST auth/password-reset/reset` のエラー定義不足**:
   - `docs/design/api_design.md` L319-322 では `400 Bad Request`（パスワード要件違反）と `403 Forbidden` のみが記載されており、現在設定されているパスワードと同一のパスワードを指定した場合の挙動が明記されていません。
   - ※なお、パスワードリセット時は現在パスワードとの同一性チェックを行うか、あるいはリセットの性質上許容するのか（またはパスワードハッシュとの比較を行うか）のポリシーを明確にする必要があります。

## 3. 推奨される修正案
1. `PATCH users/{user_id}/password` のエラー一覧に `422 Unprocessable Entity`（エラーコード: `SAME_AS_CURRENT_PASSWORD` 等、現在設定中のパスワードと同一）を追記してください。
2. `POST auth/password-reset/reset` におけるパスワード再利用の扱い（同一パスワードを許容するか、422で拒否するか）を明記してください。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 13:22:00
- **Status**: Resolved

### 実施した修正内容
- `docs/design/api_design.md` の `PATCH users/{user_id}/password` において、新パスワードが現在のパスワードと同一の場合に `422 Unprocessable Entity`（エラーコード: `SAME_AS_CURRENT_PASSWORD`）を返却するエラー仕様を明記しました。
- `POST auth/password-reset/reset` において、パスワードリセットは忘却時の復旧手段であるため、現在設定されているパスワードと同一の新パスワードが入力された場合もエラーとせずそのまま更新完了（`200 OK`）とする仕様を明記しました。

### 変更したファイル
- [api_design.md](docs/design/api_design.md)
