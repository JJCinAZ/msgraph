package msgraph

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/url"
	"os"
	"strings"
	"time"
)

type Attachment struct {
	OdataType             string    `json:"@odata.type"`
	OdataMediaContentType string    `json:"@odata.mediaContentType"`
	ID                    string    `json:"id"`
	LastModifiedDateTime  time.Time `json:"lastModifiedDateTime"`
	Name                  string    `json:"name"`
	ContentType           string    `json:"contentType"`
	Size                  int       `json:"size"`
	IsInline              bool      `json:"isInline"`
	ContentID             string    `json:"contentId"`
	ContentBytes          string    `json:"contentBytes"`
}

// Writes content of Attachment to dst. It returns the number of bytes
// written and the first error encountered while copying, if any.
func (a Attachment) Write(dst io.Writer) (int64, error) {
	decoder := base64.NewDecoder(base64.StdEncoding, strings.NewReader(a.ContentBytes))
	return io.Copy(dst, decoder)
}

func (a Attachment) AsBytesBuffer() ([]byte, error) {
	return base64.StdEncoding.DecodeString(a.ContentBytes)
}

func (a Attachment) String() string {
	if b, err := base64.StdEncoding.DecodeString(a.ContentBytes); err == nil {
		return string(b)
	}
	return ""
}

func (a Attachment) SHA256Hash() ([]byte, error) {
	decoder := base64.NewDecoder(base64.StdEncoding, strings.NewReader(a.ContentBytes))
	h := sha256.New()
	if _, err := io.Copy(h, decoder); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

func (a Attachment) WriteFile(filename string, perm os.FileMode) error {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, perm)
	if err != nil {
		return err
	}
	if _, err = a.Write(f); err != nil {
		_ = f.Close()
		os.Remove(filename)
		return err
	}
	return f.Close()
}

func (c *Client) ListAttachments(upn string, msgId string, options ...ApiOption) ([]Attachment, error) {
	var (
		err error
	)

	apiUrl, err := formatOptions("https://graph.microsoft.com/v1.0/users/"+url.PathEscape(upn)+"/messages/"+url.PathEscape(msgId)+"/attachments",
		options)
	if err != nil {
		return nil, err
	}
	max, count := getMaxItemOption(options), 0
	list := make([]Attachment, 0, 128)
	c.executeGetList(apiUrl, nil, func(body io.Reader) string {
		var (
			reply struct {
				Context  string       `json:"@odata.context"`
				Nextlink string       `json:"@odata.nextLink"`
				Data     []Attachment `json:"value"`
			}
		)
		if err = json.NewDecoder(body).Decode(&reply); err == nil {
			list = append(list, reply.Data...)
			count += len(reply.Data)
			if count >= max {
				return ""
			}
			return reply.Nextlink
		}
		return ""
	})
	return list, err
}
