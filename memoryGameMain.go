/*


Memory game gloo 13.8.2017

from Editen tutorial 2 gloo 6.8.2017

https://github.com/hajimehoshi/ebiten/wiki/Tutorial%3AScreen%2C-colors-and-squares
*/

//jj
package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/wav"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"image"
	"image/color"
	"log"
	"os"
	"path/filepath"
	//"github.com/hajimehoshi/ebiten/keys"
)

var (
	square          *ebiten.Image
	x, y            float64
	mode            int
	audioContext    *audio.Context
	audioPlayer     *audio.Player
	delayCount      int
	prevkey         string
	digitsImage     *ebiten.Image
	mouseReleased   bool
	gameover        bool
	buttonRunTimer  buttonTimerStruct
	buttonQuitTimer buttonTimerStruct
)

const ( // const
	screenwidth     = 320
	screenheight    = 320
	modeStartScreen = 1
	modeGameScreen  = 2
	sampleRate      = 44100
	startButtonCode = 801
	quitButtonCode  = 802
)

func soundinit() {
	var err error
	audioContext, err = audio.NewContext(sampleRate)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Open(filepath.Join("data", "jab.wav"))
	if err != nil {
		log.Fatal(err)
	}

	d, err := wav.Decode(audioContext, f)
	if err != nil {
		log.Fatal(err)
	}

	audioPlayer, err = audio.NewPlayer(audioContext, d)
	if err != nil {
		log.Fatal(err)
	}
}

func exitProgram() {
	os.Exit(0)
}

func digitsDraw(screen *ebiten.Image) {

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Reset()
	opts.GeoM.Translate(10, 10)
	r := image.Rect(50, 0, 100, 80)
	opts.SourceRect = &r
	screen.DrawImage(digitsImage, opts)
}

/*
check if same function name can exist in the same project
*/
func Haha1() {
	fmt.Print("Ha ha ebitut\n")
}

// where you choose to begin the game or exit
func outerloop(screen *ebiten.Image) error {

	//screen.Fill(color.NRGBA{0xff, 0xff, 128, 0xff})     //
	screen.Fill(color.NRGBA{154, 158, 5, 0xff})

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) &&
		mouseReleased {
		mx, my := ebiten.CursorPosition()
		mouseReleased = false
		if Button_Click(mx, my) == startButtonCode {
			fmt.Printf("click button\n")
			//mode = modeGameScreen
			//MG_Init()
			buttonRunTimer.set(buttonRunTimer.kButtonStart, 30) // 60 is 1 seconds\

		} else if Button_Click(mx, my) == quitButtonCode {
			buttonQuitTimer.set(buttonRunTimer.kButtonStart, 30)
		}
	} else if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		mouseReleased = true
	}

	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		exitProgram()
	}

	Button_Draw(screen)

	// timers are to allow the user to see the button press animation
	buttonRunTimer.inc()
	if buttonRunTimer.check() {
		fmt.Printf("buttonRunTimer execcute!\n")
		buttonRunTimer.reset()
		mode = modeGameScreen
		MG_Init()

	}

	buttonQuitTimer.inc()
	if buttonQuitTimer.check() {
		fmt.Printf("quit execcute!\n")
		buttonQuitTimer.reset() // unnecessary but for consistency
		exitProgram()
	}

	//digitsDraw(screen)

	if ebiten.IsKeyPressed(ebiten.KeyJ) && PushKeyOnce("J") {
		fmt.Print("key J\n")
		//PushKeyOnce("J")
	}
	if ebiten.IsKeyPressed(ebiten.KeyU) && PushKeyOnce("U") {
		fmt.Print("key U\n")
		//PushKeyOnce("U")
	}
	ebitenutil.DebugPrint(screen, "Memory Game \nin Ebiten\n")
	return nil
}

/*

need to store a map of which keys need to be non-repeating

*/

func PushKeyOnce(k string) bool {
	if k != prevkey {
		prevkey = k
		return true
	}
	prevkey = k
	return false
}

// actual game
func innerloop(screen *ebiten.Image) error {

	delayCount += 1

	if ebiten.IsRunningSlowly() {
		return nil
	}

	if ebiten.IsKeyPressed(ebiten.KeyP) && !audioPlayer.IsPlaying() {
		audioPlayer.Rewind()
		audioPlayer.Play()
	}

	if err := audioContext.Update(); err != nil {
		return err
	}

	// Fill the screen with #FF0000 color
	//screen.Fill(color.NRGBA{0xff, 0x00, 0x00, 0xff})   // red
	//screen.Fill(color.NRGBA{0xaa, 0xff, 0xff, 0xff})     // green
	screen.Fill(color.NRGBA{0x00, 0x00, 0x00, 0xff}) //black
	MG_Draw(screen)

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) &&
		mouseReleased {
		mx, my := ebiten.CursorPosition()
		gameover = mG_Click(mx, my)
		mouseReleased = false
		fmt.Printf("mouseRele false\n")
	} else if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		//fmt.Printf("mouseReleased = true\n")
		mouseReleased = true
	}

	//if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
	if gameover || ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		//mode = modeStartScreen
		gameinit()
		ebiten.SetFullscreen(false)
	}

	return nil
}

func update(screen *ebiten.Image) error {

	if mode == modeStartScreen {
		outerloop(screen)

	} else if mode == modeGameScreen {
		innerloop(screen)

	}

	return nil
}

func gameinit() {

	mode = modeStartScreen
	x = 0
	y = 25
	delayCount = 0
	mouseReleased = true
	gameover = false
}

func runOnce() {
	fmt.Printf("running ********************************\n")

	buttonOOP.init()
	b1 := new(buttonobjectData)
	b1.name = "Jack"
	b2 := new(buttonobjectData)
	b2.name = "Joey"
	b1.init()
	b2.init()

	Button_Init()
	Button_Create("run.png", startButtonCode, 100, 100, 200, 100, true)
	Button_Create("quit.png", quitButtonCode, 100, 220, 200, 100, true)
	buttonRunTimer.reset()
	buttonQuitTimer.reset()

	var err error
	var fname string
	fname = filepath.Join("data", "digits.png")
	digitsImage, _, err = ebitenutil.NewImageFromFile(
		fname,
		ebiten.FilterNearest)
	if err != nil {
		log.Fatal(err)
	}

}

func main() {
	runOnce()
	gameinit()
	soundinit()
	ebiten.SetRunnableInBackground(true)
	ebiten.SetFullscreen(false)

	scale := 2.0
	// Initialize Ebiten, and loop the update() function
	if err := ebiten.Run(update, screenwidth, screenheight, scale, "Memory Game 0.3 by George Loo"); err != nil {
		panic(err)
	}
	fmt.Printf("Program ended \n")

}
