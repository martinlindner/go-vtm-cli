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
	"os"
	"strings"
	"text/tabwriter"

	"github.com/gobwas/glob"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// getRuleStateCmd represents the getRuleState command
var getRuleStateCmd = &cobra.Command{
	Use:   "getRuleState [vserver] [target rule]",
	Short: "Get status of rules matching [target rule] on [vserver]",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("Missing argument(s)")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		getRuleState(args[0], args[1])
	},
}

func getRuleState(targetVserver, targetRule string) {
	vserverGlob := glob.MustCompile(targetVserver)
	ruleGlob := glob.MustCompile(targetRule)
	client := initClient()

	fmt.Println("Getting vserver list from", viper.Get("vtmAPIUrl"))
	serverlist, resp, err := client.ListVirtualServers()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Response:", resp.Status)

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 5, 2, ' ', 0)

	for _, vserver := range serverlist {
		if !vserverGlob.Match(vserver) {
			continue
		}

		fmt.Fprint(w, vserver, ":\n")

		r, _, err := client.GetVirtualServer(vserver)
		if err != nil {
			log.Fatal(err)
		}

		rules := *r.Basic.RequestRules

		hasRule := false

		for _, element := range rules {
			currentRule := strings.TrimPrefix(element, "/")

			if ruleGlob.Match(currentRule) {
				hasRule = true

				currentRuleState := enabledC
				if strings.HasPrefix(element, "/") {
					currentRuleState = disabledC
				}

				fmt.Fprint(w, "\t- ", currentRule, ":\t[", currentRuleState, "]\n")
			}
		}

		if !hasRule {
			fmt.Fprint(w, "\t(no matching rule)\n")
		}
	}

	w.Flush()
}

func init() {
	vserverCmd.AddCommand(getRuleStateCmd)
}
