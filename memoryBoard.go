/*

memory.go  13.8.2017

George Loo


Init
Draw
Click
Keep score
show
hide


*/

package main

import (
	//"log"
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"image/color"
	"math/rand"
	"reflect"
	"time"
)

const (
	MG_matched          = 100
	oneWidth            = 30
	oneHeight           = 30
	oneGap              = 7
	numRows             = 8
	numCols             = 8
	numCards            = numRows * numCols
	noCardClicked       = -999
	numFramesPerSec     = 60
	numSecondsToShowAll = 2 * numFramesPerSec
	numSecondsToShowTwo = 0.25 * numFramesPerSec
	numMatching         = 2
	numColors           = 7
	gameCompleteMatches = (numCols * numRows) / 2
)

type cardObjectData struct {
	//name string
	colorCode  int
	x, y, w, h float64
	clicked    bool
}

type matchingData struct {
	colorCode int
	squareIdx int
}

// type cardGraphicsData struct {
//     image *ebiten.Image
// }

var (
	aCardData        cardObjectData
	cardArr          []cardObjectData
	cardRectImage    *ebiten.Image
	frontCardImage   *ebiten.Image
	show             bool
	frameCount       int
	timerLockShowAll bool
	matcher          [10]matchingData // hack! 10 to cater for bug!
	matcherIdx       int
	ShowTwoTimerLock bool
	tries            int
	matches          int
	backToMain       int
	hgameover        bool
)

func initMatcher() {
	fmt.Printf("initMatcher \n")
	matcherIdx = 0
	for i := 0; i < numMatching; i++ {
		matcher[i].colorCode = -1
		matcher[i].squareIdx = -1
	}
}

/*

clicked on a square
    has it been clicked? - nothing happens
    not already clicked - count as one click
                show the square's contents
    2 squares clicked?
        check if they match
            yes - let them remain open
                reset Matcher
                check if all 32 are matched
                    yes - declare game over and go back
                        to main screen after user click
            no - close both squares after 2 seconds
                add to number of tries
                reset Matcher
*/
func mG_Click(mx, my int) bool {

	var x, y, w, h, x1, y1 int

	m := len(cardArr)
	for i := 0; i < m; i++ {
		x = int(cardArr[i].x)
		y = int(cardArr[i].y)
		w = int(cardArr[i].w)
		h = int(cardArr[i].h)
		x1 = x + w
		y1 = y + h
		if mx > x && mx < x1 && my > y && my < y1 {
			if !cardArr[i].clicked && matcherIdx < 2 {
				fmt.Printf("card %d clicked\n", i)
				matcher[matcherIdx].colorCode = cardArr[i].colorCode
				matcher[matcherIdx].squareIdx = i
				fmt.Printf("matcherIdx %d \n", matcherIdx)
				if matcherIdx < 2 {
					matcherIdx += 1
				}
				cardArr[i].clicked = !cardArr[i].clicked
				if matcherIdx > 1 {
					fmt.Printf("2 cards clicked\n")
					fmt.Printf("matcherIdx %d \n", matcherIdx)
					ShowTwoTimerLock = false
				}

			}
		}
	}

	if hgameover {
		backToMain++
		fmt.Printf("backToMain %d \n", backToMain)
		if backToMain > 0 {
			mouseReleased = false
			return true
		}
	}

	return false
}

func checkIfMatch() {
	var j int
	fmt.Printf("checkIfMatch \n")
	if matcher[0].colorCode != matcher[1].colorCode {
		fmt.Printf("not matching cards clicked\n")
		j = matcher[0].squareIdx
		fmt.Printf("j %d \n", j)
		cardArr[j].clicked = false
		j = matcher[1].squareIdx
		fmt.Printf("j %d \n", j)
		cardArr[j].clicked = false
		tries++
	} else {
		tries++
		matches++
	}
	initMatcher()
}

func getColorCode(c int) color.Color {
	switch c {
	case 0:
		return color.NRGBA{255, 0, 0, 0xff}
	case 1:
		return color.NRGBA{0, 255, 0, 0xff}
	case 2:
		return color.NRGBA{0, 0, 255, 0xff}
	case 3:
		return color.NRGBA{250, 250, 250, 0xff}
	case 4:
		return color.NRGBA{255, 255, 0, 0xff} // yellow
	case 5:
		return color.NRGBA{128, 128, 128, 0xff} // gray
	case 6:
		return color.NRGBA{128, 0, 255, 0xff} // purple
	default:
		return color.NRGBA{0xff, 0xff, 0xff, 0xff}
	}

}

