FROM sunmi-docker-images-registry.cn-hangzhou.cr.aliyuncs.com/public/golang:1.14.1 As go-pluginserver

ENV GOPROXY https://goproxy.cn
ENV GO111MODULE on

WORKDIR /project

# The official does not support the EXIT method, first use the local project and sunmi-OS / go-pdk
# RUN git clone https://github.com/Kong/go-pluginserver.git
COPY go-pluginserver ./go-pluginserver
RUN cd go-pluginserver && make

FROM sunmi-docker-images-registry.cn-hangzhou.cr.aliyuncs.com/public/golang:1.14.1 As go-plugins

ENV GOPROXY https://goproxy.cn
ENV GO111MODULE on

WORKDIR /go/cache
ADD go.mod .
ADD go.sum .
RUN go mod download

WORKDIR /project
ADD . .
RUN mkdir go-so
RUN make so-build

FROM sunmi-docker-images-registry.cn-hangzhou.cr.aliyuncs.com/public/kong:2.0.2-centos

USER root

# Configure time zone
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

COPY --from=go-plugins /project/go-so /usr/local/share/lua/5.1/kong/plugins/go-so
COPY --from=go-pluginserver /project/go-pluginserver/go-pluginserver /usr/local/bin/

COPY ./lua-plugins/ /usr/local/share/lua/5.1/kong/plugins/
COPY ./pkg/lualib/ /usr/local/openresty/lualib/

USER kong

COPY ./config /etc/kong/

COPY docker-entrypoint.sh /docker-entrypoint.sh

ENTRYPOINT ["/docker-entrypoint.sh"]

EXPOSE 8000 8443 8001 8444

STOPSIGNAL SIGQUIT

CMD ["kong", "docker-start","dbup"]
