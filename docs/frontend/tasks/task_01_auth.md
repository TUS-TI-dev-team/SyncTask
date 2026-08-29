# Task 01: 認証系画面・コンポーネント群の移植 (Auth Pages)

## 1. 担当概要

本タスクは、ログイン、アカウント新規登録、パスワードリセット、および関連する OTP 認証フローの画面・コンポーネントを `SyncTask-Design-Idea` から `SyncTask/frontend` に移植する作業です。

---

## 2. 移行対象ファイル一覧

| 移行元 (Design-Idea) | 移行先 (SyncTask/frontend) | 概要 |
| --- | --- | --- |
| `components/auth/otp-input.tsx` | `components/auth/otp-input.tsx` | 6桁OTP入力用マス目コンポーネント |
| `components/auth/otp-panel.tsx` | `components/auth/otp-panel.tsx` | OTP入力パネル（15分タイマー、60秒再送クールダウン、マスク表示） |
| `app/login/page.tsx` | `app/login/page.tsx` | ログイン画面 |
| `app/signup/page.tsx` | `app/signup/page.tsx` | アカウント作成（情報入力）画面 |
| `app/signup/otp/page.tsx` | `app/signup/otp/page.tsx` | アカウント作成 OTP認証画面 |
| `app/reset-password/page.tsx` | `app/reset-password/page.tsx` | パスワードリセット（メール入力）画面 |
| `app/reset-password/otp/page.tsx` | `app/reset-password/otp/page.tsx` | パスワードリセット OTP認証画面 |
| `app/reset-password/new/page.tsx` | `app/reset-password/new/page.tsx` | 新パスワード再設定画面 |

---

## 3. 仕様書・設計との整合ポイント (`screen_design.md`)

- **OTP 送信先メールアドレスのマスク仕様**:
  - 「始め4文字＋ドメイン名のみ表示、間は固定10文字のマスク（`****...`）」のフォーマットになっていること。
- **OTP タイマー & クールダウン**:
  - 15分セッションタイマー表示。
  - 再送信ボタン押下時の60秒カウントダウンおよび非活性（Disabled）化。
- **インラインエラー表示領域**:
  - パスワード不一致、必須未入力、無効なメール形式等のバリデーションメッセージが表示可能な構造になっていること。
- **レイアウト適用**:
  - 全ての認証画面が `components/layouts/auth-shell.tsx` を用いて一貫したカードデザインでラップされていること。

---

## 4. 作業手順

1. `components/auth/otp-input.tsx` および `components/auth/otp-panel.tsx` を作成・配置する。
2. `app/login/page.tsx` を作成する。
3. `app/signup/page.tsx` および `app/signup/otp/page.tsx` を作成する。
4. `app/reset-password/page.tsx`、`app/reset-password/otp/page.tsx`、`app/reset-password/new/page.tsx` を作成する。
5. 各画面間のリンク遷移（例: ログイン ↔ 新規登録 ↔ パスワードリセット）が正しく繋がっていることを確認する。

---

## 5. 完了確認チェックリスト

- [ ] `npx tsc --noEmit` で型エラーが発生しない
- [ ] ブラウザで `/login` にアクセスし、フォームおよびリンクが正常に表示される
- [ ] `/signup` -> `/signup/otp` の画面遷移・UIが確認できる
- [ ] `/reset-password` -> `/reset-password/otp` -> `/reset-password/new` の画面遷移・UIが確認できる
- [ ] OTP画面のタイマー表示および再送カウントダウンの動作が確認できる
