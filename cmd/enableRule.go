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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var enableRuleCmd = &cobra.Command{
	Use:   "enableRule [target rule]",
	Short: "Enable [target rule] on all vservers.",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("Missing argument. Please provide the [target rule] name.")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		enableRule(args[0])
	},
}

func enableRule(targetRule string) {
	client := initClient()

	fmt.Println("Getting vserver list from", viper.Get("vtmAPIUrl"))
	serverlist, resp, err := client.ListVirtualServers()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Response:", resp.Status)

	for _, vserver := range serverlist {
		fmt.Println("Processing vserver", vserver, "..")
		r, resp, err := client.GetVirtualServer(vserver)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Response:", resp.Status)

		rules := *r.Basic.RequestRules
		wasUpdated := false

		for index, element := range rules {
			if strings.HasPrefix(element, "/") && strings.HasSuffix(element, targetRule) {
				rules[index] = strings.TrimPrefix(element, "/")
				wasUpdated = true
			}
		}

		if wasUpdated {
			r.Basic.RequestRules = &rules

			resp, err = client.Set(r)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Response:", resp.Status)
		}
	}
}

func init() {
	vserverCmd.AddCommand(enableRuleCmd)

}
