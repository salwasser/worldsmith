package main

import (
	wsconfig "WorldSmith/config"
	"encoding/json"
	"gioui.org/io/event"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"image"
	"image/color"
	"io"
	"log"
	"os"
	"time"

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
			{"O", "O", "O", "O", "O", "O", "O", "O", "O", "O", "O", "O", "O", "O", "O", "O", "O", "O", "O", "O", "O"},
			{"O", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "O"},
			{"O", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "O"},
			{"O", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "O"},
			{"O", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "O"},
			{"O", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "O"},
			{"O", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "O"},
			{"O", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "O"},
			{"O", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "O"},
			{"O", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "O"},
			{"O", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "O"},
			{"O", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "O"},
			{"O", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "O"},
			{"O", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "O"},
			{"O", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "O"},
			{"O", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "O"},
			{"O", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "O"},
			{"O", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "O"},
			{"O", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "G", "O"},
			{"O", "O", "O", "O", "O", "O", "O", "O", "O", "O", "O", "O", "O", "O", "O", "O", "O", "O", "O", "O", "O"},
		},
	}
	blue  = color.NRGBA{R: 0, G: 0, B: 255, A: 255}
	green = color.NRGBA{R: 0, G: 255, B: 0, A: 255}
	black = color.NRGBA{A: 255}
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
	editorMode := false

	// Update the board 5 times per second.
	advanceGameClock := time.NewTicker(time.Second / 1)
	defer advanceGameClock.Stop()
	windowEvents := make(chan event.Event)
	acks := make(chan struct{})

	go func() {
		for {
			windowEvent := window.Event()
			//log.Printf("Window Event: %+v\n", windowEvent)
			windowEvents <- windowEvent
			//There's probably a good reason to wait here.
			<-acks
			if _, ok := windowEvent.(app.DestroyEvent); ok {
				return
			}
		}
	}()

	for {
		select {
		case windowEvent := <-windowEvents:

			switch typedWindowEvent := windowEvent.(type) {
			default:
				//log.Printf("typedWindowEvent: %+v", typedWindowEvent)
			case app.DestroyEvent:
				return typedWindowEvent.Err
			case app.FrameEvent:

				// This graphics context is used for managing the rendering state.
				//log.Printf("typedWindowEvent: %+v\n", typedWindowEvent)
				gtx := app.NewContext(&ops, typedWindowEvent)

				/*
					// register a global key listener for the escape key wrapping our entire UI.
					area := clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops)
					event.Op(gtx.Ops, window)
				*/
				for inputCheckDone := false; !inputCheckDone; {
					event, ok := gtx.Event(
						key.Filter{Name: "E"},
						key.Filter{Name: "W"})
					if !ok {
						inputCheckDone = true
					} else {
						switch typedEvent := event.(type) {
						case key.Event:
							if typedEvent.Name == "E" && typedEvent.State == key.Press {
								log.Printf("E pressed.")
								editorMode = !editorMode
							}
						}
					}
				}
				/*
					area.Pop()

				*/

				//area := clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops)
				//event.Op(gtx.Ops, window)

				// check for presses of the escape key and close the window if we find them.
				/*
					for {
						event, ok := gtx.Event(key.Filter{
							Name: key.NameEscape,
						})
						if !ok {
							break
						}
						switch event := event.(type) {
						case key.Event:
							if event.Name == key.NameEscape {
								return nil
							}
						}
					}

				*/
				// render and handle UI.
				//ui.Layout(gtx)
				//area.Pop()

				// Draw the ground layer.
				err := drawGround(&gtx, editorMode, worldMap, typedWindowEvent)
				if err != nil {
					log.Printf("%+v\n", err)
				}

				// Pass the drawing operations to the GPU.
				log.Printf("Entering draw.")
				typedWindowEvent.Frame(gtx.Ops)
				log.Printf("Draw done.")
			}
			acks <- struct{}{}
		case <-advanceGameClock.C:
			log.Printf("Tick.")
			window.Invalidate()
		}
	}
}

func drawGround(gtx *layout.Context, editorMode bool, theWorld WorldMap, event app.FrameEvent) error {
	editorWidth := 0
	editorLineWidth := 0
	if editorMode {
		editorWidth = int(float32(event.Size.X) * 0.15)
		editorLineWidth = int(0.005 * float32(event.Size.X))
		if editorLineWidth < 3 {
			editorLineWidth = 3
		} else if editorLineWidth > 10 {
			editorLineWidth = 10
		}
	}

	xJump := (event.Size.X - editorWidth) / len(theWorld.Grid[0])
	yJump := event.Size.Y / len(theWorld.Grid)
	xOrig := editorWidth
	yOrig := 0
	xOffset := xOrig
	yOffset := yOrig
	for i := range theWorld.Grid {
		for j := range theWorld.Grid[i] {
			if editorMode {
				rect := clip.Rect{Min: image.Pt(xOffset, yOffset), Max: image.Pt(xOffset+xJump, yOffset+yJump)}.Push(gtx.Ops)
				paint.ColorOp{Color: black}.Add(gtx.Ops)
				paint.PaintOp{}.Add(gtx.Ops)
				rect.Pop()
			}
			switch theWorld.Grid[i][j] {
			case "O":
				rect := clip.Rect{Min: image.Pt(xOffset+editorLineWidth, yOffset+editorLineWidth), Max: image.Pt(xOffset+xJump-editorLineWidth, yOffset+yJump-editorLineWidth)}.Push(gtx.Ops)
				paint.ColorOp{Color: blue}.Add(gtx.Ops)
				paint.PaintOp{}.Add(gtx.Ops)
				rect.Pop()
			case "G":
				rect := clip.Rect{Min: image.Pt(xOffset+editorLineWidth, yOffset+editorLineWidth), Max: image.Pt(xOffset+xJump-editorLineWidth, yOffset+yJump-editorLineWidth)}.Push(gtx.Ops)
				paint.ColorOp{Color: green}.Add(gtx.Ops)
				paint.PaintOp{}.Add(gtx.Ops)
				rect.Pop()
			}

			xOffset += xJump
		}
		xOffset = xOrig
		yOffset += yJump
	}
	return nil
}
