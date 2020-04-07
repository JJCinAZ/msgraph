package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Simply-Bits/msgraph"
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
	c.SetAPILogging(log.New(os.Stdout, "", 0))
	send1(c)
}

func send1(c *msgraph.Client) {
	msg := c.NewMessage()
	body := c.NewBody()
	body.SetText("Hello World")
	msg.SetSubject("My Subject").SetSender("John Doe", "jdoe@acme.com")
	msg.AddToRecipient("Frank Smith", "franksmith@gmail.com")
	msg.SetBody(body)
	err := msg.Send("jdoe@acme.com", false)
	if err != nil {
		fmt.Println(err)
	}
}
