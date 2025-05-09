package rpc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"VK/internal/errs"
	generated "VK/internal/generated/service"
	grpc "VK/internal/rpc/mocks"
)

func TestHandler_Publish(t *testing.T) {
	type fields struct {
		subPub        SubPubService
		streamManager StreamManager
	}
	tt := []struct {
		name       string
		setupMocks func(ctrl *gomock.Controller) fields
		expErr     error
		key        string
		data       string
	}{
		{
			name:   "successful publish",
			expErr: nil,
			key:    "123",
			data:   "kkkkkk",
			setupMocks: func(ctrl *gomock.Controller) fields {
				subpub := grpc.NewMockSubPubService(ctrl)
				subpub.EXPECT().Publish("123", "kkkkkk").Return(nil)
				return fields{
					subpub,
					nil,
				}
			},
		},
		{
			name:   "with error",
			expErr: status.Errorf(codes.NotFound, errs.ErrSubjectNotFound.Error()),
			key:    "123",
			data:   "kkkkkk",
			setupMocks: func(ctrl *gomock.Controller) fields {
				subpub := grpc.NewMockSubPubService(ctrl)
				subpub.EXPECT().Publish("123", "kkkkkk").Return(errs.ErrSubjectNotFound)
				return fields{
					subpub,
					nil,
				}
			},
		},
	}
	t.Parallel()
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := tc.setupMocks(ctrl)
			handler := New(f.subPub, f.streamManager, func(msg interface{}) {})

			_, err := handler.Publish(context.Background(), &generated.PublishRequest{
				Key:  tc.key,
				Data: tc.data,
			})
			require.ErrorIs(t, tc.expErr, err)
		})
	}

}
