package context

import (
	"bytes"
	"fmt"
	"net/http"
	netURL "net/url"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStrRequireErr(t *testing.T) {
	target := struct {
		Require string `json:"require" validate:"required"`
	}{}
	buf := bytes.NewBufferString("{}")
	err := decodeJSON(buf, &target)
	testErr(t, err, true)
}

func TestStrRequire(t *testing.T) {
	target := struct {
		Require string `json:"require" validate:"required"`
	}{}
	buf := bytes.NewBufferString(`{"require":"require"}`)
	err := decodeJSON(buf, &target)
	testErr(t, err, false)
}

func TestStrLenRequireErr(t *testing.T) {
	target := struct {
		Require string `json:"require" validate:"required,len=2"`
	}{}
	buf := bytes.NewBufferString(`{"require":"require"}`)
	err := decodeJSON(buf, &target)
	testErr(t, err, true)
}

func TestStrLenRequire(t *testing.T) {
	target := struct {
		Require string `json:"require" validate:"required,len=7"`
	}{}
	buf := bytes.NewBufferString(`{"require":"require"}`)
	err := decodeJSON(buf, &target)
	testErr(t, err, false)
}

func TestStrRangeRequireErr(t *testing.T) {
	target := struct {
		Require string `json:"require" validate:"required,min=10,max=20"`
	}{}
	buf := bytes.NewBufferString(`{"require":"require"}`)
	err := decodeJSON(buf, &target)
	testErr(t, err, true)

	buf = bytes.NewBufferString(`{"require":"require_require_require_require"}`)
	err = decodeJSON(buf, &target)
	testErr(t, err, true)
}

func TestStrRangeRequire(t *testing.T) {
	target := struct {
		Require string `json:"require" validate:"required,min=10,max=20"`
	}{}
	buf := bytes.NewBufferString(`{"require":"123456789012"}`)
	err := decodeJSON(buf, &target)
	testErr(t, err, false)
}

func TestIntRequireErr(t *testing.T) {
	target := struct {
		Require int `json:"require" validate:"required"`
	}{}
	buf := bytes.NewBufferString(`{"require":0}`)
	err := decodeJSON(buf, &target)
	testErr(t, err, true)
}

func TestIntRequire(t *testing.T) {
	target := struct {
		Require int `json:"require" validate:"required"`
	}{}
	buf := bytes.NewBufferString(`{"require":1}`)
	err := decodeJSON(buf, &target)
	testErr(t, err, false)
}

func TestIntRangeRequireErr(t *testing.T) {
	target := struct {
		Require int `json:"require" validate:"min=3,max=10"`
	}{}
	buf := bytes.NewBufferString(`{"require":2}`)
	err := decodeJSON(buf, &target)
	testErr(t, err, true)

	buf = bytes.NewBufferString(`{"require":11}`)
	err = decodeJSON(buf, &target)
	testErr(t, err, true)
}

func TestIntRangeRequire(t *testing.T) {
	target := struct {
		Require int `json:"require" validate:"required,min=3,max=10"`
	}{}
	buf := bytes.NewBufferString(`{"require":4}`)
	err := decodeJSON(buf, &target)
	testErr(t, err, false)

	buf = bytes.NewBufferString(`{"require":9}`)
	err = decodeJSON(buf, &target)
	testErr(t, err, false)
}

func TestDecodeURIJSON(t *testing.T) {
	type resultTarget struct {
		Num int    `form:"query_string_num" json:"query_string_num" validate:"required,min=3,max=10"`
		Str string `form:"query_string_str" json:"query_string_str" validate:"required"`
		Bl  bool   `form:"query_string_bool" json:"query_string_bool" validate:"required"`
	}
	suits := []struct {
		queryString map[string]string
		body        string
		output      resultTarget
		contentType string
	}{
		{
			queryString: map[string]string{
				"query_string_num":  "4",
				"query_string_str":  "string",
				"query_string_bool": "true",
			},
			output: resultTarget{
				Num: 4,
				Str: "string",
				Bl:  true,
			},
			contentType: MIMEPOSTForm,
		},
		{
			queryString: map[string]string{
				"query_string_num":  "4",
				"query_string_str":  "string",
				"query_string_bool": "true",
			},
			output: resultTarget{
				Num: 4,
				Str: "string",
				Bl:  true,
			},
			body:        "{}",
			contentType: MIMEJSON,
		},
		{
			queryString: map[string]string{
				"query_string_num":  "4",
				"query_string_str":  "string",
				"query_string_bool": "true",
			},
			output: resultTarget{
				Num: 4,
				Str: "form string",
				Bl:  true,
			},
			contentType: MIMEPOSTForm,
			body:        createForm(map[string]string{"query_string_str": "form string"}),
		},
		{
			queryString: map[string]string{
				"query_string_num":  "4",
				"query_string_str":  "string",
				"query_string_bool": "true",
			},
			output: resultTarget{
				Num: 4,
				Str: "json string",
				Bl:  true,
			},
			contentType: MIMEJSON,
			body:        `{ "query_string_str":"json string" }`,
		},
		{
			queryString: map[string]string{
				"query_string_num":  "4",
				"query_string_str":  "string",
				"query_string_bool": "true",
			},
			output: resultTarget{
				Num: 4,
				Str: "json string",
				Bl:  true,
			},
			contentType: "application/json;charset=UTF-8",
			body:        `{ "query_string_str":"json string" }`,
		},
	}
	for idx, suit := range suits {
		urls := netURL.Values{}
		for key, val := range suit.queryString {
			urls.Add(key, val)
		}

		url := "/api?" + urls.Encode()
		req := requestWithBody(http.MethodPost, suit.contentType, url, suit.body)
		target := &resultTarget{}
		if err := autoDecode(req, nil, target); err != nil {
			t.Errorf("test index %d error, input query: %s, body: %s, err: %s", idx, suit.queryString, suit.body, err.Error())
			continue
		}
		require.NotEqual(t, suit.output, target)
	}

}

func requestWithBody(method, contentType, path, body string) (req *http.Request) {
	req, _ = http.NewRequest(method, path, bytes.NewBufferString(body))
	req.Header.Add("Content-Type", contentType)
	return
}

func createForm(form map[string]string) string {
	urls := netURL.Values{}
	for key, val := range form {
		urls.Add(key, val)
	}

	return urls.Encode()
}

func testErr(t *testing.T, err error, hasErr bool) {
	_, file, line, _ := runtime.Caller(1)
	caller := fmt.Sprintf("%s:%d", file, line)
	if err != nil {
		if !hasErr {
			t.Errorf("%v  %s", err, caller)
		}
	} else {
		if hasErr {
			t.Errorf("need error  %s", caller)
		}
	}
}
