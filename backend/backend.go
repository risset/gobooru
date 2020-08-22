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

type TagOrder int

const (
	DATE TagOrder = iota
	NAME
	COUNT
)

var apiUrls = map[API]string{
	DANBOORU: "https://danbooru.donmai.us/%ss.json?",
	GELBOORU: "https://gelbooru.com/index.php?page=dapi&q=index&s=%s",
	KONACHAN: "https://konachan.com/%s.json?",
}

type JSON map[string]interface{}

type Config struct {
	GelbooruAPIKey string
	GelbooruUserID string
	DanbooruAPIKey string
	DanbooruUserID string
}

// Create config file if it doesn't exist
func InitConfig(p string) error {
	cfgPath, err := homedir.Expand(p)
	if err != nil {
		return err
	}

	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(cfgPath)
	viper.SetDefault("danbooru_api_key", " ")
	viper.SetDefault("danbooru_user_id", " ")
	viper.SetDefault("gelbooru_api_key", " ")
	viper.SetDefault("gelbooru_user_id", " ")

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
func Request(url string) ([]JSON, error) {
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
func EncodeUrl(baseUrl string, params JSON) (string, error) {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return "", err
	}

	query := u.Query()
	for k, v := range params {
		query.Set(k, fmt.Sprintf("%v", v))
	}
	u.RawQuery = query.Encode()

	return u.String(), nil
}

// Format JSON map as pretty-printed JSON string
func FormatJSON(data JSON) (string, error) {
	b, err := json.MarshalIndent(data, "", "  ")

	if err != nil {
		return " ", err
	}

	return string(b), nil
}

// Print formatted JSON to standard output
func ShowJSON(data []JSON) error {
	for _, p := range data {
		s, err := FormatJSON(p)
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
	var dir string

	if random {
		dir = "random"
	} else {
		dir = strings.Replace(tags, " ", ",", -1)
	}

	return dir
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

// Make a search of a given type (post, tag, etc.), and return JSON data
func Search(searchType string, params JSON, api int) ([]JSON, error) {
	InitConfig("~/.config/gobooru")

	switch API(api) {
	case DANBOORU:
		params["api_key"] = viper.Get("danbooru_api_key")
		params["user_id"] = viper.Get("danbooru_user_id")
	case GELBOORU:
		params["api_key"] = viper.Get("gelbooru_api_key")
		params["user_id"] = viper.Get("gelbooru_user_id")
		params["json"] = 1
		// make gelbooru wildcards same as other APIs
		s := fmt.Sprintf("%v", params["name_pattern"])
		params["name_pattern"] = strings.Replace(s, "*", "%", -1)
	}

	baseUrl := fmt.Sprintf(apiUrls[API(api)], searchType)
	searchUrl, err := EncodeUrl(baseUrl, params)
	if err != nil {
		return nil, err
	}

	data, err := Request(searchUrl)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Differentiate between JSON tag search keys for different APIs
func BuildTagParams(api int, tags string, order int) JSON {
	params := JSON{}

	switch API(api) {
	case DANBOORU:
		params["search[name_matches]"] = tags
		params["search[order]"] = order
	case GELBOORU:
		params["name_pattern"] = tags
		params["orderby"] = order
	case KONACHAN:
		params["name"] = tags
		params["order"] = order
	default:
		return params
	}

	return params
}
