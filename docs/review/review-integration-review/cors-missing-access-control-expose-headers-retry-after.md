# CORS仕様における Access-Control-Expose-Headers (Retry-After) の定義欠落

- **Status**: Resolved
- **Severity**: Major
- **Created At**: 2026-08-22 13:49:30
- **Target Files**:
  - [tech_stack.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/tech_stack.md)
  - [01_overview.md](file:///C:/Users/kazuh/Programming/repos/SyncTask/docs/design/api_design/01_overview.md)

## 1. 問題の概要

`docs/design/api_design/01_overview.md` 1.4節において、「`429 Too Many Requests` 返却時には、クライアントが再試行可能になるまでの待機時間を通知するため、標準の `Retry-After: <秒数>` レスポンスヘッダー（例: `Retry-After: 60`, `Retry-After: 900`）を必須で付与する」と規定されています。

しかし、`docs/design/tech_stack.md` 5節および `docs/design/api_design/01_overview.md` 1.2節の CORS 仕様には、`Access-Control-Expose-Headers` の定義が一切記載されていません。

W3C / Fetch Standard（CORS）の仕様上、CORSセーフリスト外のレスポンスヘッダーである `Retry-After` をフロントエンド（JavaScript / `fetch` / `axios` 等）が読み取るためには、サーバーが `Access-Control-Expose-Headers: Retry-After` をレスポンスヘッダーとして出力することが必須です。これが欠落していると、ブラウザのセキュリティ機能により JavaScript から `response.headers.get("Retry-After")` へのアクセスが遮断され、UI側で再送カウントダウンやロック待機時間の表示制御が行えなくなります。

## 2. 詳細な指摘内容

### 該当箇所の定義

1. **`docs/design/tech_stack.md` (セクション5)**:
   ```markdown
   - **CORS制御**:
     - 許可オリジン: 環境変数 `FRONTEND_URL`（例: `http://localhost:3000`, `https://synctask.app`）
     - `Access-Control-Allow-Credentials: true`（Cookie送受信用）
     - `Access-Control-Allow-Methods: GET, POST, PUT, PATCH, DELETE, OPTIONS`
     - `Access-Control-Allow-Headers: Content-Type, X-CSRF-Token, Authorization`
     - プリフライト（`OPTIONS`）レスポンス: HTTP 204 No Content
   ```

2. **`docs/design/api_design/01_overview.md` (セクション1.2)**:
   ```markdown
   - **許可オリジン (`Access-Control-Allow-Origin`)**: 環境変数 `FRONTEND_URL` で指定された単一オリジンを明示的に許可します
   - **認証情報許可 (`Access-Control-Allow-Credentials`)**: `true`
   - **許可メソッド (`Access-Control-Allow-Methods`)**: `GET, POST, PUT, PATCH, DELETE, OPTIONS`
   - **許可ヘッダー (`Access-Control-Allow-Headers`)**: `Content-Type, X-CSRF-Token, Authorization`
   - **プリフライトキャッシュ (`Access-Control-Max-Age`)**: `86400` (24時間)
   - **プリフライト（OPTIONS）レスポンス**: 条件に合致する OPTIONS リクエストに対しては `204 No Content` を返却します。
   ```

3. **影響**:
   - `01_overview.md` 1.4節で「OTP再送クールダウン中（60秒）」や「IPレートリミット超過（15分遮断）」の際に `Retry-After` を付与してクライアントに通知する設計となっていますが、`Access-Control-Expose-Headers: Retry-After` がないため、フロントエンド（Next.js / ブラウザ）の JavaScript から `Retry-After` の値を取得できず、要件定義された待機カウントダウン表示などの機能が動作しません。

## 3. 推奨される修正案

`docs/design/tech_stack.md` および `docs/design/api_design/01_overview.md` の CORS 仕様に、以下の項目を追加明記してください：

- **公開ヘッダー (`Access-Control-Expose-Headers`)**: `Retry-After`

---

## 修正完了報告

- **Resolved At**: 2026-08-22 13:55:00
- **Status**: Resolved

### 実施した修正内容
CORS 仕様に `Access-Control-Expose-Headers: Retry-After` を明記し、フロントエンドの JavaScript から 429 応答時の待機時間ヘッダーを読み取り可能にしました。

### 変更したファイル
- [tech_stack.md](file:///C:\Users\kazuh\Programming\repos\SyncTask\docs\design\tech_stack.md)
- [01_overview.md](file:///C:\Users\kazuh\Programming\repos\SyncTask\docs\design\api_design\01_overview.md)
