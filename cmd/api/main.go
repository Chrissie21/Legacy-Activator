package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	// TODO: To be replaced with absolute module name if necessary
	"legacy-activation/pkg/exploits"

	"legacy-activation/pkg/exploits/common"
	"legacy-activation/pkg/exploits/itunesstored"
)

// Constants
const (
	AppID             = "com.legacy-activation.app"
	DefaultBackendURL = "http://localhost/backend/server.php" // Points to your unported PHP backend
	DefaultAttempts   = 5
)

// AppSettings stores user preference bindings natively
type AppSettings struct {
	BackendURL string
	Attempts   int
}

// MainWindowState encapsulates UI state tracking
type MainWindowState struct {
	sync.Mutex
	currentVersion  []int
	currentSettings AppSettings
	deviceConnected bool
}

func main() {
	myApp := app.NewWithID(AppID)
	myWindow := myApp.NewWindow("legacy-activator v1.2.2")
	myWindow.Resize(fyne.NewSize(500, 200))
	myWindow.SetFixedSize(true)

	// Load stored app configurations or fall back to defaults
	prefs := myApp.Preferences()
	state := &MainWindowState{
		currentSettings: AppSettings{
			BackendURL: prefs.StringWithFallback("backend_url", DefaultBackendURL),
			Attempts:   prefs.IntWithFallback("attempts", DefaultAttempts),
		},
	}

	// UI Layout Elements
	statusLabel := widget.NewLabel("No device connected")
	activateBtn := widget.NewButton("Activate Device", nil)
	activateBtn.Disable()

	// Settings Gear Dialog Configuration
	settingsBtn := widget.NewButton("⚙", func() { // Had to port ts btw
		showSettingsDialog(myWindow, state, prefs)
	})

	// Wrap elements inside layouts
	buttonRow := container.New(layout.NewHBoxLayout(),
		container.NewGridWrap(fyne.NewSize(400, 40), activateBtn),
		layout.NewSpacer(),
		container.NewGridWrap(fyne.NewSize(40, 40), settingsBtn),
	)

	mainLayout := container.New(layout.NewVBoxLayout(),
		statusLabel,
		layout.NewSpacer(),
		buttonRow,
	)
	myWindow.SetContent(mainLayout)

	// Activation processing thread binding
	activateBtn.OnTapped = func() {
		var note string
		// Condition check logic
		if itunesstored.SupportsVersion(state.currentVersion) {
			note = "Please ensure it is connected to Wi-Fi."
		} else {
			note = "Your device will reboot during the process."
		}

		dialog.ShowInformation("Info", "Your device will now be activated. "+note, myWindow)
		activateBtn.Disable()
		statusLabel.SetText("Activating device...")

		// ActivationThread worker run
		go func() {
			// 1. Usbmuxd connection check
			lockdown, err := common.CreateUsingUsbmux()
			if err != nil {
				fyne.Do(func() {
					dialog.ShowError(fmt.Errorf("device disconnection encountered: %v", err), myWindow)
					statusLabel.SetText("Error occurred")
				})
				return
			}

			// 2. Perform activation logic verification
			if lockdown.ActivationState == "Activated" {
				fyne.Do(func() {
					dialog.ShowInformation("Success", "Device is already activated", myWindow)
					statusLabel.SetText("Device is already activated")
				})
				return
			}

			// 3. Match version parsing to extract target exploit rules
			exploit, err := exploits.GetExploit(common.ParseVersion(lockdown.ProductVersion), lockdown.UDID)
			if err != nil || exploit == nil {
				fyne.Do(func() {
					dialog.ShowError(fmt.Errorf("unsupported iOS version"), myWindow)
					statusLabel.SetText("Error occurred")
				})
				return
			}

			// 4. Fire the exploit logic using settings rules asynchronously
			err = exploit.Run(lockdown, state.currentSettings.BackendURL, state.currentSettings.Attempts, func(logMsg string) {
				fyne.Do(func() { statusLabel.SetText(logMsg) })
			})

			fyne.Do(func() {
				if err != nil {
					dialog.ShowError(err, myWindow)
					statusLabel.SetText("Error occurred")
				} else {
					dialog.ShowInformation("Success", "Done!", myWindow)
					statusLabel.SetText("Done!")
				}
				activateBtn.Enable()
			})
		}()
	}

	// Timer device poller loop running every 1000ms
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		for range ticker.C {
			lockdown, err := common.CreateUsingUsbmux()
			if err != nil {
				fyne.Do(func() {
					statusLabel.SetText("No device connected")
					activateBtn.Disable()
				})
				continue
			}

			vTuple := common.ParseVersion(lockdown.ProductVersion)
			exploit, err := exploits.GetExploit(vTuple, lockdown.UDID)
			if err != nil || exploit == nil {
				fyne.Do(func() {
					statusLabel.SetText(fmt.Sprintf("Unsupported %s iOS version: %s", lockdown.ProductType, lockdown.ProductVersion))
					activateBtn.Disable()
				})
				continue
			}

			// Call exploit back-end validator checks (ensuring PHP endpoint connection)
			supported, err := exploit.Available(lockdown.ProductType, lockdown.BuildVersion, lockdown.UDID, state.currentSettings.BackendURL)
			if err != nil {
				fyne.Do(func() {
					statusLabel.SetText("Could not reach backend. Please check your internet connection!")
					activateBtn.Disable()
				})
				continue
			}

			if !supported {
				fyne.Do(func() {
					statusLabel.SetText(fmt.Sprintf("Device variant unavailable for %s (%s)", lockdown.ProductType, lockdown.ProductVersion))
					activateBtn.Disable()
				})
				continue
			}

			// Update atomic window state variables
			state.Lock()
			state.currentVersion = vTuple
			state.deviceConnected = true
			state.Unlock()

			fyne.Do(func() {
				statusLabel.SetText(fmt.Sprintf("Connected: %s (%s)", lockdown.ProductType, lockdown.ProductVersion))
				activateBtn.Enable()
			})
		}
	}()

	myWindow.ShowAndRun()
}

// showSettingsDialog generates the custom interactive prompt box natively
func showSettingsDialog(parent fyne.Window, state *MainWindowState, prefs fyne.Preferences) {
	urlEntry := widget.NewEntry()
	urlEntry.SetText(state.currentSettings.BackendURL)
	urlEntry.PlaceHolder = DefaultBackendURL

	attemptsEntry := widget.NewEntry()
	attemptsEntry.SetText(strconv.Itoa(state.currentSettings.Attempts))

	form := widget.NewForm(
		widget.NewFormItem("Backend URL:", urlEntry),
		widget.NewFormItem("Attempts:", attemptsEntry),
	)

	dialog.ShowCustomConfirm("Settings", "Save", "Cancel", form, func(confirmed bool) {
		if confirmed {
			state.Lock()
			defer state.Unlock()

			// Sanitize input entries
			urlText := strings.TrimSpace(urlEntry.Text)
			if urlText == "" {
				urlText = DefaultBackendURL
			}

			attemptsVal, err := strconv.Atoi(attemptsEntry.Text)
			if err != nil || attemptsVal < 1 {
				attemptsVal = DefaultAttempts
			}

			state.currentSettings.BackendURL = urlText
			state.currentSettings.Attempts = attemptsVal

			// Save keys persistent across binary invocations
			prefs.SetString("backend_url", urlText)
			prefs.SetInt("attempts", attemptsVal)
		}
	}, parent)
}
