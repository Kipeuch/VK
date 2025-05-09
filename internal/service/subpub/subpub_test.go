package subpub

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"VK/internal/entity"
	"VK/internal/errs"
)

func TestSubPubBus_BasicPublishSubscribe(t *testing.T) {
	type publish struct {
		subject string
		msg     interface{}
	}

	tests := []struct {
		name       string
		subject    string
		publishes  []publish
		expectErr  error
		expectMsgs []entity.Message
	}{
		{
			name:       "no_subscribers_returns_subject_not_found",
			subject:    "topic1",
			publishes:  []publish{{subject: "topic1", msg: "msg1"}},
			expectErr:  errs.ErrSubjectNotFound,
			expectMsgs: nil,
		},
		{
			name:       "single_subscriber_receives_messages_in_order",
			subject:    "topic2",
			publishes:  []publish{{subject: "topic2", msg: "first"}, {subject: "topic2", msg: "second"}},
			expectErr:  nil,
			expectMsgs: []entity.Message{entity.Message{Subject: "topic2", Data: "first"}, entity.Message{Subject: "topic2", Data: "second"}},
		},
	}

	t.Parallel()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			bus := NewSubPub()

			var ch chan interface{}
			if tc.expectMsgs != nil {
				ch = make(chan interface{}, len(tc.expectMsgs))
				_, err := bus.Subscribe(tc.subject, func(msg interface{}) {
					ch <- msg
				})
				require.NoError(t, err)
			}

			for _, pub := range tc.publishes {
				err := bus.Publish(pub.subject, pub.msg)
				if tc.expectErr != nil {
					require.ErrorIs(t, err, tc.expectErr)
					return
				}
				require.NoError(t, err)
			}

			if tc.expectMsgs != nil {
				received := []entity.Message{}
				for i := 0; i < len(tc.expectMsgs); i++ {
					select {
					case m := <-ch:
						received = append(received, m.(entity.Message))
					case <-time.After(100 * time.Millisecond):
						t.Fatalf("timeout waiting for message %d", i)
					}
				}
				for i := range tc.expectMsgs {
					tc.expectMsgs[i].IdempontencyKey = received[i].IdempontencyKey
				}

				require.Equal(t, tc.expectMsgs, received)
			}
		})
	}
}

func TestSubPubBus_MultipleSubscribersIndependent(t *testing.T) {
	bus := NewSubPub()
	subject := "topic_multi"

	ch1 := make(chan interface{}, 1)
	ch2 := make(chan interface{}, 1)

	_, err := bus.Subscribe(subject, func(msg interface{}) {
		time.Sleep(50 * time.Millisecond)
		ch1 <- msg
	})
	require.NoError(t, err)

	_, err = bus.Subscribe(subject, func(msg interface{}) {
		ch2 <- msg
	})
	require.NoError(t, err)

	err = bus.Publish(subject, "hello")
	require.NoError(t, err)

	select {
	case m := <-ch2:
		require.Equal(t, entity.Message{IdempontencyKey: "con2", Data: "hello"}, m)
	case <-time.After(20 * time.Millisecond):
		t.Fatal("fast subscriber did not receive message in time")
	}

	select {
	case m := <-ch1:
		require.Equal(t, entity.Message{IdempontencyKey: "con1", Data: "hello"}, m)
	case <-time.After(200 * time.Millisecond):
		t.Fatal("slow subscriber did not receive message in time")
	}
}

func TestSubPubBus_CloseBehavior(t *testing.T) {
	bus := NewSubPub()
	subject := "topic_close"

	ch := make(chan interface{}, 2)
	_, err := bus.Subscribe(subject, func(msg interface{}) {
		ch <- msg
	})
	require.NoError(t, err)

	require.NoError(t, bus.Publish(subject, "before"))

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	err = bus.Close(ctx)
	require.NoError(t, err)

	err = bus.Publish(subject, "after")
	require.ErrorIs(t, err, errs.ErrSubPubClosed)

	select {
	case m := <-ch:
		require.Equal(t, entity.Message{IdempontencyKey: "con1", Data: "before"}, m)
	case <-time.After(50 * time.Millisecond):
		t.Fatal("pending message not delivered after close")
	}
}
