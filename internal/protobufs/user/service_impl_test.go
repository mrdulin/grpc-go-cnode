package user_test

import (
	"context"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mrdulin/grpc-go-cnode/internal/protobufs/user"
	"github.com/mrdulin/grpc-go-cnode/internal/utils/auth"
	"github.com/mrdulin/grpc-go-cnode/test/mocks"
	"github.com/mrdulin/grpc-go-cnode/test/testdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/metadata"
)

func TestUserServiceImpl_ValidateAccessToken(t *testing.T) {
	assert := assert.New(t)
	t.Run("should validate accesstoken correctly", func(t *testing.T) {
		var res user.ValidateAccessTokenResponse
		testHttp := new(mocks.MockedHttp)
		userServiceImpl := user.NewUserServiceImpl(testHttp, testdata.BaseURL)
		req := user.ValidateAccessTokenRequest{Accesstoken: testdata.Accesstoken}

		testHttp.On("Post", testdata.BaseURL+"/accesstoken", &user.ValidateAccessTokenRequestPayload{AccessToken: req.Accesstoken}, &res).Return(nil).Run(func(args mock.Arguments) {
			arg := args.Get(2).(*user.ValidateAccessTokenResponse)
			arg.Id = "123"
			arg.Success = true
			arg.Loginname = testdata.Loginname
			arg.AvatarUrl = testdata.AvatarURL
		})

		md := metadata.New(map[string]string{
			"authorization": auth.JWT,
		})
		ctx := metadata.NewIncomingContext(context.Background(), md)
		got, err := userServiceImpl.ValidateAccessToken(ctx, &req)
		testHttp.AssertExpectations(t)
		assert.Nil(err)
		assert.Equal(got.GetSuccess(), true)
		assert.Equal(got.GetLoginname(), testdata.Loginname)
		assert.Equal(got.GetAvatarUrl(), testdata.AvatarURL)
		assert.Equal(got.GetId(), "123")
	})

	t.Run("should return error if accesstoken is not passed in argument", func(t *testing.T) {
		testHttp := new(mocks.MockedHttp)
		userServiceImpl := user.NewUserServiceImpl(testHttp, testdata.BaseURL)
		req := user.ValidateAccessTokenRequest{Accesstoken: ""}
		md := metadata.New(map[string]string{
			"authorization": auth.JWT,
		})
		ctx := metadata.NewIncomingContext(context.Background(), md)
		res, err := userServiceImpl.ValidateAccessToken(ctx, &req)
		testHttp.AssertNotCalled(t, "Post")
		assert.Nil(res)
		assert.Error(err)
		s, ok := status.FromError(err)
		if !ok {
			t.Fatal(ok)
		}
		assert.Equal(s.Code(), codes.InvalidArgument)

	})
}
