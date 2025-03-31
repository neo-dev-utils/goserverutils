package algoUtil

import (
	"fmt"
	"gocommutils/randUtil"
	"math"
	"math/rand"
	"strconv"
)

const (
	CardsTypeBoomj_A  = 8  //J_A的豹子
	CardsTypeBoom2_10 = 7  //2-10的豹子
	CardTypeBZ        = 6  //豹子
	CardTypeSJ        = 5  //顺金
	CardTypeSZ        = 4  //顺子
	CardTypeJH        = 3  //金花
	CardTypeDZ        = 2  //对子
	CardTypeSingle    = 1  //单张
	CardTypeSingleQ_A = 0  //Q_A的单张
	CardTypeSingle9_J = -1 //9_J的单张
	CardTypeSingle7_8 = -2 //7_8的单张
	CardTypeSingle5_6 = -3 //5_6的单张

)

// 初始化手牌
func InitCards() (cards []int32) {
	count := 0
	for i := 2; i < 15; i++ {
		for j := 1; j < 5; j++ {
			pid := int32(i*100 + j)
			cards = append(cards, pid)
			count++
		}
	}

	return Shuffle(cards, count)
}

// 洗牌
func Shuffle(cards []int32, count int) []int32 {
	tmeprand := randUtil.Intn(5) + 1
	for i := 0; i < tmeprand; i++ {
		rand.Shuffle(count, func(i, j int) {
			cards[i], cards[j] = cards[j], cards[i]
		})
	}
	return cards
}

// 发牌
func Facards(cards []int32, handsnum int) ([]int32, []int32) {
	res := []int32{}
	for i := 0; i < handsnum; i++ {
		res = append(res, cards[0])
		cards = Unset(0, cards)
	}
	return res, cards
}

// 删除每个元素
func Unset(index int, A []int32) []int32 {
	if len(A)-1 <= index {
		return append(A[:index])

	}
	return append(A[:index], A[index+1:]...)
}

// 获取牌的颜色和牌值
func GetColorPrint(card int32) (cardcolor, cardprint int32) {
	return card % 100, card / 100
}

// 检测牌型
func CheckCardsType(hands []int32, changecardprint int32) (cardstype int32, card int32) {
	cards := []int32{}
	changecardnum := 0
	template := []int32{}
	for i := 0; i < len(hands); i++ {
		if _, cardsprint := GetColorPrint(hands[i]); changecardprint == cardsprint {
			changecardnum++
			continue
		}
		cards = append(cards, hands[i])
		template = append(template, hands[i])
	}
	// fmt.Printf("======开始=====cards: %v\n", cards)

	cardscolor, cardsprint := TypeColorPrint(cards)

	lencardsprint := len(cardsprint)
	lenCardsColor := len(cardscolor)
	isshunzi := CheckShunzi(cards)
	// fmt.Printf("======开始=====cards: %v\n", isshunzi)
	if lencardsprint == 1 {
		cardstype = CardTypeBZ
	} else if lencardsprint == 2 {
		cardstype = CardTypeDZ
		if isshunzi {
			cardstype = CardTypeSZ
		}

		if lenCardsColor == 1 && changecardnum > 0 {
			if isshunzi {
				cardstype = CardTypeSJ
			} else {
				cardstype = CardTypeJH
			}
		}

	} else {
		cardstype = CardTypeSingle
		if isshunzi {
			cardstype = CardTypeSZ
		}

		if lenCardsColor == 1 {
			if isshunzi {
				cardstype = CardTypeSJ
			} else {
				cardstype = CardTypeJH
			}
		}
	}

	if changecardnum == 1 {
		// fmt.Printf("======lalal=====cards: %v\n", cards)
		card = SupplyCards(template, cardstype)
	}

	return
}

// 根据颜色 、牌值进行分类
func TypeColorPrint(hands []int32) (map[int32][]int32, map[int32][]int32) {
	cardscolor := map[int32][]int32{}
	cardsprint := map[int32][]int32{}

	for i := 0; i < len(hands); i++ {
		cardcolor, cardprint := GetColorPrint(hands[i])
		cardscolor[cardcolor] = append(cardscolor[cardcolor], hands[i])
		cardsprint[cardprint] = append(cardsprint[cardprint], hands[i])
	}

	return cardscolor, cardsprint
}

// 检测顺子
func CheckShunzi(hands []int32) bool {
	// mincard, maxcard := MinMax(hands)
	// _, maxprint := GetColorPrint(maxcard)
	// _, minprint := GetColorPrint(mincard)
	SortCards(hands)
	handslen := len(hands)
	template := []int32{}
	for i := 0; i < handslen; i++ {
		_, printcards := GetColorPrint(hands[i])
		template = append(template, printcards)
	}

	if (template[0]-template[handslen-1] < 3 && handslen == 2) || (handslen == 3 && template[0]-template[1] == 1 && template[1]-template[2] == 1) {
		return true
	} // 正常的顺子

	if template[0] == 14 && template[handslen-1] <= 3 && handslen == 2 {
		return true
	} //带癞子的123
	return handslen == 3 && template[0] == 14 && template[1] == 3 && template[2] == 2
}

