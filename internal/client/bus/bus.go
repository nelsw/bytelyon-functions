package bus

import (
	"bytelyon-functions/internal/app"
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchevents"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchevents/types"
)

type Client interface {
	PutOne(string, string, string) error
	PutMany(string, string, []string) error
}

type client struct {
	*cloudwatchevents.Client
	ctx context.Context
	bus string
}

func (c *client) PutOne(s, t, d string) (err error) {
	return c.PutMany(s, t, []string{d})
}

func (c *client) PutMany(s, t string, dd []string) (err error) {
	var ee []types.PutEventsRequestEntry
	for _, d := range dd {
		ee = append(ee, types.PutEventsRequestEntry{
			EventBusName: &c.bus,
			Source:       &s,
			DetailType:   &t,
			Detail:       &d,
		})
	}
	return c.put(c.ctx, ee)
}

func (c *client) put(ctx context.Context, e []types.PutEventsRequestEntry) (err error) {
	_, err = c.PutEvents(ctx, &cloudwatchevents.PutEventsInput{Entries: e})
	return
}

// New returns a new S3 client with the provided context.
func New(ctx context.Context) Client {
	cfg, _ := config.LoadDefaultConfig(ctx)
	return &client{
		cloudwatchevents.NewFromConfig(cfg),
		ctx,
		"bytelyon-bus-" + app.Mode(),
	}
}
