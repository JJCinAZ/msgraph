package msgraph

import (
	"encoding/json"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/url"
	"strings"
)

type User struct {
	AboutMe                     string   `json:"aboutMe,omitempty"`
	AccountEnabled              bool     `json:"accountEnabled"`
	Birthday                    string   `json:"birthday,omitempty"`
	BusinessPhones              []string `json:"businessPhones,omitempty"`
	City                        string   `json:"city,omitempty"`
	CompanyName                 string   `json:"companyName,omitempty"`
	Country                     string   `json:"country,omitempty"`
	Department                  string   `json:"department,omitempty"`
	DisplayName                 string   `json:"displayName,omitempty"`
	EmployeeID                  string   `json:"employeeId,omitempty"`
	FaxNumber                   string   `json:"faxNumber,omitempty"`
	GivenName                   string   `json:"givenName,omitempty"`
	HireDate                    string   `json:"hireDate,omitempty"`
	ID                          string   `json:"id"`
	IsResourceAccount           bool     `json:"isResourceAccount"`
	JobTitle                    string   `json:"jobTitle,omitempty"`
	LastPasswordChangeDateTime  string   `json:"lastPasswordChangeDateTime,omitempty"`
	Mail                        string   `json:"mail,omitempty"`
	MailNickname                string   `json:"mailNickname,omitempty"`
	MobilePhone                 string   `json:"mobilePhone,omitempty"`
	OfficeLocation              string   `json:"officeLocation,omitempty"`
	OnPremisesDistinguishedName string   `json:"onPremisesDistinguishedName,omitempty"`
	OnPremisesDomainName        string   `json:"onPremisesDomainName,omitempty"`
	OnPremisesImmutableID       string   `json:"onPremisesImmutableId,omitempty"`
	OnPremisesLastSyncDateTime  string   `json:"onPremisesLastSyncDateTime,omitempty"`
	OnPremisesSamAccountName    string   `json:"onPremisesSamAccountName,omitempty"`
	OnPremisesSyncEnabled       bool     `json:"onPremisesSyncEnabled,omitempty"`
	OnPremisesUserPrincipalName string   `json:"onPremisesUserPrincipalName,omitempty"`
	PostalCode                  string   `json:"postalCode,omitempty"`
	PreferredDataLocation       string   `json:"preferredDataLocation,omitempty"`
	PreferredLanguage           string   `json:"preferredLanguage,omitempty"`
	ProxyAddresses              []string `json:"proxyAddresses,omitempty"`
	ShowInAddressList           bool     `json:"showInAddressList"`
	State                       string   `json:"state,omitempty"`
	StreetAddress               string   `json:"streetAddress,omitempty"`
	Surname                     string   `json:"surname,omitempty"`
	UserPrincipalName           string   `json:"userPrincipalName"`
}

type PhotoInfo struct {
	ContentType string `json:"@odata.mediaContentType"`
	Height      int    `json:"height"`
	ID          string `json:"id"`
	Width       int    `json:"width"`
}

func (c *Client) GetUserList() ([]User, error) {
	var (
		parms = []string{
			//"aboutMe",
			"accountEnabled",
			//"birthday",
			"businessPhones",
			//"city",
			"companyName",
			//"country",
			"department",
			"displayName",
			"employeeId",
			"faxNumber",
			"givenName",
			//"hireDate",
			"id",
			"isResourceAccount",
			"jobTitle",
			"lastPasswordChangeDateTime",
			"mail",
			"mailNickname",
			"mobilePhone",
			"officeLocation",
			//"onPremisesDistinguishedName",
			"onPremisesDomainName",
			//"onPremisesImmutableId",
			"onPremisesLastSyncDateTime",
			"onPremisesSamAccountName",
			"onPremisesSyncEnabled",
			"onPremisesUserPrincipalName",
			//"postalCode",
			//"preferredDataLocation",
			//"preferredLanguage",
			"proxyAddresses",
			"showInAddressList",
			"state",
			"streetAddress",
			"surname",
			"userPrincipalName",
		}
		err error
	)

	apiUrl := "https://graph.microsoft.com/v1.0/users?$select=" + strings.Join(parms, ",")
	// MS Graph API says it supports a filter option but that filter option doesn't actually
	// work for onPremisesSamAccountName (it doens't actually work for most of the fields);
	// but that's expected Microsoft quality.  I'm sure they will replace the entire API
	// with "something better" in a year anyway.  For now we have to manually filter.
	users := make([]User, 0, 128)
	c.executeGetList(apiUrl, nil, func(body io.Reader) string {
		var (
			reply struct {
				Context  string `json:"@odata.context"`
				Nextlink string `json:"@odata.nextLink"`
				Data     []User `json:"value"`
			}
		)
		if err = json.NewDecoder(body).Decode(&reply); err == nil {
			for _, u := range reply.Data {
				if len(u.OnPremisesSamAccountName) > 0 {
					users = append(users, u)
				}
			}
			return reply.Nextlink
		}
		return ""
	})
	return users, err
}

// Get profile picture given a UserPrincipalName (e.g. "bob@acme.com") or User Id (UUID)
// If successful, returns an Image and the name of the image format (jpeg, png, or gif)
func (c *Client) GetUserPhoto(upn string) (image.Image, string, error) {
	var (
		err  error
		img  image.Image
		info string
	)
	apiUrl := "https://graph.microsoft.com/v1.0/users/" + url.PathEscape(upn) + "/photo/$value"
	err = c.executeGet(apiUrl, func(reader io.Reader) error {
		var err2 error
		img, info, err2 = image.Decode(reader)
		return err2
	})
	if err == nil {
		return img, info, nil
	}
	return nil, info, err
}

// Get profile picture info given a UserPrincipalName (e.g. "bob@acme.com") or User Id (UUID)
func (c *Client) GetUserPhotoInfo(upn string) (PhotoInfo, error) {
	var pi PhotoInfo
	apiUrl := "https://graph.microsoft.com/v1.0/users/" + url.PathEscape(upn) + "/photo"
	err := c.executeGetJson(apiUrl, &pi)
	return pi, err
}
