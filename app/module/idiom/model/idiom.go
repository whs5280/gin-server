package model

const (
	FAIL           = "认输"
	EXIT           = "退出"
	FailMessage    = "You lose, game over!"
	ExitMessage    = "Bye!"
	ValidMessage   = "Invalid words, input again"
	NextMessage    = "我出：%s → 请接以【%s】开头的成语：\n"
	ComFailMessage = "我认输了，没有可以接的成语了!"
	ComNextMessage = "我接：%s → 请你接【%s】开头的成语：\n"
	IdiomDBPath    = "../data/idiom.json"
	PurifiedPath   = "../data/idiom_purified.json"
)

type Idiom struct {
	Word         string `json:"word"`
	Pinyin       string `json:"pinyin"`
	Explanation  string `json:"explanation"`
	Derivation   string `json:"derivation"`
	Example      string `json:"example"`
	Abbreviation string `json:"abbreviation"`
}

type PurifiedData struct {
	Idioms []string `json:"idioms"`
}
