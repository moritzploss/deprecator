# Deprecator

KrakenD middleware for deprecating endpoints.

## Installation

Load the middleware as part of the [KrakenD handler factory chain](https://github.com/devopsfaith/krakend-ce/blob/master/handler_factory.go):

```go
handlerFactory = deprecator.HandlerFactory(handlerFactory)
```

## Quick Start

With the following `extra_config`, KrakenD will reject requests to the
`/user/v1` endpoint starting at 08:00 on 2022-04-10 (see `deprecate` date). For
rejected requests, `status`, `body` and `headers` are set as specified in the
`response` config.

`Deprecation` and `Sunset` headers are set for all responses, both for rejected
and accepted requests, and irrespective of the time and date of the request. The
`sunset` date should be set to the officially communicated deprecation date; the
`heads_up` configuration can be used to specify time windows during which the endpoint
will reject all requests prior to its actual deprecation.

```json
"endpoints": [
  {
    "endpoint": "/user/v1",
    "output_encoding": "no-op",
    "backend": [{
      "host": [ "http://localhost:8080" ],
      "url_pattern": "/__health",
      "encoding": "no-op"
    }],
    "extra_config": {
      "moritzploss/deprecator": {
        "sunset": "2022-04-07T08:00:00+00:00",
        "deprecate": "2022-04-10T08:00:00+00:00",
        "heads_up": {
          "duration": "30m",
          "dates": [
            "2022-04-07T08:00:00+00:00",
            "2022-04-08T08:00:00+00:00",
            "2022-04-09T08:00:00+00:00"
          ]
        },
        "response": {
          "headers": { "Link": "https://myapi.com/api/spec" },
          "body": { "error": "endpoint /user/v1 is deprecated" },
          "status": 410
        }
      }
    }
  }
]
```
