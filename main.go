package main

import rl "github.com/gen2brain/raylib-go/raylib"
import (
	"fmt"
	"sync"
	"os/exec"
	"bufio"
	"strings"
	"math"
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
	isTextPredictMode := false
	textPredicIndex := 0
	cursorUI := "_"
	inputText := ""
	isShiftDown := false
	isInputMoveMode := false
	inputModeIndex := 0
	listKey := [][]string {
		{"abcd","efgh", "ijkl", "mnop", "rqst", "uvwx", "yz12", "3456", "7890", "-=[]", "\\;',", ".//` "},
		{"ABCD", "EFGH", "IJKL", "MNOP", "QRST", "UVWX", "YZ!@", "#$%^", "&*()", "_+{}", "|:\"<", ">?~ "},
	}
	autoCompleteWord := []string{}
	// rl.SetTraceLogLevel(rl.LogNone)  // disables all raylib log output
	windowId := filter(getCommandOutput("xdotool getactivewindow getwindowpid"), '\n')
	listChildProcId := strings.Split(getCommandOutput("pgrep -P "+windowId), "\n")
	if(len(listChildProcId)>0){
		windowId = filter(listChildProcId[0], '\n')
	}
	modeWindow := filter(getProcName(windowId), '\n')
	rl.InitWindow(screenW, screenH, "raylib-go keypress handler")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	// player := rl.NewVector2(float32(screenW/2), float32(screenH/2))

	handlers := map[int32]KeyHandler{
	}
	
	// listening to command
	cmd := exec.Command("bash", "-c", `nohup ./backgamepadkeyproc.sh`)
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
			// fmt.Println("output text")
			// fmt.Println(outputLines)
			if(len(outputLines)>0){
				fmt.Println(outputLines)
				// rl.CloseWindow()// for test
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
				if(xPos>3){
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
					xPos = 3;
				}
			}
		}
		if rl.IsKeyPressed(rl.KeyUp) {
			if isInputMoveMode == false {
				// if isTextPredictMode {
				// 	if((yPos-1) < 0){
				// 		return;
				// 	}
				// }
				yPos-=1;
				if (yPos < 0 && isTextPredictMode){
					textPredicIndex-=1
					if(textPredicIndex<0) { textPredicIndex = 0 }
					yPos = 0
				}else if(yPos<0 && isCommandKeySelect==false){
					yPos = 2;
				}else if (yPos<0 && isCommandKeySelect){
					yPos = 5;
				}
			}
		}
		if rl.IsKeyPressed(rl.KeyDown) {
			if isInputMoveMode == false {
				// if isTextPredictMode {
				// 	if((yPos+1) >= len(autoCompleteWord)){
				// 		return;
				// 	}
				// }
				yPos+=1;
				if (isTextPredictMode){
					if (yPos > 5) {
						textPredicIndex+=1
						if(textPredicIndex+5>len(autoCompleteWord)){
							textPredicIndex-=1
						}
						yPos = 5
					}
				}else if(yPos>2 && isCommandKeySelect==false){
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
				if isTextPredictMode {
					isCommandKeySelect = true
					isTextPredictMode = false
				}else if isInputMoveMode == false && yStartPos==0 && xStartPos==0 {
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
				xPagePos = 9
				yStartPos = 1
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
						}
					}else if yStartPos == 1 {
						if xStartPos == 0 {
							xPagePos+=3
							xStartPos+=1
							yStartPos = 0
						}else if xStartPos == 1 {
							isCommandKeySelect = true
							yStartPos = 0
							yPos = 0
						}
					}
				}else{
					if (len(autoCompleteWord) > 0) {
						isTextPredictMode = true
						yPos = 0
						isCommandKeySelect = false
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
				// SPACE, END, SHIFT, COPY, PASTE, CTRL, 
				case 0:
					tempInputText := inputText[:inputModeIndex]+string(" ")
					inputText = tempInputText+inputText[inputModeIndex:]
					autoComplete := checkAutoComplete(modeWindow, windowId, inputText)
					go func(){
						output := <-autoComplete
						func(out []string) {
							autoCompleteWord = []string{}
							autoCompleteWord = append(autoCompleteWord, out...)
						}(output)
					}()
					inputModeIndex+=1
				case 1:
					cmd := exec.Command("bash", "-c", `echo "`+inputText+`" | xclip -selection clipboard`)
					cmd.Stdout = io.Discard
    				cmd.Stderr = io.Discard
					cmd.Start();
					cmd2 := exec.Command("bash", "-c", `nohup ./command.sh paste >/dev/null 2>&1`)
					cmd2.Stdout = io.Discard
    				cmd2.Stderr = io.Discard
					cmd2.Start();
					rl.CloseWindow()
					return;
				case 2:
					isShiftDown = !isShiftDown
				case 3:
					cmd2 := exec.Command("bash", "-c", `nohup ./command.sh copy >/dev/null 2>&1`)
					cmd2.Stdout = io.Discard
    				cmd2.Stderr = io.Discard
					cmd2.Start();
					rl.CloseWindow()
					return;
				case 4:
					cmd := exec.Command("bash", "-c", `nohup ./command.sh getpaste >/dev/null 2>&1`)
					outputCmd, err := cmd.Output()
					if err == nil {
						tempInputText := inputText[:inputModeIndex]+string(outputCmd)
						inputText = tempInputText+inputText[inputModeIndex:]
						inputModeIndex += len(string(outputCmd))
					}
					rl.CloseWindow()
				case 5:
					inputText = "CTRL"
					inputModeIndex = 4
				default:
					inputText += "check"
					inputModeIndex+=4
				}
			} else if isTextPredictMode {
				addNewText := autoCompleteWord[min(yPos, len(autoCompleteWord)-1)]
				inputText+=addNewText
				inputModeIndex += len(addNewText)
			} else{
				indexShift := map[bool]int{true: 1, false: 0}[isShiftDown]
				addNewText := string(listKey[indexShift][(xPagePos)+yPos][xPos])
				tempInputText := inputText[:inputModeIndex]+addNewText
				inputText = tempInputText+inputText[inputModeIndex:]
				inputModeIndex += len(addNewText)
				autoComplete := checkAutoComplete(modeWindow, windowId, inputText)
				go func(){
					output := <-autoComplete
					func(out []string) {
						fmt.Println(out)
						autoCompleteWord = []string{}
						autoCompleteWord = append(autoCompleteWord, out...)
					}(output)
				}()
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
		if(isCommandKeySelect || isTextPredictMode){
			cursorUI = "________"
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
				isXMoreSpace += 70
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
			posYCursorTemp = 170
		} else if (isTextPredictMode) {
			posYCursorTemp = 280
		} else{
			posYCursorTemp = (xStartPos*70)+20+(15*xPos)
		}
		if isInputMoveMode == false {
			rl.DrawText(cursorUI, int32(posYCursorTemp), int32((yStartPos*80)+(45+(20*yPos))), 20, rl.Red)
		}
		
		if isShiftDown {
			rl.DrawText("--------", 170, 80, 20, rl.Red)
		}
		
		rl.DrawText("SPACE", 170, 40, 20, rl.DarkGray)
		rl.DrawText("END", 170, 60, 20, rl.DarkGray)
		rl.DrawText("SHIFT", 170, 80, 20, rl.DarkGray)
		rl.DrawText("COPY", 170, 100, 20, rl.DarkGray)
		rl.DrawText("PASTE", 170, 120, 20, rl.DarkGray)
		rl.DrawText("CTRL", 170, 140, 20, rl.DarkGray)

		// loop for 
		lastIndexPredict := min(textPredicIndex+6, len(autoCompleteWord))
		startIndexPredict := min(textPredicIndex, lastIndexPredict)
		idxPos := 0
		for idx := startIndexPredict; idx < lastIndexPredict; idx++ {
			rl.DrawText(autoCompleteWord[idx], 280, int32(40+(20*idxPos)), 20, rl.DarkGray)
			idxPos+=1
		}
		rl.EndDrawing()
	}
}