// 寻找最小值和最大值
func MinMax(c []int32) (min int32, max int32) {
	min = 9999
	for i := 0; i < len(c); i++ {
		if c[i] < min {
			min = c[i]
		}
		if c[i] > max {
			max = c[i]
		}
	}
	return
}

// 判断手牌中是否存在该牌
func HasCard(card int32, cards []int32) int {
	for key, v := range cards {
		if v == card {
			return key
		}
	}
	return -1
}

// 补充牌
func SupplyCards(cards []int32, types int32) (card int32) {
	cards = SortCards(cards)
	// fmt.Printf("cards: %v\n", cards)

	switch types {
	case CardTypeBZ:
		card = SupplyBzDz(cards)
	case CardTypeSJ:
		card = SupplySj(cards)
	case CardTypeJH:
		card = SupplyJh(cards)
	case CardTypeSZ:
		card = SupplySz(cards)
	case CardTypeDZ:
		card = SupplyBzDz(cards)
	}
	return card
}

// 将牌进行排序
func SortCards(cards []int32) []int32 {
	lencards := len(cards)
	for i := 0; i < lencards-1; i++ {
		for j := 0; j < (lencards - 1 - i); j++ {
			if (cards)[j] < (cards)[j+1] {
				cards[j], cards[j+1] = cards[j+1], cards[j]
			}
		}
	}
	return cards
}

// 升序
func RSortCards(cards []int32) []int32 {
	for i := 0; i < len(cards)-1; i++ {
		for j := 0; j < (len(cards) - 1 - i); j++ {
			if cards[j] > cards[j+1] {
				cards[j+1], cards[j] = cards[j], cards[j+1]
			}
		}
	}
	return cards
}

// 补充豹子和对子
func SupplyBzDz(cards []int32) (card int32) {
	_, cardsprint := GetColorPrint(cards[0])
	for i := 4; i > 0; i-- {
		card = cardsprint*100 + int32(i)
		if HasCard(card, cards) < 0 {
			break
		}
	}

	return card
}

// 补充顺金
func SupplySj(cards []int32) (card int32) {
	cha := cards[0] - cards[1]
	cardcolor, _ := GetColorPrint(cards[0])
	if cha == 100 {
		if cards[1] == 200+cardcolor {
			card = 1400 + cardcolor
		} else if cards[0]+100 > 1500 {
			card = cards[1] - 100
		} else {
			card = cards[0] + 100
		}
	} else if cha == 200 {
		card = cards[0] - 100
	} else {

		if cards[1] > 200 && cards[1] < 300 {
			card = 300 + cardcolor
		} else {
			card = 200 + cardcolor
		}
	}
	return card
}

// 补充金花
func SupplyJh(cards []int32) (card int32) {
	cardcolor, _ := GetColorPrint(cards[0])
	for i := 14; i > 1; i-- {
		card = int32(i)*100 + cardcolor
		if HasCard(card, cards) < 0 {
			break
		}
	}

	return card
}

// 补充顺子
func SupplySz(cards []int32) (card int32) {
	_, cardsprint := GetColorPrint(cards[0])
	_, cardsprint1 := GetColorPrint(cards[1])

	cha := cardsprint - cardsprint1
	if cha == 1 {
		if cardsprint1 == 2 {
			card = 1404
		} else if cardsprint+1 >= 15 {
			card = (cardsprint1-1)*100 + 4
		} else {
			card = (cardsprint+1)*100 + 4
		}
	} else if cha == 2 {
		card = (cardsprint-1)*100 + 4
	} else {
		if cardsprint1 > 2 {
			card = 204
		} else {
			card = 304
		}
	}
	return card
}

// 根据类型发牌
func TypeFa(types, handnum int, allcards []int32) (cards []int32) {
	switch types {
	case CardTypeBZ:
		cards = FaBz(allcards, handnum)
	case CardTypeSJ:
		cards = FaSj(allcards, handnum)
	case CardTypeSZ:
		cards = FaSz(allcards, handnum)
	case CardTypeJH:
		cards = FaJh(allcards, handnum)
	case CardTypeDZ:
		cards = FaDz(allcards, handnum)
	case CardTypeSingle:
		cards = FaSingle(allcards)
	case -3:
		cards = FaSingle6_5(allcards)
	case -2:
		cards = FaSingle8_7(allcards)
	case -1:
		cards = FaSinglej_9(allcards)
	case 0:
		cards = FaSingleA_Q(allcards)
	case 7:
		cards = FaBz2_10(allcards, handnum)
	case 8:
		cards = FaBzj_A(allcards, handnum)
	}
	return
}

// 发豹子
func FaBz(allcards []int32, handnum int) (cards []int32) {
	lenallcards := len(allcards)
	for i := 0; i < lenallcards; i++ {
		_, cardprint := GetColorPrint(allcards[i])
		cards = []int32{allcards[i]}
		for j := i + 1; j < lenallcards; j++ {
			_, cardprint1 := GetColorPrint(allcards[j])
			if cardprint1 == cardprint {
				cards = append(cards, allcards[j])
				if len(cards) >= handnum {
					break
				}
			}
		}
		if len(cards) >= handnum {
			break
		}
	}

	return
}

