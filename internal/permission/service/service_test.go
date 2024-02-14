package service

import (
	"context"
	"github.com/ananaslegend/news-crud/internal/permission/repository"
	mock_service "github.com/ananaslegend/news-crud/internal/permission/service/mocks"
	"go.uber.org/mock/gomock"
	"testing"
)

type GetAuthorPostByIDMock struct {
	returnValue int
}

func TestPermissionService_UserCanDeletePost(t *testing.T) {
	type args struct {
		ctx    context.Context
		userID int
		postID int
	}
	type mockBehavior func(s *mock_service.MockGetAuthorPostByID, postID int)

	tests := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		want         bool
	}{
		{
			name: "User can delete post, because he is author of post",
			mockBehavior: func(s *mock_service.MockGetAuthorPostByID, postID int) {
				s.EXPECT().GetAuthorPostByID(gomock.Any(), postID).Return(1, nil)
			},
			args: args{
				ctx:    context.Background(),
				userID: 1,
				postID: 1,
			},
			want: true,
		},
		{
			name: "User can not delete post, because he is not author of post",
			mockBehavior: func(s *mock_service.MockGetAuthorPostByID, postID int) {
				s.EXPECT().GetAuthorPostByID(gomock.Any(), postID).Return(2, nil)
			},
			args: args{
				ctx:    context.Background(),
				userID: 1,
				postID: 1,
			},
			want: false,
		},
		{
			name: "User can not delete post, post does not exist",
			mockBehavior: func(s *mock_service.MockGetAuthorPostByID, postID int) {
				s.EXPECT().GetAuthorPostByID(gomock.Any(), postID).Return(0, repository.ErrNoPostWasFound)
			},
			args: args{
				ctx:    context.Background(),
				userID: 1,
				postID: 1,
			},
			want: false,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockGetAuthorPostByID := mock_service.NewMockGetAuthorPostByID(c)
			testCase.mockBehavior(mockGetAuthorPostByID, testCase.args.userID)

			s := NewPermissionService(mockGetAuthorPostByID)

			if got := s.UserCanDeletePost(testCase.args.ctx, testCase.args.userID, testCase.args.postID); got != testCase.want {
				t.Errorf("UserCanDeletePost() = %v, want %v", got, testCase.want)
			}
		})
	}
}

func TestPermissionService_UserCanUpdatePost(t *testing.T) {
	type mockBehavior func(s *mock_service.MockGetAuthorPostByID, postID int)
	type args struct {
		ctx    context.Context
		userID int
		postID int
	}
	tests := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		want         bool
	}{
		{
			name: "User can update post, because he is author of post",
			mockBehavior: func(s *mock_service.MockGetAuthorPostByID, postID int) {
				s.EXPECT().GetAuthorPostByID(gomock.Any(), postID).Return(1, nil)
			},
			args: args{
				ctx:    context.Background(),
				userID: 1,
				postID: 1,
			},
			want: true,
		},
		{
			name: "User can`t update post, because he is not author of post",
			mockBehavior: func(s *mock_service.MockGetAuthorPostByID, postID int) {
				s.EXPECT().GetAuthorPostByID(gomock.Any(), postID).Return(2, nil)
			},
			args: args{
				ctx:    context.Background(),
				userID: 1,
				postID: 1,
			},
			want: false,
		},
		{
			name: "User can not delete post, post does not exist",
			mockBehavior: func(s *mock_service.MockGetAuthorPostByID, postID int) {
				s.EXPECT().GetAuthorPostByID(gomock.Any(), postID).Return(0, repository.ErrNoPostWasFound)
			},
			args: args{
				ctx:    context.Background(),
				userID: 1,
				postID: 1,
			},
			want: false,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockGetAuthorPostByID := mock_service.NewMockGetAuthorPostByID(c)
			testCase.mockBehavior(mockGetAuthorPostByID, testCase.args.userID)

			s := NewPermissionService(mockGetAuthorPostByID)

			if got := s.UserCanUpdatePost(testCase.args.ctx, testCase.args.userID, testCase.args.postID); got != testCase.want {
				t.Errorf("UserCanUpdatePost() = %v, want %v", got, testCase.want)
			}
		})
	}
}
