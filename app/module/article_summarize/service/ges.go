package service

import (
	"bufio"
	"fmt"
	"gin-server/app/module/article_summarize/model"
	"github.com/go-ego/gse"
	"os"
	"sort"
	"sync"
	"time"
)

var (
	seg       gse.Segmenter
	stopWords map[string]bool
)

func init() {
	seg, _ = gse.New()
	err := seg.LoadDict()
	if err != nil {
		return
	}
	// 当前版本 gse version: v0.80.2.705，停用词处理出现问题，手动过滤 filterStopWords

	/*stopWords := []string{
		"的", "了", "和", "是", "，", "。", "、", "“", "”", "‘", "’", "？", "！",
		"：", "；", "（", "）", "【", "】", "……", "——", "·",
	}
	seg.LoadStopArr(stopWords)*/

	stopWords = map[string]bool{
		"的": true, "了": true, "，": true, "。": true, "？": true, "、": true, "“": true, "”": true, "是": true,
	}
}

// 分块读取：避免一次性加载大文本导致内存溢出。
// 并发分词：使用多个 Goroutine 并行处理文本块。
// 合并结果：汇总各 Goroutine 的词频统计结果。

func SegmentGoroutine(filePath string, top int) []model.WordFrequency {
	startTime := time.Now()

	file, err := os.Open(filePath)
	if file == nil {
		panic(err)
	}
	defer file.Close()

	// 创建 Channel 用于分发文本行
	lines := make(chan string, 1000) // 缓存队列
	var wg sync.WaitGroup

	numberWorkers := 4
	results := make(chan map[string]int, numberWorkers)

	for i := 0; i < numberWorkers; i++ {
		wg.Add(1)
		go chuckSegment(&seg, lines, results, &wg)
	}

	// 读取文件并发送到通道
	scanner := bufio.NewScanner(file)
	go func() {
		for scanner.Scan() {
			lines <- scanner.Text()
		}
		close(lines)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	// 整合worker结果
	wordFreq := make(map[string]int)
	for freqMap := range results {
		for word, count := range freqMap {
			wordFreq[word] += count
		}
	}

	freqList := sortWordFrequency(wordFreq)
	outputResult(freqList, top)
	fmt.Println("耗时：", time.Since(startTime))

	return freqList
}

// chuckSegment 处理文本块并统计词频
func chuckSegment(seg *gse.Segmenter, lines <-chan string, results chan<- map[string]int, wg *sync.WaitGroup) {
	defer wg.Done()
	freq := make(map[string]int)

	for line := range lines {
		words := seg.Cut(line, true)
		words = filterStopWords(words)
		for _, word := range words {
			freq[word]++
		}
	}
	results <- freq
}

// Segment 分词并统计词频
func Segment(text string, top int) []model.WordFrequency {
	startTime := time.Now()

	// 分词（精确模式）
	words := seg.Cut(text, true)
	words = filterStopWords(words)

	// 词频统计
	wordFreq := make(map[string]int)
	for _, word := range words {
		wordFreq[word]++
	}

	freqList := sortWordFrequency(wordFreq)
	outputResult(freqList, top)
	fmt.Println("耗时：", time.Since(startTime))

	return freqList
}

// 排序词频（从大到小）
func sortWordFrequency(freqMap map[string]int) []model.WordFrequency {
	var freqList []model.WordFrequency
	for word, count := range freqMap {
		if count > 2 {
			freqList = append(freqList, model.WordFrequency{Word: word, Count: count})
		}
	}

	sort.Slice(freqList, func(i, j int) bool {
		return freqList[i].Count > freqList[j].Count
	})

	return freqList
}

// 输出词频统计结果
func outputResult(sortedFreq []model.WordFrequency, top int) {
	fmt.Println("词频统计（从大到小）：")
	for key, item := range sortedFreq {
		if key < top-1 {
			fmt.Printf("%s: %d\n", item.Word, item.Count)
		}
	}
}

// 过滤停用词
func filterStopWords(words []string) []string {
	var filtered []string
	for _, word := range words {
		if !stopWords[word] {
			filtered = append(filtered, word)
		}
	}
	return filtered
}
