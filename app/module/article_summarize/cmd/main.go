package main

import (
	"gin-server/app/module/article_summarize/model"
	"gin-server/app/module/article_summarize/service"
	"path/filepath"
)

func main() {
	/* content := "一把摇摇椅、一个原木色置物架，还有几个靠枕、防风烛台、铁艺花盆，再搭上一条柔软的休闲毯、几盆绿植鲜花和小木桌，你就有了城市中属于自己的一片小花园，不宽敞但却精致，感兴趣的朋友不妨看看下面的灵感图吧，在这样的阳台读读书、小酌谈天，你不想吗？"
	service.Segment(content, 5) */

	absPath, err := filepath.Abs(model.TestFilePath)
	if err != nil {
		panic(err)
	}

	_ = service.SegmentGoroutine(absPath, 5)
}
