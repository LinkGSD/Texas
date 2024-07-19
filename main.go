package main

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
)

var (
	Colors = []string{"♠", "♥", "♦", "♣"}
	Pokes  = make([]Poke, 0)
)

func init() {
	for i := 2; i <= 14; i++ {
		for _, color := range Colors {
			if i > 10 {
				if i == 11 {
					Pokes = append(Pokes, Poke{Color: color, Number: i, Val: "J"})
				} else if i == 12 {
					Pokes = append(Pokes, Poke{Color: color, Number: i, Val: "Q"})
				} else if i == 13 {
					Pokes = append(Pokes, Poke{Color: color, Number: i, Val: "K"})
				} else {
					Pokes = append(Pokes, Poke{Color: color, Number: i, Val: "A"})
				}
			} else {
				Pokes = append(Pokes, Poke{Color: color, Number: i, Val: strconv.Itoa(i)})
			}
		}
	}
}

type Poke struct {
	Color  string `json:"color"`
	Val    string `json:"value"`
	Number int    `json:"number"`
}

func (p Poke) String() string { return fmt.Sprintf("{%s %s}", p.Color, p.Val) }

func (p Poke) Compare(p2 Poke) string {
	if p.Number == 1 {
		p.Number += 13
	}
	if p2.Number == 1 {
		p2.Number += 13
	}
	return CountCompare(p.Number, p2.Number)
}

type Room struct {
	Id        string    `json:"id"`
	MaxPlayer int       `json:"max_player"`
	Players   []*Player `json:"players"`
	Mang      int       `json:"mang"`
	Pokes     []Poke    `json:"pokes"`
	MidPokers []Poke    `json:"mid_pokers"`
	MidPos    int
}

type Player struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Pokes     []Poke `json:"pokes"`
	Chips     int    `json:"chips"`
	PokeCheck int    `json:"check"`
}

func (p *Player) PokeCount(start, end int) int {
	cnt := 0
	for i := start; i < end; i++ {
		cnt += p.Pokes[i].Number
	}
	return cnt
}

func (p *Player) MaxPokes(num, ect1, ect2 int) []Poke {
	pokes := make([]Poke, 0)
	for i := len(p.Pokes) - 1; i >= 0; i-- {
		if p.Pokes[i].Number != ect1 {
			if ect2 != 0 && p.Pokes[i].Number == ect2 {
				continue
			}
			pokes = append(pokes, p.Pokes[i])
		}
		if len(pokes) == num {
			return pokes
		}
	}
	return nil
}

func CountCompare(c1, c2 int) string {
	if c1 > c2 {
		return "large"
	} else if c1 < c2 {
		return "small"
	} else {
		return "equal"
	}
}

func (p *Player) Compare(p2 *Player) string {
	if p.PokeCheck == p2.PokeCheck {
		if p.PokeCheck == 0 || p.PokeCheck == 5 {
			c1 := p.PokeCount(0, 5)
			c2 := p2.PokeCount(0, 5)
			return CountCompare(c1, c2)
		}
		if p.PokeCheck == 4 || p.PokeCheck == 6 || p.PokeCheck == 8 {
			return p.Pokes[0].Compare(p2.Pokes[0])
		}

		if p.PokeCheck == 3 || p.PokeCheck == 7 {
			compare := p.Pokes[0].Compare(p2.Pokes[0])
			if compare == "equal" {
				var c1, c2 int
				if p.PokeCheck == 3 {
					c1 = p.PokeCount(3, 5)
					c2 = p2.PokeCount(3, 5)
				} else {
					c1 = p.PokeCount(4, 5)
					c2 = p2.PokeCount(4, 5)
				}
				return CountCompare(c1, c2)
			} else {
				return compare
			}
		}
		if p.PokeCheck == 1 {
			compare := p.Pokes[0].Compare(p2.Pokes[0])
			if compare == "equal" {
				c1 := p.PokeCount(2, 5)
				c2 := p2.PokeCount(2, 5)
				return CountCompare(c1, c2)
			} else {
				return compare
			}
		}
		if p.PokeCheck == 2 {
			compare := p.Pokes[0].Compare(p2.Pokes[0])
			if compare == "equal" {
				compare = p.Pokes[2].Compare(p2.Pokes[2])
				if compare == "equal" {
					c1 := p.PokeCount(4, 5)
					c2 := p2.PokeCount(4, 5)
					return CountCompare(c1, c2)
				} else {
					return compare
				}
			} else {
				return compare
			}
		}
	}
	return CountCompare(p.PokeCheck, p2.PokeCheck)
}

