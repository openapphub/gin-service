# 变量
APP_NAME := openapphub
GO := go
GOFLAGS := -v
DOCKER := docker
DOCKER_COMPOSE := docker-compose

# 默认目标
.DEFAULT_GOAL := help

# 生成 Swagger 文档
.PHONY: swagger
swagger:
	swag init -g cmd/api/main.go -o docs --outputTypes go,json,yaml

# 帮助
.PHONY: help
help:
	@echo "可用命令："
	@echo "  install     - 安装依赖"
	@echo "  update-deps - 更新依赖"
	@echo "  run         - 本地运行应用"
	@echo "  build       - 构建应用"
	@echo "  test        - 运行测试"
	@echo "  clean       - 清理构建产物"
	@echo "  docker-up   - 启动 Docker 服务（MySQL 和 Redis）"
	@echo "  docker-down - 停止 Docker 服务"
	@echo "  docker-build- 构建 Docker 镜像"
	@echo "  docker-run  - 在 Docker 中运行应用"
	@echo "  migrate-up  - 运行数据库迁移"
	@echo "  migrate-down- 回滚数据库迁移"
	@echo "  lint        - 运行代码检查"
	@echo "  fmt         - 格式化代码"
	@echo "  swagger     - 生成 Swagger 文档"

# 安装依赖
.PHONY: install
install:
	$(GO) mod download

# 设置本地开发环境
.PHONY: dev-setup
dev-setup:
	cp .env.example .env
	@echo "请编辑 .env 文件，填入您的本地配置"

# 更新依赖
.PHONY: update-deps
update-deps:
	go get -u ./...
	go mod tidy

# 本地运行应用
.PHONY: run
run:
	$(GO) run $(GOFLAGS) ./cmd/api

# 构建应用
.PHONY: build
build:
	$(GO) build $(GOFLAGS) -o $(APP_NAME) ./cmd/api

# 运行测试
.PHONY: test
test:
	$(GO) test ./...

# 清理构建产物
.PHONY: clean
clean:
	rm -f $(APP_NAME)
	$(GO) clean

# Docker 命令
.PHONY: docker-up
docker-up:
	$(DOCKER_COMPOSE) --env-file .env up -d mysql redis

.PHONY: docker-down
docker-down:
	$(DOCKER_COMPOSE) --env-file .env down

.PHONY: docker-build
docker-build:
	$(DOCKER) build  -t $(APP_NAME) .

.PHONY: docker-run
docker-run:
	$(DOCKER) run --network host -p 3000:3000 --env-file .env --name $(APP_NAME)2 $(APP_NAME)

# 数据库迁移
.PHONY: migrate-up
migrate-up:
	@source .env && migrate -path database/migrations -database "mysql://$${MYSQL_USER}:$${MYSQL_PASSWORD}@tcp($${MYSQL_HOST}:$${MYSQL_PORT})/$${MYSQL_DATABASE}" up

.PHONY: migrate-down
migrate-down:
	@source .env && migrate -path database/migrations -database "mysql://$${MYSQL_USER}:$${MYSQL_PASSWORD}@tcp($${MYSQL_HOST}:$${MYSQL_PORT})/$${MYSQL_DATABASE}" down

# 代码质量
.PHONY: fmt
fmt:
	$(GO) fmt ./...

.PHONY: lint
lint:
	golangci-lint run
