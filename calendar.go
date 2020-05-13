package msgraph

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"time"
)

const (
	DefaultCalendarGroup = "DEFAULT"
)

type Calendar struct {
	ID                            string   `json:"id"`
	Name                          string   `json:"name"`
	Color                         string   `json:"color"`
	ChangeKey                     string   `json:"changeKey"`
	CanShare                      bool     `json:"canShare"`
	CanViewPrivateItems           bool     `json:"canViewPrivateItems"`
	CanEdit                       bool     `json:"canEdit"`
	AllowedOnlineMeetingProviders []string `json:"allowedOnlineMeetingProviders"`
	DefaultOnlineMeetingProvider  string   `json:"defaultOnlineMeetingProvider"`
	IsTallyingResponses           bool     `json:"isTallyingResponses"`
	IsRemovable                   bool     `json:"isRemovable"`
	Owner                         struct {
		Name    string `json:"name"`
		Address string `json:"address"`
	} `json:"owner"`
}

type CalendarGroup struct {
	ChangeKey string `json:"changeKey"`
	ClassID   string `json:"classId"`
	ID        string `json:"id"`
	Name      string `json:"name"`
}

type Event struct {
	ID                         string               `json:"id"`
	CreatedDateTime            time.Time            `json:"createdDateTime"`
	LastModifiedDateTime       time.Time            `json:"lastModifiedDateTime"`
	ChangeKey                  string               `json:"changeKey"`
	Categories                 []string             `json:"categories"`
	OriginalStartTimeZone      string               `json:"originalStartTimeZone"`
	OriginalEndTimeZone        string               `json:"originalEndTimeZone"`
	ICalUid                    string               `json:"iCalUid"`
	ReminderMinutesBeforeStart int                  `json:"reminderMinutesBeforeStart"`
	IsReminderOn               bool                 `json:"isReminderOn"`
	HasAttachments             bool                 `json:"hasAttachments"`
	Subject                    string               `json:"subject"`
	BodyPreview                string               `json:"bodyPreview"`
	Importance                 string               `json:"importance"`
	Sensitivity                string               `json:"sensitivity"`
	IsAllDay                   bool                 `json:"isAllDay"`
	IsCancelled                bool                 `json:"isCancelled"`
	IsOrganizer                bool                 `json:"isOrganizer"`
	ResponseRequested          bool                 `json:"responseRequested"`
	SeriesMasterID             *string              `json:"seriesMasterId,omitempty"`
	ShowAs                     string               `json:"showAs"`
	Type                       string               `json:"type"`
	WebLink                    string               `json:"webLink"`
	IsOnlineMeeting            bool                 `json:"isOnlineMeeting"`
	OnlineMeetingProvider      string               `json:"onlineMeetingProvider"`
	AllowNewTimeProposals      bool                 `json:"allowNewTimeProposals"`
	IsDraft                    bool                 `json:"isDraft"`
	Recurrence                 *PatternedRecurrence `json:"recurrence,omitempty"`
	OnlineMeeting              *OnlineMeetingInfo   `json:"onlineMeeting,omitempty"`
	ResponseStatus             ResponseStatus       `json:"responseStatus"`
	Body                       ItemBody             `json:"body"`
	Start                      DateTimeTimeZone     `json:"start"`
	End                        DateTimeTimeZone     `json:"end"`
	Location                   Location             `json:"location"`
	Locations                  []Location           `json:"locations"`
	Attendees                  []Attendee           `json:"attendees"`
	Organizer                  Recipient            `json:"organizer"`
	// From Beta API
	//TransactionID              *string       `json:"transactionId"`
	//UID                        string        `json:"uid"`
}

type PhysicalAddress struct {
	City            string `json:"city"`
	CountryOrRegion string `json:"countryOrRegion"`
	PostalCode      string `json:"postalCode"`
	State           string `json:"state"`
	Street          string `json:"street"`
}

type OutlookGeoCoordinates struct {
	Accuracy         float64 `json:"accuracy"`
	Altitude         float64 `json:"altitude"`
	AltitudeAccuracy float64 `json:"altitudeAccuracy"`
	Latitude         float64 `json:"latitude"`
	Longitude        float64 `json:"longitude"`
}

type Location struct {
	Address              *PhysicalAddress       `json:"address,omitempty"`
	Coordinates          *OutlookGeoCoordinates `json:"coordinates,omitempty"`
	DisplayName          string                 `json:"displayName"`
	LocationEmailAddress *string                `json:"locationEmailAddress,omitempty"`
	LocationURI          *string                `json:"locationUri,omitempty"`
	LocationType         string                 `json:"locationType"`
}

