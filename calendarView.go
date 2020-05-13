package msgraph

import (
	"encoding/json"
	"io"
	"net/url"
)

// Get the occurrences, exceptions, and single instances of events in a calendar view defined by a time range,
// from a user's default calendar (../me/calendarview) or some other calendar of the user's.
// Must specify a UserPrincipalName (e.g. "bob@acme.com") or User Id (UUID).
func (c *Client) GetCalendarView(upn string, options ...ApiOption) ([]Event, error) {
	var (
		err error
	)

	apiUrl, err := formatOptions("https://graph.microsoft.com/v1.0/users/"+url.PathEscape(upn)+"/calendarView",
		options)
	if err != nil {
		return nil, err
	}
	events := make([]Event, 0, 128)
	headers := make(map[string]string)
	if getTextMailBody(options) {
		headers["Prefer"] = `outlook.body-content-type="text"`
	}
	err2 := c.executeGetList(apiUrl, headers, func(body io.Reader) string {
		var (
			reply struct {
				Context  string  `json:"@odata.context"`
				Nextlink string  `json:"@odata.nextLink"`
				Data     []Event `json:"value"`
			}
		)
		if err = json.NewDecoder(body).Decode(&reply); err == nil {
			events = append(events, reply.Data...)
			return reply.Nextlink
		}
		return ""
	})
	if err == nil {
		err = err2
	}
	return events, err
}
