package utils

import (
	"log"
	"syscall"

	"golang.org/x/sys/windows/registry"
)

/*
	// Detect theme
	log.Println("[onReady] Dark Mode theme is on:", utils.IsDark())
	theme = "dark"
	if utils.IsDark() {
		theme = "white"
	}

	// Monitor theme changes
	go func() { utils.Monitor(utils.React) }()
*/

const (
	regKey  = `Software\Microsoft\Windows\CurrentVersion\Themes\Personalize` // in HKCU
	regName = `SystemUsesLightTheme`                                         // <- For taskbar & tray. Use AppsUseLightTheme for apps
)

func IsDark() bool {
	k, err := registry.OpenKey(registry.CURRENT_USER, regKey, registry.QUERY_VALUE)
	if err != nil {
		log.Println("[utils/is_dark.go][IsDark] Error opening registry key. Details:")
		log.Fatal(err)
	}
	defer k.Close()
	val, _, err := k.GetIntegerValue(regName)
	if err != nil {
		log.Println("[utils/is_dark.go][IsDark] Error getting integer value. Details:")
		log.Fatal(err)
	}
	return val == 0
}

func Monitor(fn func(bool)) {
	var regNotifyChangeKeyValue *syscall.Proc
	changed := make(chan bool)

	if advapi32, err := syscall.LoadDLL("Advapi32.dll"); err == nil {
		if p, err := advapi32.FindProc("RegNotifyChangeKeyValue"); err == nil {
			regNotifyChangeKeyValue = p
		} else {
			log.Println("[utils/is_dark.go][Monitor] Could not find function RegNotifyChangeKeyValue in Advapi32.dll. Details:")
			log.Fatal(err)
		}
	}
	if regNotifyChangeKeyValue != nil {
		go func() {
			k, err := registry.OpenKey(registry.CURRENT_USER, regKey, syscall.KEY_NOTIFY|registry.QUERY_VALUE)
			if err != nil {
				log.Println("[utils/is_dark.go][Monitor] Error opening registry key. Details:")
				log.Fatal(err)
			}
			var wasDark uint64
			for {
				regNotifyChangeKeyValue.Call(uintptr(k), 0, 0x00000001|0x00000004, 0, 0)
				val, _, err := k.GetIntegerValue(regName)
				if err != nil {
					log.Println("[utils/is_dark.go][Monitor] Error getting integer value. Details:")
					log.Fatal(err)
				}
				if val != wasDark {
					wasDark = val
					changed <- val == 0
				}
			}
		}()
	}
	for {
		val := <-changed
		fn(val)
	}

}

// React to the change
func React(isDark bool) {
	if isDark {
		log.Println("[utils/is_dark.go][React] Dark Mode ON")
	} else {
		log.Println("[utils/is_dark.go][React] Dark Mode OFF")
	}
}
