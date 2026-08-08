# Process Design (処理設計)

## 概要

システムの主要ユースケースにおける処理フローおよびシーケンス図です（Notion最新版と完全同期）。

---

## 1. アカウント作成

```mermaid
sequenceDiagram
	participant MailServer
	actor User
	participant Frontend
	participant Backend
	participant DB
	
	User->>Frontend: アカウント作成ボタンクリック
	Frontend-->>User: アカウント情報入力画面
	User->>Frontend: アカウント情報入力
	Frontend->>Backend: POST auth/register/request-otp
	Backend->>Backend: アカウント情報のバリデーション
	alt アカウント情報バリデーションエラー
		Backend-->>Frontend: アカウント情報バリデーションエラー
		Frontend-->>User: アカウント情報バリデーションエラー
	end
	Backend->>DB: メールアドレスが既に使用されていないかチェック
	DB->>Backend: 入力されたメールアドレスの使用状況

	alt メールアドレスが使用済でない
		Backend->>DB: OTPセッション作成、登録予定データ登録
		Backend->>MailServer: OTPメール送信
	end
	%% もしメールアドレスが使用済なら、成功したように見せつつメールは送らない
	Backend-->>Frontend: アカウント情報仮登録成功
	Frontend-->>User: OTP入力画面
	User->>MailServer: OTPメール確認
	MailServer-->>User: OTP
	User->>Frontend: OTP入力 & OTPセッションID
	Frontend->>Backend: OTP情報 & OTPセッションID送信
	Backend->>DB: OTP & OTPセッション情報の検索
	break OTPセッションが無い
		DB-->>Backend: OTP無し、exit
	end
	DB-->>Backend: OTP & OTPセッション情報
	Backend->>Backend: OTPの検証
	alt OTP検証失敗(5回未満)
		Backend->>DB: OTP検証失敗回数を1増加
		Backend-->>Frontend: OTP検証失敗
		Frontend-->>User: OTP検証失敗メッセージ
	else OTP検証失敗(5回)
		Backend->>DB: 現在のOTPセッションを無効化
		Backend->>DB: OTPセッション再作成、登録予定データ登録
		Backend->>MailServer: OTPメール送信
		Backend-->>Frontend: OTP検証失敗、OTP入力画面に戻る
		Frontend-->>User: <OTP入力画面に戻る>
	end
	Backend->>Backend: ログインセッション発行
	Backend->>DB: アカウントテーブルにアカウント追加、OTPセッション削除、ログインセッション追加
	note left of DB: ここでアカウント作成完了
	Backend-->>Frontend: ログインセッション
	Frontend-->>User: ホーム画面
```

---

## 2. アカウント編集

プロフィール編集画面にて「決定」ボタンがクリックされた後の処理

```mermaid
sequenceDiagram
	participant MailServer
	actor User
	participant Frontend
	participant Backend
	participant DB

	User->>Frontend: 入力されたユーザ名, メールアドレス
	Frontend->>Frontend: バリデーション検証
	alt フロントバリデーションエラー
		Frontend-->>User: 入力情報バリデーションエラー
	end
	Frontend->>Backend: PUT users/{user_id}
	activate Backend
	%%put内でユーザ名の検証 - メアド検証 - OTP飛ばす - 画面遷移
	opt PUT user
		alt バックバリデーションエラー
			Backend-->>Frontend: 入力情報バリデーションエラー
			Frontend-->>User: 入力情報バリデーションエラー
		end
		Backend->>DB: user_idに紐づくメールアドレスを要求
		DB-->>Backend: メールアドレスを送信
		Backend->>Backend: 入力情報と登録済みのメールアドレス比較
		alt メールアドレスが更新されていた場合
			Backend->>DB: OTPセッション作成、新しいユーザー情報を仮登録
			Backend->>MailServer: OTPメール送信
			Backend-->>Frontend: OTP入力画面遷移
			Frontend-->>User: OTP入力画面
		end
	end

	User->>Frontend: OTP入力
	Frontend->>Backend: OTP
	Backend->>Backend: OTP検証（諸々省略）
	alt OTP検証が成功した場合
		Backend->>DB: OTPセッション満了、アカウント情報更新
		Backend-->>Frontend: OTP検証成功
		Frontend-->>User: プロフィール表示画面
	end
	%%変更完了を通知しプロフィール画面へ
```

