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
	"text/tabwriter"

	"github.com/gobwas/glob"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// getTimeoutCmd represents the getTimeout command
var getTimeoutCmd = &cobra.Command{
	Use:   "getTimeout [vserver]",
	Short: "Get connection timeout(s) for [vserver].",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("Missing argument(s)")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		getTimeout(args[0])
	},
}

func getTimeout(targetVserver string) {
	vserverGlob := glob.MustCompile(targetVserver)
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

		r, _, err := client.GetVirtualServer(vserver)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Fprint(w, vserver, ":\t", *r.Connection.Timeout, "s\n")
	}

	w.Flush()
}

func init() {
	vserverCmd.AddCommand(getTimeoutCmd)
}