func FaBz2_10(allcards []int32, handnum int) []int32 {
	tempalcards := []int32{}
	for i := 0; i < len(allcards); i++ {
		_, print := GetColorPrint(allcards[i])
		if print <= 10 {
			tempalcards = append(tempalcards, allcards[i])
		}
	}

	return FaBz(tempalcards, handnum)
}

func FaBzj_A(allcards []int32, handnum int) []int32 {
	tempalcards := []int32{}
	for i := 0; i < len(allcards); i++ {
		_, print := GetColorPrint(allcards[i])
		if print > 10 {
			tempalcards = append(tempalcards, allcards[i])
		}
	}

	return FaBz(tempalcards, handnum)
}

// 发同花顺
func FaSj(allcards []int32, handnum int) (cards []int32) {
	for i := 0; i < len(allcards); i++ {
		color, pidprint := GetColorPrint(allcards[i])
		cards = []int32{allcards[i]}

		randcode := randUtil.Int31n(100)
		if handnum == 2 && (pidprint+2 <= 14 || pidprint == 14) && randcode >= 50 {
			pid1 := allcards[i] + 200
			if pid1 > 1400 {
				pid1 = 300 + color
			}

			if HasCard(pid1, allcards) < 0 {
				pid1 = allcards[i] + 100
				if pid1 > 1400 && pid1 < 1500 {
					pid1 = 200 + color
				}
				if HasCard(pid1, allcards) > -1 {
					cards = append(cards, pid1)
					break
				} else {
					pidprint++
				}
			} else {
				cards = append(cards, pid1)
				break
			}
		} else {
			for i := 1; i < handnum; i++ {
				pidprint1 := pidprint + int32(i)
				if pidprint1 == 15 && i == 1 {
					pidprint1 = 2
				}

				pid1 := pidprint1*100 + color
				if HasCard(pid1, allcards) > -1 {
					cards = append(cards, pid1)
				} else {
					pidprint++
					break
				}
			}
			if len(cards) >= handnum {
				break
			}
		}
	}

	return
}

// 发同花
func FaJh(allcards []int32, handnum int) (cards []int32) {
	lenallcards := len(allcards)
	for j := 0; j < lenallcards; j++ {
		cardcolor, cardprint := GetColorPrint(allcards[j])
		cards = []int32{allcards[j]}
		for i := j + 1; i < lenallcards; i++ {
			cardcolor1, cardprint1 := GetColorPrint(allcards[i])
			abs := math.Abs(float64(cardprint - cardprint1))
			// fmt.Printf("=====abs: %v %v %v\n", abs, cardprint1, cardprint)
			if cardcolor1 == cardcolor && (((cardprint1 == 14 || cardprint == 14) && abs > 2 && abs < 11) || (abs > 2 && (cardprint1 != 14 && cardprint != 14))) {
				cards = append(cards, allcards[i])
				if len(cards) >= handnum {
					break
				}
			}
		}

		if len(cards) >= handnum {
			break
		}
	}

	return
}

// 发对子
func FaDz(allcards []int32, handnum int) (cards []int32) {
	tempalcards := []int32{}
	tempalcards = append(tempalcards, allcards...)
	lenallcards := len(tempalcards)
	for i := 0; i < lenallcards; i++ {
		cardcolor, cardprint := GetColorPrint(tempalcards[i])
		cards = []int32{tempalcards[i]}
		for j := i + 1; j < lenallcards; j++ {
			cardcolor1, cardprint1 := GetColorPrint(tempalcards[j])
			abs := math.Abs(float64(cardprint - cardprint1))
			if handnum == 2 {
				if cardcolor1 != cardcolor && (abs > 2 || ((cardprint1 == 14 || cardprint == 14) && abs > 2 && abs < 11)) {
					cards = append(cards, tempalcards[j])
					break
				}
			} else {
				if len(cards) == 1 {
					if cardprint == cardprint1 {
						cards = append(cards, tempalcards[j])
					}
				} else if cardprint != cardprint1 {
					cards = append(cards, tempalcards[j])
					break
				}
			}
		}

		if len(cards) >= handnum {
			break
		}
	}

	return
}

// 发单牌
func FaSingle(allcards []int32) (cards []int32) {
	tempalcards := []int32{}
	tempalcards = append(tempalcards, allcards...)
	if len(tempalcards) <= 0 {
		fmt.Printf("==================allcards: %v\n", allcards)
		return
	}
	hand := tempalcards[0]
	tempalcards = Unset(0, tempalcards)
	cards = singleTool(tempalcards, hand)
	return
}

// 发单牌5-6
func FaSingle6_5(allcards []int32) []int32 {
	tempalcards := []int32{}
	for i := 0; i < len(allcards); i++ {
		_, print := GetColorPrint(allcards[i])
		if print <= 6 {
			tempalcards = append(tempalcards, allcards[i])
		}
	}
	return FaSingle(tempalcards)
}