func (p *Player) MulPokeCheck() (int, []Poke) {
	cnt := make([][]Poke, 13)
	for _, poke := range p.Pokes {
		cnt[poke.Number-2] = append(cnt[poke.Number-2], poke)
	}
	var mul3 []Poke
	var mul4 []Poke
	mul2 := make([]Poke, 0)
	for _, c := range cnt {
		if len(c) == 2 {
			mul2 = append(c, mul2...)
		}
		if len(c) == 3 {
			mul3 = append(c, mul3...)
		}
		if len(c) == 4 {
			mul4 = c
		}
	}
	if len(mul4) > 0 {
		fmt.Println(p.Name, "四张", mul4)
		return 7, append(mul4, p.MaxPokes(1, mul4[0].Number, 0)...)
	} else if len(mul2) > 0 && len(mul3) > 0 {
		fmt.Println(p.Name, "葫芦", mul3, mul2)
		return 6, append(mul3, mul2[:2]...)
	} else if len(mul3) > 0 {
		if len(mul3) > 3 {
			fmt.Println(p.Name, "葫芦", mul3[:5])
			return 6, mul3[:5]
		}
		fmt.Println(p.Name, "三张", mul3)
		return 3, append(mul3, p.MaxPokes(2, mul3[0].Number, 0)...)
	} else if len(mul2) > 2 {
		fmt.Println(p.Name, "两对", mul2[:4])
		return 2, append(mul2[:4], p.MaxPokes(1, mul2[0].Number, mul2[2].Number)...)
	} else if len(mul2) > 0 {
		fmt.Println(p.Name, "对子", mul2)
		return 1, append(mul2, p.MaxPokes(3, mul2[0].Number, 0)...)
	}
	return 0, nil
}

func (p *Player) FlushCheck() (int, []Poke) {
	cnt := make([][]Poke, 4)
	flush := make([]Poke, 0)
	isFlush := 1
	pre := 0
	hasA := false
	if p.Pokes[len(p.Pokes)-1].Number == 14 {
		hasA = true
	}
	for i, poke := range p.Pokes {
		switch poke.Color {
		case "♠":
			cnt[0] = append(cnt[0], poke)
		case "♥":
			cnt[1] = append(cnt[1], poke)
		case "♦":
			cnt[2] = append(cnt[2], poke)
		case "♣":
			cnt[3] = append(cnt[3], poke)
		}
		if i == 0 {
			pre = poke.Number
			flush = append(flush, poke)
			if hasA && pre == 2 {
				flush = append([]Poke{p.Pokes[len(p.Pokes)-1]}, flush...)
			}
		} else {
			if poke.Number != pre && poke.Number-pre == 1 {
				flush = append(flush, poke)
				isFlush++
			} else {
				isFlush = 1
				flush = []Poke{poke}
			}
			pre = poke.Number
		}
	}
	for _, c := range cnt {
		if len(c) == 5 {
			sort.Slice(c, func(i, j int) bool {
				return c[i].Number < c[j].Number
			})
			i := 1
			for ; i < len(c); i++ {
				if c[i].Number-c[i-1].Number != 1 {
					break
				}
			}
			if i != len(c) {
				fmt.Println(p.Name, "同花", c)
				return 5, c
			} else {
				sort.Slice(c, func(i, j int) bool {
					return c[i].Number < c[j].Number
				})
				fmt.Println(p.Name, "同花顺", c)
				return 8, c
			}
		}
	}
	if isFlush > 5 {
		fmt.Println(p.Name, "顺子", flush[len(flush)-5:])
		pokes := flush[len(flush)-5:]
		sort.Slice(pokes, func(i, j int) bool {
			return pokes[i].Number < pokes[j].Number
		})
		return 4, pokes
	}
	return 0, nil
}

