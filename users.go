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
	AboutMe                     string   `json:"aboutMe"`
	AccountEnabled              bool     `json:"accountEnabled"`
	Birthday                    string   `json:"birthday"`
	BusinessPhones              []string `json:"businessPhones"`
	City                        string   `json:"city"`
	CompanyName                 string   `json:"companyName"`
	Country                     string   `json:"country"`
	Department                  string   `json:"department"`
	DisplayName                 string   `json:"displayName"`
	EmployeeID                  string   `json:"employeeId"`
	FaxNumber                   string   `json:"faxNumber"`
	GivenName                   string   `json:"givenName"`
	HireDate                    string   `json:"hireDate"`
	ID                          string   `json:"id"`
	IsResourceAccount           bool     `json:"isResourceAccount"`
	JobTitle                    string   `json:"jobTitle"`
	LastPasswordChangeDateTime  string   `json:"lastPasswordChangeDateTime"`
	Mail                        string   `json:"mail"`
	MailNickname                string   `json:"mailNickname"`
	MobilePhone                 string   `json:"mobilePhone"`
	OfficeLocation              string   `json:"officeLocation"`
	OnPremisesDistinguishedName string   `json:"onPremisesDistinguishedName"`
	OnPremisesDomainName        string   `json:"onPremisesDomainName"`
	OnPremisesImmutableID       string   `json:"onPremisesImmutableId"`
	OnPremisesLastSyncDateTime  string   `json:"onPremisesLastSyncDateTime"`
	OnPremisesSamAccountName    string   `json:"onPremisesSamAccountName"`
	OnPremisesSyncEnabled       bool     `json:"onPremisesSyncEnabled"`
	OnPremisesUserPrincipalName string   `json:"onPremisesUserPrincipalName"`
	PostalCode                  string   `json:"postalCode"`
	PreferredDataLocation       string   `json:"preferredDataLocation"`
	PreferredLanguage           string   `json:"preferredLanguage"`
	ProxyAddresses              []string `json:"proxyAddresses"`
	ShowInAddressList           bool     `json:"showInAddressList"`
	State                       string   `json:"state"`
	StreetAddress               string   `json:"streetAddress"`
	Surname                     string   `json:"surname"`
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
	c.executeGetList(apiUrl, func(body io.Reader) string {
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
