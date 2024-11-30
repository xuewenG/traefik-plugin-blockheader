package traefik_plugin_blockheader

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type RuleConfig struct {
	Name string `json:"name,omitempty"`
	Reg  string `json:"reg,omitempty"`
}

type Config struct {
	Rules []RuleConfig `json:"regex,omitempty"`
}

func CreateConfig() *Config {
	return &Config{Rules: make([]RuleConfig, 0)}
}

type Rule struct {
	Name string
	reg  *regexp.Regexp
}

type BlockHeader struct {
	name  string
	next  http.Handler
	rules []Rule
}

func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	rules := make([]Rule, len(config.Rules))

	for i, ruleConfig := range config.Rules {
		reg, err := regexp.Compile(ruleConfig.Reg)
		if err != nil {
			return nil, fmt.Errorf("invalid rule.reg %s: %w", ruleConfig.Reg, err)
		}

		rules[i] = Rule{
			Name: ruleConfig.Name,
			reg:  reg,
		}
	}

	return &BlockHeader{
		name:  name,
		next:  next,
		rules: rules,
	}, nil
}

var forbiddenBytes = []byte(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0, user-scalable=no">
    <title>Forbiden</title>
</head>
<body>
    <div>Forbidden</div>
</body>
</html>`)

func (b *BlockHeader) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	for _, rule := range b.rules {
		value := req.Header.Get(rule.Name)
		if !rule.reg.MatchString(value) {
			continue
		}

		rw.WriteHeader(http.StatusForbidden)

		accept := req.Header.Get("Accept")
		if strings.Contains(accept, "text/html") {
			rw.Header().Set("Content-Type", "text/html")
			rw.Write(forbiddenBytes)
		}

		return
	}

	b.next.ServeHTTP(rw, req)
}
