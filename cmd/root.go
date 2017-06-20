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
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"github.com/martinlindner/go-vtm"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
)

var RootCmd = &cobra.Command{
	Use:   "go-vtm-cli",
	Short: "Brocade vTM command line tool",
	Long:  `go-vtm-cli is a cli tool to control Brocade Virtual Traffic Manager via REST API.`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.go-vtm-cli.yaml)")

	RootCmd.PersistentFlags().String("vtmAPIUrl", "http://localhost:9070/", "vTM API URL.")
	RootCmd.PersistentFlags().String("vtmAPIUser", "admin", "vTM API user.")
	RootCmd.PersistentFlags().String("vtmAPIPass", "default", "vTM API password.")

	viper.BindPFlag("vtmAPIUrl", RootCmd.PersistentFlags().Lookup("vtmAPIUrl"))
	viper.BindPFlag("vtmAPIUser", RootCmd.PersistentFlags().Lookup("vtmAPIUser"))
	viper.BindPFlag("vtmAPIPass", RootCmd.PersistentFlags().Lookup("vtmAPIPass"))
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".go-vtm-cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".go-vtm-cli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func initClient() stingray.Client {
	vtmAPIUrl := viper.Get("vtmAPIUrl").(string)
	vtmAPIUser := viper.Get("vtmAPIUser").(string)
	vtmAPIPass := viper.Get("vtmAPIPass").(string)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpclient := &http.Client{Transport: tr}

	return *stingray.NewClient(httpclient, vtmAPIUrl, vtmAPIUser, vtmAPIPass)
}
