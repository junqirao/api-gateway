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

at config.yaml

    gateway.management.password

If not configured, ignore. If set, the Authentication value is plaintext password once MD5

relationship:

    clear-text passwords --md5--> Authentication ---md5--> gateway.management.password

## GET Config

GET {entrance}/management/config

### Request parameters

| Name          | Position | Type   | Required | Description |
|---------------|----------|--------|----------|-------------|
| service_name  | query    | string | No       | none        |
| Authorization | header   | string | No       | none        |

> Response Example

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
        "interval": "30s"
      }
    }
  }
}
```

### Response Result

| Status Code | Meaning of Status Code                                  | Description | Data Model |
|-------------|---------------------------------------------------------|-------------|------------|
| 200         | [OK](https://tools.ietf.org/html/rfc7231#section-6.3.1) | none        | Inline     |

### Response Data Struct

## PUT Load Balance Config

PUT {entrance}/management/config/load_balance

> Body Request parameters

```json
{
  "service_name": "service_name",
  "config": {
    "strategy": "round-robin"
  }
}
```

supported strategy:
```round-robin```
```random```
```weighted-round-robin```
```weighted-random```
```less-load```
```hash```

### Request parameters

| Name          | Position | Type   | Required | Description |
|---------------|----------|--------|----------|-------------|
| Authorization | header   | string | No       | none        |
| body          | body     | object | No       | none        |

> Response Example

> 200 Response

```json
{
  "message": "success",
  "code": 0
}
```

### Response Result

| Status Code | Meaning of Status Code                                  | Description | Data Model |
|-------------|---------------------------------------------------------|-------------|------------|
| 200         | [OK](https://tools.ietf.org/html/rfc7231#section-6.3.1) | none        | Inline     |

### Response Data Struct

## PUT Rate Limiter Config

PUT {entrance}/management/config/rate_limiter

> Body Request parameters

```json
{
  "service_name": "test",
  "config": {
    "rate": 5,
    "peak": 10
  }
}
```

### Request parameters

| Name          | Position | Type   | Required | Description |
|---------------|----------|--------|----------|-------------|
| Authorization | header   | string | No       | none        |
| body          | body     | object | No       | none        |

> Response Example

> 200 Response

```json
{
  "message": "success",
  "code": 0
}
```

### Response Result

| Status Code | Meaning of Status Code                                  | Description | Data Model |
|-------------|---------------------------------------------------------|-------------|------------|
| 200         | [OK](https://tools.ietf.org/html/rfc7231#section-6.3.1) | none        | Inline     |

### Response Data Struct

## PUT Circuit Breaker Config

PUT {entrance}/management/config/breaker

> Body Request parameters

```json
{
  "service_name": "test",
  "config": {
    "max_failures": 5,
    "half_open_max_requests": 1,
    "open_timeout": "30s",
    "interval": "30s"
  }
}
```

### Request parameters

| Name          | Position | Type   | Required | Description |
|---------------|----------|--------|----------|-------------|
| Authorization | header   | string | No       | none        |
| body          | body     | object | No       | none        |

> Response Example

> 200 Response

```json
{
  "message": "success",
  "code": 0
}
```

### Response Result

| Status Code | Meaning of Status Code                                  | Description | Data Model |
|-------------|---------------------------------------------------------|-------------|------------|
| 200         | [OK](https://tools.ietf.org/html/rfc7231#section-6.3.1) | none        | Inline     |

### Response Data Struct

## GET Global Variables

GET {entrance}/management/program/variable

### Request parameters

| Name          | Position | Type   | Required | Description |
|---------------|----------|--------|----------|-------------|
| Authorization | header   | string | No       | none        |

> Response Example

> 200 Response

```json
{
  "message": "success",
  "code": 0
}
```

### Response Result

| Status Code | Meaning of Status Code                                  | Description | Data Model |
|-------------|---------------------------------------------------------|-------------|------------|
| 200         | [OK](https://tools.ietf.org/html/rfc7231#section-6.3.1) | none        | Inline     |

### Response Data Struct

## PUT Global Variable

PUT {entrance}/management/program/variable

> Body Request parameters

```json
{
  "key": "your_key",
  "value": "your_value"
}
```

### Request parameters

| Name          | Position | Type   | Required | Description |
|---------------|----------|--------|----------|-------------|
| Authorization | header   | string | No       | none        |
| body          | body     | object | No       | none        |

> Response Example

> 200 Response

```json
{
  "message": "success",
  "code": 0
}
```

### Response Result

| Status Code | Meaning of Status Code                                  | Description | Data Model |
|-------------|---------------------------------------------------------|-------------|------------|
| 200         | [OK](https://tools.ietf.org/html/rfc7231#section-6.3.1) | none        | Inline     |

### Response Data Struct

## DELETE Global Variable

DELETE {entrance}/management/program/variable

### Request parameters

| Name          | Position | Type   | Required | Description |
|---------------|----------|--------|----------|-------------|
| key           | query    | string | Yes      | none        |
| Authorization | header   | string | No       | none        |

> Response Example

> 200 Response

```json
{
  "message": "success",
  "code": 0
}
```

### Response Result

| Status Code | Meaning of Status Code                                  | Description | Data Model |
|-------------|---------------------------------------------------------|-------------|------------|
| 200         | [OK](https://tools.ietf.org/html/rfc7231#section-6.3.1) | none        | Inline     |

### Response Data Struct

## GET Program Info

GET {entrance}/management/program/info

### Request parameters

| Name          | Position | Type   | Required | Description |
|---------------|----------|--------|----------|-------------|
| service_name  | query    | string | No       | none        |
| Authorization | header   | string | No       | none        |

> Response Example

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

### Response Result

| Status Code | Meaning of Status Code                                  | Description | Data Model |
|-------------|---------------------------------------------------------|-------------|------------|
| 200         | [OK](https://tools.ietf.org/html/rfc7231#section-6.3.1) | none        | Inline     |

### Response Data Struct

## DELETE Program

DELETE {entrance}/management/program/info

### Request parameters

| Name          | Position | Type   | Required | Description |
|---------------|----------|--------|----------|-------------|
| service_name  | query    | string | No       | none        |
| name          | query    | string | No       | none        |
| Authorization | header   | string | No       | none        |

> Response Example

> 200 Response

```json
{
  "message": "success",
  "code": 0
}
```

### Response Result

| Status Code | Meaning of Status Code                                  | Description | Data Model |
|-------------|---------------------------------------------------------|-------------|------------|
| 200         | [OK](https://tools.ietf.org/html/rfc7231#section-6.3.1) | none        | Inline     |

### Response Data Struct

## PUT Program

PUT {entrance}/management/program/info

> Body Request parameters

```json
{
  "expr": "{base64_encoded_string}",
  "name": "test-program",
  "service_name": "test"
}
```

### Request parameters

| Name          | Position | Type   | Required | Description |
|---------------|----------|--------|----------|-------------|
| Authorization | header   | string | No       | none        |
| body          | body     | object | No       | none        |

> Response Example

> 200 Response

```json
{
  "message": "success",
  "code": 0
}
```

### Response Result

| Status Code | Meaning of Status Code                                  | Description | Data Model |
|-------------|---------------------------------------------------------|-------------|------------|
| 200         | [OK](https://tools.ietf.org/html/rfc7231#section-6.3.1) | none        | Inline     |

