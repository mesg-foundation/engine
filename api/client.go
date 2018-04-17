package api

import (
	"github.com/mesg-foundation/application/api/service"
	"github.com/mesg-foundation/application/config"
	"google.golang.org/grpc"
)

// Server is the main struct that contain the server config
type Client struct {
	connection *grpc.ClientConn
	Target     string
}

// target returns the Client's target or a default
func (c *Client) target() (target string) {
	target = c.Target
	if target == "" {
		target = config.Api.Client.Target()
	}
	return
}

// conn creates a connection if needed an return it
func (c *Client) conn() (conn *grpc.ClientConn, err error) {
	if c.connection == nil {
		c.connection, err = grpc.Dial(c.target(), grpc.WithInsecure())
	}
	conn = c.connection
	return
}

// Close closes the connection (if exist)
func (c *Client) Close() {
	if c.connection != nil {
		c.connection.Close()
		c.connection = nil
	}
}

// Service returns a Service Client
func (c *Client) ServiceClient() (client service.ServiceClient, err error) {
	conn, err := c.conn()
	if err != nil {
		return
	}
	client = service.NewServiceClient(conn)
	return
}