// 发单牌7-8
func FaSingle8_7(allcards []int32) []int32 {
	tempalcards := []int32{}
	for i := 0; i < len(allcards); i++ {
		_, print := GetColorPrint(allcards[i])
		if print <= 8 {
			tempalcards = append(tempalcards, allcards[i])
		}
	}
	var hand int32
	index := -1
	for i := 0; i < len(tempalcards); i++ {
		_, print := GetColorPrint(tempalcards[i])
		if print >= 7 {
			hand = tempalcards[i]
			index = i
			break
		}
	}
	tempalcards = Unset(index, tempalcards)

	return singleTool(tempalcards, hand)
}

// 发单牌9-j
func FaSinglej_9(allcards []int32) []int32 {
	tempalcards := []int32{}
	for i := 0; i < len(allcards); i++ {
		_, print := GetColorPrint(allcards[i])
		if print <= 11 {
			tempalcards = append(tempalcards, allcards[i])
		}
	}

	var hand int32
	index := -1
	for i := 0; i < len(tempalcards); i++ {
		_, print := GetColorPrint(tempalcards[i])
		if print >= 9 {
			hand = tempalcards[i]
			index = i
			break
		}
	}
	tempalcards = Unset(index, tempalcards)

	return singleTool(tempalcards, hand)
}

// 发单牌Q-A
func FaSingleA_Q(allcards []int32) []int32 {
	tempalcards := []int32{}
	tempalcards = append(tempalcards, allcards...)
	var hand int32
	index := -1
	for i := 0; i < len(tempalcards); i++ {
		_, print := GetColorPrint(tempalcards[i])
		if print >= 12 {
			hand = tempalcards[i]
			index = i
			break
		}
	}
	tempalcards = Unset(index, tempalcards)

	return singleTool(tempalcards, hand)
}

func singleTool(tempalcards []int32, hand int32) (cards []int32) {
	cards = append(cards, hand)
	cardcolor, cardprint := GetColorPrint(hand)
	lenallcards := len(tempalcards)
	temp := []int32{cardprint}
	for i := 0; i < lenallcards; i++ {
		cardcolor1, cardprint1 := GetColorPrint(tempalcards[i])
		if cardcolor1 != cardcolor && cardprint != cardprint1 && HasCard(cardprint1, temp) < 0 {
			cards = append(cards, tempalcards[i])
			temp = append(temp, cardprint1)
			break
		}
	}

	for i := 1; i < lenallcards; i++ {
		if HasCard(tempalcards[i], cards) < 0 {
			cardcolor1, cardprint1 := GetColorPrint(tempalcards[i])
			tempCards := []int32{tempalcards[i]}
			tempCards = append(tempCards, cards...)
			cardsType, _ := CheckCardsType(tempCards, 0)
			if HasCard(cardprint1, temp) < 0 && cardcolor1 != cardcolor && cardsType == CardTypeSingle {
				cards = append(cards, tempalcards[i])
				break
			}
		}
	}
	return
}

// 手牌权重
func HandWeight(cardstype int32, hands []int32) (Weight int64) {
	hands = SortCards(hands)
	if cardstype == CardTypeDZ {
		return DuiziWeight(hands)
	}

	res := fmt.Sprintf("%d", cardstype)
	cardcolor, _ := GetColorPrint(hands[0])
	color := fmt.Sprintf("%d", cardcolor)
	arr14 := false
	arr3 := false

	if cardstype == CardTypeSJ || cardstype == CardTypeSZ {
		for i := 0; i < len(hands); i++ {
			_, cardprint := GetColorPrint(hands[i])
			if cardprint == 14 {
				arr14 = true
			}
			if cardprint < 4 {
				arr3 = true
			}
		}
	}
	temp := []int32{}
	for i := 0; i < len(hands); i++ {
		// cardcolor1, cardprint := GetColorPrint(hands[i])
		// if cardprint == 14 && arr14 && arr3 {
		// 	temp = append(temp, 100+cardcolor1)
		// } else {
		temp = append(temp, hands[i])
		// }
	}

	temp = SortCards(temp)

	for i := 0; i < len(temp); i++ {
		cardcolor1, cardprint := GetColorPrint(temp[i])
		// if arr14 && arr3 && cardprint == 14 {
		// 	cardprint = 1
		// }
		if arr14 && arr3 && cardprint == 2 {
			cardprint = 11
		}
		if arr14 && arr3 && cardprint == 3 {
			cardprint = 13
		}

		if cardprint/10 < 1 {
			res = fmt.Sprintf("%s%d", res, 0)
		}

		res = fmt.Sprintf("%s%d", res, cardprint)
		if i > 0 {
			color = fmt.Sprintf("%s%d", color, cardcolor1)
		}
	}
	res = fmt.Sprintf("%s%s", res, color)
	Weight, _ = strconv.ParseInt(res, 10, 64)
	return Weight
}

