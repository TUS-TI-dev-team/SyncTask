# Task 05: プロフィール・アカウント管理画面群の移植 (Profile & Account Pages)

## 1. 担当概要

本タスクは、ユーザープロフィール表示・編集、メールアドレス変更OTP、パスワード変更、およびアカウント削除（再認証）に関する画面・コンポーネントを `SyncTask-Design-Idea` から `SyncTask/frontend` に移植する作業です。

---

## 2. 移行対象ファイル一覧

| 移行元 (Design-Idea) | 移行先 (SyncTask/frontend) | 概要 |
| --- | --- | --- |
| `components/profile/profile-view.tsx` | `components/profile/profile-view.tsx` | プロフィール表示ビュー |
| `components/profile/profile-edit-view.tsx` | `components/profile/profile-edit-view.tsx` | プロフィール編集ビュー（メール変更確認ダイアログ含む） |
| `app/profile/page.tsx` | `app/profile/page.tsx` | プロフィール表示画面 |
| `app/profile/edit/page.tsx` | `app/profile/edit/page.tsx` | プロフィール編集画面 |
| `app/profile/otp/page.tsx` | `app/profile/otp/page.tsx` | メールアドレス変更 OTP入力画面 |
| `app/profile/password/page.tsx` | `app/profile/password/page.tsx` | パスワード変更画面 |
| `app/profile/delete/page.tsx` | `app/profile/delete/page.tsx` | アカウント削除（パスワード再認証）画面 |

---

## 3. 仕様書・設計との整合ポイント (`screen_design.md`)

- **プロフィール表示画面 (`/profile`)**:
  - アカウント情報（ユーザー名・メールアドレス）の表示。
  - 各種設定画面（編集・パスワード変更・アカウント削除）への導線ボタン。
- **プロフィール編集画面 (`/profile/edit`)**:
  - ユーザー名・メールアドレスの入力フォーム。
  - メールアドレス変更時は、確認ダイアログを表示した上で `/profile/otp` へ誘導する挙動。
- **アカウント関連 OTP 画面 (`/profile/otp`)**:
  - 新メールアドレス宛のマスク表示（4文字＋ドメイン）。
  - 15分セッションタイマーおよび60秒再送クールダウン。
- **パスワード変更画面 (`/profile/password`)**:
  - 「現在のパスワード」「新パスワード」「新パスワード（確認）」の3つの入力欄（マスク表示）。
- **アカウント削除画面 (`/profile/delete`)**:
  - 本人確認のためのパスワード再認証入力欄。
  - 削除実行前の最終確認ポップアップ（AlertDialog）。

---

## 4. 作業手順

1. `components/profile/profile-view.tsx` および `profile-edit-view.tsx` を作成・配置する。
2. `app/profile/page.tsx` および `app/profile/edit/page.tsx` を作成する。
3. `app/profile/otp/page.tsx` を作成する。
4. `app/profile/password/page.tsx` を作成する。
5. `app/profile/delete/page.tsx` を作成する。
6. `useStore` のモックプロファイルデータとの連携および各画面遷移を確認する。

---

## 5. 完了確認チェックリスト

- [ ] `npx tsc --noEmit` で型エラーが発生しない
- [ ] ブラウザで `/profile` にアクセスし、プロファイル情報が表示される
- [ ] `/profile/edit` でユーザー名・メールアドレスの編集および保存ができる
- [ ] `/profile/password` でパスワード変更画面のUIが表示される
- [ ] `/profile/delete` でアカウント削除画面が表示され、削除確認ダイアログが機能する
