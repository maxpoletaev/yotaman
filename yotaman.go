package main

import (
	"os"
	"fmt"
	"time"
	"errors"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/briandowns/spinner"
	"github.com/zenwalker/yotaman/selfcare"
)

func newSpinner() *spinner.Spinner {
	charSet := []string{"(|)", "(/)", "(—)", "(\\)"}
	spin := spinner.New(charSet, 100 * time.Millisecond)
	spin.Suffix = " retrieving data"
	spin.Color("blue")
	spin.Stop()
	return spin
}

var (
	red = color.New(color.FgRed).SprintFunc()
	blue = color.New(color.FgBlue).SprintFunc()
	green = color.New(color.FgGreen).SprintFunc()
	spin = newSpinner()
)

var rootCmd = &cobra.Command{
	Use: "yotaman",
}

var listTariffCmd = &cobra.Command{
	Use: "list",
	Short: "Show avaliable tariffs",

	Run: func(cmd *cobra.Command, args []string) {
		spin.Start()
		err := selfcare.AutoLogin()
		if err != nil { exitWithError(err) }

		device, err := selfcare.GetCurrentDevice()
		if err != nil { exitWithError(err) }
		spin.Stop()

		for _, tariff := range device.Tariffs {
			var out string
			if device.IsCurrentTariff(tariff) {
				out = blue(fmt.Sprintf("> %4s | %s", tariff.Label(), tariff.Repr()))
			} else {
				out = fmt.Sprintf("  %4s | %s", tariff.Label(), tariff.Repr())
			}
			fmt.Println(out)
		}
	},
}

var setTariffCmd = &cobra.Command{
	Use: "set",
	Short: "Change current tariff",

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			exitWithError(errors.New("set needs a speed argument"))
		}
		spin.Start()

		newLabel := args[0]
		err := selfcare.AutoLogin()
		if err != nil { exitWithError(err) }

		device, err := selfcare.GetCurrentDevice()
		if err != nil { exitWithError(err) }
		isFound := false

		for _, tariff := range device.Tariffs {
			if tariff.Label() == newLabel {
				err = device.ChangeTariff(tariff)
				if err != nil { exitWithError(err) }

				spin.Stop()
				fmt.Println(green("Tariff changed: ", tariff.Repr()))

				isFound = true
				break
			}
		}

		if !isFound {
			exitWithError(fmt.Errorf("Tariff %s not found", newLabel))
		}
	},
}

func exitWithError(err error) {
	spin.Stop()
	fmt.Println(red(err.Error()))
	os.Exit(-1)
}

func main() {
	rootCmd.AddCommand(listTariffCmd)
	rootCmd.AddCommand(setTariffCmd)

	if err := rootCmd.Execute(); err != nil {
		exitWithError(err)
	}
}
