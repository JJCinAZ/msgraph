package msgraph

import (
	"encoding/json"
	"io"
	"net/url"
)

func (c *Client) ListMessageRules(upn string, options ...ApiOption) ([]MessageRule, error) {
	var (
		err    error
		apiUrl string
	)

	apiUrl, err = formatOptions("https://graph.microsoft.com/v1.0/users/"+url.PathEscape(upn)+"/mailFolders/inbox/messagerules",
		options)
	if err != nil {
		return nil, err
	}
	max, count := getMaxItemOption(options), 0
	rules := make([]MessageRule, 0, 256)
	headers := make(map[string]string)
	err2 := c.executeGetList(apiUrl, headers, func(body io.Reader) string {
		var (
			reply struct {
				Context  string        `json:"@odata.context"`
				Nextlink string        `json:"@odata.nextLink"`
				Data     []MessageRule `json:"value"`
			}
		)
		if err = json.NewDecoder(body).Decode(&reply); err == nil {
			rules = append(rules, reply.Data...)
			count += len(reply.Data)
			if count >= max {
				return ""
			}
			return reply.Nextlink
		}
		return ""
	})
	if err == nil {
		err = err2
	}
	return rules, err
}

func (c *Client) GetMessageRule(upn string, ruleId string) (*MessageRule, error) {
	var (
		err  error
		rule MessageRule
	)
	apiUrl := "https://graph.microsoft.com/v1.0/users/" + url.PathEscape(upn) + "/mailFolders/inbox/messagerules/" + ruleId
	err = c.executeGet(apiUrl, func(reader io.Reader) error {
		return json.NewDecoder(reader).Decode(&rule)
	})
	if err == nil {
		return &rule, nil
	}
	return nil, err
}
