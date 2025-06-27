package service

import (
	"encoding/json"
	"fmt"
	"gin-server/app/module/idiom/model"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

var idiomDB = map[string]bool{}

func InitIdiomDB() error {
	absPath, err := filepath.Abs(model.PurifiedPath)
	if err != nil {
		return err
	}

	fileContent, err := os.ReadFile(absPath)
	if err != nil {
		return err
	}

	var data struct {
		Idioms []string `json:"Idioms"`
	}
	if err := json.Unmarshal(fileContent, &data); err != nil {
		return err
	}

	idiomDB = make(map[string]bool)
	for _, idiom := range data.Idioms {
		idiomDB[idiom] = true
		/*if utf8.RuneCountInString(idiom) == 4 {
			idiomDB[idiom] = true
		} else {
			fmt.Printf("警告: 忽略非四字成语 '%s'\n", idiom)
		}*/
	}

	fmt.Printf("成功加载 %d 个成语\n", len(idiomDB))
	return nil
}

// GetLastChar 获取成语最后一个字
func GetLastChar(idiom string) string {
	r, _ := utf8.DecodeLastRuneInString(idiom)
	return string(r)
}

// IsValidIdiom 检测是否为成语
func IsValidIdiom(idiom string) bool {
	if utf8.RuneCountInString(idiom) != 4 {
		return false
	}
	_, exists := idiomDB[idiom]
	return exists
}

// FindMatchingIdioms 根据首字查找可接的成语
func FindMatchingIdioms(firstChar string, usedIdioms map[string]bool) []string {
	var matches []string
	for idiom := range idiomDB {
		if strings.HasPrefix(idiom, firstChar) && !usedIdioms[idiom] {
			matches = append(matches, idiom)
		}
	}
	return matches
}
