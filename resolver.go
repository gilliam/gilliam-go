package gilliam

import "github.com/jmcvetta/napping"
import "net/url"
import "net"
import "strings"
import "errors"
import "fmt"

type serviceRegistryClient struct {
    endpoint string
    httpClient HTTPClient
}

type resolveAdapterHTTPClient struct {
    original HTTPClient
    serviceRegistryClient *serviceRegistryClient
}

type serviceInstance struct {
    Formation string
    Service string
    Instance string
    Host string
    Ports map[string]string    `json:"ports"`
}

type serviceInstanceMap map[string]serviceInstance


func (c *serviceRegistryClient) query(formation string) (res serviceInstanceMap, err error) {
    url := c.endpoint + "/" + formation
    resp, err := c.httpClient.Get(url, nil, &res, nil)
    if err != nil {
        return
    }
    if resp.Status() != 200 {
        err = errors.New("bad registry")
    }
    return
}

// Given host and port resolve that into actual host and port.
func (c *serviceRegistryClient) resolveHostPort(host, port string) (string, string, error) {
    fmt.Println(host)
    if !strings.HasSuffix(host, ".service") {
        return host, port, nil
    }

    parts := strings.Split(host, ".")
    if len(parts) == 3 {  // <service>.<formation>.service
        res, err := c.query(parts[1])
        if err != nil {
            return "", "", err
        }
        for _, inst := range res {
            if inst.Service == parts[0] {
                host = inst.Host
                port = inst.Ports[port]
                return host, port, nil
            }
        }
    } else if len(parts) == 4 { // <instance>.<service>.<formation>.service
        res, err := c.query(parts[2])
        if err != nil {
            return "", "", nil
        }
        for _, inst := range res {
            if inst.Service == parts[1] && inst.Instance == parts[0] {
                host = inst.Host
                port = inst.Ports[port]
                return host, port, nil
            }
        }
    }
    return "", "", errors.New("No such thing")
}

func (c *resolveAdapterHTTPClient) resolveURL(urlString string) (string, error) {
    url, err := url.Parse(urlString)
    if err != nil {
        return "", err
    }
    host, port, err := net.SplitHostPort(url.Host)
    if err != nil {
        host = url.Host
        port = "80"
    }
    host, port, err = c.serviceRegistryClient.resolveHostPort(host, port)
    if err != nil {
        return "", err
    }
    url.Host = net.JoinHostPort(host, port)
    return url.String(), nil
}


func NewResolvingHTTPClient(original HTTPClient, endpoint string) (HTTPClient) {
    return &resolveAdapterHTTPClient{original, &serviceRegistryClient{endpoint, original}};
}

func (c *resolveAdapterHTTPClient) Get(url string, p *napping.Params,
        result, errMsg interface{}) (resp *napping.Response, err error) {
    url, err = c.resolveURL(url)
    if err != nil {
        return
    }
    fmt.Println(url)
    resp, err = c.original.Get(url, p, result, errMsg)
    return
}
