---
title: 默认模块
language_tabs:
  - shell: Shell
  - http: HTTP
  - javascript: JavaScript
  - ruby: Ruby
  - python: Python
  - php: PHP
  - java: Java
  - go: Go
toc_footers: []
includes: []
search: true
code_clipboard: true
highlight_theme: darkula
headingLevel: 2
generator: "@tarslib/widdershins v4.0.30"

---

# 默认模块

Base URLs:

# Authentication

# 鱼眼/user

## POST Register注册

POST /yuyanprefix/api/v1/user/register

> Body 请求参数

```json
{
  "username": "gdg",
  "password": "123456"
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|

> 返回示例

> 702 Response

```json
{
  "code": 0,
  "reason": "string",
  "message": "string",
  "metadata": {}
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|702|Unknown|none|Inline|

### 返回数据结构

状态码 **702**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» reason|string|true|none||none|
|» message|string|true|none||none|
|» metadata|object|true|none||none|

## POST Login登陆

POST /yuyanprefix/api/v1/user/login

> Body 请求参数

```json
{
  "username": "gdg",
  "password": "123456"
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|

> 返回示例

> 703 Response

```json
{
  "code": 0,
  "reason": "string",
  "message": "string",
  "metadata": {}
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|703|Unknown|none|Inline|

### 返回数据结构

状态码 **703**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» reason|string|true|none||none|
|» message|string|true|none||none|
|» metadata|object|true|none||none|

# 鱼眼/marine

## GET 获取某个user的全部历史记录

GET /yuyanprefix/api/v1/marine/history

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 否 |none|

> 返回示例

> 500 Response

```json
{
  "code": 0,
  "reason": "string",
  "message": "string",
  "metadata": {}
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|none|Inline|

### 返回数据结构

状态码 **500**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» reason|string|true|none||none|
|» message|string|true|none||none|
|» metadata|object|true|none||none|

## POST 添加识别历史记录

POST /yuyanprefix/api/v1/marine/history

> Body 请求参数

```json
{
  "file_name": "new_image22.jpg",
  "creature_id": 2,
  "result_text": "识别结果：蓝鲸"
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 否 |none|
|body|body|object| 否 |none|

> 返回示例

> 404 Response

```json
{
  "code": 0,
  "reason": "string",
  "message": "string",
  "metadata": {}
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|none|Inline|

### 返回数据结构

状态码 **404**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» reason|string|true|none||none|
|» message|string|true|none||none|
|» metadata|object|true|none||none|

## POST 添加生物

POST /yuyanprefix/api/v1/marine/creature

> Body 请求参数

```json
{
  "name": "蓝鲸",
  "description": "是地球史上现存最大的动物，属于须鲸亚目。",
  "image_url": "https://example.com/blue_whale.jpg"
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 否 |none|
|body|body|object| 否 |none|

> 返回示例

> 200 Response

```json
{
  "creatureId": 0
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» creatureId|integer|true|none||none|

## GET 获取某个user的historyid=x记录 

GET /yuyanprefix/api/v1/marine/history/1

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 否 |none|

> 返回示例

> 500 Response

```json
{
  "code": 0,
  "reason": "string",
  "message": "string",
  "metadata": {}
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|none|Inline|

### 返回数据结构

状态码 **500**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» reason|string|true|none||none|
|» message|string|true|none||none|
|» metadata|object|true|none||none|

## GET 获取生物信息By creature_id

GET /yuyanprefix/api/v1/marine/creature/2

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 否 |none|

> 返回示例

> 200 Response

```json
{
  "creature": {
    "creatureId": 0,
    "name": "string",
    "description": "string",
    "imageUrl": "string",
    "createTime": "string"
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» creature|object|true|none||none|
|»» creatureId|integer|true|none||none|
|»» name|string|true|none||none|
|»» description|string|true|none||none|
|»» imageUrl|string|true|none||none|
|»» createTime|string|true|none||none|

## PUT 更新生物信息

PUT /yuyanprefix/api/v1/marine/creature/2

> Body 请求参数

```json
{
  "id": 1,
  "creature": {
    "name": "小丑鱼（已更新）",
    "description": "更新后的描述：一种非常漂亮的热带咸水鱼。",
    "image_url": "https://example.com/new_clownfish.jpg"
  }
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 否 |none|
|body|body|object| 否 |none|

> 返回示例

> 200 Response

```json
{
  "creature": {
    "creatureId": 0,
    "name": "string",
    "description": "string",
    "imageUrl": "string",
    "createTime": "string"
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» creature|object|true|none||none|
|»» creatureId|integer|true|none||none|
|»» name|string|true|none||none|
|»» description|string|true|none||none|
|»» imageUrl|string|true|none||none|
|»» createTime|string|true|none||none|

## DELETE 删除生物信息By creature_id

DELETE /yuyanprefix/api/v1/marine/creature/1

> Body 请求参数

```json
{}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 否 |none|
|body|body|object| 否 |none|

> 返回示例

> 200 Response

```json
{}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

# 数据模型

