package curUtil

import "fmt"

// // 货币转换
// func GetScoreStr(score int64) string {
// 	yuan := score / 100
// 	remain := score % 100
// 	if remain < 0 {
// 		remain = -remain
// 	}
// 	jiao := remain / 10
// 	fen := remain % 10
// 	return fmt.Sprintf(`%d.%d%d`, yuan, jiao, fen)
// }

// 分数转货币字符串
func ScoreToStrCur(score int64) string {
	yuan := score / 100
	remain := score % 100
	if remain < 0 {
		remain = -remain
	}
	jiao := remain / 10
	fen := remain % 10
	return fmt.Sprintf(`%d.%d%d`, yuan, jiao, fen)
}
