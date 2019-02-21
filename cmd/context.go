// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/ghodss/yaml"
	config "github.com/onuryartasi/example-cli/types"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
)

var conf config.KubeConfig

// contextCmd represents the context command
var contextCmd = &cobra.Command{
	Use:   "context",
	Short: "Kubernetes context",
	Run: func(cmd *cobra.Command, args []string) {
		GetContext()
		var names = []string{}
		for _, value := range conf.Contexts {
			names = append(names, value.Name)
		}
		qs := []*survey.Question{
			{Name: "context",
				Prompt: &survey.Select{
					Message: "Choose a context:",
					Options: names,
				},
				Transform: survey.Title,
				Validate:  survey.Required,
			}}
		var selectedNames string
		err := survey.Ask(qs, &selectedNames)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Choose Context: %s", selectedNames)

	},
}

func init() {
	rootCmd.AddCommand(contextCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// contextCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// contextCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func GetContext() {
	user := os.Getenv("USER")
	dat, err := ioutil.ReadFile(fmt.Sprintf("/home/%s/.kube/config", user))
	if err != nil {
		fmt.Errorf("%s", err)
	}
	err = yaml.Unmarshal(dat, &conf)
}
