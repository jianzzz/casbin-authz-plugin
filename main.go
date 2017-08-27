// Docker RBAC & ABAC Authorization Plugin based on Casbin.
// Allows only authorized Docker operations based on access control policy.
// AUTHOR: Yang Luo <hsluoyz@gmail.com>
// Powered by Casbin: https://github.com/casbin/casbin

package main

import (
	"flag"
	"github.com/docker/go-plugins-helpers/authorization"
	"log"
        "github.com/casbin/casbin/config"
)

var (
	casbinConfig = flag.String("config", "/usr/lib/docker/casbin.conf", "Specifies the Casbin configuration file")

	TOKEN = ""
)

func init() {
	cfg, err := config.NewConfig(*casbinConfig)
        if err != nil {
                panic(err)
        }
	TOKEN = cfg.String("default::token")
}

func main() {
	// Parse command line options.
	flag.Parse()
	log.Println("Casbin config:", *casbinConfig)

	// Create Casbin authorization plugin
	plugin, err := newPlugin(*casbinConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Start service handler on the local sock
	handler := authorization.NewHandler(plugin)
	
	cfg, err := config.NewConfig(*casbinConfig)
        if err != nil {
                panic(err)
        }

        ip := cfg.String("default::app_ip")
	port := cfg.String("default::app_port")

	if err := handler.ServeTCP("casbin-authz-plugin",ip+":"+port,"",nil); err != nil {
		log.Fatal(err)
	}
}
