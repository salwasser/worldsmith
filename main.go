package main

import (
	wsconfig "WorldSmith/config"
	"encoding/json"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"image"
	"image/color"
	"io"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/op"
)

type WorldMap struct {
	Grid [][]string
}

const (
	configFile = "assets/config/config.json"
)

var (
	worldMap = WorldMap{
		Grid: [][]string{
			{"O", "O", "O", "O", "O"},
			{"O", "G", "G", "G", "O"},
			{"O", "G", "G", "G", "O"},
			{"O", "G", "G", "G", "O"},
			{"O", "O", "O", "O", "O"},
		},
	}
	blue  = color.NRGBA{R: 0, G: 0, B: 255, A: 255}
	green = color.NRGBA{R: 0, G: 255, B: 0, A: 255}
)

func main() {
	configF, errF := os.Open(configFile)
	if errF != nil {
		log.Fatalf("%+v", errF)
		return
	}

	var configB []byte
	configB, errF = io.ReadAll(configF)
	if errF != nil {
		log.Fatalf("%+v", errF)
		return
	}

	config := wsconfig.Config{}
	errF = json.Unmarshal(configB, &config)
	if errF != nil {
		log.Fatalf("%+v", errF)
	}

	go func() {
		window := new(app.Window)
		window.Option(app.Size(unit.Dp(config.Window.Width), unit.Dp(config.Window.Height)), app.Title("WorldSmith"))
		err := run(window)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(window *app.Window) error {
	var ops op.Ops
	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			// This graphics context is used for managing the rendering state.
			gtx := app.NewContext(&ops, e)

			// Draw the ground layer.
			err := drawGround(&gtx)
			if err != nil {
				log.Printf("%+v\n", err)
			}

			// Pass the drawing operations to the GPU.
			e.Frame(gtx.Ops)
		}
	}
}

func drawGround(gtx *layout.Context) error {
	xOffset := 0
	xJump := 50
	yOffset := 0
	yJump := 50
	for i := range worldMap.Grid {
		for j := range worldMap.Grid[i] {
			switch worldMap.Grid[i][j] {
			case "O":
				rect := clip.Rect{Min: image.Pt(xOffset, yOffset), Max: image.Pt(xOffset+xJump, yOffset+yJump)}.Push(gtx.Ops)
				paint.ColorOp{Color: blue}.Add(gtx.Ops)
				paint.PaintOp{}.Add(gtx.Ops)
				rect.Pop()
			case "G":
				rect := clip.Rect{Min: image.Pt(xOffset, yOffset), Max: image.Pt(xOffset+xJump, yOffset+yJump)}.Push(gtx.Ops)
				paint.ColorOp{Color: green}.Add(gtx.Ops)
				paint.PaintOp{}.Add(gtx.Ops)
				rect.Pop()
			}
			xOffset += xJump
		}
		xOffset = 0
		yOffset += yJump
	}
	return nil
}
