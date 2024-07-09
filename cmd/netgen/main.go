package netgen

import (
	"fmt"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "netgen",
		Short: "Generate a network with specified validators, commit, and chain ID",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			// // Parse the number of validators
			// numValidators, err := strconv.Atoi(args[0])
			// if err != nil {
			// 	log.Fatalf("Invalid number of validators: %s", err)
			// }

			// chainID := args[1]

			// // Call the function to generate the network with the provided arguments
			// n, err := NewNetwork(chainID, numValidators, numValidators)
			// if err != nil {
			// 	return err
			// }

			// fmt.Println("calling init")

			// return n.InitNodes("~/payload")
			return nil
		},
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
