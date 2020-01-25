package msgraph

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2"

	"golang.org/x/oauth2/clientcredentials"
)

const (
	AuthTypeClientKey = iota
)

type Client struct {
	authType    int
	oauthConfig clientcredentials.Config
	parentCtx   context.Context
	callTimeout time.Duration
	apilog      *log.Logger
}

type MsGraphError struct {
	HttpStatusCode int
	HttpStatus     string
	Message        string
}

type msGraphError struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func (e *MsGraphError) Error() string {
	if len(e.Message) > 0 {
		return e.Message
	}
	return e.HttpStatus
}

// Gets the HTTP Status code returned if there was an error
func (e *MsGraphError) StatusCode() int {
	return e.HttpStatusCode
}

func (c *Client) getHttpClient(ctx context.Context) *http.Client {
	return c.oauthConfig.Client(ctx)
}

func (c *Client) executeGetList(apiUrl string, parser func(io.Reader) string) error {
	for {
		var (
			req *http.Request
			res *http.Response
			err error
		)
		ctx, _ := context.WithTimeout(c.parentCtx, c.callTimeout)
		httpClient := c.getHttpClient(ctx)
		req, err = http.NewRequest("GET", apiUrl, nil)
		res, err = httpClient.Do(req)
		if err != nil {
			return err
		}
		if res.StatusCode == 200 {
			if c.apilog != nil {
				// If we want API logging output, read the response body, print it, and recreate the buffer
				// for the json decoder to consume
				bodyBytes, _ := ioutil.ReadAll(res.Body)
				_ = res.Body.Close() //  must close the original else we'll leak it
				c.apilog.Print("<-- ")
				c.apilog.Println(string(bodyBytes))
				res.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
			}
			apiUrl = parser(res.Body)
			_ = res.Body.Close()
			if len(apiUrl) == 0 {
				break
			}
		} else {
			// api call failure
			var mserr msGraphError
			err = json.NewDecoder(res.Body).Decode(&mserr)
			_ = res.Body.Close()
			if err == nil {
				return &MsGraphError{
					HttpStatusCode: res.StatusCode,
					HttpStatus:     res.Status,
					Message:        mserr.Error.Message,
				}
			}
			return errors.New(res.Status)
		}
	}
	return nil
}

func (c *Client) executeGetJson(apiUrl string, output interface{}) error {
	ctx, _ := context.WithTimeout(c.parentCtx, c.callTimeout)
	httpClient := c.getHttpClient(ctx)
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return err
	}
	res, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	if c.apilog != nil && res.StatusCode == 200 {
		// If we want API logging output, read the response body, print it, and recreate the buffer
		// for the json decoder to consume
		bodyBytes, _ := ioutil.ReadAll(res.Body)
		_ = res.Body.Close() //  must close the original else we'll leak it
		c.apilog.Print("<-- ")
		c.apilog.Println(string(bodyBytes))
		res.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}
	defer res.Body.Close()
	if res.StatusCode == 200 {
		err = json.NewDecoder(res.Body).Decode(output)
	} else {
		var mserr msGraphError
		if err = json.NewDecoder(res.Body).Decode(&mserr); err == nil {
			return &MsGraphError{
				HttpStatusCode: res.StatusCode,
				HttpStatus:     res.Status,
				Message:        mserr.Error.Message,
			}
		}
		return errors.New(res.Status)
	}
	return err
}

func (c *Client) executeGet(apiUrl string, parser func(io.Reader) error) error {
	ctx, _ := context.WithTimeout(c.parentCtx, c.callTimeout)
	httpClient := c.getHttpClient(ctx)
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return err
	}
	res, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	if c.apilog != nil && res.StatusCode == 200 {
		// If we want API logging output, read the response body, print it, and recreate the buffer
		// for the parser to decode
		bodyBytes, _ := ioutil.ReadAll(res.Body)
		_ = res.Body.Close() //  must close the original else we'll leak it
		c.apilog.Printf("received %d bytes", len(bodyBytes))
		x := len(bodyBytes)
		if x > 0 {
			if x > 256 {
				x = 256
			}
			c.apilog.Println(hex.Dump(bodyBytes[0:x]))
		}
		res.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}
	defer res.Body.Close()
	if res.StatusCode == 200 {
		err = parser(res.Body)
	} else {
		var mserr msGraphError
		if err = json.NewDecoder(res.Body).Decode(&mserr); err == nil {
			return &MsGraphError{
				HttpStatusCode: res.StatusCode,
				HttpStatus:     res.Status,
				Message:        mserr.Error.Message,
			}
		}
		return errors.New(res.Status)
	}
	return err
}

func NewKeyClient(ctx context.Context, TenantID string, ClientID string, ClientKey string) (*Client, error) {
	c := new(Client)
	c.authType = AuthTypeClientKey
	c.oauthConfig.ClientID = ClientID
	c.oauthConfig.ClientSecret = ClientKey
	c.oauthConfig.TokenURL = "https://login.microsoftonline.com/" + TenantID + "/oauth2/v2.0/token"
	c.oauthConfig.Scopes = append(c.oauthConfig.Scopes, "https://graph.microsoft.com/.default")
	c.oauthConfig.AuthStyle = oauth2.AuthStyleInParams
	c.parentCtx = ctx
	c.callTimeout = time.Second * 15
	return c, nil
}

func (c *Client) SetAPILogging(logger *log.Logger) {
	c.apilog = logger
}

func (c *Client) SetTimeout(timeout time.Duration) {
	if timeout > 0 {
		c.callTimeout = timeout
	}
}

func (c *Client) Close() {
}
