# Go API Template

一个基于 Go 的 API 模板，集成了常用依赖和模块化结构，适合快速开发微服务和 RESTful API。

## 技术栈

- [Fiber](https://github.com/gofiber/fiber) v2.52.9 — 轻量且高性能的 Go Web 框架
- [JWT](https://github.com/golang-jwt/jwt) v4.5.2 — JSON Web Token 认证
- [go-redis](https://github.com/redis/go-redis) v9.12.0 — Redis 客户端
- [zerolog](https://github.com/rs/zerolog) v1.34.0 — 结构化日志记录
- [viper](https://github.com/spf13/viper) v1.20.1 — 配置管理
- [GORM](https://gorm.io) v1.30.1 + MySQL 驱动 v1.6.0 — ORM 数据库操作

## 项目结构


- go mod init xxx
- go mod tidy
- air