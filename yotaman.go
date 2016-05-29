package main

import (
    "fmt"
    "os"
    "github.com/spf13/cobra"
    "./speedtest"
    "./selfcare"
)

type Config struct {
    username string
    password string
}

var RootCmd = &cobra.Command{
    Use: "yotaman",
}

var OptimizedTariffCmd = &cobra.Command{
    Use: "optimize",
    Short: "Optimize tariff speed by avaliable bandwidth",
    Run: func(cmd *cobra.Command, args []string) {
        var optimizedTariff *selfcare.Tariff

        currentSpeed := speedtest.Start()
        fmt.Println(currentSpeed)
        return

        device := selfcare.GetCurrentDevice()

        for _, tariff := range device.Tariffs {
            if currentSpeed < tariff.Speed() {
                optimizedTariff = &tariff
                break
            }
        }

        if optimizedTariff != nil {
            device.ChangeTariff(*optimizedTariff)
            fmt.Println("Tariff changed: ", optimizedTariff.Speed())
        } else {
            fmt.Println("You already have optimized tariff.")
        }
    },
}

var ListTariffCmd = &cobra.Command{
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

var SetTariffCmd = &cobra.Command{
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
    RootCmd.AddCommand(OptimizedTariffCmd)
    RootCmd.AddCommand(ListTariffCmd)
    RootCmd.AddCommand(SetTariffCmd)

    if err := RootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(-1)
    }
}
