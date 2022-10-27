package mix

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"os"
	"path"
)

type Config struct {
	ManifestFile string
	PublicPath   string
}

func FuncMap(c Config) template.FuncMap {
	manifest := loadManifest(c.ManifestFile)
	return template.FuncMap{
		"assetPath": assetPathFunc(c.PublicPath, manifest),
	}
}

func assetPathFunc(publicPath string, manifest map[string]string) func(string) (string, error) {
	return func(asset string) (string, error) {
		p, ok := manifest[asset]
		if !ok {
			return "", fmt.Errorf("asset %s not found in manifest %s", asset, manifest)
		}
		if len(publicPath) > 0 {
			p = path.Join(publicPath, p)
		}
		return p, nil
	}
}

func loadManifest(p string) map[string]string {
	data, err := os.ReadFile(p)
	if err != nil {
		log.Fatalf("couldn't read %s: %s", p, err)
	}

	manifest := make(map[string]string)
	if err = json.Unmarshal(data, &manifest); err != nil {
		log.Fatalf("couldn't parse %s: %s", p, err)
	}
	return manifest
}
