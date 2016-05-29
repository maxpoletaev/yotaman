package main

import (
	"fmt"
	"os"
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
		selfcare.AutoLogin()
		device := selfcare.GetCurrentDevice()
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

		device := selfcare.GetCurrentDevice()
		for _, tariff := range device.Tariffs {
			if tariff.Label() == newLabel {
				device.ChangeTariff(tariff)
				break
			}
		}
	},
}

func main() {
	rootCmd.AddCommand(listTariffCmd)
	rootCmd.AddCommand(setTariffCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
