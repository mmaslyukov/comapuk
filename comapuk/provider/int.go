package provider

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/comapuk/request"
	"github.com/comapuk/util"
)

type Int struct {
	_client *request.Client
	_cfg    Config
}

type IntApi struct {
	_int *Int
	_parser *request.Parser
}

type IntPageAccount struct {
	_parser *request.Parser
}

type IntPageSubscription struct {
	_parser *request.Parser
}

func NewInt(cfg Config) (c *Int) {
	client, err := request.NewClient()
	if err != nil {
		log.Fatalf("Can't create a New Client -> %#v", err)
	}
	c = &Int{_client: client, _cfg: cfg}
	return c
}

//Balance /html/body/div[1]/div/div[2]/table[1]/tbody/tr/td[2]
func (i *IntApi) Account() (acc *IntPageAccount, err error) {
	return &IntPageAccount{_parser: i._parser}, nil
}


func (i *IntApi) EstimatedDays() (val float64, err error) {
	account, err := i.Account()
	if err != nil {
		return -1, err
	}
	balance, err := account.Balance()
	if err != nil {
		return -1, err
	}
	subscription, err := i.Subscription()
	if err != nil {
		return -1, err
	}
	days, err := subscription.Days()
	if err != nil {
		return -1, err
	}
	cost, err := subscription.Cost()
	if err != nil {
		return -1, err
	}
	// log.Printf("Balance: %v, Costs:%v, Days:%v\n", balance, cost, days)
	val = days + (balance / cost) * 30
	return val, err
}



func (i *IntApi) Subscription() (sub *IntPageSubscription, err error) {
	parser, err := i._int._client.Get(i._int._cfg[PageSubscription])
	if err != nil {
		return nil, err
	}
	return &IntPageSubscription{_parser: parser}, nil
}

func (a *IntPageAccount) Balance() (val float64, err error) {
	valStr, err := a._parser.Data("/html/body/div[1]/div/div[2]/table[1]/tbody/tr/td[2]")
	if err != nil {
		return -1, err
	}
	return strconv.ParseFloat(valStr, 64)
}

func (s *IntPageSubscription) Cost() (val float64, err error) {
	valStr, err := s._parser.Data("/html/body/div[1]/div/table/tbody/tr[2]/td/table/tbody/tr/td[3]/table/tbody/tr[2]/td[2]")
	if err != nil {
		return -1, err
	}
	return strconv.ParseFloat(valStr, 64)
	
}

func (s *IntPageSubscription) Days() (val float64, err error) {
	valStr, err := s._parser.Data("/html/body/div[1]/div/table/tbody/tr[2]/td/table/tbody/tr/td[3]/table/tbody/tr[3]/td[2]")
	if err != nil {
		return -1, err
	}
	return strconv.ParseFloat(valStr, 64)
}


func (i *Int) Login() (intApi *IntApi, err error) {
	data := request.Values{
		IntLogin:    []string{i._cfg[IntLogin]},
		IntPassword: []string{i._cfg[IntPassword]},
	}
	parser, err := i._client.Post(i._cfg[PageAccount], data)
	if err != nil {
		return nil, err
	}
	if parser == nil {
		return nil, errors.New(fmt.Sprintf("Parser is nil"))
	}
	return &IntApi{_parser: parser, _int: i}, nil
}

func DoInt(cfg util.Config) {
	int := NewInt(Config{
		IntLogin:         cfg.Login,
		IntPassword:      cfg.Password,
		PageAccount:      "http://my.telegroup.ua/my_net_account.php",
		PageSubscription: "http://my.telegroup.ua/my_net_subscriptions.php",
	})
	h, err := int.Login()
	if err != nil {
		log.Fatalln(err)
	}
	_ = h
	if val, err := h.EstimatedDays(); err == nil {
		log.Printf("Days left: %v", val)
	} else {
		log.Fatalln(err)
	}
}
