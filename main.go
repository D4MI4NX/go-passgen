/* go-passgen is a graphical tool for generating random passwords.

Copyright (C) 2024  D4MI4NX

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation version 3.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>*/

package main

import (
    "crypto/rand"
    "errors"
    "flag"
    "fmt"

    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/canvas"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/dialog"
    "fyne.io/fyne/v2/driver/desktop"
    "fyne.io/fyne/v2/layout"
    "fyne.io/fyne/v2/widget"

    passwordvalidator "github.com/wagslane/go-password-validator"
    "image/color"
    "math"
    "math/big"
    mrand "math/rand"
    "strconv"
)

type toUseStruct struct {
    lowercase   bool
    uppercase   bool
    numbers     bool
    symbols     bool
}

type charStruct struct {
    lowercase   []string
    uppercase   []string
    numbers     []string
    symbols     []string
}

var length int = 16
var useableChars []string
var characters charStruct
var toUse toUseStruct
var showLicense bool

var w fyne.Window

func main() {
    fmt.Println("starting...")
    flags()
    if showLicense {
        printLicense()
    }
    setCharacters()
    getUseableCharacters()
    gui()
}

func flags() {
    flag.IntVar(&length, "l", 16, "Specify password length")
    flag.BoolVar(&showLicense, "L", false, "Print license")

    flag.Parse()
}

func setCharacters() {
    characters = charStruct{
        []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"},
        []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"},
        []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"},
        []string{"!", "#", "$", "%", "&", "+", "-", ".", "<", "=", ">", "?", "@"},
    }

    toUse = toUseStruct{true, true, true, true}
}

func getUseableCharacters() {
    useableChars = []string{}

    if toUse.lowercase {
        useableChars = append(useableChars, characters.lowercase...)
    }

    if toUse.uppercase {
        useableChars = append(useableChars, characters.uppercase...)
    }

    if toUse.numbers {
        useableChars = append(useableChars, characters.numbers...)
    }

    if toUse.symbols {
        useableChars = append(useableChars, characters.symbols...)
    }

    if len(useableChars) == 0 {
        useableChars = []string{""}
    }
}

func genPassword() (string, error) {
    var password string
    var err error
    var nBig *big.Int
    var n int64

    for i := 0; i < length; i++ {
        nBig, err = rand.Int(rand.Reader, big.NewInt(int64(len(useableChars))))
        if err != nil {
            n = int64(mrand.Intn(len(useableChars)))
        } else {
            n = nBig.Int64()
        }

        password += useableChars[n]
    }

    return password, err
}

func gui() {
    a := app.New()
    w = a.NewWindow("PassGen")
    w.SetMaster()
    w.Resize(fyne.NewSize(300, 0))

    ctrlQ := &desktop.CustomShortcut{KeyName: fyne.KeyQ, Modifier: fyne.KeyModifierControl}
    w.Canvas().AddShortcut(ctrlQ, func(shortcut fyne.Shortcut) {
		fmt.Println("exiting...")
		a.Quit()
	})


    lowercaseCheck := widget.NewCheck("lowercase (abc)", func(value bool) {
        toUse.lowercase = value
        getUseableCharacters()
    })
    lowercaseCheck.SetChecked(true)

    uppercaseCheck := widget.NewCheck("uppercase (ABC)", func(value bool) {
        toUse.uppercase = value
        getUseableCharacters()
    })
    uppercaseCheck.SetChecked(true)

    numberCheck := widget.NewCheck("numbers (123)", func(value bool) {
        toUse.numbers = value
        getUseableCharacters()
    })
    numberCheck.SetChecked(true)

    symbolCheck := widget.NewCheck("symbols (!#$)", func(value bool) {
        toUse.symbols = value
        getUseableCharacters()
    })
    symbolCheck.SetChecked(true)

    lengthLabel := widget.NewLabel("Length: ")
    lengthEntry := widget.NewEntry()
    lengthEntry.Text = strconv.Itoa(length)
    lengthEntry.OnChanged = func(number string) {
        i, err := strconv.Atoi(number)
        if err != nil {
            fmt.Println("Length not a number!")
        } else if 0 < i {
            length = i
        }
    }

    passEntry := widget.NewEntry()
    password, err := genPassword()
    if err != nil {
        fmt.Printf("Falling back to math/rand: %v\n", err)
        dialog.ShowError(errors.New("Password generation failed\nusing secure method.\nFalling back to\nless secure method."), w)
    }
    passEntry.Text = password

    strengthLabel := canvas.NewText(fmt.Sprint("strength: ", getPasswordStrength(passEntry.Text)), passwordStrengthToRGB(passEntry.Text))

    passEntry.OnChanged = func(text string) {
        updatePasswordStrength(strengthLabel, passEntry)
    }

    genPassButton := widget.NewButton("GENERATE", func() {
        passEntry.Text, _ = genPassword()
        passEntry.Refresh()
        updatePasswordStrength(strengthLabel, passEntry)
    })

    copyButton := widget.NewButton("COPY", func() {
        w.Clipboard().SetContent(passEntry.Text)
    })

    w.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
        switch(k.Name) {
            case "Return":
                passEntry.Text, _ = genPassword()
                passEntry.Refresh()
                updatePasswordStrength(strengthLabel, passEntry)

            case "C":
                w.Clipboard().SetContent(passEntry.Text)

            case "1":
                toUse.lowercase = !lowercaseCheck.Checked
                updateCheck(lowercaseCheck)

            case "2":
                toUse.uppercase = !uppercaseCheck.Checked
                updateCheck(uppercaseCheck)

            case "3":
                toUse.numbers = !numberCheck.Checked
                updateCheck(numberCheck)

            case "4":
                toUse.symbols = !symbolCheck.Checked
                updateCheck(symbolCheck)
        }
    })

    w.SetContent(
        container.NewVBox(
            lowercaseCheck,
            uppercaseCheck,
            numberCheck,
            symbolCheck,
            container.New(
                layout.NewGridLayout(2),
                lengthLabel,
                lengthEntry,
            ),
            passEntry,
            strengthLabel,
            genPassButton,
            copyButton,
        ),
    )

    w.ShowAndRun()
}

func updateCheck(check *widget.Check) {
    check.Checked = !check.Checked
    getUseableCharacters()
    check.Refresh()
}

func getPasswordStrength(password string) float64 {
    return math.Round(passwordvalidator.GetEntropy(password))
}

func interpolate(start, end, t float64) float64 {
    return start*(1-t) + end*t
}

func passwordStrengthToRGB(password string) color.Color {
    strength := int(math.Max(0, math.Min(100, float64(getPasswordStrength(password)))))

    red := 255.0
    green := 0.0

    t := float64(strength) / 100.0

    r := uint8(interpolate(red, green, t))
    g := uint8(interpolate(green, red, t))
    b := uint8(0)
    a := uint8(255)

    return color.RGBA{r, g, b, a}
}

func updatePasswordStrength(label *canvas.Text, entry *widget.Entry) {
    label.Text = fmt.Sprint("strength: ", getPasswordStrength(entry.Text))
    label.Color = passwordStrengthToRGB(entry.Text)
    label.Refresh()
}

func printLicense() {
    fmt.Println(
`---------------------------- LICENSE ----------------------------
go-passgen is a graphical tool for generating random passwords.

Copyright (C) 2024  D4MI4NX

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation version 3.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>
-----------------------------------------------------------------`)

}
