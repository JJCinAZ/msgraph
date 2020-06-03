package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/jjcinaz/msgraph"
	"github.com/jjcinaz/msgraph/examples/console"
	filecache "github.com/jjcinaz/msgraph/filecache"
	"github.com/qeesung/image2ascii/convert"
	"image"
	"log"
	"os"
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
	flag.StringVar(&userid, "u", "", "User Id or User Principal Name (email)")
	flag.BoolVar(&debugmode, "debug", false, "enable debug mode")
	flag.Parse()
	c, err = msgraph.NewClient(context.Background(), tenantid, clientid, clientsecret,
		[]string{"user.read", "Calendars.Read"}, filecache.New(), time.Minute)
	if err != nil {
		panic(err)
	}
	defer c.Close()
	if debugmode {
		c.SetAPILogging(log.New(os.Stdout, "", 0))
	}
	asciiProfilePhoto(c, userid)
}

func asciiProfilePhoto(c *msgraph.Client, id string) {
	var (
		err  error
		w, h int
		img  image.Image
	)
	w, h, err = console.InitConsole()
	if err != nil {
		fmt.Println(err)
		return
	}
	img, _, err = c.GetUserPhoto(id)
	if err == nil {
		converter := convert.NewImageConverter()
		fmt.Print(converter.Image2ASCIIString(img, &convert.Options{
			FixedWidth:  w,
			FixedHeight: h,
			Colored:     true,
		}))
	} else {
		fmt.Println(err)
	}
}
