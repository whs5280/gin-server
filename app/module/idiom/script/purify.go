package main

import (
	"encoding/json"
	"fmt"
	"gin-server/app/module/idiom/model"
	"os"
	"path/filepath"
)

func main() {
	absPath, err := filepath.Abs(model.IdiomDBPath)
	if err != nil {
		panic(err)
	}

	fileContent, err := os.ReadFile(absPath)
	if err != nil {
		fmt.Printf("读取文件失败: %v\n", err)
		os.Exit(1)
	}

	var idioms []model.Idiom
	if err := json.Unmarshal(fileContent, &idioms); err != nil {
		fmt.Printf("解析JSON失败: %v\n", err)
		os.Exit(1)
	}

	// 提纯
	var words []string
	for _, idiom := range idioms {
		words = append(words, idiom.Word)
	}
	purified := model.PurifiedData{
		Idioms: words,
	}

	// 生成新的JSON文件
	outputData, err := json.MarshalIndent(purified, "", "  ")
	if err != nil {
		fmt.Printf("生成JSON失败: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(model.PurifiedPath, outputData, 0644); err != nil {
		fmt.Printf("写入文件失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("成功提取 %d 个成语到 %s\n", len(words), model.PurifiedPath)
}