func show_all() {

	for i := 0; i < numCards; i++ {
		cardArr[i].clicked = true
	}
}
func cover_all() {
	for i := 0; i < numCards; i++ {
		cardArr[i].clicked = false
	}
}

func MG_Draw(screen *ebiten.Image) {
	var msg string
	// Create an empty option struct
	opts := &ebiten.DrawImageOptions{}
	for i := 0; i < numCards; i++ {
		x := cardArr[i].x
		y := cardArr[i].y
		opts.GeoM.Reset()
		opts.GeoM.Scale(1.0, 1.0)
		opts.GeoM.Translate(x, y)
		if cardArr[i].clicked {
			frontCardImage.Fill(getColorCode(cardArr[i].colorCode))
			screen.DrawImage(frontCardImage, opts)
		} else {
			screen.DrawImage(cardRectImage, opts)

		}

	}

	if !timerLockShowAll {
		frameCount += 1
		if frameCount > numSecondsToShowAll {
			timerLockShowAll = true
			cover_all()
			frameCount = 0
		}
	}

	if !ShowTwoTimerLock {
		frameCount += 1
		if frameCount > numSecondsToShowTwo {
			ShowTwoTimerLock = true
			frameCount = 0
			checkIfMatch()
		}

	}

	if matches == gameCompleteMatches {
		msg = fmt.Sprintf(" GAME OVER ...Tries %d Matches %d", tries, matches)
		hgameover = true
	} else {
		msg = fmt.Sprintf(" Tries %d Matches %d", tries, matches)
	}

	ebitenutil.DebugPrint(screen, msg)

	// if show {
	//     return
	// }
	//fmt.Printf("MG_Draw\n")
	//show = true
}

/*
generate n/2 random codes
copy to the other half of the n array to get pairs of cards
return array
*/
func genColorCodes(ce [numRows * numCols]int) [numRows * numCols]int {
	var (
		m, j int
	)
	//fmt.Printf("ce %d\n", ce[0])
	m = numCols * numRows
	j = m / 2
	for i := 0; i < m/2; i++ {
		ce[i] = rand.Intn(numColors)
		ce[j] = ce[i]
		j++
	}
	fmt.Println(ce)
	//ce[0] = 1965
	return ce
}

func MG_Init() {
	var (
		x     float64
		y     float64
		w     float64
		h     float64
		fudge float64
		cc    [numRows * numCols]int
	)
	tries = 0
	matches = 0
	hgameover = false
	backToMain = 0

	fmt.Printf("MG_Init\n")
	rand.Seed(time.Now().UnixNano()) // do it once during app initialization

	cc[0] = 37
	cc = genColorCodes(cc)
	fmt.Printf("cc %d\n", cc[0])
	show = false
	fudge = 9
	x = oneGap + fudge
	y = oneGap * 3
	w = oneWidth
	h = oneHeight
	r := 0
	c := 1
	for i := 0; i < numCards; i++ {
		cardArr = append(cardArr, aCardData)
		cardArr[i].x = x
		cardArr[i].y = y
		cardArr[i].w = w
		cardArr[i].h = h
		cardArr[i].clicked = false
		cardArr[i].colorCode = cc[i] //rand.Intn(numColors)
		x = x + oneWidth + oneGap
		c += 1
		if c > numCols {
			c = 1
			r += 1
			x = oneGap + fudge
			y = y + oneHeight + oneGap
		}
	}

	if cardRectImage == nil {
		// Create an w x h image
		cardRectImage, _ = ebiten.NewImage(oneWidth, oneHeight, ebiten.FilterNearest)
	}

	// Fill the square with the white color
	//cardRectImage.Fill(color.White)
	orange := color.NRGBA{255, 128, 0, 0xff}
	cardRectImage.Fill(orange)

	if frontCardImage == nil {
		// Create an w x h image
		frontCardImage, _ = ebiten.NewImage(oneWidth, oneHeight, ebiten.FilterNearest)
	}
	green := color.NRGBA{0, 128, 0, 0xff}
	frontCardImage.Fill(green)

	show_all()
	frameCount = 0
	timerLockShowAll = false
	ShowTwoTimerLock = true
	initMatcher()

	rand.Seed(time.Now().UnixNano()) // do it once during app initialization
	s := []int{1, 2, 3, 4, 5}
	Shuffle(s)
	fmt.Println(s) // Example output: [4 3 2 1 5]
}

/*
from
https://stackoverflow.com/questions/12264789/shuffle-array-in-go
*/
func Shuffle(slice interface{}) {
	rv := reflect.ValueOf(slice)
	swap := reflect.Swapper(slice)
	length := rv.Len()
	for i := length - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		swap(i, j)
	}
}