---

## 3. アカウント削除

```mermaid
sequenceDiagram
	actor User
	participant Frontend
	participant Backend
	participant DB
	
	User->>Frontend: アカウント削除ボタン押下
	Frontend->>Backend: DELETE users/{user_id} with セッション
	Backend->>DB: セッションからユーザーIDを検索
	Backend->>DB: user_id のユーザーを検索
	DB-->>Backend: 検索結果
	break user_idが、セッションのユーザーIDと一致しない
		Backend-->>Frontend: 403 権限不足
	end
	Backend->>DB: user_id に紐づいた<br>すべてのログインセッション・OTPセッションを無効化<br>タスクも削除
	Backend->>DB: user_id のユーザーを論理削除
	Backend-->>Frontend: 200 削除完了
	Frontend-->>User: ログイン画面に戻る
```

---

## 4. ログイン

```mermaid
sequenceDiagram
	actor User
	participant Frontend
	participant Backend
	participant DB

	User->>Frontend: Cookie{session_id}
	Frontend->>Backend: セッションID
	Backend->>DB: セッションID検証
	DB-->>Backend: 検証結果
	Backend-->>Frontend: 検証結果
	alt セッションが有効
		Frontend-->>User: ホーム画面へリダイレクト
	else セッションが無効 or 存在しない
		User->>Frontend: ログイン情報入力, ログインボタンクリック
		Frontend->>Frontend: 入力情報バリデーション検証
		alt フロント検証バリデーションエラー
			Frontend-->>User: 入力情報バリデーションエラー
		end
		Frontend->>Backend: POST auth/login/
		activate Backend
		Backend->>Backend: 入力情報のバリデーション
		alt バック検証バリデーションエラー
			Backend-->>Frontend: 入力情報バリデーションエラー
			Frontend-->>User: 入力情報バリデーションエラー
		end
		Backend->>DB: 入力されたユーザ情報のパスワード情報を要求
		activate DB
		DB-->>Backend: 取得したパスワード情報，<br>取得できなかった場合は空情報を送信
		deactivate DB
		Backend->>Backend: ユーザ情報認証，パスワード認証
		alt 認証エラー
			Backend-->>Frontend: 入力情報認証エラー
			Frontend-->>User: 入力情報認証エラー
		end
		Backend->>Backend: ログインセッションIDを発行
		Backend->>DB: ログインセッションIDの登録
		activate DB
		DB-->>Backend: 完了通知
		deactivate DB
		Backend-->>Frontend: ログインセッションIDをCookieに保存
		deactivate Backend
		Frontend-->>User: ホーム画面へリダイレクト
	end
```

---

## 5. ログアウト

```mermaid
sequenceDiagram
	actor User as ユーザー
	participant FE as フロントエンド
	participant BE as バックエンド
	participant DB as セッションDB

	User->>FE: ログアウトボタンをクリック

	Note over FE,DB: 正常系
	FE->>BE: POST /api/logout (Cookie: session_id)
	activate BE
	BE->>DB: session_id でセッションを検索
	activate DB
	DB-->>BE: セッション情報を返却（有効）
	deactivate DB
	BE->>DB: 該当セッションをDELETE
	activate DB
	DB-->>BE: 削除完了
	deactivate DB
	BE-->>FE: 200 OK セッションリセット
	deactivate BE
	FE->>FE: ブラウザのCookie（session_id）を削除
	FE->>FE: グローバル状態をリセット
	FE->>User: ログイン画面へリダイレクト

	Note over FE,DB: 異常系（セッション無効・期限切れ）
	FE->>BE: POST /api/logout (Cookie: session_id)
	activate BE
	BE->>DB: session_id でセッションを検索
	activate DB
	DB-->>BE: 該当セッションなし（無効/期限切れ）
	deactivate DB
	BE-->>FE: 401 Unauthorized
	deactivate BE
	FE->>FE: ブラウザのCookie（session_id）を削除
	Note right of FE: サーバー側に削除対象が無くてもクライアント側は確実にCookieを破棄する
	FE->>User: ログイン画面へリダイレクト
```

