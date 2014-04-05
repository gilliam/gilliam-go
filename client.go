package gilliam

import "github.com/jmcvetta/napping"
//import "net/http"
import "os"
import "errors"
import "encoding/json"


type HTTPClient interface {
    Get(url string, p *napping.Params, result, errMsg interface{}) (*napping.Response, error)
}

func (c *Client) queryCollection(url string, res interface{}) (string, error) {
    collection := struct {
        Items json.RawMessage
        Links struct {
            Next string
            Prev string
        }
    }{}
    resp, err := c.httpClient.Get(url, nil, &collection, nil)
    if err != nil {
        return "", err
    }
    if resp.Status() != 200 {
        return "", errors.New("Not 200 OK")
    }
    err = json.Unmarshal(collection.Items, res)
    if err != nil {
        return "", err
    }
    return collection.Links.Next, nil
}

type Client struct {
    httpClient HTTPClient
}

func New() *Client {
    httpClient := &napping.Session{}
    endpoint := os.Getenv("GILLIAM_SERVICE_REGISTRY")
    return &Client{NewResolvingHTTPClient(httpClient, endpoint)}
}
