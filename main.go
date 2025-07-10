package main

import rl "github.com/gen2brain/raylib-go/raylib"
import (
	"fmt"
	"sync"
	"os/exec"
	"bufio"
	// "log"
	"io"
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
	cursorBeatCounter := 0
	cursorBeatShow := false
	isCommandKeySelect := false
	cursorUI := "_"
	inputText := ""
	isShiftDown := false
	isInputMoveMode := false
	inputModeIndex := 0
	listKey := [][]string {
		{"abcde", "fghij", "klmno", "prqst", "uvwxy", "z1234", "56789", "0-=[]", "\\;',."},//, "/`"}
		{"ABCDE", "FGHIJ", "KLMNO", "PQRST", "UVWXY", "Z!@#$", "%^&*(", ")_+{}", "|:\"<>"},//, "?~"}
	}
	// rl.SetTraceLogLevel(rl.LogNone)  // disables all raylib log output
	rl.InitWindow(screenW, screenH, "raylib-go keypress handler")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	// player := rl.NewVector2(float32(screenW/2), float32(screenH/2))

	handlers := map[int32]KeyHandler{
	}
	
	// listening to command
	cmd := exec.Command("bash", "-c", `nohup ./backgamepadkeyproc.sh >/dev/null 2>&1`)
    stdout, err := cmd.StdoutPipe()
    if err != nil {
        panic(err)
    }

    if err := cmd.Start(); err != nil {
        panic(err)
    }

    // Channel to collect lines
    outputLines := ""
	var mu sync.Mutex

    // Read stdout in background goroutine
    go func() {
        scanner := bufio.NewScanner(stdout)
        for scanner.Scan() {
            mu.Lock()
            outputLines = scanner.Text()
			fmt.Println("output text")
			fmt.Println(outputLines)
			fmt.Println(rl.IsWindowMinimized())
			if rl.IsWindowMinimized() {
				rl.RestoreWindow()
			}else{
				rl.MinimizeWindow()
			}
            mu.Unlock()
        }
    }()
	// end command
	
	for !rl.WindowShouldClose() {
		// ── A. EDGE-triggered keys (fires once on the frame the key goes down)
		cursorBeatCounter += 1
		for key, cb := range handlers {
			if rl.IsKeyPressed(key) { // IsKeyPressed … detect a single press :contentReference[oaicite:1]{index=1}
				cb()
			}
		}

		// ── B. LEVEL-triggered keys (held down = continuous movement)
		if rl.IsKeyPressed(rl.KeyRight) {
			if isInputMoveMode {
				inputModeIndex += 1
				if inputModeIndex > len(inputText) {
					inputModeIndex -= 1
				}
				cursorBeatCounter = 0
				cursorBeatShow = true
			}else{
				xPos+=1;
				if(xPos>4){
					xPos = 0;
				}
			}
		}
		if rl.IsKeyPressed(rl.KeyLeft) {
			if isInputMoveMode {
				inputModeIndex -= 1
				if inputModeIndex < 0 {
					inputModeIndex = 0
				}
				cursorBeatCounter = 0
				cursorBeatShow = true
			}else{
				xPos-=1;
				if(xPos<0){
					xPos = 4;
				}
			}
		}
		if rl.IsKeyPressed(rl.KeyUp) {
			if isInputMoveMode == false {
				yPos-=1;
				if(yPos<0 && isCommandKeySelect==false){
					yPos = 2;
				}else if (yPos<0 && isCommandKeySelect){
					yPos = 5;
				}
			}
		}
		if rl.IsKeyPressed(rl.KeyDown) {
			if isInputMoveMode == false {
				yPos+=1;
				if(yPos>2 && isCommandKeySelect==false){
					yPos = 0;
				}else if (yPos>5 && isCommandKeySelect){
					yPos = 0;
				}
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
		if rl.IsKeyPressed(rl.KeyW) || rl.IsKeyPressed(rl.KeyPageUp) {
			if isCommandKeySelect == false {
				if isInputMoveMode == false && yStartPos==0 && xStartPos==0 {
					isInputMoveMode = true
				}else{
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
				}
			}else{
				isCommandKeySelect = false
				xPagePos = 6
				xStartPos=1
				yPos = 0
			}
		}
		// for test
		if rl.IsKeyPressed(rl.KeyP) {
			isInputMoveMode = !isInputMoveMode
		}
		if rl.IsKeyPressed(rl.KeyS) || rl.IsKeyPressed(rl.KeyPageDown) {
			if isInputMoveMode {
				isInputMoveMode = false
			}else{
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
		}
		if rl.IsKeyPressed(rl.KeyBackspace) {
			if len(inputText) > 0 {
				// inputText = inputText[:len(inputText)-1]
				inputTextTemp := ""
				if inputModeIndex-1 >= 0 {
					inputTextTemp = inputText[:inputModeIndex-1]
					// inputTextTemp = inputText[:0]
				}
				if inputModeIndex < len(inputText) {
					inputTextTemp += inputText[inputModeIndex:]
				}
				inputText = inputTextTemp
				inputModeIndex -= 1
				if inputModeIndex < 0 {
					inputModeIndex = 0
				}
			}
		}
		if rl.IsKeyPressed(rl.KeyEnter) {
			if isCommandKeySelect {
				switch(yPos){
				// COPY, PASTE, SHIFT, CTRL, SPACE	
				case 0:
					inputText += "COPY"
					inputModeIndex = 4
				case 1:
					cmd2 := exec.Command("bash", "-c", `nohup ./command.sh >/dev/null 2>&1`)
					cmd2.Stdout = io.Discard
    				cmd2.Stderr = io.Discard
					cmd2.Start();
					rl.MinimizeWindow()
					return;
				case 2:
					isShiftDown = !isShiftDown
				case 3:
					inputText = "CTRL"
					inputModeIndex = 4
				case 4:
					inputText += " "
					inputModeIndex+=1
				case 5:
					cmd := exec.Command("bash", "-c", `echo "`+inputText+`" | xclip -selection clipboard`)
					cmd.Stdout = io.Discard
    				cmd.Stderr = io.Discard
					cmd.Start();
					cmd2 := exec.Command("bash", "-c", `nohup ./command.sh >/dev/null 2>&1`)
					cmd2.Stdout = io.Discard
    				cmd2.Stderr = io.Discard
					cmd2.Start();
					rl.MinimizeWindow()
					return;
				default:
					inputText += "check"
					inputModeIndex+=4
				}
			}else{
				indexShift := map[bool]int{true: 1, false: 0}[isShiftDown]
				addNewText := string(listKey[indexShift][(xPagePos)+yPos][xPos])
				tempInputText := inputText[:inputModeIndex]+addNewText
				inputText = tempInputText+inputText[inputModeIndex:]
				inputModeIndex += len(addNewText)
			}
			
		}
		// ── DRAW
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		inputTextCursor := ""
		if inputModeIndex >= 0 {
			inputTextCursor = inputText[:inputModeIndex]
		}else{
			inputTextCursor = inputText
		}
		if cursorBeatCounter == 50 {
			cursorBeatCounter = 0
			cursorBeatShow = !cursorBeatShow
		}
		if cursorBeatShow {
			inputTextCursor += "|"
		}
		if inputModeIndex < len(inputText) {
			inputTextCursor += inputText[inputModeIndex:]
		}
		rl.DrawText(inputTextCursor, 20, 10, 20, rl.DarkGray)
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
		if isInputMoveMode == false {
			rl.DrawText(cursorUI, int32(posYCursorTemp), int32((yStartPos*80)+(45+(20*yPos))), 20, rl.Red)
		}
		
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
