---
title: API-Gateway
language_tabs:
  - shell: Shell
  - http: HTTP
  - javascript: JavaScript
  - ruby: Ruby
  - python: Python
  - php: PHP
  - java: Java
  - go: Go
toc_footers: [ ]
includes: [ ]
search: true
code_clipboard: true
highlight_theme: darkula
headingLevel: 2
generator: "@tarslib/widdershins v4.0.23"

---

# API-Gateway


## Authentication
见config.yaml

    gateway.management.password

如没有配置则忽略，若有设置，Authentication值为明文密码一次MD5

关系为：

    明文密码 --md5--> Authentication ---md5--> gateway.management.password

## GET 获取配置

GET {entrance}/management/config

### 请求参数

| 名称            | 位置     | 类型     | 必选 | 说明   |
|---------------|--------|--------|----|------|
| service_name  | query  | string | 否  | none |
| Authorization | header | string | 否  | none |

> 返回示例

> 200 Response

```json
{
  "message": "success",
  "code": 0,
  "data": {
    "default": false,
    "config": {
      "reverse_proxy": {
        "trim_routing_key_prefix": true,
        "retry_count": 1,
        "dial_timeout": "",
        "tls_handshake_timeout": "",
        "scheme": ""
      },
      "load_balance": {
        "strategy": "round_robin"
      },
      "rate_limiter": {
        "rate": 0,
        "peak": 0
      },
      "breaker": {
        "name": "",
        "max_failures": 5,
        "half_open_max_requests": 1,
        "open_timeout": "30s",
        "Interval": "30s"
      }
    }
  }
}
```

### 返回结果

