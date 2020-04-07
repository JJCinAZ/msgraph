package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jjcinaz/msgraph"
)

func main() {
	var (
		err                           error
		c                             *msgraph.Client
		tenantid, clientid, clientkey string
	)
	tenantid = os.Getenv("AZURE_TENANTID")
	clientid = os.Getenv("AZURE_CLIENTID")
	clientkey = os.Getenv("AZURE_CLIENTKEY")
	if len(tenantid) == 0 {
		fmt.Println("Missing environment variable AZURE_TENANTID")
	}
	if len(clientid) == 0 {
		fmt.Println("Missing environment variable AZURE_CLIENTID")
	}
	if len(clientkey) == 0 {
		fmt.Println("Missing environment variable AZURE_CLIENTKEY")
	}
	c, err = msgraph.NewKeyClient(context.Background(), tenantid, clientid, clientkey)
	if err != nil {
		panic(err)
	}
	defer c.Close()
	upn := "jdoe@acme.com"
	msgs, err := c.ListMessages(upn, msgraph.OptionTextMailBody(),
		msgraph.OptionFilter("hasAttachments eq true and isRead eq false"),
		msgraph.OptionMaxItems(100), msgraph.OptionSelect("sender"), msgraph.OptionSelect("createdDateTime"),
		msgraph.OptionSelect("subject"), msgraph.OptionSelect("hasAttachments"))
	if err != nil {
		fmt.Println(err)
	} else {
		for _, m := range msgs {
			fmt.Println(m.CreatedDateTime, m.Sender, m.Subject)
			attachments, err := c.ListAttachments(upn, m.ID)
			if err != nil {
				fmt.Println(err)
			} else {
				for _, a := range attachments {
					fmt.Println(a.Name, a.Size)
					a.WriteFile(filepath.Base(a.Name), 0666)
				}
			}
		}
	}
}