---

## 6. パスワードリセット

### 1. パスワード再設定要求フェーズ

```mermaid
sequenceDiagram
	autonumber
	actor User as ユーザー
	participant FE as Next.js (Frontend)
	participant BE as Gin (Backend)
	participant DB as DB (Userテーブル, OTPセッションテーブル)
	participant Mail as メールサーバー/メールAPI

	Note over User,Mail: フェーズ1: OTP発行・送信
	User->>FE: 「パスワードを忘れた」をクリック
	FE->>User: メールアドレス入力フォーム表示
	User->>FE: メールアドレスを入力・送信

	FE->>BE: POST /api/auth/password-reset/request-otp { email }
	BE->>DB: SELECT * FROM users WHERE EMAIL = ?

	alt ユーザーが存在しない
		DB-->>BE: 該当ユーザーなし
		Note right of BE: メール存在の有無を秘匿するため<br/>成功時と同じレスポンスを返す
		BE-->>FE: 200 OK (汎用メッセージ)
		FE-->>User: 「メールをご確認ください」とOTP入力画面を表示
	else ユーザーが存在する
		DB-->>BE: ユーザー情報を返す

		BE->>BE: OTP(8桁英数字)を生成
		BE->>DB: 既存の有効なOTPセッションがあれば失効
		BE->>DB: INSERT INTO otp_sessions<br/>(UNAME, otp_hash, expires_at, category, <br>send_count=0, attempt_count=0, status='active')
		DB-->>BE: 保存完了

		BE->>Mail: OTPメール送信リクエスト(同期)
		Mail-->>BE: 送信結果

		alt 送信失敗
			BE-->>FE: 500 Internal Server Error
			FE-->>User: 「送信に失敗しました。再度お試しください」
		else 送信成功
			BE-->>FE: 200 OK (汎用メッセージ)
			FE-->>User: OTP入力画面を表示
		end
	end
```

### 2. OTP入力・検証フェーズ

```mermaid
sequenceDiagram
	autonumber
	actor User as ユーザー
	participant FE as Next.js (Frontend)
	participant BE as Gin (Backend)
	participant DB as DB (Userテーブル, OTPセッションテーブル)
	participant Mail as メールサーバー/メールAPI

	Note over User,DB: フェーズ2: OTP検証(リトライ・再送信・期限切れ)
	User->>FE: メール記載のOTPを入力フォームへ

	loop OTP検証(最大試行回数まで)
		User->>FE: OTPを入力・送信
		FE->>BE: POST /api/auth/verify-otp { otp } <br>Cookie: OTPセッション
		BE->>DB: CookieのOTPセッションIDから<br>対応するOTPセッションを取得

		alt OTPセッションが存在しない/失効済み
			DB-->>BE: セッションなし or status != 'active'
			BE-->>FE: 410 Gone (セッション無効)
			FE-->>User: 「再度パスワード再設定をお試しください」
			Note over User,FE: フェーズ1へ戻る
		else 有効期限切れ
			DB-->>BE: expires_at < 現在時刻
			BE->>DB: UPDATE otp_sessions SET status='expired'
			BE-->>FE: 410 Gone (期限切れ)
			FE-->>User: フェーズ1へ戻る
		else 有効なセッションあり
			DB-->>BE: セッション情報(otp_hash, attempt_count)を返す
			BE->>BE: 入力OTPとotp_hashを比較

			alt OTP不一致
				BE->>DB: UPDATE otp_sessions SET attempt_count += 1

				alt リトライ回数が制限を超えた
					BE->>DB: UPDATE otp_sessions SET status='locked'
					BE-->>FE: 429 Too Many Requests (リトライ上限)
					FE-->>User: 「試行回数の上限に達しました。<br/>最初からやり直してください」
					Note over User,FE: フェーズ1へ戻る(セッション失効)
				else リトライ可能
					BE-->>FE: 400 Bad Request (不一致, 残り回数)
					FE-->>User: 「OTPが正しくありません(残りN回)」
				end
			else OTP一致
				BE->>DB: UPDATE otp_sessions SET status='verified', expire_at=今から5分後
				BE-->>FE: 200 OK
				FE-->>User: 新しいパスワード入力画面へ遷移
			end
		end
	end

	opt ユーザーが「再送信」ボタンを押す
		User->>FE: 再送信をクリック
		FE->>BE: POST /api/auth/resend-otp<br>Cookie: OTPセッションID
		BE->>DB: 既存OTPセッションを失効(status='expired')
		BE->>BE: 新しいOTP(8桁英数字)を生成
		BE->>DB: INSERT INTO otp_sessions<br/>(新しいレコード, attempt_count=0, status='active')
		BE->>Mail: OTPメール送信リクエスト(同期)
		Mail-->>BE: 送信結果
		BE-->>FE: 200 OK
		FE-->>User: 「再送信しました」と表示
		Note over User,FE: OTP検証ループへ戻る
	end
```

