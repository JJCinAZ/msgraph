package msgraph

import (
	"encoding/json"
	"io"
	"net/url"
	"time"
)

type EmailAddress struct {
	//OdataType string `json:"@odata.type"`
	Address string `json:"address"`
	Name    string `json:"name"`
}

type Recipient struct {
	EmailAddress EmailAddress `json:"emailAddress"`
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

type MailFolder struct {
	ChildFolderCount              int                       `json:"childFolderCount"`
	DisplayName                   string                    `json:"displayName"`
	ID                            string                    `json:"id"`
	ParentFolderID                string                    `json:"parentFolderId"`
	TotalItemCount                int                       `json:"totalItemCount"`
	UnreadItemCount               int                       `json:"unreadItemCount"`
	WellKnownName                 string                    `json:"wellKnownName"`
	ChildFolders                  []MailFolder              `json:"childFolders"`
	MessageRules                  []MessageRule             `json:"messageRules"`
	Messages                      []Message                 `json:"messages"`
	MultiValueExtendedProperties  []MultiValueExtendedProp  `json:"multiValueExtendedProperties"`
	SingleValueExtendedProperties []SingleValueExtendedProp `json:"singleValueExtendedProperties"`
}

type MessageRule struct {
	Actions     MessageRuleActions    `json:"actions"`
	Conditions  MessageRulePredicates `json:"conditions"`
	DisplayName string                `json:"displayName"`
	Exceptions  MessageRulePredicates `json:"exceptions"`
	HasError    bool                  `json:"hasError"`
	ID          string                `json:"id"`
	IsEnabled   bool                  `json:"isEnabled"`
	IsReadOnly  bool                  `json:"isReadOnly"`
	Sequence    int                   `json:"sequence"`
}

type MessageRuleActions struct {
	AssignCategories      []string    `json:"assignCategories"`
	CopyToFolder          string      `json:"copyToFolder"`
	Delete                bool        `json:"delete"`
	ForwardAsAttachmentTo []Recipient `json:"forwardAsAttachmentTo"`
	ForwardTo             []Recipient `json:"forwardTo"`
	MarkAsRead            bool        `json:"markAsRead"`
	MarkImportance        string      `json:"markImportance"`
	MoveToFolder          string      `json:"moveToFolder"`
	PermanentDelete       bool        `json:"permanentDelete"`
	RedirectTo            []Recipient `json:"redirectTo"`
	StopProcessingRules   bool        `json:"stopProcessingRules"`
}

type MessageRulePredicates struct {
	BodyContains           []string    `json:"bodyContains"`
	BodyOrSubjectContains  []string    `json:"bodyOrSubjectContains"`
	Categories             []string    `json:"categories"`
	FromAddresses          []Recipient `json:"fromAddresses"`
	HasAttachments         bool        `json:"hasAttachments"`
	HeaderContains         []string    `json:"headerContains"`
	Importance             string      `json:"importance"`
	IsApprovalRequest      bool        `json:"isApprovalRequest"`
	IsAutomaticForward     bool        `json:"isAutomaticForward"`
	IsAutomaticReply       bool        `json:"isAutomaticReply"`
	IsEncrypted            bool        `json:"isEncrypted"`
	IsMeetingRequest       bool        `json:"isMeetingRequest"`
	IsMeetingResponse      bool        `json:"isMeetingResponse"`
	IsNonDeliveryReport    bool        `json:"isNonDeliveryReport"`
	IsPermissionControlled bool        `json:"isPermissionControlled"`
	IsReadReceipt          bool        `json:"isReadReceipt"`
	IsSigned               bool        `json:"isSigned"`
	IsVoicemail            bool        `json:"isVoicemail"`
	MessageActionFlag      string      `json:"messageActionFlag"`
	NotSentToMe            bool        `json:"notSentToMe"`
	RecipientContains      []string    `json:"recipientContains"`
	SenderContains         []string    `json:"senderContains"`
	Sensitivity            string      `json:"sensitivity"`
	SentCcMe               bool        `json:"sentCcMe"`
	SentOnlyToMe           bool        `json:"sentOnlyToMe"`
	SentToAddresses        []Recipient `json:"sentToAddresses"`
	SentToMe               bool        `json:"sentToMe"`
	SentToOrCcMe           bool        `json:"sentToOrCcMe"`
	SubjectContains        []string    `json:"subjectContains"`
	WithinSizeRange        SizeRange   `json:"withinSizeRange"`
}

type SizeRange struct {
	MaximumSize int `json:"maximumSize"`
	MinimumSize int `json:"minimumSize"`
}

func (c *Client) ListMessages(upn string, options ...ApiOption) ([]Message, error) {
	return c.ListMessagesInFolder(upn, "", options...)
}

func (c *Client) ListMessagesInFolder(upn string, folderId string, options ...ApiOption) ([]Message, error) {
	var (
		err    error
		apiUrl string
	)
	url := "https://graph.microsoft.com/v1.0/users/" + url.PathEscape(upn)
	if len(folderId) == 0 {
		url = url + "/messages"
	} else {
		url = url + "/mailFolders/" + folderId + "messages"
	}
	if apiUrl, err = formatOptions(url, options); err != nil {
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
