---
title: API v
language_tabs:
  - shell: cURL
  - go: Go
  - python: Python
language_clients:
  - shell: ""
  - go: ""
  - python: ""
toc_footers: []
includes: []
search: true
highlight_theme: darkula
headingLevel: 2

---

<!-- Generator: Widdershins v4.0.1 -->

<h1 id=""> Effective Mobile Тестовое</h1>

> Запуск через yaml кинфиг через **yq** 

```shell
sudo ./start.sh
```


> Scroll down for code samples, example requests and responses. Select a language for code samples from the tabs above or the mobile navigation menu.

<h1 id="-subscriptions">subscriptions</h1>

## Delete subscription

> Code samples

```shell
# You can also use wget
curl -X DELETE /api/v1/subscriptions?id=550e8400-e29b-41d4-a716-446655440000 \
  -H 'Accept: application/json'

```

```go
package main

import (
       "bytes"
       "net/http"
)

func main() {

    headers := map[string][]string{
        "Accept": []string{"application/json"},
    }

    data := bytes.NewBuffer([]byte{jsonReq})
    req, err := http.NewRequest("DELETE", "/api/v1/subscriptions", data)
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}

```

```python
import requests
headers = {
  'Accept': 'application/json'
}

r = requests.delete('/api/v1/subscriptions', params={
  'id': '550e8400-e29b-41d4-a716-446655440000'
}, headers = headers)

print(r.json())

```

`DELETE /api/v1/subscriptions`

Удаление подписки по её id

<h3 id="delete-subscription-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|id|query|string(uuid)|true|Subscription ID (UUID format)|

> Example responses

> 200 Response

```json
{
  "message": "Subscription deleted successfully"
}
```

<h3 id="delete-subscription-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Success message|[handlers.SuccessResponse](#schemahandlers.successresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Bad request - missing ID|[handlers.ErrorResponse](#schemahandlers.errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Subscription not found|[handlers.ErrorResponse](#schemahandlers.errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal server error|[handlers.ErrorResponse](#schemahandlers.errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## Get subscriptions

> Code samples

```shell
# You can also use wget
curl -X GET /api/v1/subscriptions \
  -H 'Accept: application/json'

```

```go
package main

import (
       "bytes"
       "net/http"
)

func main() {

    headers := map[string][]string{
        "Accept": []string{"application/json"},
    }

    data := bytes.NewBuffer([]byte{jsonReq})
    req, err := http.NewRequest("GET", "/api/v1/subscriptions", data)
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}

```

```python
import requests
headers = {
  'Accept': 'application/json'
}

r = requests.get('/api/v1/subscriptions', headers = headers)

print(r.json())

```

`GET /api/v1/subscriptions`

Получение подписки по её id. Если id не указано, возвращаются все

<h3 id="get-subscriptions-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|id|query|string(uuid)|false|Subscription ID (UUID format)|

> Example responses

> 200 Response

```json
{
  "created_at": "string",
  "end_date": "string",
  "id": "string",
  "price": 0,
  "service_name": "string",
  "start_date": "string",
  "updated_at": "string",
  "user_id": "string"
}
```

<h3 id="get-subscriptions-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Single subscription when ID provided|[subscription.Subscription](#schemasubscription.subscription)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Subscription not found|[handlers.ErrorResponse](#schemahandlers.errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal server error|[handlers.ErrorResponse](#schemahandlers.errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## Create a new subscription

> Code samples

```shell
# You can also use wget
curl -X POST /api/v1/subscriptions \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

```go
package main

import (
       "bytes"
       "net/http"
)

func main() {

    headers := map[string][]string{
        "Content-Type": []string{"application/json"},
        "Accept": []string{"application/json"},
    }

    data := bytes.NewBuffer([]byte{jsonReq})
    req, err := http.NewRequest("POST", "/api/v1/subscriptions", data)
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}

```

```python
import requests
headers = {
  'Content-Type': 'application/json',
  'Accept': 'application/json'
}

r = requests.post('/api/v1/subscriptions', headers = headers)

print(r.json())

```

`POST /api/v1/subscriptions`

Создание новой подписки

> Body parameter

```json
{
  "created_at": "string",
  "end_date": "string",
  "id": "string",
  "price": 0,
  "service_name": "string",
  "start_date": "string",
  "updated_at": "string",
  "user_id": "string"
}
```

<h3 id="create-a-new-subscription-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|[subscription.Subscription](#schemasubscription.subscription)|true|Body of the request|

> Example responses

> 200 Response

```json
{
  "created_at": "string",
  "end_date": "string",
  "id": "string",
  "price": 0,
  "service_name": "string",
  "start_date": "string",
  "updated_at": "string",
  "user_id": "string"
}
```

<h3 id="create-a-new-subscription-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|[subscription.Subscription](#schemasubscription.subscription)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Bad Request|[handlers.ErrorResponse](#schemahandlers.errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|[handlers.ErrorResponse](#schemahandlers.errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## Update subscription

> Code samples

```shell
# You can also use wget
curl -X PUT /api/v1/subscriptions?id=497f6eca-6276-4993-bfeb-53cbbbba6f08 \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

```go
package main

import (
       "bytes"
       "net/http"
)

func main() {

    headers := map[string][]string{
        "Content-Type": []string{"application/json"},
        "Accept": []string{"application/json"},
    }

    data := bytes.NewBuffer([]byte{jsonReq})
    req, err := http.NewRequest("PUT", "/api/v1/subscriptions", data)
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}

```