### 3. パスワード更新フェーズ

```mermaid
sequenceDiagram
	autonumber
	actor User as ユーザー
	participant FE as Next.js (Frontend)
	participant BE as Gin (Backend)
	participant DB as DB (Userテーブル, OTPセッションテーブル)
	participant Mail as メールサーバー/メールAPI

	Note over User,Mail: フェーズ3: パスワード更新・自動ログイン
	User->>FE: 新しいパスワードを入力・送信
	FE->>BE: POST /api/auth/reset-password<br/>{ newPassword } Cookie: OTPセッション
 
	BE->>DB: OTPセッションの有効性を確認(verified状態のOTPセッションがあるか)
 
	alt セッションが無効/期限切れ
		DB-->>BE: 該当なし
		BE-->>FE: 401 Unauthorized
		FE-->>User: 「セッションが無効です。再度お試しください」
		Note over User,FE: フェーズ1へ戻る
	else セッションが有効
		DB-->>BE: ユーザー情報を返す
		BE->>BE: 新パスワードをハッシュ化
		BE->>DB: UPDATE users SET password_hash = ?
		BE->>DB: UPDATE otp_sessions SET status='completed'
		DB-->>BE: 更新完了
		
		BE->>DB: 既存のログインセッションを無効化
		BE->>BE: 新しいログインセッションを生成
		BE->>DB: INSERT INTO sessions (user_id, ...)
		BE->>Mail: 「パスワード更新完了」メール送信(同期)
		Mail-->>BE: 送信結果
 
		BE-->>FE: 200 OK<br/>Set-Cookie: session_id<br/>{ user情報 }
		FE-->>User: タスク管理画面(ログイン済み)へ自動遷移
	end
```

---

## 7. パスワード変更

```mermaid
sequenceDiagram
	actor User
	participant Frontend
	participant Backend
	participant DB
	
	User->>Frontend: 現在のパスワード＆新のパスワード入力
	Frontend->>Backend: PATCH users/{user_id}/password<br>現在パスワード＆新パスワード
	opt PATCH users/{user_id}/password
		Backend->>DB: ユーザーのパスワードハッシュを取得
		DB-->>Backend: パスワードハッシュを返却
		Backend->>Backend: 現在のパスワードでパスワード認証
		alt 認証エラー
			Backend->>DB: パスワード変更失敗回数取得
			DB-->>Backend: パスワード変更失敗回数，<br>最後に失敗した日時
			alt 失敗回数が5回目
				Backend->>DB: セッションを破棄
				Backend-->>Frontend: ログイン画面に遷移せよ
				Frontend-->>User: ログイン画面に遷移
			else 5回未満
				Backend-->>Frontend: 認証失敗
				Frontend-->>User: 認証失敗，再入力要求
			end
		else 認証成功
			Backend->>DB: 新しいパスワードのハッシュに更新
			Backend->>DB: パスワード変更失敗回数をリセット
			Backend-->>Frontend: パスワード変更完了
			Frontend-->>User: パスワード変更完了
		end
	end
```
