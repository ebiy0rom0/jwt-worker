package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/pelletier/go-toml/v2"
)

type supabaseConfig struct {
	Secret string `toml:"secret"`
}
type config struct {
	Supabase supabaseConfig
}

var C config

func loadConfig() {
	path := filepath.Join(rootDir(), "config.toml")
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	decoder := toml.NewDecoder(file)
	if err := decoder.Decode(&C); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("load config: %+v", C)
}

func rootDir() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Dir(file)
}

func init() {
	loadConfig()
}