// 对子权重
func DuiziWeight(hands []int32) (Weight int64) {
	_, cardsprint := GetColorPrint(hands[1])
	res := fmt.Sprintf("%d", CardTypeDZ)
	for i := 0; i < 2; i++ {
		if cardsprint/10 < 1 {
			res = fmt.Sprintf("%s%d", res, 0)
		}
		res = fmt.Sprintf("%s%d", res, cardsprint)
	}

	var color string
	var dan int32

	for i := 0; i < len(hands); i++ {
		cardcolor1, cardprint1 := GetColorPrint(hands[i])
		if cardprint1 == cardsprint {
			color = fmt.Sprintf("%s%d", color, cardcolor1)
		} else {
			dan = hands[i]
		}
	}
	var cardcolor1 int32
	cardcolor1, cardsprint = GetColorPrint(dan)
	if cardsprint/10 < 1 {
		res = fmt.Sprintf("%s%d", res, 0)
	}
	res = fmt.Sprintf("%s%d", res, cardsprint)
	res = fmt.Sprintf("%s%s%d", res, color, cardcolor1)
	Weight, _ = strconv.ParseInt(res, 10, 64)
	return
}

// 发一副大牌
func BigCards(hands, allcards []int32, types int32, changecardprint int32, handnum int) []int32 {
	template := []int32{}
	for i := 0; i < len(hands); i++ {
		_, cardprint := GetColorPrint(hands[i])
		if cardprint != changecardprint {
			template = append(template, hands[i])
		}
	}
	_, max := MinMax(template)
	_, max = GetColorPrint(max)
	templateallcards := []int32{}
	for i := 0; i < len(allcards); i++ {
		_, cardprint := GetColorPrint(allcards[i])
		if cardprint > max {
			templateallcards = append(templateallcards, allcards[i])
		}
	}
	randcardstype := randUtil.RandInt(int(types), CardTypeBZ)
	return TypeFa(randcardstype, handnum, templateallcards)
}

// 发一副小牌
func MinCards(hands, allcards []int32, types int32, changecardprint int32, handnum int) []int32 {
	template := []int32{}
	for i := 0; i < len(hands); i++ {
		_, cardprint := GetColorPrint(hands[i])
		if cardprint != changecardprint {
			template = append(template, hands[i])
		}
	}
	min, _ := MinMax(template)
	_, min = GetColorPrint(min)
	templateallcards := []int32{}
	for i := 0; i < len(allcards); i++ {
		_, cardprint := GetColorPrint(allcards[i])
		if cardprint < min {
			templateallcards = append(templateallcards, allcards[i])
		}
	}
	randcardstype := randUtil.RandInt(1, int(types))
	return TypeFa(randcardstype, handnum, templateallcards)
}

func MaxCardsType(cards []int32, num int) (res []int32) {
	cardscolor, cardsprint := TypeColorPrint(cards)
	if temp := FoundMaxBoom(cardsprint, num); len(temp) >= num {
		res = temp
	} else if temp := FoundMaxShunJi(cardscolor, num); len(temp) >= num {
		res = temp
	} else if temp := FoundMaxShunzi(cards, num); len(temp) >= num {
		res = temp
	} else if temp := FoundMaxTongHua(cardscolor, num); len(temp) >= num {
		res = temp
	} else if temp := FoundMaxBoom(cardsprint, num-1); len(temp) >= num-1 {
		res = temp
		cards = SortCards(cards)
		for i := 0; i < len(cards); i++ {
			if HasCard(cards[i], res) < 0 {
				res = append(res, cards[i])
				break
			}
		}
	} else {
		res = FoundMaxDan(cards, num)
	}

	// fmt.Printf("==========res: %v\n", res)
	return
}

func MinCardsType(cards []int32, num int) (res []int32) {
	cardscolor, cardsprint := TypeColorPrint(cards)
	if temp := FoundMinDan(cards, num); len(temp) >= num {
		res = temp
	} else if temp := FoundMinBoom(cardsprint, num-1); len(temp) >= num-1 {
		res = temp
		cards = SortCards(cards)
		for i := 0; i < len(cards); i++ {
			if HasCard(cards[i], res) < 0 {
				res = append(res, cards[i])
				break
			}
		}
	} else if temp := FoundMinTongHua(cardscolor, num); len(temp) >= num {
		res = temp
	} else if temp := FoundMinShunzi(cards, num); len(temp) >= num {
		res = temp
	} else if temp := FoundMinShunJi(cardscolor, num-1); len(temp) >= num {
		res = temp
	} else {
		res = FoundMinBoom(cardsprint, num)

	}

	return
}

// 最大对子和炸弹
func FoundMaxBoom(cards map[int32][]int32, num int) []int32 {
	max := 0
	res := []int32{}
	for k, v := range cards {
		v = SortCards(v)
		if k > int32(max) && len(v) >= num {
			res = make([]int32, 0)
			for i := 0; i < num; i++ {
				res = append(res, v[i])
			}
			max = int(k)
		}
	}
	// fmt.Printf("res: %v num %v  cards %v\n", res, num, cards)
	return res
}

