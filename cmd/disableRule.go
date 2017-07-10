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
	"strings"

	"github.com/gobwas/glob"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var disableRuleCmd = &cobra.Command{
	Use:   "disableRule [vserver] [target rule]",
	Short: "Disable [target rule] on [vserver]",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("Missing argument(s)")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		disableRule(args[0], args[1])
	},
}

func disableRule(targetVserver string, targetRule string) {
	if dryRun {
		fmt.Println(dryRunC)
	}

	vserverGlob := glob.MustCompile(targetVserver)
	client := initClient()

	fmt.Println("Getting vserver list from", viper.Get("vtmAPIUrl"))
	serverlist, resp, err := client.ListVirtualServers()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Response:", resp.Status)

	for _, vserver := range serverlist {
		if !vserverGlob.Match(vserver) {
			continue
		}

		r, _, err := client.GetVirtualServer(vserver)
		if err != nil {
			log.Fatal(err)
		}

		rules := *r.Basic.RequestRules
		hasUpdates := false
		hasRule := false

		for index, element := range rules {
			if strings.HasSuffix(element, targetRule) {
				hasRule = true
				if !strings.HasPrefix(element, "/") {
					rules[index] = "/" + element
					hasUpdates = true
				}
			}
		}

		if !hasRule {
			fmt.Print(vserver, ":\t(no rule)\n")
			continue
		}

		if hasUpdates {
			fmt.Print(vserver, ":\t", targetRule, " [enabled] -> [", disabledC, "]\n")
			if !dryRun {
				r.Basic.RequestRules = &rules

				_, err = client.Set(r)
				if err != nil {
					log.Fatal(err)
				}
			}
		} else {
			fmt.Print(vserver, ":\t", targetRule, " [", disabledC, "] (no change)\n")
		}
	}
}

func init() {
	vserverCmd.AddCommand(disableRuleCmd)
}
