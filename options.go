package msgraph

import (
	"fmt"
	"math"
	"net/url"
	"strconv"
	"strings"
)

type ApiOption interface{}

type ApiOptions []ApiOption

type optSearch struct {
	value string
	prop  string // "attachment", "bcc", "body", "cc", "from", "participants", "received", "sent", "recipients", "to", "subject"
}
type optSelect struct {
	field string
}
type optFilter struct {
	filter string
}

type optPageSize struct {
	n int
}

type optMax struct {
	n int
}

type optTextMailBody struct {
}

func OptionPageSize(n int) ApiOption {
	return optPageSize{n: n}
}

func OptionTextMailBody() ApiOption {
	return optTextMailBody{}
}

func OptionSearch(value string) ApiOption {
	return optSearch{value: value}
}

func OptionSelect(field string) ApiOption {
	return optSelect{field: field}
}

func OptionFilter(filter string) ApiOption {
	return optFilter{filter: filter}
}

func OptionMaxItems(n int) ApiOption {
	return optMax{n: n}
}

func getMaxItemOption(options []ApiOption) int {
	for _, o := range options {
		if t, ok := o.(optMax); ok {
			return t.n
		}
	}
	return math.MaxInt32
}

func getTextMailBody(options []ApiOption) bool {
	for _, o := range options {
		if _, ok := o.(optTextMailBody); ok {
			return true
		}
	}
	return false
}

func formatOptions(apiUrl string, options []ApiOption) (string, error) {
	var (
		sel                strings.Builder
		nSrch, nFilt, nSel int
	)
	baseUrl, err := url.ParseRequestURI(apiUrl)
	if err != nil {
		return "", err
	}
	for _, o := range options {
		switch o.(type) {
		case optSearch:
			nSrch++
		case optFilter:
			nFilt++
		}
	}
	if nSrch > 1 {
		return "", fmt.Errorf("cannot have more than one search")
	} else if nSrch > 0 {
		if nFilt > 0 {
			return "", fmt.Errorf("cannot use filter with a search request")
		}
	}
	params := url.Values{}
	for _, o := range options {
		switch x := o.(type) {
		case optSearch:
			if len(x.prop) > 0 {
				params.Add("$search", fmt.Sprintf(`"%s:%s"`, x.prop, x.escapeValue()))
			} else {
				params.Add("$search", fmt.Sprintf(`"%s"`, x.escapeValue()))
			}
		case optSelect:
			if nSel > 0 {
				sel.WriteByte(',')
			}
			sel.WriteString(x.field)
			nSel++
		case optFilter:
			params.Add("$filter", x.filter)
		case optPageSize:
			params.Add("top", strconv.Itoa(x.n))
		}
	}
	if nSel > 0 {
		params.Add("$select", sel.String())
	}
	baseUrl.RawQuery = params.Encode()
	return baseUrl.String(), nil
}

func (o optSearch) escapeValue() string {
	var b strings.Builder
	for _, r := range o.value {
		switch r {
		case '"':
			b.WriteString(`\"`)
		case '\'':
			b.WriteString("''")
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}
