# API Gateway

[中文文档](/README_CN.md) | [English Doc](/README.md)

支持配置热更新的可编程API网关

-------

* [API Gateway](#api-gateway)
    * [总览](#总览)
    * [功能模块](#功能模块)
        * [支持配置热更](#支持配置热更)
        * [待完成](#待完成)
    * [快速开始](#快速开始)
        * [API文档](#api文档)
        * [服务配置](#服务配置)
            * [负载均衡](#负载均衡)
            * [限流器](#限流器)
            * [断路器](#断路器)
        * [可编程脚本](#可编程脚本)
            * [示例](#示例)
            * [支持的方法和值](#支持的方法和值)
                * [request](#request)
                * [response](#response)
                * [global](#global)
                * [ctx](#ctx)
                * [logger](#logger)
                * [upstream](#upstream)
                * [terminate_if](#terminate_if)
    * [配置示例](#配置示例)
    * [依赖项目](#依赖项目)

## 总览

![overview](/doc/overview.png)

## 功能模块

* 负载均衡
    * 加权轮询 (Weighted Round Robin)
    * 加权随机 (Weighted Random)
* 限流器
* 断路器
* 可插拔脚本
* 请求重试
* 配置热更新

### 支持配置热更

* 负载均衡器
* 限流器
* 断路器
* 可插拔脚本程序

### 待完成

* Prometheus指标接入

## 快速开始

### API文档

[api reference](/doc/api_reference_cn.md)

### 服务配置

更新服务配置 详见 [api reference](/doc/api_reference_cn.md)

可用的API:

```
PUT {entrance}/management/config/load_balance
PUT {entrance}/management/config/rate_limiter
PUT {entrance}/management/config/breaker
```

#### 负载均衡

请求参数:

```json
{
  "service_name": "...",
  "config": {
    "strategy": "weighted-round-robin"
  }
}
```

| 参数       | 描述     | 默认值         |
|----------|--------|-------------|
| strategy | 负载均衡策略 | round-robin |

可用的负载均衡策略:

```round-robin```
```random```
```weighted-round-robin```
```weighted-random```
```less-load```
```hash```

#### 限流器

请求参数:

```json
{
  "service_name": "...",
  "config": {
    "rate": 5,
    "peak": 10
  }
}
```

| 参数   | 描述   | 默认值 |
|------|------|-----|
| rate | 每秒限制 | -   |
| peak | 突发值  | -   |

#### 断路器

请求参数:

```json
{
  "service_name": "...",
  "config": {
    "max_failures": 5,
    "half_open_max_requests": 1,
    "open_timeout": "30s",
    "interval": "30s"
  }
}
```

| 参数                     | 描述                                                | 默认值 |
|------------------------|---------------------------------------------------|-----|
| max_failures           | 当错误达到阈值时，断路器状态从闭合变为半开                             | 5   |
| half_open_max_requests | 当断路器状态为半开时达到阈值后，断路器状态从半开变为完全断开                    | 1   |
| open_timeout           | 断路器断开状态持续时间，默认1分钟                                 | 1m  |
| interval               | 清除max_failures和half_open_max_requests计数器的间隔，默认1分钟 | 1m  |

### 可编程脚本

利用 [go-expr](https://github.com/expr-lang/expr) 作为基础执行引擎
额外支持了通过";"作为结尾的多行脚本。

#### 示例

```
logger.Info(ctx, "hello %s", "world");
response.Header.Set("X-Test-Header",upstream.HostName);
terminate_if(request.Header.Get("Authorization") == "","authorization is empty");
```

第一行: 输出 info 级别的日志

第二行: 将响应头 "X-Test-Header" 的值设置为变量 upstream.HostName

第三行: 如果请求头 "Authorization" 的值为空, 则终止请求

#### 支持的方法和值

##### request

request object

| Name               | Type     | Description | Usage                                             |
|--------------------|----------|-------------|---------------------------------------------------|
| request.RemoteAddr | String   | 请求客户端地址     | ```request.RemoteAddr```                          |
| request.Host       | String   | 请求主机        | ```request.Host```                                |
| request.URL        | String   | 请求url       | ```request.URL```                                 |
| request.Method     | String   | 请求方法        | ```request.Method```                              |
| request.Header     | Object   | 请求头         | -                                                 |
| request.Header.Add | Function | 添加请求头       | ```request.Header.Add(key string,value string)``` |
| request.Header.Set | Function | 设置请求头       | ```request.Header.Set(key string,value string)``` |
| request.Header.Get | Function | 获取请求头的值     | ```request.Header.Get(key string)```              |

##### response

response object

| Name                | Type     | Description | Usage                                              |
|---------------------|----------|-------------|----------------------------------------------------|
| response.Header.Add | Function | 添加响应头值      | ```response.Header.Add(key string,value string)``` |
| response.Header.Set | Function | 设置响应头值      | ```response.Header.Set(key string,value string)``` |
| response.Header.Get | Function | 获取响应头值      | ```response.Header.Get(key string)```              |

##### global

全局变量，这个值在所有注册到同一注册中心的网关服务中是共享的

| Name       | Type     | Description | Usage                                     |
|------------|----------|-------------|-------------------------------------------|
| global     | Map      | 全局变量map实例   | ```global.key / global["key"]```          |
| set_global | Function | 设置全局变量      | ```set_global(key string,value string)``` |

##### ctx

请求上下文，一般用作于函数参数

##### logger

logger object

| Name          | Type     | Description     | Usage                                                      |
|---------------|----------|-----------------|------------------------------------------------------------|
| logger.Info   | Function | 记录info级别日志      | ```logger.Info(ctx context, value any)```                  |
| logger.Warn   | Function | 记录warn级别日志      | ```logger.Warn(ctx context, value any)```                  |
| logger.Error  | Function | 记录error级别日志     | ```logger.Error(ctx context, value any)```                 |
| logger.Infof  | Function | 记录格式化的info级别日志  | ```logger.Infof(ctx context, format string, value any)```  |
| logger.Warnf  | Function | 记录格式化的warn级别日志  | ```logger.Warnf(ctx context, format string, value any)```  |
| logger.Errorf | Function | 记录格式化的error级别日志 | ```logger.Errorf(ctx context, format string, value any)``` |

##### upstream

| Name                 | Type    | Description | Usage                      |
|----------------------|---------|-------------|----------------------------|
| upstream.Id          | String  | 上游实例id      | ```upstream.Id```          |
| upstream.Host        | String  | 上游主机        | ```upstream.Host```        |
| upstream.HostName    | String  | 上游主机名       | ```upstream.HostName```    |
| upstream.Port        | Integer | 上游端口        | ```upstream.Port```        |
| upstream.ServiceName | String  | 上游服务名       | ```upstream.ServiceName``` |
| upstream.Meta        | Map     | 上游元数据       | ```upstream.Meta```        |

##### terminate_if

| Name         | Type     | Description | Usage                                                       |
|--------------|----------|-------------|-------------------------------------------------------------|
| terminate_if | Function | 条件终止        | ```terminate_if(condition bool, reason [optional]string)``` |

## 配置文件示例

```yaml
# goframe 服务器相关设置
server:
  address: ":8000"
  # 如果debug=true 同时会打开pprof，访问路径为：/debug/pprof
  # 对所有请求的响应头中添加debug信息
  debug: true

# goframe 日志相关设置
logger:
  level: "all"
  stdout: true

# 网关配置
gateway:
  # 代理url前缀
  # 默认 "/api/"
  prefix: "/api/"
  # 管理功能相关配置
  management:
    # 管理入口
    # "/management/*" -> "{entrance}/management/*"
    entrance: ""
    # 是否开启管理api
    enable: true
    # md5两次加密的密码
    # md5("password")->"5f4dcc3b5aa765d61d8327deb882cf99"->md5("5f4dcc3b5aa765d61d8327deb882cf99") = "696d29e0940a4957748fe3fc9efd22a3"
    # 如果这个字段不为空，在调用管理API时，需要在请求头""Authorization"中添加明文密码一次哈希的值
    password: "696d29e0940a4957748fe3fc9efd22a3"

# 注册中心配置
registry:
  # 注册中心类型
  type: "etcd"
  # 注册的服务名
  service_name: "api-gateway"
  # 数据库设置
  database:
    # 数据库入口
    endpoints:
      - "endpoint-1:2379"
      - "endpoint-2:2380"
      - "endpoint-3:2381"
    # 用户名密码相关
    username: "your username"
    password: "your password"
```

## 依赖项目

* [GoFrame](https://github.com/gogf/gf)
* [go-expr](https://github.com/expr-lang/expr)
* [simple-registry](https://github.com/junqirao/simple-registry)
