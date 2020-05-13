package msgraph

import (
	"encoding/json"
	"time"
)

type DateTimeTimeZone struct {
	DateTime string    `json:"dateTime"`
	TimeZone string    `json:"timeZone"`
	Native   time.Time `json:"-"`
}

func (d *DateTimeTimeZone) UnmarshalJSON(b []byte) error {
	var internal struct {
		DateTime string `json:"dateTime"`
		TimeZone string `json:"timeZone"`
	}
	err := json.Unmarshal(b, &internal)
	if err == nil && len(internal.DateTime) > 0 {
		d.DateTime = internal.DateTime
		d.TimeZone = internal.TimeZone
		if loc, err2 := time.LoadLocation(internal.TimeZone); err2 == nil {
			d.Native, _ = time.ParseInLocation("2006-01-02T15:04:05.999999999", internal.DateTime, loc)
		}
	}
	return err
}
