package main

import (
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

const winWidth, winHeight int = 800, 600

func die(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

type color struct {
	r, g, b byte
}

type pos struct {
	x, y float32
}
type ball struct {
	pos
	radius int
	vx     float32
	vy     float32
	color  color
}
type paddle struct {
	pos
	width  int
	height int
	color  color
}

func setPixel(x, y int, c color, pixels []byte) {
	index := (y*winWidth + x) * 4

	if index < len(pixels)-4 && index >= 0 {
		pixels[index] = c.r
		pixels[index+1] = c.g
		pixels[index+2] = c.b
	}
}

func drawRectangle(x, y, a, b int, c color, pixels []byte) {
	for j := y; j < y+b; j++ {
		for i := x; i < x+a; i++ {
			setPixel(i, j, c, pixels)
		}
	}
}

func drawCircle(x, y, r int, c color, pixels []byte) {
	for j := y - r; j < y+r; j++ {
		for i := x - r; i < x+r; i++ {
			if (i-x)*(i-x)+(j-y)*(j-y) <= r*r {
				setPixel(i, j, c, pixels)
			}
		}
	}
}

func drawBackground(c color, pixels []byte) {
	for y := 0; y < winHeight; y++ {
		for x := 0; x < winWidth; x++ {
			setPixel(x, y, c, pixels)
		}
	}
}

func (paddle *paddle) draw(pixels []byte) {
	startX := int(paddle.x) - paddle.width/2
	startY := int(paddle.y) - paddle.height/2
	drawRectangle(startX, startY, paddle.width, paddle.height,
		paddle.color, pixels)
}

func (ball *ball) draw(pixels []byte) {
	drawCircle(int(ball.x), int(ball.y), ball.radius,
		ball.color, pixels)
}

func (ball *ball) update() {
	ball.x += ball.vx
	ball.y += ball.vy

	if int(ball.y)-ball.radius < 0 || int(ball.y)+ball.radius > winHeight {
		ball.vy = -ball.vy
	}

	if int(ball.x)-ball.radius < 0 || int(ball.x)+ball.radius > winWidth {
		ball.vx = -ball.vx
	}
}

func (paddle *paddle) update(keyState []uint8) {
	if keyState[sdl.SCANCODE_UP] != 0 {
		paddle.y -= 3
	}
	if keyState[sdl.SCANCODE_DOWN] != 0 {
		paddle.y += 3
	}
}

func (paddle *paddle) aiUpdate(ball *ball) {
	paddle.y = ball.y
}

func main() {

	window, err := sdl.CreateWindow("Testing SDL", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(winWidth), int32(winHeight), sdl.WINDOW_SHOWN)
	die(err)
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	die(err)
	defer renderer.Destroy()

	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING,
		int32(winWidth), int32(winHeight))
	die(err)
	defer tex.Destroy()

	pixels := make([]byte, winWidth*winHeight*4)

	drawBackground(color{150, 100, 200}, pixels)

	player1 := paddle{pos{100, 100}, 10, 100, color{255, 255, 255}}
	player2 := paddle{pos{700, 100}, 10, 100, color{255, 255, 255}}
	ball := ball{pos{300, 300}, 10, 5, 5, color{255, 255, 255}}

	keyState := sdl.GetKeyboardState()
	// Event Loop
	for {

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		drawBackground(color{150, 100, 200}, pixels)
		player1.update(keyState)
		player1.draw(pixels)
		player2.aiUpdate(&ball)
		player2.draw(pixels)
		ball.update()
		ball.draw(pixels)
		tex.Update(nil, pixels, winWidth*4)
		renderer.Copy(tex, nil, nil)
		renderer.Present()

		sdl.Delay(16)
	}
}

/*
	drawBackground
	----This can be modified to make more efficient. We just need to erase the objects.

*/
