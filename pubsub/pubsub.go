package pubsub

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
)

func CreateTopic(ctx context.Context, projectID, topicID string) error {
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("%w: failed to setup pubsub client", err)
	}

	topic := client.TopicInProject(topicID, projectID)
	exists, err := topic.Exists(ctx)
	if err != nil {
		return fmt.Errorf("%w: failed to check if topic exists", err)
	}
	if exists {
		fmt.Printf("Topic %s already exists\n", topicID)
		return nil
	}

	_, err = client.CreateTopic(ctx, topicID)
	if err != nil {
		return fmt.Errorf("%w: failed to create topic", err)
	}

	fmt.Printf("Topic created: %s\n", topicID)
	return nil
}
