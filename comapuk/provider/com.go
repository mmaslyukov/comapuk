package provider

import (
	// "errors"
	// "fmt"
	"log"
	"strconv"
	"strings"

	"github.com/comapuk/request"
	"github.com/comapuk/util"
)

type Com struct {
	_client *request.Client
	_cfg    Config
}

func NewCom(cfg Config) (c *Com) {
	client, err := request.NewClient()
	if err != nil {
		log.Fatalf("Can't create a New Client -> %#v", err)
	}
	c = &Com{_client: client, _cfg: cfg}
	return c
}

func (z *Com) Read() (val float64, err error) {
	parser, err := z._client.Get(z._cfg[PageCalc])
	if err != nil {
		return val, err
	}
	data := request.Values{
		ComEmail:    []string{z._cfg[ComEmail]},
		ComPassword: []string{z._cfg[ComPassword]},
	}
	parser, err = z._client.Post(z._cfg[PageLogin], data)
	if err != nil {
		return val, err
	}
	parsed, err := parser.Data("/html/body/div[2]/div[2]/div/div[2]/div[2]/div[1]/p[2]")
	if err != nil {
		return val, err
	}
	
	split := strings.Split(parsed, " ")
	valStr := strings.Replace(split[0], ",", ".", 1)
	return strconv.ParseFloat(valStr, 64)
}

func DoCom(cfg util.Config) {
	com := NewCom(Config{
		ComEmail:    cfg.Login,
		ComPassword: cfg.Password,
		PageLogin:   "https://komunalka.ua/post/cabinet/login/",
	})
	if val, err := com.Read(); err == nil {
		log.Printf("City Utility Read value: %v", val)
	} else {
		log.Fatalln(err)
	}
}
