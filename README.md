# API Gateway

[中文文档](/README_CN.md) | [English Doc](/README.md)

configuration hot updates supported programmable api gateway

-------

* [API Gateway](#api-gateway)
    * [Overview](#overview)
    * [Features](#features)
        * [Hot Updates Supported](#hot-updates-supported)
    * [Quick Start](#quick-start)
        * [API Reference](#api-reference)
        * [Config](#config)
            * [LoadBalance](#loadbalance)
            * [Rate Limiter](#rate-limiter)
            * [Circuit Breaker](#circuit-breaker)
        * [Program](#program)
            * [Example](#example)
            * [Supported Functions](#supported-functions)
                * [request](#request)
                * [response](#response)
                * [global](#global)
                * [ctx](#ctx)
                * [logger](#logger)
                * [upstream](#upstream)
                * [terminate_if](#terminate_if)
        * [Metric](#metric)
            * [Basic](#basic)
            * [Extra](#extra)
                * [Fields](#fields)
                * [Example](#example-1)
        * [Request Replicate](#request-replicate)
            * [How it works](#how-it-works)
            * [Config](#config-1)
                * [Common](#common)
                * [Server](#server)
                * [Client](#client)
    * [Config Example](#config-example)
    * [Dependencies](#dependencies)

## Overview

![overview](/doc/overview.png)

## Features

* Load Balance
    * Round-Robin
    * Random
    * Weighted Round-Robin
    * Weighted Random
    * Less Load
    * Hash
* Rate Limiter
* Circuit Breaker
* Programmable
* Retry
* Configuration Hot Updates
* Metrics for prometheus
* Request Replicate (Mirror)

### Hot Updates Supported

* Load Balance
* Rate Limiter
* Circuit Breaker
* Program

## Quick Start

### API Reference

[api reference](/doc/api_reference_en.md)

### Config

to update config, see [api reference](/doc/api_reference_en.md)

available apis:

```
PUT {entrance}/management/config/load_balance
PUT {entrance}/management/config/rate_limiter
PUT {entrance}/management/config/breaker
```

#### LoadBalance

Parameters:

```json
{
  "service_name": "...",
  "config": {
    "strategy": "weighted-round-robin"
  }
}
```

| Parameter | Description            | Default     |
|-----------|------------------------|-------------|
| strategy  | load balancer strategy | round-robin |

Supported Strategies:

```round-robin```
```random```
```weighted-round-robin```
```weighted-random```
```less-load```
```hash```

#### Rate Limiter

Parameters:

```json
{
  "service_name": "...",
  "config": {
    "rate": 5,
    "peak": 10
  }
}
```

| Parameter | Description     | Default |
|-----------|-----------------|---------|
| rate      | rate pre second | -       |
| peak      | burst           | -       |

#### Circuit Breaker

Parameters:

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

| Parameter              | Description                                                                                                   | Default |
|------------------------|---------------------------------------------------------------------------------------------------------------|---------|
| max_failures           | when reached max_failures breaker<br/> changes status from close to half_open                                 | 5       |
| half_open_max_requests | when reached half_open_max_requests <br/>at half_open state breaker changes state<br/> from half_open to open | 1       |
| open_timeout           | the period of the open state, default 1m                                                                      | 1m      |
| interval               | for the CircuitBreaker to clear the internal <br/>Counts, default 1m                                          | 1m      |

### Program

using [go-expr](https://github.com/expr-lang/expr) as basic execution engine

additionally support multiline statements, end line with ";"

#### Example

```
logger.Info(ctx, "hello %s", "world");
response.Header.Set("X-Test-Header",upstream.HostName);
terminate_if(request.Header.Get("Authorization") == "","authorization is empty");
```

line1. print info level log to server

line2. set variable upstream.HostName to response header "X-Test-Header"

line3. check request.Header "Authorization" is empty, if true, terminate request with error message "authorization is
empty"

#### Supported Functions

##### request

request object

| Name               | Type     | Description              | Usage                                             |
|--------------------|----------|--------------------------|---------------------------------------------------|
| request.RemoteAddr | String   | request remote address   | ```request.RemoteAddr```                          |
| request.Host       | String   | request host             | ```request.Host```                                |
| request.URL        | String   | request url              | ```request.URL```                                 |
| request.Method     | String   | request method           | ```request.Method```                              |
| request.Header     | Object   | request header           | -                                                 |
| request.Header.Add | Function | add request header value | ```request.Header.Add(key string,value string)``` |
| request.Header.Set | Function | set request header value | ```request.Header.Set(key string,value string)``` |
| request.Header.Get | Function | get request header value | ```request.Header.Get(key string)```              |

##### response

response object

| Name                | Type     | Description               | Usage                                              |
|---------------------|----------|---------------------------|----------------------------------------------------|
| response.Header.Add | Function | add response header value | ```response.Header.Add(key string,value string)``` |
| response.Header.Set | Function | set response header value | ```response.Header.Set(key string,value string)``` |
| response.Header.Get | Function | get response header value | ```response.Header.Get(key string)```              |

##### global

global variable, this value is shared across all api-gateway instances.

| Name       | Type     | Description         | Usage                                     |
|------------|----------|---------------------|-------------------------------------------|
| global     | Map      | global variable map | ```global.key / global["key"]```          |
| set_global | Function | set global variable | ```set_global(key string,value string)``` |

##### ctx

context from request, use as parameter of function

##### logger

logger object

| Name          | Type     | Description                   | Usage                                                      |
|---------------|----------|-------------------------------|------------------------------------------------------------|
| logger.Info   | Function | log with info level           | ```logger.Info(ctx context, value any)```                  |
| logger.Warn   | Function | log with warn level           | ```logger.Warn(ctx context, value any)```                  |
| logger.Error  | Function | log with error level          | ```logger.Error(ctx context, value any)```                 |
| logger.Infof  | Function | formated log with info level  | ```logger.Infof(ctx context, format string, value any)```  |
| logger.Warnf  | Function | formated log with warn level  | ```logger.Warnf(ctx context, format string, value any)```  |
| logger.Errorf | Function | formated log with error level | ```logger.Errorf(ctx context, format string, value any)``` |

##### upstream

| Name                 | Type    | Description           | Usage                      |
|----------------------|---------|-----------------------|----------------------------|
| upstream.Id          | String  | upstream instance id  | ```upstream.Id```          |
| upstream.Host        | String  | upstream host         | ```upstream.Host```        |
| upstream.HostName    | String  | upstream host name    | ```upstream.HostName```    |
| upstream.Port        | Integer | upstream port         | ```upstream.Port```        |
| upstream.ServiceName | String  | upstream service name | ```upstream.ServiceName``` |
| upstream.Meta        | Map     | upstream metadata     | ```upstream.Meta```        |

##### terminate_if

| Name         | Type     | Description                      | Usage                                                       |
|--------------|----------|----------------------------------|-------------------------------------------------------------|
| terminate_if | Function | terminate request with condition | ```terminate_if(condition bool, reason [optional]string)``` |

### Metric

```
GET {entrance}/management/metrics
```

_notice: if gateway.management.password not empty, request header "Authorization" is required._

#### Basic

for detail, see [prometheus/client_golang](https://github.com/prometheus/client_golang)

#### Extra

##### Fields

| Field                          | Description       | Label Description                          |
|--------------------------------|-------------------|--------------------------------------------|
| gateway_http_request_total     | total requests    | -                                          |
| gateway_http_service_time_cost | request time cost | service:"service name" quantile:"quantile" |
| gateway_http_status_total      | request status    | status:"HTTP status code"                  |

##### Example

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

### Request Replicate

#### How it works

> Register the client to the registry through the server-side proxy
> then server create a reverse proxy handler based on it. after the
> request proxy response is completed, the request will be replicated
> and push to the target server asynchronously.

#### Config

##### Common

| Field       | Description  | Value                                  | Default |
|-------------|--------------|----------------------------------------|---------|
| mirror.mode | running mode | ```"server"``` ```"client"``` ```""``` | ""      |

##### Server

| Field                            | Description                                                                                       | Value  | Default |
|----------------------------------|---------------------------------------------------------------------------------------------------|--------|---------|
| mirror.server.ch_buffer_size     | channel buffer size, if set to 0 push operation<br/> will be blocked till it is done.             | int    | 100     |
| mirror.server.worker_count       | reverse proxy worker count. when request costs <br/>time more worker may improve block situation. | int    | 10      |
| mirror.server.heartbeat_interval | client heartbeat interval seconds.                                                                | int    | 30      |
| mirror.server.white_list         | access white list, accept domain name or ip address.<br/> allow all when value not set.           | array  | []      |
| auth.secret                      | access secret.                                                                                    | string | ""      |

##### Client

| Field                        | Description                                               | Value  | Default |
|------------------------------|-----------------------------------------------------------|--------|---------|
| mirror.client.server_address | server address. e.g. "some.host:8001"                     | string | ""      |
| mirror.client.secret         | same as server side auth.secret                           | string | ""      |
| mirror.client.filter         | service filter, if not set all request will be replicated | array  | []      |

## Config Example

```yaml
# goframe server config
server:
  address: ":8000"
  # if debug mode enabled
  # pprof will also enable and all response proxied will add debug header
  debug: true

# goframe logger config
logger:
  level: "all"
  stdout: true

# gateway config
gateway:
  # proxy prefix
  # default "/api/"
  prefix: "/api/"
  # management config
  management:
    # management entrance
    # "/management/*" -> "{entrance}/management/*"
    entrance: ""
    # enable management api
    enable: true
    # md5 twice hashed password
    # md5("password")->"5f4dcc3b5aa765d61d8327deb882cf99"->md5("5f4dcc3b5aa765d61d8327deb882cf99") = "696d29e0940a4957748fe3fc9efd22a3"
    # if password not empty, calling management api, "Authorization" header is required, value is md5("password")
    password: "696d29e0940a4957748fe3fc9efd22a3"

# request mirror config
mirror:
  # working mode "","client","server"
  mode: "server"
  server:
    # channel buffer size, default 100
    ch_buffer_size: 100
    # async worker count, default 10
    worker_count: 10
    # client heartbeat interval, default 30
    heartbeat_interval: 30
    # access white list, accept domain name or ip address
    white_list:
      - "127.0.0.1"
  client:
    server_address: "127.0.0.1:8001"
    secret: "P@sswOrd"
    # service_name going to replicate to,
    # if not set will replicate to all
    filter:
      - "service_name_1"
      - "service_name_2"

# registry config
registry:
  # registry type
  type: "etcd"
  # register service_name
  service_name: "api-gateway"
  # registry database config
  database:
    # database endpoints
    endpoints:
      - "endpoint-1:2379"
      - "endpoint-2:2380"
      - "endpoint-3:2381"
    username: "your username"
    password: "your password"
```

## Dependencies

* [GoFrame](https://github.com/gogf/gf)
* [go-expr](https://github.com/expr-lang/expr)
* [simple-registry](https://github.com/junqirao/simple-registry)