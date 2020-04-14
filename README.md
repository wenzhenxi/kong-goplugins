# kong-goplugin

Based on the go expansion capability provided by kong official, provide basic projects and development environment

The initialization and update of the integrated kong database will be executed automatically before starting, without manual execution

## Reference

https://github.com/Kong/kong

https://github.com/Kong/go-pdk

https://github.com/Kong/go-pluginserver

https://github.com/Kong/go-plugins



## Depend

- Docker > 17.03
- Docker-composer > 0.21


## RUN

```
make build
make start
```

## Update Config

```editorconfig

plugins = bundled,go-log,go-hello,go-exit,go-token              # Comma-separated list of plugins this node

go_plugins_dir = /usr/local/share/lua/5.1/kong/plugins/go-so            # Directory for installing Kong plugins


# Nginx Worker Open multiple will have bugs, Waiting for official repair
nginx_worker_processes = 1   # Determines the number of worker processes

```


## ADD Plugin

```
vim Makefile

...
so-build: test.so

...

test.so:
	go build -o go-so/test.so -buildmode=plugin ./app/test.go

...

```


## Note

The official does not support the EXIT() method, first use the local project and sunmi-OS / go-pdk

```
github.com/Kong/go-pluginserver -> local

github.com/Kong/go-pdk -> github.com/sunmi-OS/go-pdk

```
