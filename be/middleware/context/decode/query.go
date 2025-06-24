package decode

import "net/http"

func Query(req *http.Request, obj interface{}) error {
	values := req.URL.Query()
	if len(values) == 0 {
		return nil
	}
	if err := mapForm(obj, values); err != nil {
		return err
	}
	if err := mapFormByTag(obj, values, "query"); err != nil {
		return err
	}
	return nil
}

func HTTPHeader(header http.Header, obj interface{}) error {
	if len(header) == 0 {
		return nil
	}
	if err := mapFormByTag(obj, header, "header"); err != nil {
		return err
	}
	return nil

}

func HTTPUri(values map[string][]string, obj interface{}) error {
	if len(values) == 0 {
		return nil
	}
	if err := mapFormByTag(obj, values, "uri"); err != nil {
		return err
	}
	return nil
}
