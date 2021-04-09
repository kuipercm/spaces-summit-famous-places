package pubsub

import (
	"context"
	"fmt"
	"google.golang.org/api/iterator"
	"os"

	"cloud.google.com/go/pubsub"
)

type TopicDetails struct {
	ProjectId string
	TopicId   string
}

func (t TopicDetails) CreateTopic() error {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, t.ProjectId)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}

	topics := client.Topics(ctx)
	for {
		topic, err := topics.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("client.Topics(%q).Next: %v", t.ProjectId, err)
		}
		if topic.String() == t.TopicId {
			fmt.Fprintf(os.Stdout, "Topic %v already exists\n", t.TopicId)
			return nil
		}
	}

	_, err = client.CreateTopic(ctx, t.TopicId)
	if err != nil {
		return fmt.Errorf("CreateTopic: %v", err)
	}
	fmt.Fprintf(os.Stdout, "Topic created: %v\n", t)

	return nil
}
