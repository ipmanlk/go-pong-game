package main

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

// settings
const winWidth, winHeight int32 = 800, 600

// game state
type gameState int

const (
	start gameState = iota
	play
)

var state = start

// supported numbers for score (0-3)
var nums = [][]byte{
	{
		1, 1, 1,
		1, 0, 1,
		1, 0, 1,
		1, 0, 1,
		1, 1, 1,
	},
	{
		1, 1, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		1, 1, 1,
	},
	{
		1, 1, 1,
		0, 0, 1,
		1, 1, 1,
		1, 0, 0,
		1, 1, 1,
	},
	{
		1, 1, 1,
		0, 0, 1,
		0, 1, 1,
		0, 0, 1,
		1, 1, 1,
	},
}

// structs
type color struct {
	r, g, b byte
}

type pos struct {
	x, y float32
}

type ball struct {
	pos    // composistion <- can be used instead of inheritance in oop
	radius float32
	xv     float32
	yv     float32
	color  color
}

type paddle struct {
	pos
	w     float32
	h     float32
	speed float32
	score int
	color color
}

// draw functions
func (paddle *paddle) draw(pixels []byte) {
	startX := int(paddle.x - paddle.w/2)
	startY := int(paddle.y - paddle.h/2)

	for y := 0; y < int(paddle.h); y++ {
		for x := 0; x < int(paddle.w); x++ {
			setPixcel(startX+x, startY+y, paddle.color, pixels)
		}
	}

	numX := lerp(paddle.x, getCenter().x, 0.2)
	drawNumber(pos{numX, 35}, paddle.color, 10, paddle.score, pixels)
}

func (ball *ball) draw(pixels []byte) {
	//  YAGNI- Ya Aint Gonna Need It
	for y := -ball.radius; y < ball.radius; y++ {
		for x := -ball.radius; x < ball.radius; x++ {
			if x*x+y*y < ball.radius*ball.radius {
				setPixcel(int(ball.x+x), int(ball.y+y), ball.color, pixels)
			}
		}
	}
}

func drawNumber(pos pos, color color, size int, num int, pixels []byte) {
	startX := int(pos.x) - (size*3)/2
	startY := int(pos.y) - (size*5)/2

	for i, v := range nums[num] {
		if v == 1 {
			for y := startY; y < startY+size; y++ {
				for x := startX; x < startX+size; x++ {
					setPixcel(x, y, color, pixels)
				}
			}
		}

		startX += size

		if (i+1)%3 == 0 {
			startY += size
			startX -= size * 3
		}
	}
}

// update functions
func (ball *ball) update(leftPaddle *paddle, rightPaddle *paddle, elapsedTime float32) {
	ball.x += ball.xv * elapsedTime
	ball.y += ball.yv * elapsedTime

	if (ball.y-ball.radius) < 0 || (ball.y+ball.radius) > float32(winHeight) {
		ball.yv = -ball.yv
	}

	if ball.x < 0 {
		rightPaddle.score++
		ball.pos = getCenter()
		state = start
	} else if int(ball.x) > int(winWidth) {
		leftPaddle.score++
		ball.pos = getCenter()
		state = start
	}

	if ball.x-ball.radius < leftPaddle.x+leftPaddle.w/2 {
		if (ball.y) > (leftPaddle.y)-leftPaddle.h/2 && (ball.y) < (leftPaddle.y)+leftPaddle.h/2 {
			ball.xv = -ball.xv
			ball.x = leftPaddle.x + leftPaddle.w/2.0 + ball.radius
		}
	}

	if ball.x+ball.radius > rightPaddle.x-rightPaddle.w/2 {
		if ball.y > rightPaddle.y-rightPaddle.h/2 && ball.y < rightPaddle.y+rightPaddle.h/2 {
			ball.xv = -ball.xv
			ball.x = rightPaddle.x - rightPaddle.w/2.0 - ball.radius
		}
	}

}

func (paddle *paddle) update(keyState []uint8, elapsedTime float32) {
	if keyState[sdl.SCANCODE_UP] != 0 {
		paddle.y -= paddle.speed * elapsedTime
	}
	if keyState[sdl.SCANCODE_DOWN] != 0 {
		paddle.y += paddle.speed * elapsedTime
	}
}

func (paddle *paddle) aiUpdate(ball *ball, elapsedTime float32) {
	paddle.y = ball.y
}

// utility functions
func setPixcel(x, y int, c color, pixcels []byte) {
	index := (y*int(winWidth) + x) * 4

	if index < len(pixcels)-4 && index >= 0 {
		pixcels[index] = c.r
		pixcels[index+1] = c.g
		pixcels[index+2] = c.b
	}

}

func clear(pixels []byte) {
	for i := range pixels {
		pixels[i] = 0
	}
}

func lerp(a float32, b float32, pct float32) float32 {
	return a + pct*(b-a)
}

func getCenter() pos {
	return pos{float32(winWidth / 2), float32(winHeight / 2)}
}

// driver
func main() {
	err := sdl.Init(sdl.INIT_EVERYTHING)

	if err != nil {
		fmt.Println(err)
		return
	}

	window, err := sdl.CreateWindow("pong-game", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, winWidth, winHeight, sdl.WINDOW_SHOWN)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer renderer.Destroy()

	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, winWidth, winHeight)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer tex.Destroy()

	// make byte array according to the window height and width
	pixels := make([]byte, winWidth*winHeight*4)

	// game components
	manualPlayer := paddle{pos{50, 100}, 20, 100, 800, 0, color{255, 255, 255}}
	aiPlayer := paddle{pos{float32(winWidth) - 50, 100}, 20, 100, 300, 0, color{255, 255, 255}}

	ball := ball{pos{300, 300}, 20, 400, 400, color{255, 255, 255}}

	// represents the current state of keyboard
	keyState := sdl.GetKeyboardState()

	var frameStart time.Time
	var elapsedTime float32

	// game loop
	for {
		frameStart = time.Now()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		if state == play {
			manualPlayer.update(keyState, elapsedTime)
			aiPlayer.aiUpdate(&ball, elapsedTime)
			ball.update(&manualPlayer, &aiPlayer, elapsedTime)
		} else if state == start {
			if keyState[sdl.SCANCODE_SPACE] != 0 {
				if aiPlayer.score == 3 || manualPlayer.score == 3 {
					aiPlayer.score = 0
					manualPlayer.score = 0
				}
				state = play
			}
		}

		// clear screen before drawing
		clear(pixels)

		manualPlayer.draw(pixels)
		aiPlayer.draw(pixels)
		ball.draw(pixels)

		tex.Update(nil, pixels, int(winWidth)*4)
		renderer.Copy(tex, nil, nil)
		renderer.Present()

		// maximum limit for fps based on how much time passed
		elapsedTime = float32(time.Since(frameStart).Seconds())

		if elapsedTime < .005 {
			sdl.Delay(5 - uint32(elapsedTime/1000.0))
			elapsedTime = float32(time.Since(frameStart).Seconds())
		}
	}
}
