package imap

import "github.com/emersion/go-imap/client"
import "github.com/emersion/go-imap-id"

type Client struct {
	*client.Client
	IDCli *id.Client
}

func WrapEXClient(c *client.Client) *Client {
	return &Client{
		Client: c,
		IDCli:  id.NewClient(c),
	}
}
