package traefik_plugin_blockheader_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	traefik_plugin_blockheader "github.com/xuewenG/traefik-plugin-blockheader"
)

type HeaderConfig struct {
	Name  string
	Value string
}

func TestCreateConfig(t *testing.T) {
	t.Run("create config", func(t *testing.T) {
		config := traefik_plugin_blockheader.CreateConfig()
		if config == nil {
			t.Fatal("create config failed")
		}
	})
}

func TestNew(t *testing.T) {
	tests := []struct {
		desc   string
		rules  []traefik_plugin_blockheader.RuleConfig
		expErr bool
	}{
		{
			desc:   "should return no error",
			expErr: false,
			rules: []traefik_plugin_blockheader.RuleConfig{
				{
					Name: "bar",
					Reg:  "foo",
				},
				{
					Name: "bar",
					Reg:  "foo",
				},
			},
		},
		{
			desc:   "should return an error",
			expErr: true,
			rules: []traefik_plugin_blockheader.RuleConfig{
				{
					Name: "bar",
					Reg:  "*",
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			config := &traefik_plugin_blockheader.Config{
				Rules: test.rules,
			}

			_, err := traefik_plugin_blockheader.New(context.Background(), nil, config, "traefik_plugin_blockheader")

			if !test.expErr && err != nil {
				t.Fatal("test failed")
			}

			if test.expErr && err == nil {
				t.Fatal("test failed")
			}
		})
	}
}

func TestServeHTTP(t *testing.T) {
	tests := []struct {
		desc      string
		rules     []traefik_plugin_blockheader.RuleConfig
		headers   []HeaderConfig
		forbidden bool
	}{
		{
			desc:      "should be forbidden",
			forbidden: true,
			rules: []traefik_plugin_blockheader.RuleConfig{
				{
					Name: "User-Agent",
					Reg:  "MicroMessenger",
				},
			},
			headers: []HeaderConfig{
				{
					Name:  "User-Agent",
					Value: "Mozilla/5.0 (iPhone; CPU iPhone OS 17_6_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/8.0.50(0x1800323d) NetType/WIFI Language/zh_CN",
				},
			},
		},
		{
			desc:      "should not be forbidden",
			forbidden: false,
			rules: []traefik_plugin_blockheader.RuleConfig{
				{
					Name: "User-Agent",
					Reg:  "MicroMessenger",
				},
			},
			headers: []HeaderConfig{
				{
					Name:  "User-Agent",
					Value: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36",
				},
			},
		},
		{
			desc:      "should not be forbidden",
			forbidden: false,
			rules: []traefik_plugin_blockheader.RuleConfig{
				{
					Name: "",
					Reg:  "MicroMessenger",
				},
			},
			headers: []HeaderConfig{
				{
					Name:  "User-Agent",
					Value: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			config := &traefik_plugin_blockheader.Config{
				Rules: test.rules,
			}

			next := func(rw http.ResponseWriter, req *http.Request) {
			}

			blockHeader, err := traefik_plugin_blockheader.New(context.Background(), http.HandlerFunc(next), config, "traefik_plugin_blockheader")
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Accept", "text/html")
			for _, header := range test.headers {
				req.Header.Set(header.Name, header.Value)
			}

			blockHeader.ServeHTTP(recorder, req)

			statusCode := recorder.Result().StatusCode
			if test.forbidden && statusCode != http.StatusForbidden {
				t.Error("should be forbidden")
			}

			if !test.forbidden && statusCode == http.StatusForbidden {
				t.Error("should not be forbidden")
			}
		})
	}
}
