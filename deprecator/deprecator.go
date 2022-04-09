package deprecator

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/luraproject/lura/config"
)

type Config struct {
	Start    time.Time              `json:"start"`
	Complete time.Time              `json:"complete"`
	Status   int                    `json:"status"`
	Body     map[string]interface{} `json:"body"`
	Headers  map[string]string      `json:"headers"`
}

const Namespace = "github_com/moritzploss/deprecator"

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
	if cfg.Complete.Before(cfg.Start) {
		panic("[deprecator]: Error: time `start` greater than time `complete`")
	}
	if _, ok := json.Marshal(cfg.Body); ok != nil {
		panic("[deprecator]: Error: cannot parse response body. Invalid JSON.")
	}

	return cfg, true
}
