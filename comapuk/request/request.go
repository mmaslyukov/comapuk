package request

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

type Values = url.Values

type Client struct {
	_client *http.Client
}
type Node = html.Node
type Parser struct {
	_body *Node
}

func NewClient(cert ...string) (c *Client, err error) {
	var tlsConfig *tls.Config = nil
	if len(cert) == 1 {
		caCert, err := ioutil.ReadFile(cert[0])
		if err != nil {
			return c, err
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		tlsConfig = &tls.Config{
			RootCAs: caCertPool,
		}
	}
	// jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	jar, err := cookiejar.New(nil)
	if err != nil {
		return c, err
	}
	client := &http.Client{
		Jar: jar,
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}
	c = &Client{_client: client}
	return c, err
}

func (c *Client) Get(url string) (p *Parser, e error) {

	resp, err := c._client.Get(url)
	if err != nil {
		return p, e
	}
	defer resp.Body.Close()
	
	return c.handleResponse(resp)
}

func (c *Client) Post(url string, data Values) (p *Parser, err error) {
	// log.Printf("%#v\n", data)
	resp, err := c._client.PostForm(url, data)
	if err != nil {
		return p, err
	}
	defer resp.Body.Close()
	
	// bytes, err := httputil.DumpResponse(resp, true)
	// log.Printf("Response %v\n",string(bytes))
	
	// return nil,nil
	return c.handleResponse(resp)
}

func (c *Client) handleResponse(resp *http.Response) (p *Parser, err error) {

	if resp.StatusCode != 200 {
		return p, errors.New(fmt.Sprintf("Status code isn't 200\n"))
	}

	node, err := htmlquery.Parse(resp.Body)
	if err != nil {
		return p, err
	}

	return &Parser{_body: node}, err
}

func (p *Parser) Attr(xpath string, name string) (val string, err error) {
	node, err := htmlquery.Query(p._body, xpath)
	if err != nil {
		return "", err
	}
	for _, a := range node.Attr {
		if a.Key == name {
			val = a.Val
		}
		// log.Printf("attr: %#v\n", a)
	}
	// p._body.FirstChild
	return val, err
}

func (p *Parser) Data(xpath string) (val string, err error) {
	node, err := htmlquery.Query(p._body, xpath)
	if err != nil {
		return "", err
	}
	if node == nil {
		return "", errors.New(fmt.Sprintf("Node is nil"))
	}
	if node.FirstChild == nil {
		return "", errors.New(fmt.Sprintf("Node.FirstChild is nil"))
	}
	if len(node.FirstChild.Data) == 0  {
		return "", err
	}

	// log.Printf("%v\n", val)
	return node.FirstChild.Data, err
}

