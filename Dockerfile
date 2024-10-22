# 构建：使用golang:1.16版本
FROM acr.cqrb.cn/develop/golang:1.22.6 as build
ARG CI_PROJECT_NAME
ARG GIT_TOKEN
ARG ENVER


# 容器环境变量添加，会覆盖默认的变量值
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct
ENV  GOPRIVATE="gitlab.cqrb.cn"
RUN git config --global url."http://oauth2:${GIT_TOKEN}@gitlab.cqrb.cn/".insteadOf "http://gitlab.cqrb.cn/"
# 设置工作区
WORKDIR /go/src/${CI_PROJECT_NAME}

# 把全部文件添加到/go/release目录
ADD . .
RUN rm -f go.mod go.sum
RUN go mod init
#RUN if [ "${ENVER}" = "develop" ];then\
#    echo 'require gitlab.cqrb.cn/shangyou_mic/shangyou-api-pb future_test'>>go.mod;\
#    elif [ "${ENVER}" = "testing" ];then\
#    echo 'require gitlab.cqrb.cn/shangyou_mic/shangyou-api-pb ceshi_test'>>go.mod;\
#    fi
RUN go mod tidy
# 编译：把cmd/main.go编译成可执行的二进制文件，命名为app
RUN cd ./cmd/${CI_PROJECT_NAME}&& GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-s -w" -installsuffix cgo -o ${CI_PROJECT_NAME} ./...
# 运行：使用alpine作为基础镜像
FROM acr.cqrb.cn/develop/shangyou_go_image
ARG CI_PROJECT_NAME
# 在build阶段复制时区到
ENV ZONEINFO=/app/zoneinfo.zip
WORKDIR /app
COPY --from=build /usr/local/go/lib/time/zoneinfo.zip /app
# 在build阶段复制可执行的go二进制文件app
COPY --from=build /go/src/${CI_PROJECT_NAME}/cmd/${CI_PROJECT_NAME} /app
# 在build阶段复制配置文件
#COPY --from=build /go/src/${CI_PROJECT_NAME}/conf ./conf
RUN addgroup -S nonroot && adduser -u 65530 -S nonroot -G nonroot
RUN chown nonroot:nonroot /app
USER 65530
# 启动服务
CMD ["/bin/bash", "-c","./${CI_PROJECT_NAME}"]