// 最大同花顺子
func FoundMaxShunJi(cards map[int32][]int32, num int) []int32 {
	res := []int32{}
	var Weight int64
	// fmt.Printf("cards: %v\n", cards)
	for _, v := range cards {
		v = SortCards(v)
		tempres := []int32{}
		for i := 1; i < len(v); i++ {
			index1 := HasCard(v[i]+100, v)
			pid := v[i] + 200
			// fmt.Printf("pid: %v\n", pid)
			if num == 2 {
				index2 := HasCard(pid, v)
				if index1 > -1 && HasCard(pid, tempres) < 0 {
					tempres = append(tempres, v[i]+100)
				} else if index2 > -1 && HasCard(pid, tempres) < 0 {
					tempres = append(tempres, v[i]+200)
				}

				if index2 < 0 && index1 < 0 {
					tempres = make([]int32, 0)
				}

				tempres = append(tempres, v[i])

				// fmt.Printf("tempres: %v\n", tempres)
			} else {
				if index1 > -1 && HasCard(pid, tempres) < 0 {
					tempres = append(tempres, v[i]+100)
				} else {
					if index1 < 0 {
						tempres = make([]int32, 0)
					}
					tempres = append(tempres, v[i])

				}

			}
			if len(res)+1 == num {
				temp := []int32{}
				temp = append(temp, res...)
				for j := 1401; j < 1405; j++ {
					if HasCard(int32(j), v) > -1 {
						temp = append(temp, int32(j))
						changeprint := 0
						if num < 3 {
							changeprint = 16
						}
						cardstyps, _ := CheckCardsType(temp, int32(changeprint))
						if cardstyps == CardTypeJH {
							res = append(res, int32(j))
							break
						}
					}
				}
			}

			if len(tempres) >= num {
				break
			}
		}

		if len(tempres) >= num {
			tempWeight := HandWeight(CardTypeSJ, tempres)
			if tempWeight > Weight {
				res = make([]int32, 0)
				res = tempres
				Weight = tempWeight
			}

		}
	}
	// fmt.Printf("=====顺金==res: %v\n", res)
	return res
}

// 最大顺子
func FoundMaxShunzi(cards []int32, num int) []int32 {
	res := []int32{}
	cards = SortCards(cards)
	for i := 1; i < len(cards); i++ {
		index := -1
		_, print := GetColorPrint(cards[i])
		var pid int32
		// fmt.Printf("cards: %v\n", cards[i])
		for j := 1; j <= 4; j++ {
			pid = (print+1)*100 + int32(j)
			tempindex := HasCard(pid, cards)
			// fmt.Printf("==========================/=====pid: %v %vres %v\n", pid, res, tempindex)
			index = tempindex
			if tempindex > -1 {
				// fmt.Printf("========pid: %v\n", pid)
				if HasCard(pid, res) < 0 {
					res = append(res, pid)
				}
				break
			}
		}

		if num == 2 {
			if index == -1 {
				for j := 1; j <= 4; j++ {
					pid = (print+2)*100 + int32(j)
					tempindex := HasCard(pid, cards)
					index = tempindex
					if tempindex > -1 {
						// fmt.Printf("========pid====: %v\n", pid)
						if HasCard(pid, res) < 0 {
							res = append(res, pid)
						}

						break
					}

				}
			}
		}
		if index < 0 {
			res = make([]int32, 0)
		}
		res = append(res, cards[i])

		if len(res)+1 == num {
			temp := []int32{}
			temp = append(temp, res...)
			for j := 1401; j < 1405; j++ {
				if HasCard(int32(j), cards) > -1 {
					temp = append(temp, int32(j))
					changeprint := 0
					if num < 3 {
						changeprint = 16
					}
					cardstyps, _ := CheckCardsType(temp, int32(changeprint))
					if cardstyps == CardTypeSZ {
						res = append(res, int32(j))
						break
					}
				}
			}
		}

		if len(res) >= num {
			break
		}
	}

	// fmt.Printf("===========最大顺子==res: %v\n", res)
	return res
}

// 最大同花
func FoundMaxTongHua(Cards map[int32][]int32, num int) []int32 {
	var Weight int64
	res := []int32{}
	for _, v := range Cards {
		v = SortCards(v)
		// fmt.Printf("v: %v\n", v)
		if len(v) >= num {
			template := []int32{}
			for i := 0; i < num; i++ {
				template = append(template, v[i])
			}
			tempWeight := HandWeight(CardTypeJH, template)
			// fmt.Printf("tempWeight: %v\n", tempWeight)
			if tempWeight > Weight {
				res = make([]int32, 0)
				res = template
				Weight = tempWeight
			}
		}
	}

	// fmt.Printf("=====同花==res: %v\n", res)

	return res
}

// 最大单张
func FoundMaxDan(cards []int32, num int) []int32 {
	cards = SortCards(cards)
	res := []int32{}
	for i := 0; i < num; i++ {
		res = append(res, cards[i])
	}
	return res
}

// 最小对子和炸弹
func FoundMinBoom(cards map[int32][]int32, num int) []int32 {
	max := 9999999999999
	res := []int32{}
	for k, v := range cards {
		v = RSortCards(v)
		if k < int32(max) && len(v) >= num {
			res = make([]int32, 0)
			for i := 0; i < num; i++ {
				res = append(res, v[i])
			}
			max = int(k)
		}
	}
	return res
}

