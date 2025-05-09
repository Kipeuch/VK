package stream_manager

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	generated "VK/internal/generated/service"
)

type mockPubServer struct {
	grpc.ServerStream
}

func (m *mockPubServer) Send(e *generated.Event) error {
	return nil
}

func TestStreamManager_AddConn(t *testing.T) {
	mockConn := &mockPubServer{}
	tt := []struct {
		name     string
		subject  string
		conn     generated.PubSub_SubscribeServer
		beginMap map[string][]generated.PubSub_SubscribeServer
		waitMap  map[string][]generated.PubSub_SubscribeServer
	}{
		{
			name:     "Add new conn, success",
			subject:  "cinema",
			conn:     mockConn,
			beginMap: map[string][]generated.PubSub_SubscribeServer{},
			waitMap:  map[string][]generated.PubSub_SubscribeServer{"cinema": {mockConn}},
		},
		{
			name:     "Add existing conn",
			conn:     mockConn,
			subject:  "cinema",
			beginMap: map[string][]generated.PubSub_SubscribeServer{"cinema": {mockConn}},
			waitMap:  map[string][]generated.PubSub_SubscribeServer{"cinema": {mockConn}},
		},
	}
	t.Parallel()
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			m := New()

			m.conns = tc.beginMap

			m.AddConn(tc.subject, tc.conn)

			require.EqualValues(t, m.conns, tc.waitMap)

		})
	}

}

func TestStreamManager_RemoveConn(t *testing.T) {
	mockConn := &mockPubServer{}
	tt := []struct {
		name     string
		subject  string
		conn     generated.PubSub_SubscribeServer
		beginMap map[string][]generated.PubSub_SubscribeServer
		waitMap  map[string][]generated.PubSub_SubscribeServer
	}{
		{
			name:     "Delete exists conn, success",
			subject:  "cinema",
			conn:     mockConn,
			beginMap: map[string][]generated.PubSub_SubscribeServer{"cinema": {mockConn}},
			waitMap:  map[string][]generated.PubSub_SubscribeServer{"cinema": {}},
		},
		{
			name:     "Delete no exist con",
			subject:  "cinema",
			conn:     mockConn,
			beginMap: map[string][]generated.PubSub_SubscribeServer{"1": {mockConn}},
			waitMap:  map[string][]generated.PubSub_SubscribeServer{"1": {mockConn}},
		},
	}
	t.Parallel()
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			m := New()

			m.conns = tc.beginMap

			m.RemoveConn(tc.subject, tc.conn)

			require.EqualValues(t, m.conns, tc.waitMap)

		})
	}

}
