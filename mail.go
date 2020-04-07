package msgraph

import (
	"encoding/json"
	"io"
	"net/url"
	"time"
)

type Recipient struct {
	EmailAddress struct {
		//OdataType string `json:"@odata.type"`
		Address string `json:"address"`
		Name    string `json:"name"`
	} `json:"emailAddress"`
}

type ItemBody struct {
	ContentType string `json:"contentType"`
	Content     string `json:"content"`
}

type FollowUpFlag struct {
	CompletedDateTime DateTimeTimeZone `json:"completedDateTime"`
	DueDateTime       DateTimeTimeZone `json:"dueDateTime"`
	FlagStatus        string           `json:"flagStatus"` // Possible values are notFlagged, complete, and flagged.
	StartDateTime     DateTimeTimeZone `json:"startDateTime"`
}

type InternetMessageHeader struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Message struct {
	client                     *Client
	fromHasValue               bool
	BccRecipients              []Recipient             `json:"bccRecipients,omitempty"`
	Body                       ItemBody                `json:"body"`
	BodyPreview                string                  `json:"bodyPreview,omitempty"`
	Categories                 []string                `json:"categories,omitempty"`
	CcRecipients               []Recipient             `json:"ccRecipients,omitempty"`
	ChangeKey                  string                  `json:"changeKey,omitempty"`
	ConversationID             string                  `json:"conversationId,omitempty"`
	ConversationIndex          string                  `json:"conversationIndex,omitempty"`
	CreatedDateTime            string                  `json:"createdDateTime"`
	Flag                       *FollowUpFlag           `json:"flag,omitempty"`
	From                       Recipient               `json:"from"`
	HasAttachments             bool                    `json:"hasAttachments"`
	ID                         string                  `json:"id,omitempty"`
	Importance                 string                  `json:"importance,omitempty"`
	InferenceClassification    string                  `json:"inferenceClassification,omitempty"`
	InternetMessageHeaders     []InternetMessageHeader `json:"internetMessageHeaders,omitempty"`
	InternetMessageID          string                  `json:"internetMessageId,omitempty"`
	IsDeliveryReceiptRequested bool                    `json:"isDeliveryReceiptRequested"`
	IsDraft                    bool                    `json:"isDraft"`
	IsRead                     bool                    `json:"isRead"`
	IsReadReceiptRequested     bool                    `json:"isReadReceiptRequested"`
	LastModifiedDateTime       string                  `json:"lastModifiedDateTime,omitempty"`
	ParentFolderID             string                  `json:"parentFolderId,omitempty"`
	ReceivedDateTime           string                  `json:"receivedDateTime,omitempty"`
	ReplyTo                    []Recipient             `json:"replyTo,omitempty"`
	Sender                     Recipient               `json:"sender"`
	SentDateTime               string                  `json:"sentDateTime,omitempty"`
	Subject                    string                  `json:"subject"`
	ToRecipients               []Recipient             `json:"toRecipients,omitempty"`
	UniqueBody                 *ItemBody               `json:"uniqueBody,omitempty"`
	WebLink                    string                  `json:"webLink,omitempty"`
}

func (c *Client) ListMessages(upn string, options ...ApiOption) ([]Message, error) {
	var (
		err error
	)

	apiUrl, err := formatOptions("https://graph.microsoft.com/v1.0/users/"+url.PathEscape(upn)+"/messages",
		options)
	if err != nil {
		return nil, err
	}
	max, count := getMaxItemOption(options), 0
	msgs := make([]Message, 0, 1024)
	headers := make(map[string]string)
	if getTextMailBody(options) {
		headers["Prefer"] = `outlook.body-content-type="text"`
	}
	err2 := c.executeGetList(apiUrl, headers, func(body io.Reader) string {
		var (
			reply struct {
				Context  string    `json:"@odata.context"`
				Nextlink string    `json:"@odata.nextLink"`
				Data     []Message `json:"value"`
			}
		)
		if err = json.NewDecoder(body).Decode(&reply); err == nil {
			msgs = append(msgs, reply.Data...)
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
	return msgs, err
}

func (c *Client) DeleteMessage(upn, msgid string) error {
	return c.executeDelete("https://graph.microsoft.com/v1.0/users/" + url.PathEscape(upn) + "/messages/" + url.PathEscape(msgid))
}

func (m Message) Send(upn string, saveToSentItems bool) error {
	var (
		data struct {
			Msg             Message `json:"message"`
			SaveToSentItems bool    `json:"saveToSentItems"`
		}
	)
	if len(upn) == 0 {
		upn = m.Sender.EmailAddress.Address
	}
	if !m.fromHasValue {
		m.From = m.Sender
	}
	data.Msg = m
	data.SaveToSentItems = saveToSentItems
	return m.client.executePost("https://graph.microsoft.com/v1.0/users/"+url.PathEscape(upn)+"/sendMail",
		data, nil)
}

func (c *Client) NewMessage() Message {
	return Message{
		client:          c,
		CreatedDateTime: time.Now().Format(time.RFC3339),
	}
}

func (m *Message) SetSubject(subject string) *Message {
	m.Subject = subject
	return m
}

func (m *Message) SetSender(name, address string) *Message {
	m.Sender.EmailAddress.Name = name
	m.Sender.EmailAddress.Address = address
	return m
}

func (m *Message) SetFrom(name, address string) *Message {
	m.From.EmailAddress.Name = name
	m.From.EmailAddress.Address = address
	return m
}

func (m *Message) SetBody(body ItemBody) *Message {
	m.Body = body
	return m
}

func (m *Message) AddToRecipient(name, address string) *Message {
	var r Recipient
	r.EmailAddress.Name = name
	r.EmailAddress.Address = address
	m.ToRecipients = append(m.ToRecipients, r)
	return m
}

func (m *Message) AddCcRecipient(name, address string) *Message {
	var r Recipient
	r.EmailAddress.Name = name
	r.EmailAddress.Address = address
	m.CcRecipients = append(m.CcRecipients, r)
	return m
}

func (m *Message) AddBccRecipient(name, address string) *Message {
	var r Recipient
	r.EmailAddress.Name = name
	r.EmailAddress.Address = address
	m.BccRecipients = append(m.BccRecipients, r)
	return m
}

func (m *Message) AddReplyTo(name, address string) *Message {
	var r Recipient
	r.EmailAddress.Name = name
	r.EmailAddress.Address = address
	m.ReplyTo = append(m.ReplyTo, r)
	return m
}

func (c *Client) NewBody() ItemBody {
	return ItemBody{}
}

func (i *ItemBody) SetText(text string) *ItemBody {
	i.ContentType = "text"
	i.Content = text
	return i
}

func (i *ItemBody) SetHtml(html string) *ItemBody {
	i.ContentType = "html"
	i.Content = html
	return i
}
