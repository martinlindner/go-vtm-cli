// Copyright © 2017 Martin Lindner <mlindner@gaba.co.jp>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/gobwas/glob"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// setMaxReplyTimeCmd represents the setMaxReplyTime command
var setMaxReplyTimeCmd = &cobra.Command{
	Use:   "setMaxReplyTime [pool] [maxReplyTime]",
	Short: "Set maximum reply time for [pool] to [maxReplyTime].",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("Missing argument(s)")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		maxReplyTime, err := strconv.Atoi(args[1])

		if err != nil {
			log.Fatal(err)
		}

		setMaxReplyTime(args[0], maxReplyTime)
	},
}

func setMaxReplyTime(targetPool string, maxReplyTime int) {
	if dryRun {
		fmt.Println("Note: Dry-Run!")
	}

	poolGlob := glob.MustCompile(targetPool)
	client := initClient()

	fmt.Println("Getting pool list from", viper.Get("vtmAPIUrl"))
	poollist, resp, err := client.ListPools()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Response:", resp.Status)

	for _, pool := range poollist {
		if !poolGlob.Match(pool) {
			continue
		}

		r, _, err := client.GetPool(pool)
		if err != nil {
			log.Fatal(err)
		}

		currentMaxReplyTime := *r.Connection.MaxReplyTime

		if currentMaxReplyTime != maxReplyTime {
			fmt.Print(pool, ":\t", currentMaxReplyTime, "s -> ", maxReplyTime, "s\n")

			if !dryRun {
				*r.Connection.MaxReplyTime = maxReplyTime
				_, err = client.Set(r)
				if err != nil {
					log.Fatal(err)
				}
			}
		} else {
			fmt.Print(pool, ":\t", currentMaxReplyTime, "s (no change)\n")
		}

	}
}

func init() {
	poolCmd.AddCommand(setMaxReplyTimeCmd)
}
