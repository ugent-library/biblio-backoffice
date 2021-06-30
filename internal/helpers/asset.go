package helpers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"path"
)

func Asset() template.FuncMap {
	manifest := loadAssetManifest("static/mix-manifest.json")

	return template.FuncMap{
		"assetPath": assetPathFunc(manifest),
	}
}

func assetPathFunc(manifest map[string]string) func(string) (string, error) {
	return func(asset string) (string, error) {
		p, ok := manifest[asset]
		if !ok {
			err := fmt.Errorf("Asset %s not found in manifest %s", asset, manifest)
			log.Println(err)
			return "", err
		}
		return path.Join("/static/", p), nil
	}
}

func loadAssetManifest(p string) map[string]string {
	data, err := ioutil.ReadFile(p)
	if err != nil {
		log.Fatalf("Couldn't read %s: %s", p, err)
	}
	manifest := make(map[string]string)
	if err = json.Unmarshal(data, &manifest); err != nil {
		log.Fatalf("Couldn't parse %s: %s", p, err)
	}
	return manifest
}
