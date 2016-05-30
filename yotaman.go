package main

import (
	"os"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zenwalker/yotaman/selfcare"
)

var rootCmd = &cobra.Command{
	Use: "yotaman",
}

var listTariffCmd = &cobra.Command{
	Use: "list",
	Short: "Show avaliable tariffs",

	Run: func(cmd *cobra.Command, args []string) {
		err := selfcare.AutoLogin()
		if err != nil { exitWithError(err) }

		device, err := selfcare.GetCurrentDevice()
		if err != nil { exitWithError(err) }

		for _, tariff := range device.Tariffs {
			var currentFlag string
			if device.IsCurrentTariff(tariff) {
				currentFlag = ">"
			}
			output := fmt.Sprintf("%1s %4s | %s", currentFlag, tariff.Label(), tariff.Repr())
			fmt.Println(output)
		}
	},
}

var setTariffCmd = &cobra.Command{
	Use: "set",
	Short: "Change current tariff",

	Run: func(cmd *cobra.Command, args []string) {
		selfcare.AutoLogin()
		newLabel := args[0] // TODO: check exists

		device, err := selfcare.GetCurrentDevice()
		if err != nil { exitWithError(err) }
		isFound := false

		for _, tariff := range device.Tariffs {
			if tariff.Label() == newLabel {
				err = device.ChangeTariff(tariff)
				if err != nil { exitWithError(err) }
				fmt.Println("Tariff changed:", tariff.Repr())
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
	fmt.Println(err)
	os.Exit(-1)
}

func main() {
	rootCmd.AddCommand(listTariffCmd)
	rootCmd.AddCommand(setTariffCmd)

	if err := rootCmd.Execute(); err != nil {
		exitWithError(err)
	}
}
