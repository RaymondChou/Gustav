package rest

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Client struct {
	conn     *httputil.ClientConn
	resource *url.URL
}

// 给指定资源建立一个client
//
// 资源跟地址
// http://127.0.0.1:3000/snips/)
func NewClient(resource string) (*Client, error) {
	var client = new(Client)
	var err error

	// 设置 host
	if client.resource, err = url.Parse(resource); err != nil {
		return nil, err
	}

	// 设置 conn
	var tcpConn net.Conn
	if tcpConn, err = net.Dial("tcp", client.resource.Host); err != nil {
		return nil, err
	}
	client.conn = httputil.NewClientConn(tcpConn, nil)

	return client, nil
}

// 关闭client
func (client *Client) Close() {
	client.conn.Close()
}

// 针对请求按照特定方法建立
func (client *Client) newRequest(method string, id string) (*http.Request, error) {
	request := new(http.Request)
	var err error

	request.ProtoMajor = 1
	request.ProtoMinor = 1
	request.TransferEncoding = []string{"chunked"}

	request.Method = method

	// Generate Resource-URI and parse it
	uri := client.resource.String() + id
	if request.URL, err = url.Parse(uri); err != nil {
		return nil, err
	}

	return request, nil
}

// 发送请求
func (client *Client) Request(request *http.Request) (*http.Response, error) {
	var err error
	var response *http.Response

	// Send the request
	if err = client.conn.Write(request); err != nil {
		return nil, err
	}

	// Read the response
	if response, err = client.conn.Read(request); err != nil {
		return nil, err
	}

	return response, nil
}

// GET /resource/
func (client *Client) Index() (*http.Response, error) {
	var request *http.Request
	var err error

	if request, err = client.newRequest("GET", ""); err != nil {
		return nil, err
	}

	return client.Request(request)
}

// GET /resource/id
func (client *Client) Find(id string) (*http.Response, error) {
	var request *http.Request
	var err error

	if request, err = client.newRequest("GET", id); err != nil {
		return nil, err
	}

	return client.Request(request)
}

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error {
	return nil
}

// POST /resource
func (client *Client) Create(body string) (*http.Response, error) {
	var request *http.Request
	var err error

	if request, err = client.newRequest("POST", ""); err != nil {
		return nil, err
	}

	request.Body = nopCloser{bytes.NewBufferString(body)}

	return client.Request(request)
}

// PUT /resource/id
func (client *Client) Update(id string, body string) (*http.Response, error) {
	var request *http.Request
	var err error
	if request, err = client.newRequest("PUT", id); err != nil {
		return nil, err
	}

	request.Body = nopCloser{bytes.NewBufferString(body)}

	return client.Request(request)
}

// 解析返回数据获得参数
func (client *Client) IdFromURL(urlString string) (string, error) {
	var uri *url.URL
	var err error
	if uri, err = url.Parse(urlString); err != nil {
		return "", err
	}

	return string(uri.Path[len(client.resource.Path):]), nil
}

// DELETE /resource/id
func (client *Client) Delete(id string) (*http.Response, error) {
	var request *http.Request
	var err error
	if request, err = client.newRequest("DELETE", id); err != nil {
		return nil, err
	}

	return client.Request(request)
}
