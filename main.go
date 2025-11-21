package main

import (
	"encoding/base64"
	"fmt"
	"image/color"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"golang.design/x/hotkey"

	"github.com/Minsuh1204/PLATiNA-ARCHiVE-Go-Client/client"
)

const VERSION = "1.0.0"

var APIKey string
var b64APIKey string
var decoderName string
var logLabel *widget.Label
var songTitleLabel *widget.Label
var songLevelLabel *widget.Label
var judgeLabel *widget.Label
var scoreLabel *widget.Label
var patchLabel *widget.Label
var jacketContainer *fyne.Container
var analyzeButton *widget.Button

func main() {
	a := app.New()
	w := a.NewWindow(fmt.Sprintf("PLATiNA-ARCHiVE v%s", VERSION))
	w.Resize(fyne.NewSize(800, 600))
	w.SetFixedSize(true)

	jacketPlaceholder := canvas.NewRectangle(color.Black)
	jacketPlaceholder.SetMinSize(fyne.NewSize(200, 200))
	jacketContainer = container.NewStack(jacketPlaceholder)

	songTitleLabel = widget.NewLabel("")
	songLevelLabel = widget.NewLabel("")
	judgeLabel = widget.NewLabel("")
	scoreLabel = widget.NewLabel("")
	patchLabel = widget.NewLabel("")
	statsContainer := container.NewVBox(songTitleLabel, songLevelLabel, judgeLabel, scoreLabel, patchLabel)
	topContainer := container.NewHBox(jacketContainer, layout.NewSpacer(), statsContainer)
	paddedTopContainer := container.NewPadded(topContainer)

	analyzeButton = widget.NewButton("Analyze", startAnalyze)
	buttonContainer := container.New(layout.NewCenterLayout(), analyzeButton)

	logLabel = widget.NewLabel("")
	logScroll := container.NewScroll(logLabel)
	logScroll.SetMinSize(fyne.NewSize(0, 400))
	mainContainer := container.NewVBox(paddedTopContainer, canvas.NewLine(color.Gray{}), buttonContainer, canvas.NewLine(color.Gray{}), logScroll)
	w.SetContent(mainContainer)

	go registerHotkeys()
	APIKey = client.LoadAPIKey()
	if APIKey == "" {
		showWelcomeDialog(w)
	} else {
		decoderName = strings.Split(APIKey, "::")[0]
		b64APIKey = base64.StdEncoding.EncodeToString([]byte(APIKey))
		logMessage(fmt.Sprintf("환영합니다, %s님!", decoderName))
	}

	w.ShowAndRun()
}

func showWelcomeDialog(w fyne.Window) {
	var d dialog.Dialog
	transitioning := false

	loginBtn := widget.NewButton("Login", func() {
		transitioning = true
		d.Hide()
		showLoginDialog(w)
	})
	registerBtn := widget.NewButton("Register", func() {
		transitioning = true
		d.Hide()
		showRegisterDialog(w)
	})

	content := container.NewVBox(
		widget.NewLabel("Please login or register to continue."),
		loginBtn,
		registerBtn,
	)

	d = dialog.NewCustom("Welcome", "Quit", content, w)
	d.SetOnClosed(func() {
		if !transitioning && APIKey == "" {
			w.Close()
		}
	})
	d.Show()
}

func showLoginDialog(w fyne.Window) {
	nameEntry := widget.NewEntry()
	nameEntry.PlaceHolder = "Username"
	passEntry := widget.NewPasswordEntry()
	passEntry.PlaceHolder = "Password"

	d := dialog.NewForm("Login", "Login", "Back", []*widget.FormItem{
		widget.NewFormItem("Username", nameEntry),
		widget.NewFormItem("Password", passEntry),
	}, func(confirm bool) {
		if !confirm {
			showWelcomeDialog(w)
			return
		}
		name := nameEntry.Text
		pass := passEntry.Text
		result, err := client.Login(name, pass)
		if err != nil {
			errDialog := dialog.NewError(fmt.Errorf("login failed: %v", err), w)
			errDialog.SetOnClosed(func() {
				showLoginDialog(w)
			})
			errDialog.Show()
			return
		}

		handleAuthSuccess(result.APIKey)
	}, w)
	d.Resize(fyne.NewSize(300, 200))
	d.Show()
}

func showRegisterDialog(w fyne.Window) {
	nameEntry := widget.NewEntry()
	nameEntry.PlaceHolder = "Username"
	passEntry := widget.NewPasswordEntry()
	passEntry.PlaceHolder = "Password"

	d := dialog.NewForm("Register", "Register", "Back", []*widget.FormItem{
		widget.NewFormItem("Username", nameEntry),
		widget.NewFormItem("Password", passEntry),
	}, func(confirm bool) {
		if !confirm {
			showWelcomeDialog(w)
			return
		}
		name := nameEntry.Text
		pass := passEntry.Text
		result, err := client.Register(name, pass)
		if err != nil {
			errDialog := dialog.NewError(fmt.Errorf("register failed: %v", err), w)
			errDialog.SetOnClosed(func() {
				showRegisterDialog(w)
			})
			errDialog.Show()
			return
		}

		handleAuthSuccess(result.APIKey)
	}, w)
	d.Resize(fyne.NewSize(300, 200))
	d.Show()
}

func handleAuthSuccess(key string) {
	err := client.SaveAPIKey(key)
	if err != nil {
		logMessage(fmt.Sprintf("Failed to save API key: %v", err))
	}

	APIKey = key
	decoderName = strings.Split(APIKey, "::")[0]
	b64APIKey = base64.StdEncoding.EncodeToString([]byte(APIKey))
	logMessage(fmt.Sprintf("Success! Welcome, %s!", decoderName))
}

func logMessage(msg string) {
	logLabel.SetText(logLabel.Text + fmt.Sprintf("[%v] %v\n", client.FormatCurrentTime(), msg))
}

func registerHotkeys() {
	keyInsertWin := hotkey.Key(0x2D) // Insert key for Windows
	// keyInsertMac := hotkey.Key0 // Testing key for Mac
	hk := hotkey.New([]hotkey.Modifier{hotkey.ModAlt}, keyInsertWin)
	if err := hk.Register(); err != nil {
		logMessage(fmt.Sprintf("Failed to register hotkey: %v", err))
		return
	}
	for range hk.Keydown() {
		fyne.Do(func() {
			analyzeButton.Tapped(&fyne.PointEvent{})
		})
	}
}

func startAnalyze() {
	if analyzeButton.Disabled() {
		return
	}
	analyzeButton.Disable()
	logMessage("Analyze started...")

	go func() {
		defer func() {
			fyne.Do(func() {
				analyzeButton.Enable()
			})
		}()

		// Simulate heavy background work (OCR)
		time.Sleep(2 * time.Second)

		fyne.Do(func() {
			logMessage("Analyze finished!")
		})
	}()
}
