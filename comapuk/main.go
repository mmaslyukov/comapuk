package main

import (
	"log"

	"github.com/BurntSushi/toml"
	"github.com/comapuk/provider"
	"github.com/comapuk/util"
)



// type Config struct {
// 	Login string `toml:"login"`
// 	Password string `toml:"password"`
// 	// Certificate string
// }

type Provider struct {
	// config map[string]Config `toml:"config"`
	// config Config `toml:"config"`
}


// type Configs map[string]Provider
type Configs map[string]util.Config



func main() {
	
	var cfg Configs
	_, err := toml.DecodeFile("config.toml", &cfg); 
	if err != nil {
		log.Fatalln(err)
	}
	if len(cfg["zkh-plus"].Login) == 0 {
		_, err = toml.DecodeFile("secret/config.toml", &cfg); 
		if err != nil {
			log.Fatalln(err)
		}
	} 

	provider.DoZkh(cfg["zkh-plus"])
	provider.DoCom(cfg["komunalka"])
	provider.DoInt(cfg["telegroup"])
}
 
