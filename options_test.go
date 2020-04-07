package msgraph

import (
	"net/url"
	"testing"
)

func TestOptionSearch_escapeValue(t *testing.T) {
	type fields struct {
		value string
		prop  string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"test1", fields{value: "abc"}, "abc"},
		{"test2", fields{value: "ab'c'"}, "ab''c''"},
		{"test3", fields{value: `abc "def"`}, `abc \"def\"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := optSearch{
				value: tt.fields.value,
				prop:  tt.fields.prop,
			}
			if got := o.escapeValue(); got != tt.want {
				t.Errorf("escapeValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_formatOptions(t *testing.T) {
	type args struct {
		apiUrl  string
		options []ApiOption
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"test1", args{
			apiUrl: "https://graph.microsoft.com/v1.0/me/messages",
			options: []ApiOption{
				optFilter{"(from/emailAddress/address) eq 'jdoe@acme.com'"},
			},
		}, "https://graph.microsoft.com/v1.0/me/messages?$filter=(from/emailAddress/address) eq 'jdoe@acme.com'", false},
		{"test2", args{
			apiUrl:  "https://graph.microsoft.com/v1.0/me/messages",
			options: []ApiOption{},
		}, "https://graph.microsoft.com/v1.0/me/messages", false},
		{"test2", args{
			apiUrl: "https://graph.microsoft.com/v1.0/me/messages",
			options: []ApiOption{
				optSelect{"sender"},
				optSelect{"body"},
				optFilter{"(from/emailAddress/address) eq 'jdoe@acme.com'"},
			},
		}, "https://graph.microsoft.com/v1.0/me/messages?$filter=(from/emailAddress/address) eq 'jdoe@acme.com'&$select=sender,body", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := formatOptions(tt.args.apiUrl, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("formatOptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotUrl, _ := url.ParseRequestURI(got)             // Parse https://host.com/path?escapedquery
			gotQuery, _ := url.QueryUnescape(gotUrl.RawQuery) // get the unescaped version of thr query
			gotUrl.RawQuery = ""                              // remove the query from the Url
			gotBase := gotUrl.String()                        // get everything except the query, https://host.com/path
			if len(gotQuery) > 0 {                            // if we have a query, add it back to the path
				gotBase += "?" + gotQuery
			}
			if gotBase != tt.want {
				t.Errorf("formatOptions() got = %v, want %v", got, tt.want)
			}
		})
	}
}
