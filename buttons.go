

/*

buttons
gloo 9.8.2017


create
draw
click
return codes

*/

package main

import (
    "log"
    "fmt"
    "image/color"
    "path/filepath"
    "github.com/hajimehoshi/ebiten"
    "github.com/hajimehoshi/ebiten/ebitenutil"
)

type buttonobjectData struct {
    name string
    code int        // return this code if clicked
    x,y,w,h float64
    active bool
    image *ebiten.Image
    clicked bool
    count int  // count frames
}




var (
    aButtonData buttonobjectData 
    buttonArr []buttonobjectData
   	buttonRectImage *ebiten.Image
    buttonOOP buttonobjectData
)

const (
    noButtonPressed = 0
)

func (b *buttonobjectData ) init() {
    fmt.Print("button init test \n")
    //b.name = "George"
    fmt.Print("button string test ",b.name," is very good \n")
} 

func Button_Init() {

    buttonRectImage, _ = ebiten.NewImage(100, 100, ebiten.FilterNearest)
}

func Button_Create( n string, 
                    code int, 
                    x float64, 
                    y float64, 
                    w float64, 
                    h float64, 
                    a bool) {
                    
                    
    fmt.Print("button create\n")
 
    buttonArr = append(buttonArr, aButtonData)
    i := len(buttonArr)
    i = i - 1   // zero index
    buttonArr[i].name = n
    buttonArr[i].code = code
    buttonArr[i].x = x
    buttonArr[i].y = y
    buttonArr[i].w = w
    buttonArr[i].h = h
    buttonArr[i].active = a
    buttonArr[i].clicked = false
    buttonArr[i].count = 0
    
    var err error
    var fname string
	fname = filepath.Join("data", n)
	img, _, err := ebitenutil.NewImageFromFile(
        fname, 
        ebiten.FilterNearest)
	if err != nil {
		log.Fatal(err)
	}
	buttonArr[i].image = img

}

/*
check if same function name can exist in the same project
*/
func Haha() {
    fmt.Print("Ha ha buttons\n")
}


func Button_Draw(screen *ebiten.Image) {
    //fmt.Print("button draw\n")
    // Create an empty option struct
    opts := &ebiten.DrawImageOptions{}
    
    m := len(buttonArr)
    blue := color.NRGBA{0x00, 0x00, 0xff, 0xff}
    for i := 0; i < m; i++ {
        buttonRectImage.Fill(blue)
        opts.GeoM.Reset()
        //opts.GeoM.Scale( buttonArr[i].w, buttonArr[i].h )
        if buttonArr[i].clicked {
            opts.GeoM.Translate(buttonArr[i].x+5, buttonArr[i].y+5)
            buttonArr[i].count += 1
            if buttonArr[i].count > 30 {
                buttonArr[i].clicked = false
                buttonArr[i].count = 0
            }
            
        } else {
            opts.GeoM.Translate(buttonArr[i].x, buttonArr[i].y)
        }
        
        //screen.DrawImage(buttonRectImage, opts)
        screen.DrawImage(buttonArr[i].image, opts)
        
        //ebitenutil.DebugPrint(buttonRectImage, buttonArr[i].name)
        //ebitenutil.DebugPrint(screen, buttonArr[i].name)
    }

}

/*

start timer
how long
what to do when timer runs out
reset timer

button timer

*/


//const kButtonStart = 2222

type buttonTimerStruct struct {
    counter int 
    numFrames int 
    code int 
    kButtonStart int
}

func (bt *buttonTimerStruct) set(codef, numF int) {
    fmt.Print("buttontimer set! \n")
    bt.code = codef
    bt.numFrames = numF
}

func (bt *buttonTimerStruct) inc() {
    if bt.code == bt.kButtonStart {
        //fmt.Print("inc\n")
        bt.counter++
    }
}

func (bt *buttonTimerStruct) check() bool {
    if bt.counter > bt.numFrames {
        return true
    } 
    return false
}

func (bt *buttonTimerStruct) reset() {
    fmt.Print("reset buttonTimer \n")
    bt.counter = 0
    bt.code = 0
    bt.numFrames = 0
    bt.kButtonStart = 2222
}

func Button_Click(mx, my int) int {
    //fmt.Print("button click\n")
    var x, y, w, h, x1, y1 int
    
    m := len(buttonArr)
    for i := 0; i < m; i++ {
        x = int (buttonArr[i].x)
        y = int (buttonArr[i].y)
        w = int (buttonArr[i].w)
        h = int (buttonArr[i].h)
        x1 = x + w
        y1 = y + h
        code := buttonArr[i].code
        if mx > x && mx < x1 && my > y && my < y1 {
            buttonArr[i].clicked = true
            buttonArr[i].count = 0
            return code
        }
    }

    return noButtonPressed
}
