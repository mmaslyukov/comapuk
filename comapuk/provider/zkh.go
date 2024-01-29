package provider

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/comapuk/request"
	"github.com/comapuk/util"
)

type Zkh struct {
	_client *request.Client
	_cache  map[string]string
	_cfg    Config
	// _form   int
}

func NewZkh(cfg Config) (z *Zkh) {
	client, err := request.NewClient(
		cfg[ZkhCertPath],
	)

	if err != nil {
		log.Fatalf("Can't create a New Client -> %#v", err)
	}
	z = &Zkh{_client: client, _cache: make(util.Cache), _cfg: cfg}
	z._cache[ZkhViewState] = ""
	z._cache[ZkhRenderer] = "paidForm"
	z._cache[ZkhSource] = "calculationsHeadForm:j_idt46"
	z._cache[ZkhExecute] = z._cache[ZkhSource]
	// z._cache[FormId] = ""//string(z._form)
	return z
}

func (z *Zkh) Login() error {
	parser, err := z._client.Get(z._cfg[PageLogin])
	if err != nil {
		return err
	}
	err = z.updateViewState(parser)
	if err != nil {
		return err
	}
	data := request.Values{
		ZkhViewState:        []string{z._cache[ZkhViewState]},
		"mf":                []string{"mf"},
		"mf:" + ZkhLogin:    []string{z._cfg[ZkhLogin]},
		"mf:" + ZkhPassword: []string{z._cfg[ZkhPassword]},
		"mf:btnLogin":       []string{""},
	}
	parser, err = z._client.Post(z._cfg[PageLogin], data)
	if err != nil {
		return err
	}
	err = z.updateViewState(parser)
	if err != nil {
		return err
	}
	_, err = parser.Attr("//*[@id=\"j_idt18:btnExit\"]", "id")
	return err
}

// Possible values of 'id' are 0, 1
func (z *Zkh) Read(id uint) (val float64, err error) {
	if id > 1 {
		return val, errors.New(
			fmt.Sprintf("Form id(%v) is wrong, expect 0 or 1", id))
	}
	parser, err := z._client.Get(z._cfg[PageCalc])
	if err != nil {
		return val, err
	}
	err = z.updateViewState(parser)
	if err != nil {
		return val, err
	}
	data := request.Values{
		ZkhViewState: []string{z._cache[ZkhViewState]},
		ZkhSource:    []string{z._cache[ZkhSource]},
		ZkhExecute:   []string{z._cache[ZkhExecute]},
		ZkhRenderer:  []string{z._cache[ZkhRenderer]},
		ZkhFormId:    []string{fmt.Sprint(id)},
	}
	parser, err = z._client.Post(z._cfg[PageCalc], data)
	if err != nil {
		return val, err
	}
	valStr, err := parser.Attr("//*[@id=\"paidForm:amount\"]", "value")
	if err != nil {
		return val, err
	}
	return strconv.ParseFloat(valStr, 64)
}

func (z *Zkh) updateViewState(parser *request.Parser) (err error) {
	value, err := parser.Attr(
		"//*[@id=\"j_id1:"+ZkhViewState+":0\"]",
		"value")
	if err != nil {
		return err
	}
	// log.Printf("%v: %v", ViewState, value)
	z._cache[ZkhViewState] = value
	return err
}

func DoZkh(cfg util.Config) {
	var err error = nil
	zkh := NewZkh(Config{
		ZkhLogin:    cfg.Login,
		ZkhPassword: cfg.Password,
		ZkhCertPath: cfg.Certificate,
		PageLogin:   "https://www.zkh-plus.com.ua/ZkhPlusWebClient/pages/login.jsf",
		PageCalc:    "https://www.zkh-plus.com.ua/ZkhPlusWebClient/pages/user/pageCalculations.jsf",
	})
	err = zkh.Login()
	if err != nil {
		log.Fatalln(err)
	}
	if val, err := zkh.Read(0); err == nil {
		log.Printf("Zkh Utility Read value: %v", val)
	} else {
		log.Fatalln(err)
	}
	if val, err := zkh.Read(1); err == nil {
		log.Printf("Zkh Heating Read value: %v", val)
	} else {
		log.Fatalln(err)
	}
}