func (r *Room) GetPokes() {
	for i := 0; i < len(r.Pokes); i++ {
		n := rand.Intn(len(r.Pokes) - i)
		r.Pokes[i], r.Pokes[n+i] = r.Pokes[n+i], r.Pokes[i]
	}
}
func (r *Room) Start() bool {
	r.Reset()
	num := 0
	for i := 0; i < 2; i++ {
		for _, player := range r.Players {
			player.Pokes = append(player.Pokes, r.Pokes[num])
			num++
		}
	}
	for i := 0; i < 5; i++ {
		r.MidPokers = append(r.MidPokers, r.Pokes[num])
		num += 2
	}
	fmt.Println(r.MidPokers)
	r.NextMidPoker()
	r.NextMidPoker()
	player, err := r.Finish()
	if err != nil {
		return false
	}
	win := ""
	for _, p := range player {
		sort.Slice(p.Pokes, func(i, j int) bool {
			return p.Pokes[i].Number < p.Pokes[j].Number
		})
		win = fmt.Sprintf("%s,\n%s,%v", win, p.Name, p.Pokes)
	}
	fmt.Println(win[1:], "获胜")
	return false
}

func (r *Room) Finish() ([]*Player, error) {
	if r.MidPos == 0 {
		return nil, errors.New("游戏未开始")
	}
	for _, player := range r.Players {
		player.Pokes = append(player.Pokes, r.MidPokers[:r.MidPos]...)
		sort.Slice(player.Pokes, func(i, j int) bool {
			return player.Pokes[i].Number < player.Pokes[j].Number
		})
	}
	fmt.Println("=====")
	for _, player := range r.Players {
		fmt.Println(player.Pokes)
		mul, mulPokes := player.MulPokeCheck()
		flush, flushPokes := player.FlushCheck()
		if mul > flush {
			player.Pokes = mulPokes
			player.PokeCheck = mul
		} else if flush != 0 {
			player.Pokes = flushPokes
			player.PokeCheck = flush
		} else {
			player.Pokes = player.Pokes[len(player.Pokes)-5:]
		}
	}
	p := []*Player{r.Players[0]}
	for i := 1; i < len(r.Players); i++ {
		if len(r.Players[i].Pokes) != 5 {
			fmt.Println(r.Players[i].Name, len(r.Players[i].Pokes))
			panic(nil)
		}
		if p[0].Compare(r.Players[i]) == "small" {
			p = []*Player{r.Players[i]}
		} else if p[0].Compare(r.Players[i]) == "equal" {
			p = append(p, r.Players[i])
		}
	}
	return p, nil
}

func (r *Room) NextMidPoker() {
	if r.MidPos == 0 {
		r.MidPos = 3
	} else {
		r.MidPos++
	}
}

func (r *Room) Reset() {
	r.MidPos = 0
	r.MidPokers = nil
	for _, player := range r.Players {
		player.Pokes = nil
		player.PokeCheck = 0
	}
	r.GetPokes()
}

func main() {
	room := Room{Id: "123", MaxPlayer: 5, Pokes: Pokes, Players: []*Player{{Id: "1", Name: "test1"}, {Id: "2", Name: "test2"}, {Id: "3", Name: "test3"}, {Id: "4", Name: "test4"}, {Id: "5", Name: "test5"}}}
	room.Start()
}