```python
import requests
headers = {
  'Content-Type': 'application/json',
  'Accept': 'application/json'
}

r = requests.put('/api/v1/subscriptions', params={
  'id': '497f6eca-6276-4993-bfeb-53cbbbba6f08'
}, headers = headers)

print(r.json())

```

`PUT /api/v1/subscriptions`

Обновление подписки по её id

> Body parameter

```json
{
  "end_date": "string",
  "price": 0,
  "service_name": "string",
  "start_date": "string",
  "user_id": "string"
}
```

<h3 id="update-subscription-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|id|query|string(uuid)|true|Subscription ID (UUID format)|
|body|body|[handlers.UpdateSubscriptionRequest](#schemahandlers.updatesubscriptionrequest)|true|Subscription object with updated fields|

> Example responses

> 200 Response

```json
{
  "created_at": "string",
  "end_date": "string",
  "id": "string",
  "price": 0,
  "service_name": "string",
  "start_date": "string",
  "updated_at": "string",
  "user_id": "string"
}
```

<h3 id="update-subscription-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Updated subscription|[subscription.Subscription](#schemasubscription.subscription)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Bad request - missing ID or invalid body|[handlers.ErrorResponse](#schemahandlers.errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Subscription not found|[handlers.ErrorResponse](#schemahandlers.errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## Calculate total cost of subscriptions

> Code samples

```shell
# You can also use wget
curl -X GET /api/v1/subscriptions/calculate \
  -H 'Accept: application/json'

```

```go
package main

import (
       "bytes"
       "net/http"
)

func main() {

    headers := map[string][]string{
        "Accept": []string{"application/json"},
    }

    data := bytes.NewBuffer([]byte{jsonReq})
    req, err := http.NewRequest("GET", "/api/v1/subscriptions/calculate", data)
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}

```

```python
import requests
headers = {
  'Accept': 'application/json'
}

r = requests.get('/api/v1/subscriptions/calculate', headers = headers)

print(r.json())

```

`GET /api/v1/subscriptions/calculate`

Подсчет стоимости подписки по фильтрам

<h3 id="calculate-total-cost-of-subscriptions-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|user_id|query|string(uuid)|false|Filter by user ID (UUID format)|
|service_name|query|string|false|Filter by service name|
|start_date|query|string(MM-YYYY)|false|Filter by start date (MM-YYYY format)|
|end_date|query|string(MM-YYYY)|false|Filter by end date (MM-YYYY format)|

> Example responses

> 200 Response

```json
{
  "total": 0
}
```

<h3 id="calculate-total-cost-of-subscriptions-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Total cost calculation result|[handlers.SuccessCostResponse](#schemahandlers.successcostresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal server error|[handlers.ErrorResponse](#schemahandlers.errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

# Schemas

<h2 id="tocS_handlers.ErrorResponse">handlers.ErrorResponse</h2>
<!-- backwards compatibility -->
<a id="schemahandlers.errorresponse"></a>
<a id="schema_handlers.ErrorResponse"></a>
<a id="tocShandlers.errorresponse"></a>
<a id="tocshandlers.errorresponse"></a>

```json
{
  "error": "string",
  "message": "string"
}

```

Error response object

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|error|string|false|none|none|
|message|string|false|none|none|

<h2 id="tocS_handlers.SuccessCostResponse">handlers.SuccessCostResponse</h2>
<!-- backwards compatibility -->
<a id="schemahandlers.successcostresponse"></a>
<a id="schema_handlers.SuccessCostResponse"></a>
<a id="tocShandlers.successcostresponse"></a>
<a id="tocshandlers.successcostresponse"></a>

```json
{
  "total": 0
}

```

Success calculate object

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|total|integer|false|none|none|

<h2 id="tocS_handlers.SuccessResponse">handlers.SuccessResponse</h2>
<!-- backwards compatibility -->
<a id="schemahandlers.successresponse"></a>
<a id="schema_handlers.SuccessResponse"></a>
<a id="tocShandlers.successresponse"></a>
<a id="tocshandlers.successresponse"></a>

```json
{
  "message": "Subscription deleted successfully"
}

```

Success response object

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|message|string|false|none|none|

<h2 id="tocS_handlers.UpdateSubscriptionRequest">handlers.UpdateSubscriptionRequest</h2>
<!-- backwards compatibility -->
<a id="schemahandlers.updatesubscriptionrequest"></a>
<a id="schema_handlers.UpdateSubscriptionRequest"></a>
<a id="tocShandlers.updatesubscriptionrequest"></a>
<a id="tocshandlers.updatesubscriptionrequest"></a>

```json
{
  "end_date": "string",
  "price": 0,
  "service_name": "string",
  "start_date": "string",
  "user_id": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|end_date|string|false|none|none|
|price|integer|false|none|none|
|service_name|string|false|none|none|
|start_date|string|false|none|none|
|user_id|string|false|none|none|

<h2 id="tocS_subscription.Subscription">subscription.Subscription</h2>
<!-- backwards compatibility -->
<a id="schemasubscription.subscription"></a>
<a id="schema_subscription.Subscription"></a>
<a id="tocSsubscription.subscription"></a>
<a id="tocssubscription.subscription"></a>

```json
{
  "created_at": "string",
  "end_date": "string",
  "id": "string",
  "price": 0,
  "service_name": "string",
  "start_date": "string",
  "updated_at": "string",
  "user_id": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|created_at|string|false|none|none|
|end_date|string|false|none|none|
|id|string|false|none|none|
|price|integer|true|none|none|
|service_name|string|true|none|none|
|start_date|string|true|none|none|
|updated_at|string|false|none|none|
|user_id|string|true|none|none|

