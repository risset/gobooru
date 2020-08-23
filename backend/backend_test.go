package backend

import (
	"log"
	"testing"
)

// Test post search for all APIs
func TestPostSearch(t *testing.T) {
	tags := "blue_sky 1girl"
	limit := 1
	random := false

	for _, api := range []API{DANBOORU, GELBOORU, KONACHAN} {
		s := BuildPostSearch(api, tags, limit, random)
		_, err := GetData(s)
		if err != nil {
			log.Fatal(err)
		}
	}

}

// Test tag search for all APIs
func TestTagSearch(t *testing.T) {
	tag := "1girl"
	limit := 1
	order := 1

	for _, api := range []API{DANBOORU, GELBOORU, KONACHAN} {
		s := BuildTagSearch(api, tag, limit, order)
		_, err := GetData(s)
		if err != nil {
			log.Fatal(err)
		}
	}

}
