Build
```
$ go get github.com/jianzzz/casbin-authz-plugin
$ cd $GOPATH/src/github.com/jianzzz/casbin-authz-plugin
$ make
$ make install
```

Run the plugin as a systemd service
```
$ systemctl daemon-reload
$ systemctl enable casbin-authz-plugin
$ systemctl start casbin-authz-plugin
```

See whether the plugin starts correctly:
```
$ journalctl -xe -u casbin-authz-plugin -f
```

Enable the authorization plugin on docker engine
Step-1: Add authorization plugin to the docker engine configuration
```
--authorization-plugin casbin-authz-plugin
```
Step-2: Restart docker engine
```
$ systemctl daemon-reload
$ systemctl restart docker
```

My expand:
1. Use tcp socket but not unix socket，so multi docker daemons can map at one plugin.
2. Service install will not install systemd/casbin-authz-plugin.socket, cause we don't use unix socket.
3. Add token auth. You should add token in request header when you send api requests matching the policy.
4. See systemd/casbin-authz-plugin.service, when you restart docker service, casbin-authz-plugin service will be restarted.
5. When you restart casbin-authz-plugin service, it will automatically write /etc/docker/plugins/casbin-authz-plugin.spec.
   Change the listening port in casbin.conf and run "make && make install", then "systemctl daemon-reload && systemctl restart casbin-authz-plugin".


