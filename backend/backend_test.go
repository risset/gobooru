package backend

import (
	"fmt"
	"log"
	"testing"
)

// Test all APIs and verify they make the correct response
func SearchAll(searchType string, tags string) error {
	apiList := []API{
		DANBOORU,
		GELBOORU,
		KONACHAN,
	}
	params := JSON{}

	for _, api := range apiList {
		if searchType == "tag" {
			params = BuildTagParams(int(api), tags, 0)

		} else {
			params["tags"] = tags
		}

		params["limit"] = 1

		data, err := Search(searchType, params, int(api))

		if err != nil {
			return err
		}

		if len(data) == 0 {
			return fmt.Errorf("API %d: No data in response.", api)
		}
	}

	return nil
}

func TestTagSearches(t *testing.T) {
	err := SearchAll("tag", "black_hair")
	if err != nil {
		log.Fatal(err)
	}
}

func TestPostSearches(t *testing.T) {
	err := SearchAll("post", "black_hair")
	if err != nil {
		log.Fatal(err)
	}
}

func TestFormatJSON(t *testing.T) {
	data := JSON{"limit": 1}
	_, err := FormatJSON(data)
	if err != nil {
		log.Fatal(err)
	}
}

func TestEncodeUrl(t *testing.T) {
	baseUrl := fmt.Sprintf(apiUrls[DANBOORU], "post")
	params := JSON{"limit": 1}
	_, err := EncodeUrl(baseUrl, params)
	if err != nil {
		log.Fatal(err)
	}
}

func TestInitConfig(t *testing.T) {
	err := InitConfig("~/.config/gobooru")
	if err != nil {
		log.Fatal(err)
	}
}
