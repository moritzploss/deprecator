# Deprecator

KrakenD middleware for deprecating endpoints.

## Installation

Load the middleware as part of the [KrakenD handler factory chain](https://github.com/devopsfaith/krakend-ce/blob/master/handler_factory.go):

```go
handlerFactory = deprecator.HandlerFactory(handlerFactory)
```

## Quick Start

With the following `extra_config`, KrakenD will reject requests to the
`/user/v1` endpoint starting at midnight on 2022-04-10. The probability
of a request beeing rejected increases linearly between the `start` time (0 %)
and the `complete` time (100 %). For rejected requests, `status`, `body` and
`headers` are set as specified.

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
      "github_com/moritzploss/deprecator": {
        "start": "2022-04-10T00:00:00+00:00",
        "complete": "2022-04-17T00:00:00+00:00",
        "status": 410,
        "body": { "error": "endpoint /user/v1 is deprecated" },
        "headers": { "Link": "https://myapi.com/api/spec" }
      }
    }
  }
]
```

The calculation of the rejection probability is stateless and evaluated in
isolation for each incoming request. `Deprecation` and `Sunset` headers are
set for all responses, both for rejected and accepted requests.
