// Docker RBAC & ABAC Authorization Plugin based on Casbin.
// Allows only authorized Docker operations based on access control policy.
// AUTHOR: Yang Luo <hsluoyz@gmail.com>
// Powered by Casbin: https://github.com/casbin/casbin

package main

import (
	"github.com/casbin/casbin"
	"github.com/docker/go-plugins-helpers/authorization"
	"log"
	"net/url"
	"github.com/casbin/casbin/config"
)

// CasbinAuthZPlugin is the Casbin Authorization Plugin
type CasbinAuthZPlugin struct {
	// Casbin enforcer
	enforcer *casbin.Enforcer
}

// newPlugin creates a new casbin authorization plugin
func newPlugin(casbinConfig string) (*CasbinAuthZPlugin, error) {
	plugin := &CasbinAuthZPlugin{}

	//modify by jianzhuang at 2017/08/19
//	plugin.enforcer = casbin.NewEnforcer(casbinConfig)
	cfg, err := config.NewConfig(casbinConfig)
	if err != nil {
		panic(err)
	}
	modelPath := cfg.String("default::model_path")
	log.Println("Casbin model path: ",modelPath)
	fileOrDB := cfg.String("default::policy_backend")
	if fileOrDB == "file"{
	        policyPath := cfg.String("file::policy_path")
        	log.Println("Casbin policy path: ",policyPath)
	        plugin.enforcer = casbin.NewEnforcer(modelPath,policyPath)
	}else{
		//todo
	}
	//end modify	

	return plugin, nil
}

// AuthZReq authorizes the docker client command.
// The command is allowed only if it matches a Casbin policy rule.
// Otherwise, the request is denied!
func (plugin *CasbinAuthZPlugin) AuthZReq(req authorization.Request) authorization.Response {
	// Parse request and the request body
	reqURI, _ := url.QueryUnescape(req.RequestURI)
	reqURL, _ := url.ParseRequestURI(reqURI)
	obj := reqURL.String()
	act := req.RequestMethod

	// Modify by jianzhuang
//	if plugin.enforcer.Enforce(obj, act) {
//		log.Println("obj:", obj, ", act:", act, "res: allowed")
//		return authorization.Response{Allow: true}
//	}
//
//	log.Println("obj:", obj, ", act:", act, "res: denied")
//	return authorization.Response{Allow: false, Msg: "Access denied by casbin plugin"}

	log.Println("user: ",req.User)
	log.Println("header: ",req.RequestHeaders)
//	log.Println("body: ",string(req.RequestBody))
        if plugin.enforcer.Enforce(obj, act) {
		if TOKEN == "" {
			log.Println("obj:", obj, ", act:", act, "res: denied")
		        return authorization.Response{Allow: false, Msg: "Access denied by casbin plugin, token is empty"}
		}
		if req.RequestHeaders["Token"] == TOKEN {
	                log.Println("obj:", obj, ", act:", act, "res: allowed")
                	return authorization.Response{Allow: true}
		} else {
			log.Println("obj:", obj, ", act:", act, "res: denied")
		        return authorization.Response{Allow: false, Msg: "Access denied by casbin plugin, you need to add token in request header"}
		}
        } else {
		log.Println("obj:", obj, ", act:", act, "res: allowed")
	        return authorization.Response{Allow: true}
	}
	// End mofify
}

// AuthZRes authorizes the docker client response.
// All responses are allowed by default.
func (plugin *CasbinAuthZPlugin) AuthZRes(req authorization.Request) authorization.Response {
	// Allowed by default.
	return authorization.Response{Allow: true}
}
