# booruget
A simple CLI program for interacting with the REST API of various
booru sites, for downloading images, listing tags, etc. 

Currently supports:

[Danbooru](https://danbooru.donmai.us)

[Gelbooru](https://gelbooru.com)

[Konachan](https://konachan.com)

## Usage
gobooru <post/tag> tags flags

- Tags can be separated by spaces without quotes in the argument, queries using wildcard characters must be in quotes however
- Danbooru has a limit of 2 tags per request, unless you have an upgraded account
- Gelbooru sometimes requires authentication for requests. API keys and User IDs for accounts can be placed in ~/.config/gobooru/config.yaml

https://danbooru.donmai.us/wiki_pages/help%3Aapi

https://gelbooru.com/index.php?page=wiki&s=view&id=18780

https://konachan.com/help/api

## Examples
### Download 200 images that match the specified tags from gelbooru
```bash
gobooru post patchouli_knowledge rating:safe 1girl -a 1 -n 200 -d
```

### Download 20 random images from danbooru
```bash
gobooru post -r -n 20 -d
```

### Look up a given tag pattern on konachan, sort by post count
```bash
gobooru tag "patch*" -a 2 -o 2
```
