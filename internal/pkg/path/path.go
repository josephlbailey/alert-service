package path

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func Determine(name string) (string, error) {

	for i := 0; i <= 4; i++ {
		path, err := os.Getwd()

		if err != nil {
			return "", err
		}

		dirContents, err := os.ReadDir(path)

		for _, dir := range dirContents {
			if dir.Name() == name {
				path, _ = strings.CutSuffix(path, "/")
				log.Printf("found %v, path: %v/%v\n", name, path, name)
				return path, nil
			}
		}

		if err := os.Chdir("../"); err != nil {
			return "", err
		}
	}

	return "", fmt.Errorf("%v directory not found in parent directories", name)
}