// 最小同花顺子
func FoundMinShunJi(cards map[int32][]int32, num int) []int32 {
	res := []int32{}
	var Weight int64 = 99999999999
	for _, v := range cards {
		v = RSortCards(v)
		tempres := []int32{}
		for i := 1; i < len(v); i++ {
			index1 := HasCard(v[i]+100, v)
			if num == 2 {
				index2 := HasCard(v[i]+200, v)
				if index1 > -1 {
					tempres = append(tempres, v[i]+100)
				} else if index2 > -1 {
					tempres = append(tempres, v[i]+200)

				}
				if len(tempres) < num-1 {
					tempres = make([]int32, 0)
				}

				tempres = append(tempres, v[i])

			} else {
				if index1 > -1 {
					tempres = append(tempres, v[i]+100)
				} else {
					if len(tempres) < num-1 {
						tempres = make([]int32, 0)
					}

					tempres = append(tempres, v[i])

				}

			}

			if len(tempres) >= num {
				break
			}
		}

		if len(tempres) >= num {
			tempWeight := HandWeight(CardTypeSJ, res)
			if tempWeight < Weight {
				res = make([]int32, 0)
				res = tempres
				Weight = tempWeight
			}

		}
	}

	return res
}

// 最小顺子
func FoundMinShunzi(cards []int32, num int) []int32 {
	res := []int32{}
	cards = RSortCards(cards)
	for i := 1; i < len(cards); i++ {
		index := -1
		_, print := GetColorPrint(cards[i])
		var pid int32
		for j := 1; j <= 4; j++ {
			pid = (print+1)*100 + int32(j)
			tempindex := HasCard(pid, cards)
			if tempindex > -1 && HasCard(pid, res) < 0 {
				res = append(res, pid)
			}
			index = tempindex
		}

		if num == 2 {
			if index == -1 {
				for j := 1; j <= 4; j++ {
					pid = (print+2)*100 + int32(j)
					tempindex := HasCard(pid, cards)
					if tempindex > -1 && HasCard(pid, res) < 0 {
						res = append(res, pid)
					}
					index = tempindex
				}
			}
		}

		if len(res) < num-1 {
			res = make([]int32, 0)
		}

		res = append(res, cards[i])
		if len(res) >= num {
			break
		}
	}

	return res
}

// 最小同花
func FoundMinTongHua(Cards map[int32][]int32, num int) []int32 {
	var Weight int64 = 999999999999
	res := []int32{}
	for _, v := range Cards {
		v = RSortCards(v)
		if len(v) >= num {
			template := []int32{}
			for i := 0; i < num; i++ {
				template = append(template, v[i])
			}
			tempWeight := HandWeight(CardTypeJH, template)
			if tempWeight < Weight {
				res = make([]int32, 0)
				res = template
				Weight = tempWeight
			}
		}
	}
	return res
}

// 最小单张
func FoundMinDan(cards []int32, num int) []int32 {
	cards = RSortCards(cards)
	res := []int32{}
	changeprint := 0
	if num < 3 {
		changeprint = 16
	}
	lencards := len(cards)
	for i := 1; i < lencards; i++ {
		res = []int32{cards[0], cards[i]}
		code := false
		for j := i + 1; j < lencards; j++ {
			temp := make([]int32, 0)
			temp = append(res, cards[j])
			if num < 3 {
				temp = append(temp, 1605)
			}
			template, _ := CheckCardsType(temp, int32(changeprint))
			if template == CardTypeSingle {
				res = make([]int32, 0)
				res = temp
				code = true
				break
			}
		}

		if code {
			break
		}
	}
	return res
}

func DetailedSingleCardsType(cards []int32) (types int32) {
	var max int32
	for i := 0; i < len(cards); i++ {
		if cards[i] > max {
			max = cards[i]
		}
	}
	_, max = GetColorPrint(max)
	if max > 11 {
		types = CardTypeSingleQ_A
	} else if max > 8 {
		types = CardTypeSingle9_J
	} else if max > 6 {
		types = CardTypeSingle7_8
	} else {
		types = CardTypeSingle5_6
	}
	return
}

func DetailedBoomCardsType(cards []int32) (types int32) {
	var max int32
	for i := 0; i < len(cards); i++ {
		if cards[i] > max {
			max = cards[i]
		}
	}
	_, max = GetColorPrint(max)
	if max > 10 {
		types = CardsTypeBoomj_A
	} else {
		types = CardsTypeBoom2_10
	}
	return
}

func GetWeight(hands []int32) int64 {
	hands = SortCards(hands)
	res := ""
	for i := 0; i < len(hands); i++ {
		_, cardsprint := GetColorPrint(hands[i])
		if cardsprint/10 < 1 {
			res = fmt.Sprintf("%s%d", res, 0)
		}
		res = fmt.Sprintf("%s%d", res, cardsprint)
	}
	res1, _ := strconv.ParseInt(res, 10, 64)
	return res1
}

