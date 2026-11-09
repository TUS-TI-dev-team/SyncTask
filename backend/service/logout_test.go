package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"synctask/backend/model"
)

type mockLogoutRepository struct {
	attemptFunc func(context.Context, *model.LogoutAttempt) (*model.LogoutAttemptResult, error)
}

func (m *mockLogoutRepository) AttemptLogout(ctx context.Context, attempt *model.LogoutAttempt) (*model.LogoutAttemptResult, error) {
	return m.attemptFunc(ctx, attempt)
}

func TestLogoutService(t *testing.T) {

	t.Run("正常系: セッションID が正常なものでありrepository 内でもエラーがない際にnilを返すこと", func(t *testing.T) {
		input := model.LogoutInput{SessionID: "4013350b-9e2b-4ef8-9f98-b07b7341418f"}
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

	t.Run("準正常系: セッションID が空文字の際にnil を返しrepository にいかないこと", func(t *testing.T) {
		input := model.LogoutInput{}
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
		assert.False(t, called)
	})

	t.Run("準正常系: セッションID がnil の際にnil を返しrepository にいかないこと", func(t *testing.T) {
		called := false
		repo := &mockLogoutRepository{attemptFunc: func(_ context.Context, attempt *model.LogoutAttempt) (*model.LogoutAttemptResult, error) {
			called = true
			return &model.LogoutAttemptResult{
				Status: model.LogoutStatusSuccess,
			}, nil
		}}

		svc := NewLogoutService(repo)
		err := svc.Logout(context.Background(), nil)
		require.NoError(t, err)
		assert.False(t, called)
	})

	t.Run("準正常系: セッションID がUUIDv7 の際にnil を返しrepository にいかないこと", func(t *testing.T) {
		uuidv7, err := uuid.NewV7()
		require.NoError(t, err)
		id, err := uuid.Parse(uuidv7.String())
		require.NoError(t, err)
		assert.Equal(t, id.Version(), uuid.Version(7))

		input := model.LogoutInput{SessionID: id.String()}
		called := false
		repo := &mockLogoutRepository{attemptFunc: func(_ context.Context, attempt *model.LogoutAttempt) (*model.LogoutAttemptResult, error) {
			called = true
			assert.Equal(t, input.SessionID, attempt.SessionID)
			return &model.LogoutAttemptResult{
				Status: model.LogoutStatusSuccess,
			}, nil
		}}

		svc := NewLogoutService(repo)
		err = svc.Logout(context.Background(), &input)
		require.NoError(t, err)
		assert.False(t, called)
	})

	t.Run("異常系: セッションID が正常かつrepository にてdb のエラーが帰ってきた際にerrをそのまま伝播させること", func(t *testing.T) {
		dbError := errors.New("database failed")
		called := false
		repo := &mockLogoutRepository{attemptFunc: func(context.Context, *model.LogoutAttempt) (*model.LogoutAttemptResult, error) {
			called = true
			return nil, dbError
		}}

		svc := NewLogoutService(repo)
		err := svc.Logout(context.Background(), &model.LogoutInput{SessionID: uuid.New().String()})
		require.ErrorIs(t, err, dbError)
		assert.True(t, called)
	})

	t.Run("準正常系: セッションID が正常かつrepository にてsid が見つからなかった際にnil を返すこと", func(t *testing.T) {
		input := model.LogoutInput{SessionID: "4013350b-9e2b-4ef8-9f98-b07b7341418f"}
		called := false
		repo := &mockLogoutRepository{attemptFunc: func(context.Context, *model.LogoutAttempt) (*model.LogoutAttemptResult, error) {
			called = true
			return &model.LogoutAttemptResult{
				Status: model.LogoutStatusNotFound,
			}, nil
		}}

		svc := NewLogoutService(repo)
		err := svc.Logout(context.Background(), &input)
		require.NoError(t, err)
		assert.True(t, called)
	})

	t.Run("異常系: セッションID が 正常かつrepository からのステータスが想定外の値であった際にerr で返しかつその err が(*AppErr)でキャストできないこと", func(t *testing.T) {
		input := model.LogoutInput{SessionID: "4013350b-9e2b-4ef8-9f98-b07b7341418f"}
		called := false
		repo := &mockLogoutRepository{attemptFunc: func(context.Context, *model.LogoutAttempt) (*model.LogoutAttemptResult, error) {
			called = true
			return &model.LogoutAttemptResult{
				Status: model.LogoutStatus("unknown"),
			}, nil
		}}

		svc := NewLogoutService(repo)
		err := svc.Logout(context.Background(), &input)
		require.ErrorContains(t, err, "unknown logout status")
		_, ok := err.(*model.AppError)
		assert.False(t, ok)
		assert.True(t, called)
	})

	t.Run("異常系: セッションID が 正常かつrepository からの戻り値が(nil, nil)際にerrで返答しかつerr が(*AppErr)でキャストできないこと", func(t *testing.T) {
		input := model.LogoutInput{SessionID: "4013350b-9e2b-4ef8-9f98-b07b7341418f"}
		called := false
		repo := &mockLogoutRepository{attemptFunc: func(context.Context, *model.LogoutAttempt) (*model.LogoutAttemptResult, error) {
			called = true
			return nil, nil
		}}

		svc := NewLogoutService(repo)
		err := svc.Logout(context.Background(), &input)
		require.ErrorContains(t, err, "logout result is nil")
		_, ok := err.(*model.AppError)
		assert.False(t, ok)
		assert.True(t, called)
	})
}
