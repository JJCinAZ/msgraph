package msgraph

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"os/exec"
	"runtime"
	"time"
)

type callbackReturn struct {
	code    string
	err     error
	errdesc string // longer description some providers return
}

const closeWindowJS = `<!doctype html><html><body><script>close();</script></body></html>`

func (c *Client) startCallbackServer(ctx context.Context, state string, userWait time.Duration) chan callbackReturn {
	returnChan := make(chan callbackReturn)
	codeChan := make(chan callbackReturn)
	mux := http.NewServeMux()
	mux.HandleFunc("/authcb", func(w http.ResponseWriter, r *http.Request) {
		var x callbackReturn
		defer r.Body.Close()
		if v := r.FormValue("error"); len(v) > 0 {
			x.err = fmt.Errorf("%s", v)
			x.errdesc = r.FormValue("error_description")
			if len(x.errdesc) > 0 {
				w.Write([]byte(x.errdesc))
			} else {
				w.Write([]byte(v))
			}
		} else {
			if r.FormValue("state") != state {
				x.err = fmt.Errorf("invalid state during callback")
				http.Error(w, "state doesn't match", http.StatusBadRequest)
			} else {
				x.code = r.FormValue("code")
				w.Write([]byte(closeWindowJS))
			}
		}
		codeChan <- x
	})
	srv := &http.Server{Addr: "localhost:8001", Handler: mux}
	go func() {
		serveErr := srv.ListenAndServe()
		if serveErr == http.ErrServerClosed {
			serveErr = nil
		} else {
			codeChan <- callbackReturn{err: serveErr}
		}
	}()
	go func() {
		var ret callbackReturn
		select {
		case ret = <-codeChan:
			c, _ := context.WithTimeout(ctx, 3*time.Second)
			srv.Shutdown(c)
			srv.Close()
		case <-ctx.Done():
			ret.err = fmt.Errorf("context closed")
		case <-time.After(userWait):
			ret.err = fmt.Errorf("timeout, no callback received")
		}
		returnChan <- ret
	}()
	return returnChan
}

func generateNonce(len int) (string, error) {
	nonceBytes := make([]byte, len)
	_, err := rand.Read(nonceBytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(nonceBytes), nil
}

// openUrl opens a browser window to the specified location.
// This code originally appeared at:
//   http://stackoverflow.com/questions/10377243/how-can-i-launch-a-process-that-is-not-a-file-in-go
func openUrl(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("cannot open URL %s on this platform", url)
	}
	return err
}
