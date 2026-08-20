# パスワード変更APIにおけるINVALID_PASSWORD_CONTENTエラーコードの定義漏れおよび評価順序の不整合

- **Status**: Resolved
- **Severity**: Major
- **Created At**: 2026-08-17 17:45:00
- **Target Files**:
  - [03_users.md](docs/design/api_design/03_users.md)
  - [01_overview.md](docs/design/api_design/01_overview.md)
  - [02_auth.md](docs/design/api_design/02_auth.md)

## 1. 問題の概要
`03_users.md` の 3.2.4 節 (`PATCH users/{user_id}/password`) において、新パスワード（`new_password`）にユーザー名やメールアドレスのローカル部が含まれる場合のバリデーションが「リクエスト評価順序」のステップ2（`400 Bad Request` / code: `"BAD_REQUEST"`）に分類されており、`01_overview.md`（1.3 節）および `02_auth.md`（3.1.9 節）で定義されている `422 Unprocessable Entity`（code: `"INVALID_PASSWORD_CONTENT"`）の仕様と直接矛盾しています。
また、ユーザー名やメールアドレスローカル部の含有チェックはログインセッション/DBからユーザー情報を参照して行うドメインバリデーションであるにもかかわらず、認可チェック（ステップ3 `404 Not Found`）やパスワード再認証（ステップ4 `401 REAUTH_FAILED`）より前のステップ2で評価する記述となっており、リクエスト評価順序の設計上も不整合が発生しています。

## 2. 詳細な指摘内容
1. **エラーコードおよびHTTPステータスコードのドキュメント間不一致**:
   - `01_overview.md` 1.3 節「代表的なエラーコード一覧」では、`422 INVALID_PASSWORD_CONTENT`（「パスワード変更/リセット時のユーザー名・メールアドレスローカル部含有違反」）が明記されています。
   - `02_auth.md` 3.1.9 節 (`POST auth/password-reset/reset`) においても、新パスワードにユーザー名やメールアドレスローカル部が含まれる場合は `422 Unprocessable Entity`（code: `"INVALID_PASSWORD_CONTENT"`）を返却する仕様となっています。
   - しかし `03_users.md` 3.2.4 節では、`new_password` へのユーザー名/メールローカル部含有違反をステップ2の構文エラー（`400 BAD_REQUEST`）として扱い、`##### Errors` 一覧からも `422 INVALID_PASSWORD_CONTENT` が完全に漏れています。

2. **リクエスト評価順序の不備（認可・再認証前のコンテキスト参照）**:
   - 3.2.4 節のステップ2（`400 Bad Request`）で「`new_password` の文字数・文字種・ユーザー名/メールローカル部含有検証を最優先で実施」と記載されています。
   - しかし、ユーザー名やメールアドレスローカル部の含有チェックは単なる文字数・文字種（構文）チェックではなく、認証セッション/DBから取得したユーザーコンテキストに依存する検証です。
   - これをステップ3（認可・IDORチェック `404`）やステップ4（パスワード再認証 `401`）の前に評価すると、未認可のリクエストや再認証失敗リクエストに対してもユーザー情報依存のバリデーションが先行評価されることになり、評価順序の責務分離に反します。

## 3. 推奨される修正案

1. `03_users.md` 3.2.4 節の「リクエスト評価順序」を以下のように修正し、構文検証（ステップ2: 8〜128文字、英大文字/小文字/数字/記号のうち3種以上）と、ドメインビジネスルール検証（ステップ5: ユーザー名/メールローカル部含有チェックおよび同一パスワードチェック）を適切に分離してください。

   **評価順序の修正案**:
   - **ステップ2（リクエスト構文・入力バリデーション `400 Bad Request`）**: リクエストボディの JSON 形式、必須フィールド（`current_password`, `new_password`）の有無、および `new_password` の文字数（8〜128文字）・使用可能文字種（英大文字/小文字/数字/記号のうち3種以上）を検証。不備時は即座に 400 `BAD_REQUEST` を返却（遅延・カウンター加算なし）。
   - **ステップ5（ビジネスルール検証 `422 Unprocessable Entity`）**: 再認証成功後、`new_password` に対象ユーザーのユーザー名・メールアドレスのローカル部（4文字以上の場合、大文字・小文字を区別せず Case-Insensitive）が含まれていないか検証（含有時は 422 `INVALID_PASSWORD_CONTENT`）、および照合済みの `current_password` と同一でないか検証（同一時は 422 `SAME_AS_CURRENT_PASSWORD`）。

2. `03_users.md` 3.2.4 節の `##### Errors` に `422 INVALID_PASSWORD_CONTENT` を追加し、`400 Bad Request` の説明からユーザー名/メール含有違反の記述を削除・移動してください。

---

## 修正完了報告

- **Resolved At**: 2026-08-17 17:50:00
- **Status**: Resolved

### 実施した修正内容
`03_users.md` 3.2.4 節 (`PATCH users/{user_id}/password`) の「リクエスト評価順序」において、ステップ2では新パスワード自体の構文チェック（文字数8〜128文字、3種以上の文字種）のみを行い、ユーザー名・メールアドレスローカル部の含有チェックを認可・再認証成功後のステップ5（ビジネスルール検証）に移動しました。含有違反時は `422 Unprocessable Entity` (code: `"INVALID_PASSWORD_CONTENT"`) を返却するよう修正し、Errors 一覧にも追加して全体仕様書間（`01_overview.md`, `02_auth.md`）の整合性を確保しました。

### 変更したファイル
- [03_users.md](docs/design/api_design/03_users.md)