| 状态码 | 状态码含义                                                   | 说明   | 数据模型   |
|-----|---------------------------------------------------------|------|--------|
| 200 | [OK](https://tools.ietf.org/html/rfc7231#section-6.3.1) | none | Inline |

### 返回数据结构

## PUT 设置负载均衡配置

PUT {entrance}/management/config/load_balance

> Body 请求参数

```json
{
  "service_name": "service_name",
  "config": {
    "strategy": "round_robin/random"
  }
}
```

### 请求参数

| 名称            | 位置     | 类型     | 必选 | 说明   |
|---------------|--------|--------|----|------|
| Authorization | header | string | 否  | none |
| body          | body   | object | 否  | none |

> 返回示例

> 200 Response

```json
{
  "message": "success",
  "code": 0
}
```

### 返回结果

| 状态码 | 状态码含义                                                   | 说明   | 数据模型   |
|-----|---------------------------------------------------------|------|--------|
| 200 | [OK](https://tools.ietf.org/html/rfc7231#section-6.3.1) | none | Inline |

### 返回数据结构

## PUT 设置限流器配置

PUT {entrance}/management/config/rate_limiter

> Body 请求参数

```json
{
  "service_name": "test",
  "config": {
    "rate": 5,
    "peak": 10
  }
}
```

### 请求参数

| 名称            | 位置     | 类型     | 必选 | 说明   |
|---------------|--------|--------|----|------|
| Authorization | header | string | 否  | none |
| body          | body   | object | 否  | none |

> 返回示例

> 200 Response

```json
{
  "message": "success",
  "code": 0
}
```

### 返回结果

| 状态码 | 状态码含义                                                   | 说明   | 数据模型   |
|-----|---------------------------------------------------------|------|--------|
| 200 | [OK](https://tools.ietf.org/html/rfc7231#section-6.3.1) | none | Inline |

### 返回数据结构

## PUT 设置断路器配置

PUT {entrance}/management/config/breaker

> Body 请求参数

```json
{
  "service_name": "test",
  "config": {
    "max_failures": 5,
    "half_open_max_requests": 1,
    "open_timeout": "30s",
    "Interval": "30s"
  }
}
```

### 请求参数

| 名称            | 位置     | 类型     | 必选 | 说明   |
|---------------|--------|--------|----|------|
| Authorization | header | string | 否  | none |
| body          | body   | object | 否  | none |

> 返回示例

> 200 Response

```json
{
  "message": "success",
  "code": 0
}
```

### 返回结果

| 状态码 | 状态码含义                                                   | 说明   | 数据模型   |
|-----|---------------------------------------------------------|------|--------|
| 200 | [OK](https://tools.ietf.org/html/rfc7231#section-6.3.1) | none | Inline |

### 返回数据结构

## GET 获取全局变量

GET {entrance}/management/program/variable

### 请求参数

| 名称            | 位置     | 类型     | 必选 | 说明   |
|---------------|--------|--------|----|------|
| Authorization | header | string | 否  | none |

> 返回示例

> 200 Response

```json
{
  "message": "success",
  "code": 0
}
```

### 返回结果

| 状态码 | 状态码含义                                                   | 说明   | 数据模型   |
|-----|---------------------------------------------------------|------|--------|
| 200 | [OK](https://tools.ietf.org/html/rfc7231#section-6.3.1) | none | Inline |

### 返回数据结构

## PUT 设置全局变量

PUT {entrance}/management/program/variable

> Body 请求参数

```json
{
  "key": "your_key",
  "value": "your_value"
}
```

### 请求参数

| 名称            | 位置     | 类型     | 必选 | 说明   |
|---------------|--------|--------|----|------|
| Authorization | header | string | 否  | none |
| body          | body   | object | 否  | none |

> 返回示例

> 200 Response

```json
{
  "message": "success",
  "code": 0
}
```

### 返回结果

| 状态码 | 状态码含义                                                   | 说明   | 数据模型   |
|-----|---------------------------------------------------------|------|--------|
| 200 | [OK](https://tools.ietf.org/html/rfc7231#section-6.3.1) | none | Inline |

### 返回数据结构

## DELETE 移除全局变量

DELETE {entrance}/management/program/variable


### 请求参数

| 名称            | 位置     | 类型     | 必选 | 说明   |
|---------------|--------|--------|----|------|
| key           | query  | string | 是  | none |
| Authorization | header | string | 否  | none |

> 返回示例

> 200 Response

```json
{
  "message": "success",
  "code": 0
}
```

### 返回结果

| 状态码 | 状态码含义                                                   | 说明   | 数据模型   |
|-----|---------------------------------------------------------|------|--------|
| 200 | [OK](https://tools.ietf.org/html/rfc7231#section-6.3.1) | none | Inline |

### 返回数据结构

## GET 获取程序

GET {entrance}/management/program/info

### 请求参数

| 名称            | 位置     | 类型     | 必选 | 说明   |
|---------------|--------|--------|----|------|
| service_name  | query  | string | 否  | none |
| Authorization | header | string | 否  | none |

> 返回示例

> 200 Response

```json
{
  "message": "success",
  "code": 0,
  "data": {
    "test": [
      {
        "expr": "{base64_encoded_string}",
        "name": "test-program",
        "service_name": "test"
      }
    ]
  }
}
```

### 返回结果

| 状态码 | 状态码含义                                                   | 说明   | 数据模型   |
|-----|---------------------------------------------------------|------|--------|
| 200 | [OK](https://tools.ietf.org/html/rfc7231#section-6.3.1) | none | Inline |

### 返回数据结构

## DELETE 移除程序

DELETE {entrance}/management/program/info

### 请求参数

| 名称            | 位置     | 类型     | 必选 | 说明   |
|---------------|--------|--------|----|------|
| service_name  | query  | string | 否  | none |
| name          | query  | string | 否  | none |
| Authorization | header | string | 否  | none |

> 返回示例

> 200 Response

```json
{
  "message": "success",
  "code": 0
}
```

### 返回结果

| 状态码 | 状态码含义                                                   | 说明   | 数据模型   |
|-----|---------------------------------------------------------|------|--------|
| 200 | [OK](https://tools.ietf.org/html/rfc7231#section-6.3.1) | none | Inline |

### 返回数据结构

## PUT 设置程序

PUT {entrance}/management/program/info

> Body 请求参数

```json
{
  "expr": "{base64_encoded_string}",
  "name": "test-program",
  "service_name": "test"
}
```

### 请求参数

| 名称            | 位置     | 类型     | 必选 | 说明   |
|---------------|--------|--------|----|------|
| Authorization | header | string | 否  | none |
| body          | body   | object | 否  | none |

> 返回示例

> 200 Response

```json
{
  "message": "success",
  "code": 0
}
```

### 返回结果

| 状态码 | 状态码含义                                                   | 说明   | 数据模型   |
|-----|---------------------------------------------------------|------|--------|
| 200 | [OK](https://tools.ietf.org/html/rfc7231#section-6.3.1) | none | Inline |

