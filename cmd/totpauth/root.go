package totpauth

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"

	totp "github.com/Appelfeldt/totpauth/pkg/totpauth"
	"github.com/spf13/cobra"
)

var BuildVersion string

var rootCmd = &cobra.Command{
	Use:     "totpauth",
	Version: BuildVersion,
	Short:   "totpauth - Generates time-based one-time passwords",
	Long:    "totpauth is a CLI tool that generates time-based one-time passwords using a supplied key",
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var input string
		var inputReader io.Reader = cmd.InOrStdin()

		if hasPipedInput() {
			var err error
			input, err = bufio.NewReader(inputReader).ReadString('\n')
			if err != nil {
				fmt.Fprintf(os.Stderr, "error reading input: %v\n", err)
				os.Exit(1)
			}
		} else {
			if len(args) > 0 && args[0] != "-" {
				file, err := os.Open(args[0])
				if err != nil {
					fmt.Fprintf(os.Stderr, "failure opening file: %v", err)
					os.Exit(1)
				}
				inputReader = file
			}

			var buffer bytes.Buffer
			io.Copy(&buffer, inputReader)

			input = buffer.String()

		}
		start, _ := cmd.Flags().GetUint64("timestart")
		step, _ := cmd.Flags().GetUint64("timestep")
		if step == 0 {
			fmt.Fprint(os.Stderr, "timestep value cannot be zero")
			os.Exit(1)
		}

		password := totp.Auth(input, start, step)
		fmt.Println(password)
	},
}

func init() {
	rootCmd.PersistentFlags().Uint64("timestart", 0, "")
	rootCmd.PersistentFlags().Uint64("timestep", 30, "")
}

func hasPipedInput() bool {
	fileInfo, _ := os.Stdin.Stat()
	return fileInfo.Mode()&os.ModeCharDevice == 0
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred '%s'", err)
		os.Exit(1)
	}
}
