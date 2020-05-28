package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	ctx := context.Background()

	var size string
	var up, middle, down bool
	middle, down = true, true

	cmd := &cobra.Command{
		Use: "eeemo",
		RunE: func(cmd *cobra.Command, args []string) error {
			if size != "mini" && size != "normal" && size != "maxi" {
				return errors.New("invalid size")
			}

			text := HECOMES("...", size, up, middle, down)

			fmt.Print(text)

			return nil
		},
	}

	cmd.Flags().StringVarP(&size, "size", "", "mini", "Valid sizes are 'mini', 'normal' and 'maxi'")
	cmd.Flags().BoolVarP(&up, "up", "", true, "")
	cmd.Flags().BoolVarP(&middle, "middle", "", true, "")
	cmd.Flags().BoolVarP(&down, "down", "", false, "")

	err := cmd.ExecuteContext(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func HECOMES(s string, size string, up bool, middle bool, down bool) string {
	return s
}
