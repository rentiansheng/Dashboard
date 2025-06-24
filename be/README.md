# Dashboard Backend

一个基于 Go 语言开发的数据仪表板后端服务，提供灵活的数据源管理和图表查询功能。

## 🚀 功能特性

- **多数据源支持**: 支持 MySQL、Elasticsearch 等多种数据源
- **灵活查询**: 支持复杂的聚合查询、过滤条件和时间范围查询
- **图表生成**: 支持柱状图、折线图、饼图等多种图表类型
- **字段枚举**: 支持动态字段枚举和API枚举
- **分组管理**: 支持数据分组和分组键管理
- **RESTful API**: 提供完整的 RESTful API 接口
- **配置管理**: 支持 YAML 配置文件管理
- **日志记录**: 完整的请求日志和错误日志记录

## 📋 系统要求

- Go 1.24.1 或更高版本
- MySQL 5.7+ 或 MySQL 8.0+
- Elasticsearch 7.x+ (可选)

## 🛠️ 安装和配置

### 1. 克隆项目

```bash
git clone <repository-url>
cd dashboard/be
```

### 2. 安装依赖

```bash
go mod download
```

### 3. 配置数据库

复制配置文件模板：

```bash
cp config/config_example.yaml config/config.yaml
```

编辑 `config/config.yaml` 文件，配置数据库连接信息：

```yaml
# Server configuration
log_level: "info"

server:
  port: "8080"
  read_timeout: "100s"
  write_timeout: "100s"
  shutdown_timeout: "15s"
  log_level: "info"

engine:
  engine_type: "mysql"
  is_all: true
  default_engine_config:
    dsn: "user:password@tcp(127.0.0.1:3306)/dashboard?charset=utf8mb4&parseTime=True&loc=Local"

es:
  addrs:
    - "http://127.0.0.1:9200"
  username: ""
  password: ""
```

### 4. 初始化数据库

运行数据库初始化脚本：

```bash
# 创建数据库表
mysql -u your_username -p your_database < app/metrics/repository/data_source/impl/mysql/model/create.sql
mysql -u your_username -p your_database < app/metrics/repository/group_key/impl/mysql/model/create.sql
```

### 5. 运行服务

```bash
go run cmd/main.go
```

服务将在 `http://localhost:8080` 启动。

 
```
POST /api/v1/data-sources/enum
Content-Type: application/json

{
  "department_id": 1,
  "data_source_id": 1,
  "date": {
    "start": "2024-01-01",
    "end": "2024-01-31"
  },
  "field_name": "status"
}
```

 

## 🏗️ 项目结构

```
be/
├── app/                    # 应用层
│   └── metrics/           # 指标模块
│       ├── datasource/    # 数据源处理
│       ├── define/        # 数据定义
│       ├── repository/    # 数据访问层
│       ├── schema/        # API 模式层
│       └── service/       # 业务逻辑层
├── cmd/                   # 命令行入口
├── config/                # 配置文件
├── middleware/            # 中间件
├── pkg/                   # 公共包
└── plugin/                # 插件
```

## 🔧 开发指南

### 添加新的数据源类型

1. 在 `app/metrics/datasource/` 下创建新的数据源处理器
2. 实现相应的查询和聚合逻辑
3. 在 `app/metrics/repository/` 下添加数据访问层
4. 更新配置和路由

### 添加新的图表类型

1. 在 `app/metrics/define/` 中定义新的图表类型
2. 在 `app/metrics/datasource/format/` 中实现格式化逻辑
3. 更新查询处理逻辑

### 运行测试

```bash
go test ./...
```

## 📝 日志

服务会记录以下类型的日志：

- **请求日志**: 记录所有 HTTP 请求的详细信息
- **错误日志**: 记录系统错误和异常
- **业务日志**: 记录重要的业务操作

日志级别可通过配置文件中的 `log_level` 参数调整。

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

[许可证信息]

## 📞 联系方式

如有问题，请联系开发团队。 