func checkAutoComplete(mode string, id string, inputword string) <-chan[]string {
	// incompleteWord := inputword
	autoCompleteListAsync := make(chan []string)
	autoCompleteList := []string{}
	go func() {
		fmt.Println("get mode,", mode)
		defer close(autoCompleteListAsync)
		if(mode == "terminal" || strings.Contains(mode, "zsh") || mode == "konsole" || mode == "bash"){
			separatorCmd := []string{";", "&&", "||", "&",  "|"}
			arrInput := strings.Fields(inputword)
			lastIndexCmd := lastIndexOf(arrInput, separatorCmd)
			commandNow := arrInput[min(lastIndexCmd+1, len(arrInput)-1)]
			// TODO: get history terminal
			if (commandNow == "go"){
				autoCompleteList = append(autoCompleteList, []string{
					"run", "build",
				}...)
			}
			// check what folder is terminal, and list possible file or folder in dir
			if (string(inputword[len(inputword)-1]) == " ") {
				autoCompleteList = append(autoCompleteList, getListFolder(filter(getCurrDirProcId(id), '\n'))...)
			}else{
				// get word before /
				lastCmd := string(arrInput[len(arrInput)-1])
				fmt.Println("get cmd", lastCmd)
				splitLastCmd := strings.Split(lastCmd, "/")
				dirCheck := ""
				if(len(splitLastCmd)>1){
					dirCheck+="/"
				}
				dirCheck+=strings.Join(splitLastCmd[:len(splitLastCmd)-1], "/")
				autoCompleteTarget := splitLastCmd[len(splitLastCmd)-1]
				autoCompleteListTempp := getListFolder(filter(dirCheck, '\n'))
				// append(autoCompleteList, ...)
				for _, item := range autoCompleteListTempp {
					checkCharLen := int(math.Min(float64(len(autoCompleteTarget)), float64(len(item))))
					if string(item[:checkCharLen]) == autoCompleteTarget {
						autoCompleteList = append(autoCompleteList, string(item[checkCharLen:]))
					}
				}
			}
			autoCompleteList = append(autoCompleteList, getClipboardList()...)
		}else if(mode == "browser"){ fmt.Println("sementara" )}
		autoCompleteListAsync <- autoCompleteList
		// fmt.Println(autoCompleteListAsync)
	}()
	return autoCompleteListAsync
}
func lastIndexOf(slice []string, target []string) int {
    for i := len(slice) - 1; i >= 0; i-- {
        if contains(target, slice[i]) {
			return i
		}
    }
    return 0
}
func contains[T comparable](slice []T, value T) bool {
    for _, v := range slice {
        if v == value {
            return true
        }
    }
    return false
}
func getListFolder(curr string) []string {
	// currDir := getCommandOutput("cd "+curr+" | ls | sed 's#/##'")
	currDir := getCommandOutput("ls "+curr+" | sed 's#/##'")
	return strings.Split(currDir, "\n")
}
func getCommandOutput(cmd string) string {
	cmdOut := exec.Command("bash", "-c", cmd)
	// fmt.Println("hehehe :", cmd)
	outputCmd, err := cmdOut.Output()
	if err == nil {
		// fmt.Println("hehehe :", string(outputCmd))
		return string(outputCmd)
	}
	return ""
}
func max(a int, b int) int {
	if a>b { return a } else { return b}
}
func getCurrDirProcId(procId string) string {
	return getCommandOutput(`lsof -p `+procId+` 2>/dev/null | awk '$4 == "cwd" { print $9 }'`)
}
func getProcName(procId string) string {
	return getCommandOutput(`lsof -p `+procId+` 2>/dev/null | awk '$4 == "cwd" { print $1 }'`)
}
func filter(str string, filterStr rune) string {
	output := ""
	for _, ch := range str {
		if ch != filterStr {
			output += string(ch)
		}
	}
	return output
}
func getClipboardList() []string {
	return strings.Split(getCommandOutput("gpaste-client list"), "\n")
}
func getCommandOutputAsync(cmd string) <-chan string {
	result := make(chan string)
	go func() {
		defer close(result)
		cmdOut := exec.Command("bash", "-c", cmd)
		outputCmd, err := cmdOut.Output()
		if err == nil {
			result <- string(outputCmd)
		} else {
			result <- "" // or send error message if needed
		}
	}()

	return result
}
