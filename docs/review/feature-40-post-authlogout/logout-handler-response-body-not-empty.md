# ログアウト成功時（200 OK）のレスポンスボディが空になっていない

- **Status**: Resolved
- **Severity**: Major
- **Created At**: 2026-10-30 18:52:00
- **Target Files**:
  - [logout.go](backend/handler/logout.go)
  - [logout_test.go](backend/handler/logout_test.go)

## 1. 問題の概要
`POST /api/auth/logout` の成功レスポンス（`200 OK`）は、仕様上「レスポンスボディなし（空）」と定義されていますが、`c.JSON(http.StatusOK, nil)` を呼び出していたため、レスポンスボディに JSON の `"null"`（4バイト）が出力されていました。

## 2. 詳細な指摘内容
1. API設計書（`docs/design/api_design/02_auth.md` 3.1.5節）および業務設計書（`docs/design/process_design/05_logout.md` 4.1節）では、成功応答について「レスポンスボディなし（JSONオブジェクトなし）」と規定されています。
2. `backend/handler/logout.go` の `c.JSON(http.StatusOK, nil)` を `c.Status(http.StatusOK)` に変更する必要があります。

---

## 修正完了報告

- **Resolved At**: 2026-10-30 19:27:00
- **Status**: Resolved

### 実施した修正内容
- `backend/handler/logout.go` (L60) において `c.JSON(http.StatusOK, nil)` から `c.Status(http.StatusOK)` へ修正されたことを確認しました。

### 変更したファイル
- [logout.go](backend/handler/logout.go)
