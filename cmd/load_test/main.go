package main

import (
	"context"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"

	"VK/internal/generated/service"
)

const (
	target          = "localhost:4122"
	subscriptionKey = "load-test-key"
	numSubscribers  = 3
	numMessages     = 1_000_000
	receiveTimeout  = 5 * time.Second
)

func main() {
	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("connection failed: %v", err)
	}
	defer conn.Close()
	client := service.NewPubSubClient(conn)

	var counters [numSubscribers]atomic.Int32
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	ready := make(chan struct{})

	for i := 0; i < numSubscribers; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			stream, err := client.Subscribe(ctx, &service.SubscribeRequest{Key: subscriptionKey})
			if err != nil {
				log.Fatalf("subscriber %d: %v", idx, err)
			}
			ready <- struct{}{}

			for {
				_, err := stream.Recv()
				if err != nil {
					return
				}
				counters[idx].Add(1)
			}
		}(i)
	}

	for i := 0; i < numSubscribers; i++ {
		<-ready
	}

	log.Println("Starting load test...")
	start := time.Now()

	var publishWg sync.WaitGroup
	publishWg.Add(numMessages)

	for i := 0; i < numMessages; i++ {
		go func(msgNum int) {
			defer publishWg.Done()
			_, err := client.Publish(ctx, &service.PublishRequest{
				Key:  subscriptionKey,
				Data: string(rune(msgNum)),
			})
			if err != nil {
				log.Printf("publish error: %v", err)
			}
		}(i)
	}

	publishWg.Wait()
	publishDuration := time.Since(start)
	log.Printf("Published %d messages in %v", numMessages, publishDuration)

	time.Sleep(receiveTimeout)
	cancel()

	wg.Wait()

	totalReceived := 0
	for idx := range counters {
		count := counters[idx].Load()
		totalReceived += int(count)
		log.Printf("Subscriber %d received %d messages", idx, count)
	}

	expected := numMessages * numSubscribers
	log.Printf("\nTotal received: %d/%d (%.2f%%)",
		totalReceived,
		expected,
		100*float64(totalReceived)/float64(expected),
	)
}
