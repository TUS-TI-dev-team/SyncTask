# VS Code での Go / Gin バックエンド デバッグ手順書

本ドキュメントは、VS Code 上で Go (Gin) バックエンドおよびテストコードをブレークポイント付きでデバッグ実行するための環境構築と実行手順のメモです。他環境や新規セットアップ時にも本手順を参照することで再現可能です。

---

## 1. 前提ツールの準備

### 1-1. VS Code 拡張機能「Go」
VS Code の拡張機能マーケットプレイスで以下をインストールします。
- **名前**: Go
- **提供元**: Go Team at Google (`golang.go`)

### 1-2. デバッガ Delve (`dlv`) のインストール
Go のデバッグには公式デバッガ **Delve** が必要です。

ターミナルで以下のコマンドを実行してインストールします：
```bash
go install github.com/go-delve/delve/cmd/dlv@latest
```

> **PATH の確認**:
> インストールされた `dlv` は通常 `$HOME/go/bin` に配置されます。`dlv version` で確認できない場合は、シェル設定（`~/.bashrc` や `~/.zshrc`）に以下が設定されているか確認してください。
> ```bash
> export PATH=$PATH:$(go env GOPATH)/bin
> ```
> *(VS Code 内で `Ctrl+Shift+P` -> `Go: Install/Update Tools` を実行し、`dlv` にチェックを入れてインストールすることも可能です)*

---

## 2. VS Code デバッグ設定 (`.vscode/launch.json`)

プロジェクトルート直下に `.vscode/launch.json` を配置します。

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Debug Backend Server",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceFolder}/backend/main.go",
      "cwd": "${workspaceFolder}/backend",
      "env": {
        "GIN_MODE": "debug"
      }
    },
    {
      "name": "Debug Current Test File",
      "type": "go",
      "request": "launch",
      "mode": "test",
      "program": "${fileDirname}",
      "cwd": "${workspaceFolder}/backend",
      "env": {
        "GIN_MODE": "test"
      }
    }
  ]
}
```

### 設定項目のポイント
- **`program`**:
  - サーバー起動: エントリーポイントとなる `${workspaceFolder}/backend/main.go` を指定。
  - テスト実行: 現在アクティブなファイルのディレクトリ `${fileDirname}` を指定することで、開いているテストパッケージを実行可能。
- **`cwd`**:
  - `${workspaceFolder}/backend` に指定します。これにより、相対パスでのファイル参照や設定読み込みが `backend/` 基準で正しく動作します。
- **`env`**:
  - 必要に応じて環境変数（`GIN_MODE`, `DB_HOST`, `DB_PORT` 等）を上書きできます。

---

## 3. 事前準備 (DB 起動)

バックエンドサーバーは起動時に PostgreSQL に接続してマイグレーションを実行するため、事前に DB コンテナを起動しておきます。

```bash
docker compose up -d db
```

> ローカル実行時のデフォルト接続先（`localhost:5432`, ユーザー: `synctask`, DB: `synctask_dev`）と `docker-compose.yml` のポートフォワーディング（`5432:5432`）が一致しているため、追加の環境変数設定は不要です。

---

## 4. デバッグの実行手順

### 4-1. バックエンドサーバーのデバッグ
1. **ブレークポイントの設置**:
   - `backend/main.go` や `backend/handler/*.go` の任意のコード行番号の左側をクリックして赤丸を付けます。
2. **デバッグの開始**:
   - VS Code 左のアクティビティバーから **「実行とデバッグ」**（`Ctrl+Shift+D`）を開きます。
   - 上部ドロップダウンから **「Debug Backend Server」** を選択します。
   - **`F5`** キー（または緑の再生ボタン）を押下します。
3. **リクエスト送信と調査**:
   - 下部デバッグコンソールに `Listening and serving HTTP on :8080` と表示されたら起動完了です。
   - ブラウザ、curl、Postman などからエンドポイント（例: `http://localhost:8080/health`）にリクエストを送ると、ブレークポイントで一時停止します。
   - **変数パネル**での値確認、**ウォッチ式**の評価、**ステップオーバー (`F10`)**、**ステップイン (`F11`)** などのデバッグ操作が可能です。

### 4-2. 単体テストのデバッグ
1. デバッグしたいテストコード（例: `backend/handler/*_test.go`）を開きます。
2. 止めたい箇所にブレークポイントを設置します。
3. 以下のいずれかの方法でデバッグを開始します：
   - テスト関数の上に表示される **`debug test`** リンク（CodeLens）をクリックする。
   - 構成ドロップダウンから **「Debug Current Test File」** を選択して **`F5`** を押す。

---

## 5. よくあるトラブルシューティング

| 症状 | 原因 | 対処法 |
|---|---|---|
| `Failed to launch: could not launch process: exec: "dlv": executable file not found in $PATH` | Delve がインストールされていない、または PATH に含まれていない | `go install github.com/go-delve/delve/cmd/dlv@latest` を実行し、`export PATH=$PATH:$(go env GOPATH)/bin` を確認 |
| `DB接続に失敗しました: dial tcp 127.0.0.1:5432: connect: connection refused` | PostgreSQL が起動していない | `docker compose up -d db` で DB コンテナを起動する |
| `bind: address already in use` (8080) | すでに別のプロセス（またはコンテナ）が 8080 を使用している | `docker compose stop backend` や `lsof -i :8080` で既存プロセスを停止する |
| `import cycle` や パッケージが見つからないエラー | `cwd` または `program` のパス指定がずれている | `launch.json` の `cwd` が `${workspaceFolder}/backend` になっているか確認 |
