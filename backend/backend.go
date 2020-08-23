package backend

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

type API int

const (
	DANBOORU API = iota
	GELBOORU
	KONACHAN
)

type searchType int

const (
	POST searchType = iota
	TAG
)

type JSON map[string]interface{}

type search struct {
	baseUrl string
	params  JSON
}

func init() {
	initConfig("~/.config/gobooru")
}

// Create config file if it doesn't exist
func initConfig(p string) error {
	cfgPath, err := homedir.Expand(p)
	if err != nil {
		return err
	}

	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(cfgPath)
	viper.SetDefault("danbooru_username", " ")
	viper.SetDefault("danbooru_api_key", " ")
	viper.SetDefault("gelbooru_user_id", " ")
	viper.SetDefault("gelbooru_api_key", " ")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			os.Mkdir(cfgPath, 0755)
			viper.WriteConfigAs(path.Join(cfgPath, "config.yaml"))
		} else {
			return err
		}
	}

	return nil
}

// Get array of JSON objects from URL
func request(url string) ([]JSON, error) {
	data := []JSON{}

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Get formatted URL from API URL and search parameters
func encodeUrl(s search) (string, error) {
	u, err := url.Parse(s.baseUrl)
	if err != nil {
		return "", err
	}

	query := u.Query()

	for k, v := range s.params {
		query.Set(k, fmt.Sprintf("%v", v))
	}
	u.RawQuery = query.Encode()

	return u.String(), nil
}

// Format JSON map as pretty-printed JSON string
func formatJSON(data JSON) (string, error) {
	b, err := json.MarshalIndent(data, "", "  ")

	if err != nil {
		return " ", err
	}

	return string(b), nil
}

// Print formatted JSON to standard output
func ShowJSON(data []JSON) error {
	for _, p := range data {
		s, err := formatJSON(p)
		if err != nil {
			return err
		}
		fmt.Println(s)
	}

	return nil
}

// Download image from URL and save to local file
func GetImg(post JSON, dir string, ch chan<- string) {
	start := time.Now()
	fileUrl := fmt.Sprintf("%v", post["file_url"])
	resp, err := http.Get(fileUrl)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}
	defer resp.Body.Close()

	file, err := os.Create(fmt.Sprintf("%s/%s", dir, path.Base(fileUrl)))
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}
	defer file.Close()

	bytes, err := io.Copy(file, resp.Body)
	if err != nil {
		ch <- fmt.Sprintf("While reading %s: %v", fileUrl, err)
		return
	}

	secs := time.Since(start).Seconds()
	mb := float64(bytes) / 1e+6
	ch <- fmt.Sprintf("%.2fs %7fmiB %s", secs, mb, fileUrl)
}

// Return directory name for given parameters
func GetImgDirName(random bool, tags string) string {
	if random {
		return "random"
	} else {
		return strings.Replace(tags, " ", ",", -1)
	}
}

// For array of JSON entries, concurrently download all images
func GetAllImages(data []JSON, dir string) error {
	err := os.Mkdir(dir, 0755)
	if err != nil {
		return err
	}

	ch := make(chan string)
	for _, p := range data {
		go GetImg(p, dir, ch)
	}

	for range data {
		fmt.Println(<-ch)
	}

	return nil
}

// Fill base API url template with search type
func convertBaseUrl(baseUrl string, st searchType) string {
	switch st {
	case POST:
		return fmt.Sprintf(baseUrl, "post")
	case TAG:
		return fmt.Sprintf(baseUrl, "tag")
	default:
		return " "
	}
}

// Initialize search object with common and danbooru-specific values
func danbooruSearch(limit int, st searchType) search {
	s := search{
		baseUrl: convertBaseUrl("https://danbooru.donmai.us/%ss.json?", st),
		params: JSON{
			"login":   viper.Get("danbooru_username"),
			"api_key": viper.Get("danbooru_api_key"),
			"limit":   limit,
		},
	}

	return s
}

// Initialize search object with common and gelbooru-specific values
func gelbooruSearch(limit int, st searchType) search {
	s := search{
		baseUrl: convertBaseUrl("https://gelbooru.com/index.php?page=dapi&q=index&s=%s", st),
		params: JSON{
			"user_id": viper.Get("gelbooru_user_id"),
			"api_key": viper.Get("gelbooru_api_key"),
			"limit":   limit,
			"json":    1,
		},
	}

	return s
}

// Initialize search object with common and konachan-specific values
func konachanSearch(limit int, st searchType) search {
	return search{
		baseUrl: convertBaseUrl("https://konachan.com/%s.json?", st),
		params: JSON{
			"limit": limit,
		},
	}
}

// Return a new search object for a given API and search type
func newSearch(a API, limit int, st searchType) search {
	switch a {
	case DANBOORU:
		return danbooruSearch(limit, st)
	case GELBOORU:
		return gelbooruSearch(limit, st)
	case KONACHAN:
		return konachanSearch(limit, st)
	default:
		panic("Invalid API value")
	}
}

// Build post search object for given API and parameters
func BuildPostSearch(api API, tags string, limit int, random bool) search {
	s := newSearch(api, limit, POST)
	s.params["random"] = random
	return s
}

// Build tag search object for given API and parameters
func BuildTagSearch(api API, tag string, limit int, order int) search {
	s := newSearch(api, limit, TAG)

	switch api {
	case DANBOORU:
		s.params["search[name_matches]"] = tag
		s.params["search[order]"] = order
	case GELBOORU:
		s.params["name_pattern"] = strings.Replace(tag, "*", "%", -1)
		s.params["order"] = order
	case KONACHAN:
		s.params["name"] = tag
		s.params["order"] = order
	}

	return s
}

// Make a search of a given type (post, tag, etc.), and return JSON data
func GetData(s search) ([]JSON, error) {
	searchUrl, err := encodeUrl(s)
	if err != nil {
		return nil, err
	}

	data, err := request(searchUrl)
	if err != nil {
		return nil, err
	}

	return data, nil
}
