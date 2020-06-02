package tokencache

import (
	"encoding/gob"
	"fmt"
	"golang.org/x/oauth2"
	"hash/fnv"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/jjcinaz/msgraph"
)

type FileCache struct {
}

func New() *FileCache {
	return &FileCache{}
}

func (fc *FileCache) Load(c *msgraph.Client) (*oauth2.Token, error) {
	file, err := tokenCacheFile(&c.OauthConfig)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := new(oauth2.Token)
	err = gob.NewDecoder(f).Decode(t)
	return t, err
}

func (fc *FileCache) Save(c *msgraph.Client, token *oauth2.Token) error {
	file, err := tokenCacheFile(&c.OauthConfig)
	if err != nil {
		return err
	}
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	return gob.NewEncoder(f).Encode(token)
}

func tokenCacheFile(config *oauth2.Config) (string, error) {
	var (
		dir string
		err error
	)
	hash := fnv.New32a()
	hash.Write([]byte(config.ClientID))
	hash.Write([]byte(config.ClientSecret))
	hash.Write([]byte(strings.Join(config.Scopes, " ")))
	fn := fmt.Sprintf("oauth-tok%v", hash.Sum32())
	if dir, err = os.UserCacheDir(); err != nil {
		return "", err
	}
	dir = filepath.Join(dir, "msgraph")
	if err = os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}
	return filepath.Join(dir, url.QueryEscape(fn)), nil
}
