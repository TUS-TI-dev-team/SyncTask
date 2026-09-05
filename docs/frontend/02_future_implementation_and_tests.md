# 02. 今後の機能追加およびテスト実装計画 (Future Roadmap & Tests)

本ドキュメントでは、UI移行（Phase 1 & 2）完了後に実施する **バックエンドAPI連携・状態管理の刷新・バリデーション実装**、および **テスト自動化** の方針と主要ToDoリストを定義します。

---

## 1. 今後のフロントエンド・アーキテクチャ方針

### 1.1 API クライアント層 & データフェッチ設計
- **推奨技術**: `TanStack Query (React Query)` または Next.js App Router ネイティブの `fetch` + Server Actions
- **APIエンドポイント対応**: `docs/design/api_design/`（Ginバックエンド）に準拠したクライアントメソッド群（`lib/api/`）を定義。
- **ベースURL**: 環境変数 `NEXT_PUBLIC_API_BASE_URL`（開発時: `http://localhost:8080/api/v1`）。

### 1.2 認証・セッション管理アーキテクチャ
- **Cookie セッション連携**:
  - バックエンド発行の `session_id`（HttpOnly, Secure, SameSite=Lax）を用いたCookie認証。
  - Next.js の `middleware.ts` によるルーティングガード（未ログイン時の `/login` リダイレクト、ログイン時の `/home` 自動遷移）。
- **OTP認証セッション**:
  - 新規登録、パスワードリセット、メールアドレス変更の一時セッション管理。

### 1.3 フォームバリデーション & エラーハンドリング
- **フォームバリデーション**: `react-hook-form` + `zod` による宣言的スキーマバリデーション。
- **APIエラー通知**: `sonner`（Toast）およびインラインエラーバナーによるエラーコード別メッセージ表示（例: `AUTH_001`, `VALIDATION_ERROR` 等）。

---

## 2. テスト実装方針

### 2.1 採用推奨テストツール
| カテゴリ | ツール | 目的・対象 |
| --- | --- | --- |
| **単体 / コンポーネントテスト** | [Vitest](https://vitest.dev/) + [React Testing Library](https://testing-library.com/) | 状態ロジック、UIコンポーネントの表示・ユーザー操作の検証 |
| **E2E (End-to-End) テスト** | [Playwright](https://playwright.dev/) | ログイン〜タスク作成・編集〜ログアウトまでの一連のユーザーシナリオ検証 |

### 2.2 テスト自動化方針
- CI（GitHub Actions）上で `npm run lint`、`tsc --noEmit`、`vitest run` を自動実行。

---

## 3. 機能別 実装 & テスト ToDo リスト

### ① 認証・アカウントフロー (Auth & Account)
- [ ] **APIクライアント実装**:
  - [ ] `POST /api/v1/auth/signup`（登録要求・OTP発行）
  - [ ] `POST /api/v1/auth/signup/verify`（OTP検証・セッション発行）
  - [ ] `POST /api/v1/auth/login`（ログイン・セッション発行）
  - [ ] `POST /api/v1/auth/logout`（ログアウト）
  - [ ] `POST /api/v1/auth/password-reset/request` & `/verify` & `/confirm`（パスワードリセット一連）
- [ ] **画面連携**:
  - [ ] 各フォームへの Zod バリデーション組み込み
  - [ ] エラーコード別（無効なOTP、ロックアウト等）のToast/エラー表示
- [ ] **テスト項目**:
  - [ ] 正常ログイン / 失敗時エラー表示の単体テスト
  - [ ] OTPタイマー失効および再送信クールダウンの動作テスト
  - [ ] 新規登録からログイン完了までのPlaywright E2Eテスト

### ② ホーム画面・ダッシュボード (Home Dashboard)
- [ ] **APIクライアント実装**:
  - [ ] `GET /api/v1/tasks?priority=high`（高優先度タスク取得）
  - [ ] `GET /api/v1/tasks?deadline_near=true`（締切間近タスク取得）
  - [ ] `GET /api/v1/tasks?pinned=true`（ピン止めタスク取得）
- [ ] **画面連携**:
  - [ ] タブ切り替え時のデータフェッチ・キャッシュ制御
  - [ ] 完了タスク表示/非表示トグルとAPIクエリ連動
  - [ ] ページネーション（$N > 10$ 省略ルール）のサーバー連動
- [ ] **テスト項目**:
  - [ ] タブ切り替えによる表示タスク一覧の切り替えテスト
  - [ ] ページネーション計算ロジックの単体テスト

### ③ タスク管理・カレンダー (Tasks & Calendar)
- [ ] **APIクライアント実装**:
  - [ ] `GET /api/v1/tasks`（キーワード検索・フィルタ・ページネーション対応）
  - [ ] `POST /api/v1/tasks`（タスク作成・繰り返し一括作成）
  - [ ] `PATCH /api/v1/tasks/:id`（ステータス更新・個別更新）
  - [ ] `DELETE /api/v1/tasks/:id`（タスク削除）
  - [ ] `POST /api/v1/tasks/:id/pin`（ピン留め切り替え）
- [ ] **画面連携**:
  - [ ] 検索入力デバウンス処理
  - [ ] カレンダー月/週表示切り替えおよび該当日タスクマッピング
  - [ ] カレンダー上での直接ステータス変更（楽観的UI更新）
  - [ ] タスク作成・編集・削除ダイアログのAPI連携
- [ ] **テスト項目**:
  - [ ] 検索・絞り込みフィルターの結合テスト
  - [ ] カレンダーグリッドの日付計算・タスク配置テスト
  - [ ] 繰り返しタスク作成（曜日・期間指定）のロジック検証
  - [ ] タスクCRUD操作のPlaywright E2Eテスト

### ④ プロフィール・アカウント管理 (Profile & Settings)
- [ ] **APIクライアント実装**:
  - [ ] `GET /api/v1/users/me`（プロフィール取得）
  - [ ] `PATCH /api/v1/users/me`（ユーザー名変更）
  - [ ] `POST /api/v1/users/me/email/request` & `/verify`（メールアドレス変更OTP）
  - [ ] `POST /api/v1/users/me/password`（パスワード変更）
  - [ ] `DELETE /api/v1/users/me`（アカウント削除）
- [ ] **画面連携**:
  - [ ] メールアドレス変更時の確認モーダルおよびログアウトリダイレクト処理
  - [ ] パスワード変更/アカウント削除成功時のセッション破棄・リダイレクト処理
- [ ] **テスト項目**:
  - [ ] プロフィール編集のバリデーションテスト
  - [ ] パスワード再認証および削除実行の結合テスト
