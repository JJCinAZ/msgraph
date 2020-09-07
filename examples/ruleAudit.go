package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gookit/color"
	"github.com/jjcinaz/msgraph"
	filecache "github.com/jjcinaz/msgraph/filecache"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	var (
		err                              error
		c                                *msgraph.Client
		userid                           string
		debugmode                        bool
		tenantid, clientid, clientsecret string
	)
	tenantid = os.Getenv("AZURE_TENANTID")
	clientid = os.Getenv("AZURE_CLIENTID")
	clientsecret = os.Getenv("AZURE_CLIENTSECRET")
	if len(tenantid) == 0 {
		fmt.Println("Missing environment variable AZURE_TENANTID")
	}
	if len(clientid) == 0 {
		fmt.Println("Missing environment variable AZURE_CLIENTID")
	}
	if len(clientsecret) == 0 {
		fmt.Println("Missing environment variable AZURE_CLIENTSECRET")
	}
	flag.StringVar(&userid, "u", "", "Email address of user to check (if not supplied, all users are checked)")
	flag.BoolVar(&debugmode, "debug", false, "enable debug mode")
	flag.Parse()
	if len(userid) == 0 {
		c, err = msgraph.NewKeyClient(context.Background(), tenantid, clientid, clientsecret)
		if err != nil {
			panic(err)
		}
		defer c.Close()
		users, err := c.GetUserList()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		for _, u := range users {
			if len(u.Mail) > 0 {
				doAudit(c, u.Mail)
			}
		}
	} else {
		c, err = msgraph.NewClient(context.Background(), tenantid, clientid, clientsecret,
			[]string{"User.Read.All", "Mail.ReadBasic", "MailboxSettings.Read"},
			filecache.New(), time.Minute)
		if err != nil {
			panic(err)
		}
		defer c.Close()
		if debugmode {
			c.SetAPILogging(log.New(os.Stdout, "", 0))
		}
		doAudit(c, userid)
	}
}

func doAudit(c *msgraph.Client, userid string) {
	var (
		err   error
		rules []msgraph.MessageRule
	)
	fmt.Printf("Checking user %s\n-------------------------------\n", userid)
	rules, err = c.ListMessageRules(userid)
	if err != nil {
		fmt.Println(err)
	} else {
		for i, m := range rules {
			if i > 0 {
				fmt.Println("")
			}
			color.New(color.BgWhite, color.FgBlack).Print(m.DisplayName)
			fmt.Println(":")
			if notSameDomain(userid, m.Actions.ForwardTo...) {
				fmt.Println(m.Actions.ForwardTo)
			}
			if notSameDomain(userid, m.Actions.ForwardAsAttachmentTo...) {
				fmt.Println(m.Actions.ForwardAsAttachmentTo)
			}
			if notSameDomain(userid, m.Actions.RedirectTo...) {
				fmt.Println(m.Actions.RedirectTo)
			}
			if m.HasError {
				color.New(color.FgRed).Printf("Rule has an error")
			}
		}
	}
}

func notSameDomain(email string, dests ...msgraph.Recipient) bool {
	var (
		domain string
	)
	at := strings.LastIndex(email, "@")
	if at >= 0 {
		domain = email[at+1:]
	} else {
		return false
	}
	for _, dest := range dests {
		at = strings.LastIndex(dest.EmailAddress.Address, "@")
		if at >= 0 {
			if strings.EqualFold(dest.EmailAddress.Address[at+1:], domain) == false {
				return true
			}
		}
	}
	return false
}
