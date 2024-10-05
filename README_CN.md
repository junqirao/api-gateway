# API Gateway

[中文文档](/README_CN.md) | [English Doc](/README.md)

支持配置热更新的可编程API网关

-------

* [API Gateway](#api-gateway)
    * [总览](#总览)
    * [功能模块](#功能模块)
        * [支持配置热更](#支持配置热更)
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
                * [jwt](#jwt)
                * [ipgeo](#ipgeo)
        * [Metric](#metric)
            * [基础指标](#基础指标)
            * [额外指标](#额外指标)
                * [释义](#释义)
                * [示例](#示例-1)
        * [流量复制](#流量复制)
            * [如何工作](#如何工作)
            * [配置](#配置)
                * [通用](#通用)
                * [服务端](#服务端)
                * [客户端](#客户端)
    * [配置文件示例](#配置文件示例)
    * [依赖项目](#依赖项目)

## 总览

![overview](/doc/overview.png)

## 功能模块

* 负载均衡
    * 轮询 (Round-Robin)
    * 随机 (Random)
    * 加权轮询 (Weighted Round-Robin)
    * 加权随机 (Weighted Random)
    * 最少负载 (Less Load)
    * 客户端地址哈希 (Hash)
* 限流器
* 断路器
* 可插拔脚本
* 请求重试
* 配置热更新
* Prometheus指标接入
* 流量复制 (镜像)

### 支持配置热更

* 负载均衡器
* 限流器
* 断路器
* 可插拔脚本程序

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
| request.ClientIP   | String   | 请求客户端IP地址   | ```request.ClientIP```                            |
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

| Name          | Type     | Description     | Usage                                         |
|---------------|----------|-----------------|-----------------------------------------------|
| logger.Info   | Function | 记录info级别日志      | ```logger.Info(value any)```                  |
| logger.Warn   | Function | 记录warn级别日志      | ```logger.Warn(value any)```                  |
| logger.Error  | Function | 记录error级别日志     | ```logger.Error(value any)```                 |
| logger.Infof  | Function | 记录格式化的info级别日志  | ```logger.Infof(format string, value any)```  |
| logger.Warnf  | Function | 记录格式化的warn级别日志  | ```logger.Warnf(format string, value any)```  |
| logger.Errorf | Function | 记录格式化的error级别日志 | ```logger.Errorf(format string, value any)``` |

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

##### jwt

| Name              | Type     | Description   | Usage                                               |
|-------------------|----------|---------------|-----------------------------------------------------|
| jwt               | Function | jwt对象         | ```jwt(header_key [optinal]string)```               |
| jwt().MustSuccess | Function | 当解析jwt失败时终止请求 | ```jwt(header_key [optinal]string).MustSuccess()``` |
| jwt().Claims      | Map      | jwt claims 对象 | ```jwt(header_key [optinal]string).Claims```        |

Example:

```
// header["Authentication"] = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJiYXoiOiJxdXgiLCJmb28iOiJiYXIifQ.UyL2QgWoIz5Y3IdCeBk98_2-W26RTc5SCi7k0PCXbQw"
// payload = {"baz": "qux", "foo": "bar"}, key = "test_key_123456"

jwt().MustSuccess();                                      // nil
logger.Infof("jwt claims %v", jwt().Claims);              // jwt claims map[baz:qux foo:bar]
logger.Infof("jwt claims-foo %v", jwt().Claims["foo"]);   // jwt claims-foo bar
```

Config:

```yaml
program:
  extra:
    jwt:
      # jwt key
      key: "test_key_123456"
      # 用于jwt header中的kid映射密钥，可选
      keys_mapping: {
        "key1": "123456",
        "key2": "abcdefg"
      }
```

##### ipgeo

> 提供ISO 3166国家代码和城市代码查询，基于本地mmdb数据库。

| Name          | Type     | Description       | Usage                 |
|---------------|----------|-------------------|-----------------------|
| ipgeo.Country | Function | 获取国家代码            | ```ipgeo.Country()``` |
| ipgeo.City    | Function | 获取城市代码（需要数据库支持）   | ```ipgeo.City()```    |
| ipgeo.CityEN  | Function | 获取英文城市代码（需要数据库支持） | ```ipgeo.CityEN()```  |

Example:

```
logger.Infof("country %s", ipgeo.Country());  // country CN
logger.Infof("city %s", ipgeo.City());        // city map[de:Hangzhou en:Hangzhou es:Hangzhou fr:Hangzhou ja:杭州市 pt-BR:Hangzhou ru:Ханчжоу zh-CN:杭州]
logger.Infof("city_en %s", ipgeo.CityEN());   // city_en Hangzhou
```

Config:

```yaml
program:
  extra:
    ipgeo:
      # 数据库文件下载 https://www.maxmind.com/
      database: "GeoLite2-City.mmdb"
```

### Metric

```
GET {entrance}/management/metrics
```

_注意：如果配置中 ```auth.secret``` 不为空，请求头需要添加 "Authorization" 字段。_

#### 基础指标

见 [prometheus/client_golang](https://github.com/prometheus/client_golang)

#### 额外指标

##### 释义

| 字段                             | 含义     | 标签含义                          |
|--------------------------------|--------|-------------------------------|
| gateway_http_request_total     | 请求总量   | -                             |
| gateway_http_service_time_cost | 请求时间成本 | service:"服务名" quantile:"百分位数" |
| gateway_http_status_total      | 状态总量   | status:"HTTP响应状态码"            |

##### 示例

```
# HELP gateway_http_request_total Request count.
# TYPE gateway_http_request_total counter
gateway_http_request_total 5236
# HELP gateway_http_service_time_cost Request time cost by service.
# TYPE gateway_http_service_time_cost summary
gateway_http_service_time_cost{service="test",quantile="0.5"} 59
gateway_http_service_time_cost{service="test",quantile="0.9"} 93
gateway_http_service_time_cost{service="test",quantile="0.99"} 107
gateway_http_service_time_cost_sum{service="test"} 300712
gateway_http_service_time_cost_count{service="test"} 5236
# HELP gateway_http_status_total Request status count.
# TYPE gateway_http_status_total counter
gateway_http_status_total{status="200"} 5234
gateway_http_status_total{status="502"} 2
```

### 流量复制

#### 如何工作

> 通过服务端代理注册客户端到注册中心并创建对应的反向代理处理器，在流量代理响应结束后异步推送到客户端注册的目标主机。

#### 配置

##### 通用

| 字段          | 描述   | 值                                      | 默认 |
|-------------|------|----------------------------------------|----|
| mirror.mode | 工作模式 | ```"server"``` ```"client"``` ```""``` | "" |

##### 服务端

| 字段                               | 描述                                  | 值      | 默认  |
|----------------------------------|-------------------------------------|--------|-----|
| mirror.server.ch_buffer_size     | 推送缓冲区大小，如果设置为 0 会阻塞直到推送完成，影响代理资源回收。 | int    | 100 |
| mirror.server.worker_count       | 代理工作线程数，当请求耗时高时提高工作线程数量会改善阻塞时长。     | int    | 10  |
| mirror.server.heartbeat_interval | 客户端心跳间隔秒数                           | int    | 30  |
| mirror.server.white_list         | 访问白名单，允许访问的域名或 IP 地址。为空时允许所有。       | array  | []  |
| auth.secret                      | 访问密钥                                | string | ""  |

##### 客户端

| 字段                           | 描述                       | 值      | 默认 |
|------------------------------|--------------------------|--------|----|
| mirror.client.server_address | 服务端地址，如 "some.host:8001" | string | "" |
| mirror.client.secret         | 访问密钥，与服务端 auth.secret 相同 | string | "" |
| mirror.client.filter         | 服务过滤，如果不设置所有请求都会被复制      | array  | [] |

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

# 鉴权相关
auth:
  secret: "P@sswOrd"

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

# 流量复制配置
mirror:
  # 工作模式 "","client","server"
  mode: "server"
  server:
    # 缓冲区大小, 默认 100
    ch_buffer_size: 100
    # 异步处理worker数量, 默认 10
    worker_count: 10
    # 客户端心跳间隔（秒）, 默认 30
    heartbeat_interval: 30
    # 白名单，允许访问的域名或ip
    white_list:
      - "127.0.0.1"
  client:
    server_address: "127.0.0.1:8001"
    secret: "P@sswOrd"
    # 需要复制流量的服务，不设置时复制全部
    filter:
      - "service_name_1"
      - "service_name_2"

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
