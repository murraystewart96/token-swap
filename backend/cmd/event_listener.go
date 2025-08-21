package cmd

import (
	"github.com/spf13/cobra"
)

func createEventListenerCmd() *cobra.Command { //nolint:gocognit
	clientCmd := &cobra.Command{
		Use:   "event-listener",
		Short: "event listener",
		Long:  `Listens for transaction events and published to event queue`,

		Run: func(cmd *cobra.Command, _ []string) {

		},
	}

	return clientCmd
}
