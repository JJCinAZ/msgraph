package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/kofoworola/godate"

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
	//c.SetAPILogging(log.New(os.Stdout, "", 0))
	weekstart := godate.Create(time.Now())
	events, err := c.GetCalendarView("SES-LgConf@simplybits.com", msgraph.OptionTextMailBody(),
		msgraph.OptionStartDateTime(weekstart.Time), msgraph.OptionEndDateTime(weekstart.EndOfMonth().Time))
	if err != nil {
		fmt.Println(err)
	} else {
		sort.Slice(events, func(i, j int) bool {
			return events[i].Start.Native.Before(events[j].Start.Native)
		})
		for _, e := range events {
			if e.Type == "occurrence" {
				fmt.Print("🔁\t")
			} else {
				fmt.Print("\t")
			}
			fmt.Printf("%-20.20s\t%s\t%s\n",
				e.Subject, e.Start.Native.Local().Format(time.RFC850), e.End.Native.Local().Format(time.RFC850))
		}
	}
}
