## cli

```
$ ./cli
Usage: cli {cart|order|stock|product} [options ...]

$ ./cli product
Usage:
  product list
  product info <sku>

$ ./cli product list
HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 21 Feb 2025 15:22:26 GMT
Content-Length: 90

{
    "skus": [
        1076963,
        1148162,
        1625903,
        2618151,
        2956315,
        2958025,
        3596599,
        3618852,
        4288068,
        4465995
    ]
}

$ ./cli product info 2618151
HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 21 Feb 2025 15:23:11 GMT
Content-Length: 59

{
    "name": "Пора снимать бикини",
    "price": 452
}

$ ./cli product info 3618852
HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 21 Feb 2025 15:23:23 GMT
Content-Length: 198

{
    "name": "Защитный код. Как выжить в нашем городе | Майоров Олег Вячеславович, Степанов Максим Викторович",
    "price": 4645
}

$ ./cli stock info 2618151
HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 21 Feb 2025 15:23:54 GMT
Content-Length: 11

{
    "count": 9
}

$ ./cli stock info 3618852
HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 21 Feb 2025 15:24:04 GMT
Content-Length: 11

{
    "count": 8
}

$ ./cli cart add 2618151 10
HTTP/1.1 412 Precondition Failed
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Fri, 21 Feb 2025 15:24:35 GMT
Content-Length: 20

precondition failed

$ ./cli cart add 2618151 9
HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 21 Feb 2025 15:24:42 GMT
Content-Length: 2

{}

$ ./cli cart add 3618852
HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 21 Feb 2025 15:24:55 GMT
Content-Length: 2

{}

$ ./cli cart list
HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 21 Feb 2025 15:25:04 GMT
Content-Length: 336

{
    "items": [
        {
            "sku": 2618151,
            "count": 9,
            "name": "Пора снимать бикини",
            "price": 452
        },
        {
            "sku": 3618852,
            "count": 1,
            "name": "Защитный код. Как выжить в нашем городе | Майоров Олег Вячеславович, Степанов Максим Викторович",
            "price": 4645
        }
    ],
    "totalPrice": 8713
}

$ ./cli cart checkout
HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 21 Feb 2025 15:25:22 GMT
Content-Length: 13

{
    "orderID": 6
}

$ ./cli cart list
HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 21 Feb 2025 15:25:29 GMT
Content-Length: 27

{
    "items": [],
    "totalPrice": 0
}

$ ./cli order info 6
HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 21 Feb 2025 15:25:43 GMT
Content-Length: 102

{
    "status": "awaiting payment",
    "user": 123,
    "items": [
        {
            "sku": 2618151,
            "count": 9
        },
        {
            "sku": 3618852,
            "count": 1
        }
    ]
}

$ ./cli stock info 2618151
HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 21 Feb 2025 15:26:04 GMT
Content-Length: 11

{
    "count": 0
}

$ ./cli stock info 3618852
HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 21 Feb 2025 15:26:20 GMT
Content-Length: 11

{
    "count": 7
}

$ ./cli order pay 6
HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 21 Feb 2025 15:26:36 GMT
Content-Length: 2

{}

$ ./cli order info 6
HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 21 Feb 2025 15:26:56 GMT
Content-Length: 91

{
    "status": "payed",
    "user": 123,
    "items": [
        {
            "sku": 2618151,
            "count": 9
        },
        {
            "sku": 3618852,
            "count": 1
        }
    ]
}
```
