# API Gateway

[中文文档](/README_CN) | [English Doc](/README)

configuration hot updates supported programmable api gateway

-------

<!-- TOC -->
* [API Gateway](#api-gateway)
  * [Overview](#overview)
  * [Features](#features)
    * [Hot Updates Supported](#hot-updates-supported)
    * [Todo List](#todo-list)
  * [Quick Start](#quick-start)
    * [API Reference](#api-reference)
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
  * [Config Example](#config-example)
  * [Dependencies](#dependencies)
<!-- TOC -->

## Overview

![overview](/doc/overview.png)

## Features

* Load Balance
    * Weighted Round-Robin
    * Weighted Random
* Rate Limiter
* Circuit Breaker
* Programmable
* Retry
* Configuration Hot Updates

### Hot Updates Supported

* Load Balance
* Rate Limiter
* Circuit Breaker
* Program

### Todo List

* Metrics for prometheus

## Quick Start

### API Reference

[api reference](/doc/api_reference_en.md)

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

| Name               | Type     | Description              | Usage                                       |
|--------------------|----------|--------------------------|---------------------------------------------|
| request.RemoteAddr | String   | request remote address   | request.RemoteAddr                          |
| request.Host       | String   | request host             | request.Host                                |
| request.URL        | String   | request url              | request.URL                                 |
| request.Method     | String   | request method           | request.Method                              |
| request.Header     | Object   | request header           | -                                           |
| request.Header.Add | Function | add request header value | request.Header.Add(key string,value string) |
| request.Header.Set | Function | set request header value | request.Header.Set(key string,value string) |
| request.Header.Get | Function | get request header value | request.Header.Get(key string)              |

##### response

response object

| Name                | Type     | Description               | Usage                                        |
|---------------------|----------|---------------------------|----------------------------------------------|
| response.Header.Add | Function | add response header value | response.Header.Add(key string,value string) |
| response.Header.Set | Function | set response header value | response.Header.Set(key string,value string) |
| response.Header.Get | Function | get response header value | response.Header.Get(key string)              |

##### global

global variable, this value is shared across all api-gateway instances.

| Name       | Type     | Description         | Usage                               |
|------------|----------|---------------------|-------------------------------------|
| global     | Map      | global variable map | global.key / global["key"]          |
| set_global | Function | set global variable | set_global(key string,value string) |

##### ctx

context from request, use as parameter of function

##### logger

logger object

| Name          | Type     | Description                   | Usage                                                |
|---------------|----------|-------------------------------|------------------------------------------------------|
| logger.Info   | Function | log with info level           | logger.Info(ctx context, value any)                  |
| logger.Warn   | Function | log with warn level           | logger.Warn(ctx context, value any)                  |
| logger.Error  | Function | log with error level          | logger.Error(ctx context, value any)                 |
| logger.Infof  | Function | formated log with info level  | logger.Infof(ctx context, format string, value any)  |
| logger.Warnf  | Function | formated log with warn level  | logger.Warnf(ctx context, format string, value any)  |
| logger.Errorf | Function | formated log with error level | logger.Errorf(ctx context, format string, value any) |

##### upstream

| Name                 | Type    | Description           | Usage                |
|----------------------|---------|-----------------------|----------------------|
| upstream.Id          | String  | upstream instance id  | upstream.Id          |
| upstream.Host        | String  | upstream host         | upstream.Host        |
| upstream.HostName    | String  | upstream host name    | upstream.HostName    |
| upstream.Port        | Integer | upstream port         | upstream.Port        |
| upstream.ServiceName | String  | upstream service name | upstream.ServiceName |
| upstream.Meta        | Map     | upstream metadata     | upstream.Meta        |

##### terminate_if

| Name         | Type     | Description                      | Usage                                                 |
|--------------|----------|----------------------------------|-------------------------------------------------------|
| terminate_if | Function | terminate request with condition | terminate_if(condition bool, reason [optional]string) |

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