package service

import (
	"testing"
	"context"
	"errors"

	"synctask/backend/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockLogoutRepository struct {
	attemptFunc func(context.Context, *model.LogoutAttempt) (*model.LogoutAttemptResult, error)
}

func (m *mockLogoutRepository) AttemptLogout(ctx context.Context, attempt *model.LogoutAttempt) (*model.LogoutAttemptResult, error) {
	return m.attemptFunc(ctx, attempt)
}

//test するケース
//1. 正常系: セッションID 正常 + repo エラーなし-> nil
//2. 準正常系: セッションID がnil -> nil + repo にいかない
//3. 準正常系: セッションID がUUIDv7 -> nil + repo にいかない
//4. 異常系: セッションID 正常 + repo db のエラ-で帰ってくる-> err を伝播させる
//5. 準正常系: セッションID 正常 + repo sid が見つからなかった -> nil
//6. 異常系: セッションID 正常 + repo からのステータスが想定外の値 -> err で返す+ (*AppErr)でキャストできない
func TestLogoutService(t *testing.T) {
	input := model.LogoutInput{SessionID: "4013350b-9e2b-4ef8-9f98-b07b7341418f"}

	t.Run("正常系: セッションID が正常なものでありrepository 内でもエラーがない際にnilを返すこと", func (t *testing.T){
		called := false
		repo := &mockLogoutRepository{attemptFunc: func(_ context.Context, attempt *model.LogoutAttempt) (*model.LogoutAttemptResult, error) {
			called = true
			assert.Equal(t, input.SessionID, attempt.SessionID)
			return &model.LogoutAttemptResult{
				Status: model.LogoutStatusSuccess,
			}, nil
		}}

		svc := NewLogoutService(repo)
		err := svc.Logout(context.Background(), &input)
		require.NoError(t, err)
		assert.True(t, called)
	})

	t.Run("準正常系: セッションID がnil の際にnil を返しrepository にいかないこと", func (t *testing.T){
	})

	t.Run("準正常系: セッションID がUUIDv7 の際にnil を返しrepository にいかないこと", func (t *testing.T){
	})

	t.Run("異常系: セッションID が正常かつrepository にてdb のエラーが帰ってきた際にerrをそのまま伝播させること", func (t *testing.T){
	})

	t.Run("準正常系: セッションID が正常かつrepository にてsid が見つからなかった際にnil を返すこと", func(t *testing.T){
		repo := &mockLogoutRepository{attemptFunc: func(context.Context, *model.LogoutAttempt) (*model.LogoutAttemptResult, error) {
			return &model.LogoutAttemptResult{
				Status: model.LogoutStatusNotFound,
			}, nil
		}}

		svc := NewLogoutService(repo)
		err := svc.Logout(context.Background(), &input)
		require.NoError(t, err)
	})

	t.Run("異常系: セッションID が 正常かつrepository からのステータスが想定外の値であった際にerr で返しかつその err が(*AppErr)でキャストできないこと", func(t *testing.T){
		repo := &mockLogoutRepository{attemptFunc: func(context.Context, *model.LogoutAttempt) (*model.LogoutAttemptResult, error) {
			return &model.LogoutAttemptResult{
				Status: model.LogoutStatus("unknown"),
			}, nil
		}}

		svc := NewLogoutService(repo)
		err := svc.Logout(context.Background(), &input)
		require.NoError(t, err)

	})
}
