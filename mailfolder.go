package msgraph

import (
	"encoding/json"
	"io"
	"net/url"
)

func (c *Client) ListMailFolders(upn string, options ...ApiOption) ([]MailFolder, error) {
	var (
		err    error
		apiUrl string
	)

	apiUrl, err = formatOptions("https://graph.microsoft.com/v1.0/users/"+url.PathEscape(upn)+"/mailFolders",
		options)
	if err != nil {
		return nil, err
	}
	max, count := getMaxItemOption(options), 0
	folders := make([]MailFolder, 0, 256)
	headers := make(map[string]string)
	err2 := c.executeGetList(apiUrl, headers, func(body io.Reader) string {
		var (
			reply struct {
				Context  string       `json:"@odata.context"`
				Nextlink string       `json:"@odata.nextLink"`
				Data     []MailFolder `json:"value"`
			}
		)
		if err = json.NewDecoder(body).Decode(&reply); err == nil {
			folders = append(folders, reply.Data...)
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
	return folders, err
}

// Get a Folder object identified by folderId for a user.
// Must specify a UserPrincipalName (e.g. "bob@acme.com") or User Id (UUID).
func (c *Client) GetFolder(upn string, folderId string) (*MailFolder, error) {
	var (
		err    error
		folder MailFolder
	)
	apiUrl := "https://graph.microsoft.com/v1.0/users/" + url.PathEscape(upn) + "/mailFolders/" + folderId
	err = c.executeGet(apiUrl, func(reader io.Reader) error {
		return json.NewDecoder(reader).Decode(&folder)
	})
	if err == nil {
		return &folder, nil
	}
	return nil, err
}
