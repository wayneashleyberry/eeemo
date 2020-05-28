// Zalgo Text Generator by Tchouky
// http://www.eeemo.net/
package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"github.com/wayneashleyberry/eeemo/pkg/zalgo"
)

func main() {
	ctx := context.Background()

	var size string
	var up, middle, down bool
	middle, down = true, true

	cmd := &cobra.Command{
		Use:  "eeemo [text]",
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if size != "mini" && size != "normal" && size != "maxi" {
				return errors.New("invalid size")
			}

			var input string

			file := os.Stdin
			fi, err := file.Stat()
			if err != nil {
				return err
			}
			fsize := fi.Size()
			if fsize == 0 {
				if len(args) > 0 {
					input = args[0]
				}
			}

			if fsize > 0 {
				b, err := ioutil.ReadAll(file)
				if err != nil {
					return err
				}

				input = string(b)
			}

			if input == "" {
				return errors.New("couldn't read input")
			}

			fmt.Print(
				zalgo.Generate(input, size, up, middle, down),
			)

			return nil
		},
	}

	cmd.Flags().StringVarP(&size, "size", "", "mini", "'mini', 'normal' or 'maxi'")
	cmd.Flags().BoolVarP(&up, "up", "", true, "")
	cmd.Flags().BoolVarP(&middle, "middle", "", true, "")
	cmd.Flags().BoolVarP(&down, "down", "", false, "")

	err := cmd.ExecuteContext(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
