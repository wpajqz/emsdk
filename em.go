package emsdk

import "sync"

var (
	client *Client
	once   sync.Once
)

type Client struct {
	clientID     string
	clientSecret string
	baseURL      string
	adminToken   adminTokenResponse
}

func New(orgName, appName, clientID, clientSecret string) (*Client, error) {
	var err error
	once.Do(func() {
		client = &Client{
			baseURL:      "https://a1.easemob.com/" + orgName + "/" + appName,
			clientID:     clientID,
			clientSecret: clientSecret,
		}

		client.adminToken, err = client.getAccessToken()
	})

	return client, err
}
