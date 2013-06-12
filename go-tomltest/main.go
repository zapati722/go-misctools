package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
)

func test() map[string]interface{} {
	return map[string]interface{}{
		"pkg": map[string]interface{}{
			"title": "jQuery",
			"desc":  "jQuery is a fast, small, and feature-rich JavaScript library. It makes things like HTML document traversal and manipulation, event handling, animation, and Ajax much simpler with an easy-to-use API that works across a multitude of browsers.",
			"www":   "http://jquery.com",
		},
		"default": map[string]interface{}{
			"js": []interface{}{"jquery"},
		},
		"legacy": map[string]interface{}{
			"title":      "Legacy",
			"desc":       "Serves 1.x (instead of newer versions) to MSIE versions older than 9.0",
			"version_lt": 2,
			"auto": map[string]interface{}{
				"header_regexp": []interface{}{"User-Agent", "^Mozilla/\\d+\\.\\d+ \\([^(]*MSIE [1-8]\\.[0-9][^)]*\\).*"}}}}
}

func main() {
	m := map[string]interface{}{}
	_, err := toml.DecodeFile("Q:\\d\\src\\github.com\\openbase\\ob-build\\hive-default\\dist\\pkg\\webuilib-jquery\\jquery.webuilib.obpkg", m)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v", m)
}