// 发顺子
// func FaSz(allcards []int32, handnum int) (cards []int32) {
// 	for w := 0; w < len(allcards); w++ {
// 		color, pidprint := GetColorPrint(allcards[w])
// 		cards = []int32{allcards[w]}
// 		var i, j int32
// 		randcode := rr.Int31n(100)
// 		color1 := color + 1
// 		if handnum == 2 && (pidprint+2 <= 14 || pidprint == 14) && randcode <= 5 {
// 			pidprint1 := pidprint + 2
// 			if pidprint1 == 15 {
// 				pidprint1 = 3
// 			}
// 			code := false
// 			for j = 1; j < 4; j++ {
// 				pid1 := pidprint1*100 + color1
// 				if HasCard(pid1, allcards) > -1 {
// 					code = true
// 					cards = append(cards, pid1)
// 					break
// 				}
// 				color1++
// 				if color1 > 4 {
// 					color1 = 1
// 				}
// 			}

// 			if code {
// 				break
// 			}
// 		} else {
// 			for i = 1; i < int32(handnum); i++ {
// 				pidprint1 := pidprint + i
// 				if i == 1 && pidprint1 == 15 {
// 					pidprint1 = 2
// 					pidprint = 1
// 				}

// 				for j = 1; j < 4; j++ {
// 					pid1 := pidprint1*100 + color1
// 					if HasCard(pid1, allcards) > -1 {
// 						cards = append(cards, pid1)
// 						break
// 					}
// 					color1++
// 					if color1 > 4 {
// 						color1 = 1
// 					}
// 				}
// 				if len(cards) >= handnum {
// 					break
// 				}
// 			}

// 			if len(cards) < handnum {
// 				for i = 1; i <= int32(handnum-len(cards)); i++ {
// 					pidprint1 := pidprint - i
// 					code := false

// 					for j = 1; j < 4; j++ {
// 						pid1 := pidprint1*100 + color1
// 						if HasCard(pid1, allcards) > -1 {
// 							code = true
// 							cards = append(cards, pid1)
// 							break
// 						}
// 						color1++
// 						if color1 > 4 {
// 							color1 = 1
// 						}
// 					}

// 					if !code {
// 						break
// 					}
// 				}

// 			} else {
// 				break
// 			}
// 		}
// 	}

// 	return
// }

func FaSz(allcards []int32, handnum int) (cards []int32) {
	var color, print int32
	isSort := -1
	cards, isSort = tempFaSz(allcards, cards, color, print, isSort, handnum)
	// fmt.Printf("cards: %v\n", cards)
	if len(cards) < handnum {
		min, _ := MinMax(cards)
		color, print = GetColorPrint(min)
		// fmt.Printf("=============================cards: %v min %v\n", cards, min)
		isSort = 1
		if print == 2 {
			print = 3
			isSort = 2
		}

		cards, _ = tempFaSz(allcards, cards, color, print, isSort, handnum)
	}
	return
}
func tempFaSz(allcards, cards []int32, color, print int32, isSort, handnum int) ([]int32, int) {
	// fmt.Printf("isSort: %v\n", isSort)
	for i := 0; i < len(allcards); i++ {
		color1, print1 := GetColorPrint(allcards[i])
		lencards := len(cards)
		if lencards <= 0 {
			cards = append(cards, allcards[i])
			color, print = color1, print1
			continue
		}

		if (print1+1 == print && (isSort == 1 || isSort == -1)) || (print1-1 == print && (isSort == -1 || isSort == 2)) {
			if print1+1 == print && (isSort == 1 || isSort == -1) {
				isSort = 1
			} else if print1-1 == print && (isSort == -1 || isSort == 2) {
				isSort = 2
			}

			if color1 == color && (handnum-lencards == 1) {
				continue
			}
			cards = append(cards, allcards[i])
			color, print = color1, print1
			if len(cards) >= handnum {
				// fmt.Println("============================================1")
				break
			}
			continue
		}

		// fmt.Printf("isSort: %v\n", isSort)
		if isSort == -1 && ((print == 14 && (print1 == 2 || print1 == 3)) || (print1 == 14 && (print == 2 || print == 3))) {
			if handnum == 2 {
				if color1 != color {
					cards = append(cards, allcards[i])
				} else {
					for j := 0; j < len(allcards); j++ {
						tempColor1, tempPrint1 := GetColorPrint(allcards[j])
						if tempColor1 != color && tempPrint1 == print1 {
							cards = append(cards, allcards[j])
							break
						}
					}
				}
				// fmt.Println("============================================2")

				break
			} else if handnum == 3 {
				cards = append(cards, allcards[i])
				temp := []int32{1602}
				temp = append(temp, cards...)
				_, tempcard := CheckCardsType(temp, 16)
				_, tempPrint := GetColorPrint(tempcard)
				for j := 0; j < len(allcards); j++ {
					tempColor1, tempPrint1 := GetColorPrint(allcards[j])
					if tempColor1 != color && tempPrint1 == tempPrint {
						cards = append(cards, allcards[j])
						break
					}
				}
			}
			break
		}
	}
	return cards, isSort
}
