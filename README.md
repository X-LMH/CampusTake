# CampusTake Backend

CampusTake 后端服务，基于 Go、go-zero、MySQL、Redis 和 RabbitMQ。

## 目录说明

- `campusservice.go`: 服务启动入口
- `etc/config.yaml`: 本地开发配置
- `deploy/docker-config.yaml`: Docker 部署配置
- `sql/campus.sql`: 数据库初始化脚本
- `sql/performance_indexes.sql`: 性能索引脚本

## 本地运行

```bash
go mod download
go run campusservice.go -f etc/config.yaml
```

服务默认端口为 `8888`。

## Docker 镜像构建

在项目根目录执行：

```bash
docker build -t campustake-backend:latest .
```

## Docker 部署

### 1. 创建 Docker 网络

```bash
docker network create campustake-net
```

### 2. 启动 MySQL

```bash
docker run -d --name campustake-mysql --network campustake-net -p 3306:3306 -e MYSQL_ROOT_PASSWORD=123123 -e MYSQL_DATABASE=campus -v campustake-mysql-data:/var/lib/mysql mysql:8.4
```

初始化数据库：

PowerShell:

```powershell
Get-Content .\sql\campus.sql | docker exec -i campustake-mysql mysql -uroot -p123123 campus
Get-Content .\sql\performance_indexes.sql | docker exec -i campustake-mysql mysql -uroot -p123123 campus
```

Bash:

```bash
docker exec -i campustake-mysql mysql -uroot -p123123 campus < sql/campus.sql
docker exec -i campustake-mysql mysql -uroot -p123123 campus < sql/performance_indexes.sql
```

### 3. 启动 Redis

```bash
docker run -d --name campustake-redis --network campustake-net -p 6379:6379 redis:7-alpine
```

### 4. 启动 RabbitMQ

```bash
docker run -d --name campustake-rabbitmq --network campustake-net -p 5672:5672 -p 15672:15672 -e RABBITMQ_DEFAULT_USER=admin -e RABBITMQ_DEFAULT_PASS=123123 rabbitmq:3-management
```

RabbitMQ 管理后台地址：`http://localhost:15672`，账号 `admin`，密码 `123123`。

### 5. 启动后端服务

```bash
docker run -d --name campustake-backend --network campustake-net -p 8888:8888 -v campustake-upload:/app/upload campustake-backend:latest
```

访问地址：

```text
http://localhost:8888
```

## 使用自定义配置

如果需要修改数据库、Redis、RabbitMQ、JWT 或上传地址，可以复制一份配置文件：

PowerShell:

```powershell
Copy-Item .\deploy\docker-config.yaml .\deploy\prod-config.yaml
```

Bash:

```bash
cp deploy/docker-config.yaml deploy/prod-config.yaml
```

修改后通过挂载方式启动：

PowerShell:

```powershell
docker run -d --name campustake-backend --network campustake-net -p 8888:8888 -v campustake-upload:/app/upload -v "$PWD/deploy/prod-config.yaml:/app/deploy/prod-config.yaml:ro" campustake-backend:latest -f /app/deploy/prod-config.yaml
```

Bash:

```bash
docker run -d --name campustake-backend --network campustake-net -p 8888:8888 -v "$PWD/deploy/prod-config.yaml:/app/deploy/prod-config.yaml:ro" -v campustake-upload:/app/upload campustake-backend:latest -f /app/deploy/prod-config.yaml
```

生产环境建议至少修改：

- `JwtAuth.SecretKey`
- `MySQLConfig.Password`
- `RabbitMQConfig.Password`
- `UploadConfig.UrlPrefix`

## 查看日志

```bash
docker logs -f campustake-backend
```

## 停止和删除容器

```bash
docker stop campustake-backend campustake-mysql campustake-redis campustake-rabbitmq
docker rm campustake-backend campustake-mysql campustake-redis campustake-rabbitmq
```