type ResponseStatus struct {
	Response string    `json:"response"`
	Time     time.Time `json:"time"`
}

type TimeSlot struct {
	Start DateTimeTimeZone `json:"start"`
	End   DateTimeTimeZone `json:"end"`
}

type Attendee struct {
	Type            string         `json:"type"`
	Status          ResponseStatus `json:"status"`
	EmailAddress    EmailAddress   `json:"emailAddress"`
	ProposedNewTime *TimeSlot      `json:"proposedNewTime,omitempty"`
}

type OnlineMeetingInfo struct {
	ConferenceID    *string  `json:"conferenceId"`
	JoinURL         *string  `json:"joinUrl"`
	Phones          []Phone  `json:"phones,omitempty"`
	QuickDial       *string  `json:"quickDial,omitempty"`
	TollFreeNumbers []string `json:"tollFreeNumbers,omitemptu"`
	TollNumber      *string  `json:"tollNumber,omitempty"`
}

type Phone struct {
	Number string `json:"number"`
	Type   string `json:"type"`
}

type PatternedRecurrence struct {
	Pattern RecurrencePattern `json:"pattern"`
	Range   RecurrenceRange   `json:"range"`
}

type RecurrencePattern struct {
	DayOfMonth     *int     `json:"dayOfMonth,omitempty"`
	DaysOfWeek     []string `json:"daysOfWeek,omitempty"`
	FirstDayOfWeek *string  `json:"firstDayOfWeek,omitempty"`
	Index          *string  `json:"index,omitempty"`
	Interval       int      `json:"interval"`
	Month          *int     `json:"month,omitempty"`
	Type           string   `json:"type"`
}

type RecurrenceRange struct {
	EndDate             *string `json:"endDate,omitempty"`
	NumberOfOccurrences int     `json:"numberOfOccurrences"`
	RecurrenceTimeZone  *string `json:"recurrenceTimeZone,omitempty"`
	StartDate           string  `json:"startDate"`
	Type                string  `json:"type"`
}

func (c *Client) ListCalendars(upn string, options ...ApiOption) ([]Calendar, error) {
	var (
		err error
	)

	apiUrl, err := formatOptions("https://graph.microsoft.com/v1.0/users/"+url.PathEscape(upn)+"/calendars",
		options)
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

// Get a Calendar object for a user or the default calendar of an Office 365 Group.
// Must specify a UserPrincipalName (e.g. "bob@acme.com") or User Id (UUID).
// If neither calendarId nor calendarGroupId are specified, this routine gets the user or groups
// default calendar.  You can get the default calendar in the default calendarGroup
// by passing an empty calendarId and the msgraph.DefaultCalendarGroup constant
func (c *Client) GetCalendar(upn string, calendarId string, calendarGroupId string) (*Calendar, error) {
	var (
		err error
		cal Calendar
	)
	apiUrl := "https://graph.microsoft.com/v1.0/users/" + url.PathEscape(upn)
	if len(calendarGroupId) > 0 {
		// Specified a calendarGroup ID, so we either get the default calendar in the default calendarGroup
		// or we get a specific calendar from a specific calendarGroup
		if len(calendarId) == 0 {
			return nil, fmt.Errorf("must specify a calendar Id when targeting a calendarGroup")
		}
		if calendarGroupId == DefaultCalendarGroup {
			apiUrl = apiUrl + "/calendarGroup/calendars/" + calendarId
		} else {
			apiUrl = apiUrl + "/calendarGroup/" + calendarGroupId + "/calendars/" + calendarId
		}
	} else if len(calendarId) == 0 {
		apiUrl = apiUrl + "/calendar"
	} else {
		apiUrl = apiUrl + "/calendars/" + calendarId
	}
	err = c.executeGet(apiUrl, func(reader io.Reader) error {
		return json.NewDecoder(reader).Decode(&cal)
	})
	if err == nil {
		return &cal, nil
	}
	return nil, err
}

// Get the default Calendar object for an Office 365 Group.
// Must specify a UserPrincipalName (e.g. "bob@acme.com") or Id (UUID).
func (c *Client) GetGroupCalendar(upn string) (*Calendar, error) {
	var (
		err error
		cal Calendar
	)
	apiUrl := "https://graph.microsoft.com/v1.0/groups/" + url.PathEscape(upn) + "/calendar"
	err = c.executeGet(apiUrl, func(reader io.Reader) error {
		return json.NewDecoder(reader).Decode(&cal)
	})
	if err == nil {
		return &cal, nil
	}
	return nil, err
}
