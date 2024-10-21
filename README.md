# openapphub-go

> 项目基于大佬的 [singo](https://github.com/gourouting/singo) 二开修改, 调整了项目结构，加入了docker、makefile、AUTH_MODE（jwt or session）等内容。从头学习的话推荐跟着大佬的项目进行：

- [让我们写个G站吧！Golang全栈编程实况](https://space.bilibili.com/10/channel/detail?cid=78794)

- [仿B站的G站](https://github.com/Gourouting/giligili)

- [Singo框架为移动端提供Token登录的案例](https://github.com/bydmm/singo-token-exmaple)

-----

openapphub: Simple Single Golang Web Service

openapphub: 用最简单的架构，实现够用的框架，服务海量用户

https://github.com/Gourouting/openapphub

## 更新日志

1. 已支持接口测试
2. 已经支持go1.20，请安装这个版本的golang使用本项目
3. 新增JWT认证支持
4. 新增Swagger文档
5. 新增logs处理

## 目的

本项目采用了一系列Golang中比较流行的组件，可以以本项目为基础快速搭建Restful Web API

## 特色

本项目已经整合了许多开发API所必要的组件：

1. [Gin](https://github.com/gin-gonic/gin): 轻量级Web框架，自称路由速度是golang最快的 
2. [GORM](https://gorm.io/index.html): ORM工具。本项目需要配合Mysql使用 
3. [Gin-Session](https://github.com/gin-contrib/sessions): Gin框架提供的Session操作工具
4. [Go-Redis](https://github.com/go-redis/redis): Golang Redis客户端
5. [godotenv](https://github.com/joho/godotenv): 开发环境下的环境变量工具，方便使用环境变量
6. [Gin-Cors](https://github.com/gin-contrib/cors): Gin框架提供的跨域中间件
7. [httpexpect](https://github.com/gavv/httpexpect): 接口测试工具
8. [JWT-Go](https://github.com/golang-jwt/jwt): JWT认证支持
9. [Swagger](https://github.com/swaggo/gin-swagger): API文档生成工具
10. 自行实现了国际化i18n的一些基本功能
11. 本项目支持基于cookie的session和JWT两种认证方式

本项目已经预先实现了一些常用的代码方便参考和复用:

1. 创建了用户模型
2. 实现了`/api/v1/user/register`用户注册接口
3. 实现了`/api/v1/user/login`用户登录接口
4. 实现了`/api/v1/user/me`用户资料接口(需要登录后获取session)
5. 实现了`/api/v1/user/logout`用户登出接口(需要登录后获取session)
6. 实现了`/api/v1/user/refresh`刷新JWT token接口

本项目已经预先创建了一系列文件夹划分出下列模块:

1. `cmd/api`: 主程序入口
2. `database/migrations`: 数据库建表相关
2. `internal/api`: MVC框架的controller，负责协调各部件完成任务
3. `internal/model`: 数据库模型和数据库操作相关的代码
4. `internal/service`: 负责处理比较复杂的业务，把业务代码模型化可以有效提高业务代码的质量
5. `internal/serializer`: 储存通用的json模型，把model得到的数据库模型转换成api需要的json对象
6. `pkg/cache`: redis缓存相关的代码
7. `internal/auth`: 权限控制相关的代码
8. `internal/util`: 一些通用的小工具
9. `internal/config`: 配置文件和配置加载相关的代码
10. `internal/middleware`: 中间件相关的代码
11. `test`: 测试用例

## 环境变量

项目在启动的时候依赖以下环境变量，但是在也可以在项目根目录创建.env文件设置环境变量便于使用(建议开发环境使用)

```shell
MYSQL_DSN="db_user:db_password@/db_name?charset=utf8&parseTime=True&loc=Local" # Mysql连接地址
REDIS_ADDR="127.0.0.1:6379" # Redis端口和地址
REDIS_PW="" # Redis连接密码
REDIS_DB="" # Redis库从0到10
SESSION_SECRET="setOnProducation" # Seesion密钥，必须设置而且不要泄露
GIN_MODE="debug"
LOG_LEVEL="debug"
AUTH_MODE="session" # 认证模式，可选值：session 或 jwt
JWT_SECRET="setOnProducation" # JWT密钥，使用JWT认证模式时必须设置
PORT="3000" # 服务端口号
```
## Godotenv

本项目使用[Godotenv](https://github.com/joho/godotenv)加载环境变量，在使用和部署项目的时候可以配置环境变量增加灵活性。

## Go Mod

本项目使用[Go Mod](https://github.com/golang/go/wiki/Modules)管理依赖。

```shell
go mod download
```

## 运行

```shell
go run cmd/api/main.go
// or user make
make install
make dev-setup
make docker-up
make run
```
项目运行后启动在3000端口（可以通过PORT环境变量修改）

## 编译

```shell
go build -o openapphub cmd/api/main.go
```

## 接口测试

本项目使用`httpexpect`进行接口测试，测试文件位于`test`目录下。运行测试：

```
go test -v ./test
```

## Swagger文档

本项目使用Swagger自动生成API文档。

1. 安装swag
   ```
   go install github.com/swaggo/swag/cmd/swag@latest
   ```

2. 生成文档
   ```
   swag init -g cmd/api/main.go -o docs
   ```

3. 访问文档：启动服务后，访问 `http://localhost:3000/swagger/index.html`

## Makefile

项目根目录下的Makefile文件包含了常用的操作命令，使用make命令即可执行。

## Docker支持

本项目支持使用Docker进行部署，根目录下的Dockerfile文件已经配置好了相应的构建流程。

## 贡献

如果你有好的意见或建议，欢迎给我们提 issue 或 pull request。

## 版权

Copyright (c) 2024 Gourouting

Licensed under the [MIT license](LICENSE)
