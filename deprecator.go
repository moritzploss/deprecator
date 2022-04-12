package deprecator

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/luraproject/lura/config"
)

type Response struct {
	Status  int                    `json:"status"`
	Body    map[string]interface{} `json:"body"`
	Headers map[string]string      `json:"headers"`
}

type HeadsUp struct {
	Duration Duration    `json:"duration"`
	Dates    []time.Time `json:"dates"`
}

type Config struct {
	Sunset    time.Time `json:"sunset"`
	Deprecate time.Time `json:"deprecate"`
	HeadsUp   HeadsUp   `json:"heads_up"`
	Response  Response  `json:"response"`
}

const Namespace = "moritzploss/deprecator"

func ConfigGetter(e config.ExtraConfig) (*Config, bool) {
	cfg := new(Config)

	tmp, ok := e[Namespace]
	if !ok {
		return cfg, false
	}

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(tmp); err != nil {
		panic("[deprecator]: Error: failed to parse config")
	}
	if err := json.NewDecoder(buf).Decode(cfg); err != nil {
		panic("[deprecator]: Error: failed to parse config")
	}
	if cfg.Deprecate.Before(cfg.Sunset) {
		panic("[deprecator]: Error: time `sunset` greater than time `deprecate`")
	}
	if _, ok := json.Marshal(cfg.Response.Body); ok != nil {
		panic("[deprecator]: Error: cannot parse response body. Invalid JSON.")
	}

	return cfg, true
}
