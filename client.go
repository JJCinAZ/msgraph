package msgraph

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"golang.org/x/oauth2/microsoft"
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
	AuthTypeAuthCode
)

type Client struct {
	authType    int
	OauthConfig oauth2.Config
	token       *oauth2.Token
	ccConfig    clientcredentials.Config
	parentCtx   context.Context
	callTimeout time.Duration
	apilog      *log.Logger
}

type TokenCache interface {
	Load(c *Client) (*oauth2.Token, error)
	Save(c *Client, token *oauth2.Token) error
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
	if c.authType == AuthTypeClientKey {
		return c.ccConfig.Client(ctx)
	}
	return c.OauthConfig.Client(ctx, c.token)
}

// GetList is specialized in that we get back paged results from MSGraph API
// We need to detect this and keep calling back for the next page.
// This function wraps all that logic with the supplied parser controlling
// and early stop or continuation by returning the NextLink parameter
func (c *Client) executeGetList(apiUrl string, headers map[string]string, parser func(io.Reader) string) error {
	for {
		var (
			req *http.Request
			res *http.Response
			err error
		)
		ctx, _ := context.WithTimeout(c.parentCtx, c.callTimeout)
		httpClient := c.getHttpClient(ctx)
		req, err = http.NewRequest("GET", apiUrl, nil)
		if headers != nil {
			for k, v := range headers {
				req.Header.Add(k, v)
			}
		}
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

func (c *Client) executePost(apiUrl string, body interface{}, parser func(io.Reader) error) error {
	ctx, _ := context.WithTimeout(c.parentCtx, c.callTimeout)
	httpClient := c.getHttpClient(ctx)
	pr, pw := io.Pipe()
	go func() {
		_ = json.NewEncoder(pw).Encode(body)
		pw.Close()
	}()
	req, err := http.NewRequest("POST", apiUrl, pr)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	if res, err := httpClient.Do(req); err != nil {
		return err
	} else {
		return c.executeProcessResult(res, parser)
	}
}

func (c *Client) executeGetJson(apiUrl string, output interface{}) error {
	return c.executeMethod("GET", apiUrl,
		func(reader io.Reader) error {
			return json.NewDecoder(reader).Decode(output)
		})
}

func (c *Client) executeGet(apiUrl string, parser func(io.Reader) error) error {
	return c.executeMethod("GET", apiUrl, parser)
}

func (c *Client) executeDelete(apiUrl string) error {
	return c.executeMethod("DELETE", apiUrl, func(reader io.Reader) error {
		return nil
	})
}

func (c *Client) executeMethod(method string, apiUrl string, parser func(io.Reader) error) error {
	ctx, _ := context.WithTimeout(c.parentCtx, c.callTimeout)
	httpClient := c.getHttpClient(ctx)
	req, err := http.NewRequest(method, apiUrl, nil)
	if err != nil {
		return err
	}
	if res, err := httpClient.Do(req); err != nil {
		return err
	} else {
		return c.executeProcessResult(res, parser)
	}
}

func (c *Client) executeProcessResult(res *http.Response, parser func(io.Reader) error) error {
	var err error
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
		if parser != nil {
			err = parser(res.Body)
		}
		return err
	} else if res.StatusCode >= 201 && res.StatusCode <= 299 {
		return nil
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

// Creates a new MSGraph API Client using the Client Credentials Grant flow [see https://oauth.net/2/grant-types/client-credentials/]
// One needs to create the Client ID and Key in Azure AD prior to calling this.
// All permissions must be added to the Configured Permissions in the App Registration within Azure AD console
// [see https://docs.microsoft.com/en-us/azure/active-directory/develop/v2-permissions-and-consent#the-default-scope]
func NewKeyClient(ctx context.Context, TenantID string, ClientID string, ClientKey string) (*Client, error) {
	c := new(Client)
	c.authType = AuthTypeClientKey
	c.ccConfig.ClientID = ClientID
	c.ccConfig.ClientSecret = ClientKey
	c.ccConfig.TokenURL = "https://login.microsoftonline.com/" + TenantID + "/oauth2/v2.0/token"
	c.ccConfig.Scopes = append(c.ccConfig.Scopes, "https://graph.microsoft.com/.default")
	c.ccConfig.AuthStyle = oauth2.AuthStyleInParams
	c.parentCtx = ctx
	c.callTimeout = time.Second * 180
	return c, nil
}

// Create a new MSGraph API Client using the Authorization Code Grant flow [see https://oauth.net/2/grant-types/authorization-code/]
// One needs to create the Client ID and Secret in Azure AD prior to calling this.
// One or more scopes must be supplied to indicate the type of access being requested.
// This authentication flow launches a web browser for the OAuth2 work with a callback to
// a web server at localhost:8001.
func NewClient(ctx context.Context, TenantID string, ClientID string, ClientSecret string, scopes []string,
	cache TokenCache, timeout time.Duration) (*Client, error) {
	var (
		state string
		err   error
	)
	c := new(Client)
	c.authType = AuthTypeAuthCode
	c.OauthConfig.ClientID = ClientID
	c.OauthConfig.ClientSecret = ClientSecret
	c.OauthConfig.Endpoint = microsoft.AzureADEndpoint(TenantID)
	c.OauthConfig.Scopes = append(c.ccConfig.Scopes, "offline_access")
	c.OauthConfig.Scopes = append(c.ccConfig.Scopes, scopes...)
	c.parentCtx = ctx
	c.callTimeout = time.Second * 180

	// Load any token which was previously cached
	if cache == nil {
		cache = NullCache{}
	}
	c.token, err = cache.Load(c)
	if err == nil {
		return c, nil
	}
	state, err = generateNonce(16)
	if err == nil {
		err = openUrl(c.OauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline))
	}
	if err != nil {
		return nil, err
	}
	cbChan := c.startCallbackServer(ctx, state, timeout)
	cb := <-cbChan
	if cb.err != nil {
		return nil, cb.err
	}
	c.token, err = c.OauthConfig.Exchange(ctx, cb.code)
	if err != nil {
		return nil, err
	}
	err = cache.Save(c, c.token)
	return c, nil
}
