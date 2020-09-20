package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

func getAndParseHTML(u string) (kgData, error) {
	data := kgData{}
	resp, err := http.Get(u)
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return data, fmt.Errorf("http status %s", resp.Status)
	}

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return data, err
	}

	rx := regexp.MustCompile(`window.__DATA__\s?=\s?(.*)?;\s?</script>`)
	results := rx.FindAllSubmatch(res, 1)
	if results == nil || len(results) == 0 {
		return data, fmt.Errorf("playurl not found")
	}
	bytesData := results[0][1]
	err = json.Unmarshal(bytesData, &data)
	if err != nil {
		return data, fmt.Errorf("failed to decode data \n%s \n%v", string(bytesData), err)
	}

	return data, nil
}
