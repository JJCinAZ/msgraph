package msgraph

import (
	"encoding/json"
	"io"
	"net/url"
)

func (c *Client) ListCalendarGroups(upn string, options ...ApiOption) ([]CalendarGroup, error) {
	var (
		err error
	)

	apiUrl, err := formatOptions("https://graph.microsoft.com/v1.0/users/"+url.PathEscape(upn)+"/calendarGroups",
		options)
	if err != nil {
		return nil, err
	}
	list := make([]CalendarGroup, 0, 8)
	headers := make(map[string]string)
	err2 := c.executeGetList(apiUrl, headers, func(body io.Reader) string {
		var (
			reply struct {
				Context  string          `json:"@odata.context"`
				Nextlink string          `json:"@odata.nextLink"`
				Data     []CalendarGroup `json:"value"`
			}
		)
		if err = json.NewDecoder(body).Decode(&reply); err == nil {
			list = append(list, reply.Data...)
			return reply.Nextlink
		}
		return ""
	})
	if err == nil {
		err = err2
	}
	return list, err
}

func (c *Client) ListCalendarsInGroup(upn string, calendarGroupId string, options ...ApiOption) ([]Calendar, error) {
	var (
		err error
	)

	urlbase := "https://graph.microsoft.com/v1.0/users/" + url.PathEscape(upn)
	if calendarGroupId == DefaultCalendarGroup {
		urlbase = urlbase + "/calendarGroup/calendars"
	} else {
		urlbase = urlbase + "/calendarGroups/" + calendarGroupId + "/calendars"
	}
	apiUrl, err := formatOptions(urlbase, options)
	if err != nil {
		return nil, err
	}
	list := make([]Calendar, 0, 8)
	headers := make(map[string]string)
	err2 := c.executeGetList(apiUrl, headers, func(body io.Reader) string {
		var (
			reply struct {
				Context  string     `json:"@odata.context"`
				Nextlink string     `json:"@odata.nextLink"`
				Data     []Calendar `json:"value"`
			}
		)
		if err = json.NewDecoder(body).Decode(&reply); err == nil {
			list = append(list, reply.Data...)
			return reply.Nextlink
		}
		return ""
	})
	if err == nil {
		err = err2
	}
	return list, err
}
