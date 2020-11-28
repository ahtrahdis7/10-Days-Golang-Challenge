package urlshort

import (
	yaml "gopkg.in/yaml.v2"
	"net/http"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}
		// if we can match a path ...
		//  redirect
		// else
		fallback.ServeHTTP(w, r)
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yamlBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	// 1: parse the YAML
	var pathUrls []pathUrl
	err := yaml.Unmarshal(yamlBytes, &pathUrls)

	if err != nil {
		return nil, err
	}
	// 2. Convert YAML array into map
	pathsToUrls := make(map[string]string)

	for _, pu := range pathUrls {
		pathsToUrls[pu.Path] = pu.URL
	}
	// 3. return map handler using map

	// pathsToUrls := map[string]string{}
	return MapHandler(pathsToUrls, fallback), nil

}

type pathUrl struct {
	Path string `yaml:"path"`
	URL string `yaml:"url"`
}