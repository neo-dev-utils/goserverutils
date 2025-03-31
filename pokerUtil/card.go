package pokerUtil

import (
	"fmt"
	"gocommutils/randUtil"
	"gocommutils/timeUtil"
)

// 黑红梅方
var Deck = []byte{
	0x21, 0x31, 0x41, 0x51, 0x61, 0x71, 0x81, 0x91, 0xa1, 0xb1, 0xc1, 0xd1, 0xe1,
	0x22, 0x32, 0x42, 0x52, 0x62, 0x72, 0x82, 0x92, 0xa2, 0xb2, 0xc2, 0xd2, 0xe2,
	0x23, 0x33, 0x43, 0x53, 0x63, 0x73, 0x83, 0x93, 0xa3, 0xb3, 0xc3, 0xd3, 0xe3,
	0x24, 0x34, 0x44, 0x54, 0x64, 0x74, 0x84, 0x94, 0xa4, 0xb4, 0xc4, 0xd4, 0xe4,
}

type GamePoker struct {
	Cards []byte
}

func (gp *GamePoker) InitPoker() {
	for _, v := range Deck {
		gp.Cards = append(gp.Cards, v)
	}
}

func (gp *GamePoker) ShuffleCards() {
	for i := 0; i < len(gp.Cards); i++ {
		index := randUtil.Intn(len(gp.Cards))
		gp.Cards[i], gp.Cards[index] = gp.Cards[index], gp.Cards[i]
	}
}

func (gp *GamePoker) DealCards() byte {
	card := gp.Cards[0]
	gp.Cards = append(gp.Cards[:0], gp.Cards[1:]...)
	return card
}

func (gp *GamePoker) GetCardsCount() int {
	return len(gp.Cards)
}

// 测试牌型编码效率
func testCardEncode() {
	var cards []byte
	start := timeUtil.Time.NowTime()
	count := 0
	for a := 0; a < 52-4; a++ {
		for b := a + 1; b < 52-3; b++ {
			for c := b + 1; c < 52-2; c++ {
				for d := c + 1; d < 52-1; d++ {
					for e := d + 1; e < 52; e++ {

						cards = []byte{Deck[a], Deck[b], Deck[c], Deck[d], Deck[e]}
						//encodeCard(getCardType(cards))
						count++
					}
				}
			}
		}
	}
	end := timeUtil.Time.NowTime()
	fmt.Println("cards : ", cards)
	fmt.Println(end.Sub(start), count)
}

// 随便输入一个数字，得出它的牌值和花色
func GetCardValueAndColor(value byte) (cardValue, cardColor byte) {
	cardValue = (value & 0xf0) //byte的高4位总和是240
	cardColor = value & 0xf    //byte的低4位总和是15
	return
}

// 将牌 进行牌型编码并返回
func GetEncodeCard(cardType int, cards []byte) (cardEncode int) {
	cardEncode = (cardType) << 20
	if cardType != CardTypeDZ {
		num := 0
		for i, card := range cards {
			cardEncode |= (int(card) >> 4) << uint((5-i-1)*4)
			_, color := GetCardValueAndColor(card)
			temp := 100
			if i == 1 {
				temp = 10
			} else if i == 2 {
				temp = 1
			}

			num += (int(color) * temp)
		}

		cardEncode += num
	} else {
		dui := cards[0]
		dan := cards[0]
		duizi := cards[0]
		if cards[1]&0xf0 == cards[2]&0xf0 {
			dui = cards[1]
			duizi = cards[1]
		} else {
			dan = cards[2]
		}

		cardEncode |= (int(dui) >> 4) << uint((4)*4)
		cardEncode |= (int(dui) >> 4) << uint((3)*4)
		cardEncode |= (int(dan) >> 4) << uint((2)*4)
		tempnum, _ := GetCardValueAndColor(duizi)
		num := 0
		j := 0
		for i := 0; i < len(cards); i++ {
			tempnum1, color1 := GetCardValueAndColor(cards[i])
			if tempnum1 == tempnum {
				j++
				temp := 100
				if j == 1 {
					temp = 10
				}

				num += (temp * int(color1))
			}
		}
		cardEncode += num
	}

	return
}

// 生成3张牌 TODO 暂时先随便生成3张牌，后面要根据相应的策略来生成牌
func GenerateCards() (cards []byte) {
	cards = make([]byte, 3)
	c0 := Deck[randUtil.Intn(52)]
	c1 := Deck[timeUtil.Time.NowUnixNano()%52]
	c2 := Deck[(timeUtil.Time.NowUnixNano()+123)%52]
	cards[0] = c0
	cards[1] = c1
	cards[2] = c2
	return
}

func Test(cards []byte) {
	sortRes := sortCards(cards)
	fmt.Println("sortRes : ", sortRes)
	t, _ := GetCardTypeJH(sortRes)
	fmt.Println("cardsType : ", t)
}
