package main

import rl "github.com/gen2brain/raylib-go/raylib"
import (
	"fmt"
	"os/exec"
)

// KeyHandler is the signature we’ll use for one-shot callbacks.
type KeyHandler func()

func main() {
	// 1 – Setup a window and some state.
	const (
		screenW, screenH = 400, 225
		moveSpeed        = 200.0 // pixels-per-second
	)
	xPos := 0
	yPos := 0
	xStartPos := 0
	xPagePos := 0
	yStartPos := 0
	isCommandKeySelect := false
	cursorUI := "_"
	inputText := ""
	isShiftDown := false
	listKey := [][]string {
		{"abcde", "fghij", "klmno", "prqst", "uvwxy", "z1234", "56789", "0-=[]", "\\;',."},//, "/`"}
		{"ABCDE", "FGHIJ", "KLMNO", "PQRST", "UVWXY", "Z!@#$", "%^&*(", ")_+{}", "|:\"<>"},//, "?~"}
	}
	
	rl.InitWindow(screenW, screenH, "raylib-go keypress handler")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	// player := rl.NewVector2(float32(screenW/2), float32(screenH/2))

	// //-----------------------------------------------------------------
	// // 2 – Make a tiny “dispatcher”: map desired keys to callback logic.
	// //-----------------------------------------------------------------
	handlers := map[int32]KeyHandler{
	}

	//-----------------------------------------------------------------
	// 3 – Main loop.
	//-----------------------------------------------------------------
	for !rl.WindowShouldClose() {

		// ── A. EDGE-triggered keys (fires once on the frame the key goes down)
		for key, cb := range handlers {
			if rl.IsKeyPressed(key) { // IsKeyPressed … detect a single press :contentReference[oaicite:1]{index=1}
				cb()
			}
		}

		// ── B. LEVEL-triggered keys (held down = continuous movement)
		if rl.IsKeyPressed(rl.KeyRight) {
			xPos+=1;
			if(xPos>4){
				xPos = 0;
			}
		}
		if rl.IsKeyPressed(rl.KeyLeft) {
			xPos-=1;
			if(xPos<0){
				xPos = 4;
			}
		}
		if rl.IsKeyPressed(rl.KeyUp) {
			yPos-=1;
			if(yPos<0 && isCommandKeySelect==false){
				yPos = 2;
			}else if (yPos<0 && isCommandKeySelect){
				yPos = 5;
			}
		}
		if rl.IsKeyPressed(rl.KeyDown) {
			yPos+=1;
			if(yPos>2 && isCommandKeySelect==false){
				yPos = 0;
			}else if (yPos>5 && isCommandKeySelect){
				yPos = 0;
			}
		}

		if rl.IsKeyPressed(rl.KeyD) {
			if isCommandKeySelect == false {
				if((xPagePos+6)<len(listKey[0])){
					xStartPos+=1
					xPagePos+=6
				}else{
					isCommandKeySelect = true
					yStartPos = 0
					yPos = 0
				}
				// fmt.Println(xPagePos)
			}
		}
		if rl.IsKeyPressed(rl.KeyA) {
			if isCommandKeySelect {
				xPagePos = 6
				xStartPos=1
				yPos = 0
			}else{
				if (xPagePos-6) >= 0 {
					xPagePos-=6
					xStartPos-=1
				}
			}
			isCommandKeySelect = false
			// xPos=0
		}
		if rl.IsKeyPressed(rl.KeyW) {
			if isCommandKeySelect == false {
				if(yStartPos<=0){
					if xStartPos > 0 {
						xPagePos-=3
						xStartPos-=1
						yStartPos = 1
					}else{
						xPagePos = 0
					}
				}else{
					yStartPos-=1
					xPagePos-=3
				}
			}else{
				isCommandKeySelect = false
				xPagePos = 6
				xStartPos=1
				yPos = 0
			}
		}
		if rl.IsKeyPressed(rl.KeyS) {
			if isCommandKeySelect == false {
				if yStartPos == 0 {
					if((xPagePos+3)+1<len(listKey[0])){
						xPagePos+=3
						yStartPos+=1
					} else if xStartPos == 1 {
						isCommandKeySelect = true
						yStartPos = 0
						yPos = 0
					}
				}else if yStartPos == 1 {
					xPagePos+=3
					xStartPos+=1
					yStartPos = 0
				}
			}
		}
		if rl.IsKeyPressed(rl.KeyBackspace) {
			if len(inputText) > 0 {
				inputText = inputText[:len(inputText)-1]
			}
		}
		if rl.IsKeyPressed(rl.KeyEnter) {
			if isCommandKeySelect {
				switch(yPos){
				// COPY, PASTE, SHIFT, CTRL, SPACE	
				case 0:
					inputText += "CTRL"
				case 2:
					isShiftDown = !isShiftDown
				case 3:
					inputText += "CTRL"
				case 4:
					inputText += " "
				case 5:
					exec.Command("echo \""+inputText+"\" | xclip -selection clipboard")
					cmd := exec.Command("nohup", "./command.sh", ">/dev/null 2>&1")
					_, err := cmd.Output()
					if err != nil {
						fmt.Println("Error:", err)
						return
					}
					rl.CloseWindow()
					return;
				default:
					inputText += "check"
				}
			}else{
				indexShift := map[bool]int{true: 1, false: 0}[isShiftDown]
				inputText += string(listKey[indexShift][(xPagePos)+yPos][xPos])
			}
			
		}
		// ── DRAW
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		
		rl.DrawText(inputText, 20, 0, 20, rl.DarkGray)
		isYMoreSpace := 0
		isXMoreSpace := 0
		idxTemp := 0
		if(isCommandKeySelect){
			cursorUI = "__________"
		}else{
			// yPos = 0
			cursorUI = "_"
		}
		indexShift := map[bool]int{true: 1, false: 0}[isShiftDown]
		for idx, item := range listKey[indexShift] {
			if(idx%3==0 && idx!=0){
				isYMoreSpace += 20;
			}
			if(idx == 6){
				isYMoreSpace = 0
				isXMoreSpace += 110
			}
			if(idxTemp==6){
				idxTemp = 0
			}
			for validx, strVal := range item {
				rl.DrawText(string(strVal), int32(isXMoreSpace+20+(validx*15)), int32(isYMoreSpace+40+(20*idxTemp)), 20, rl.DarkGray)
			}
			idxTemp+=1
		}
		posYCursorTemp := 0
		if(isCommandKeySelect){
			posYCursorTemp = 260
		}else{
			posYCursorTemp = (xStartPos*110)+20+(15*xPos)
		}
		rl.DrawText(cursorUI, int32(posYCursorTemp), int32((yStartPos*80)+(45+(20*yPos))), 20, rl.Red)
		if isShiftDown {
			rl.DrawText("--------", 260, 80, 20, rl.Red)
		}
		
		rl.DrawText("COPY", 260, 40, 20, rl.DarkGray)
		rl.DrawText("PASTE", 260, 60, 20, rl.DarkGray)
		rl.DrawText("SHIFT", 260, 80, 20, rl.DarkGray)
		rl.DrawText("CTRL", 260, 100, 20, rl.DarkGray)
		rl.DrawText("SPACE", 260, 120, 20, rl.DarkGray)
		rl.DrawText("END", 260, 140, 20, rl.DarkGray)
		rl.EndDrawing()
	}
}
