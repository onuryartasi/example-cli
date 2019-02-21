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
	"regexp"

	"github.com/ghodss/yaml"
	config "github.com/onuryartasi/example-cli/types"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
)

var conf config.KubeConfig

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
				Validate: survey.Required,
			}}
		var selectedNames string
		err := survey.Ask(qs, &selectedNames)
		fmt.Println(selectedNames)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Choose Context: %s", selectedNames)
		ReplaceContext(selectedNames)

	},
	Aliases: []string{"ctx", "cx"},
}

func GetContext() {
	home := os.Getenv("HOME")
	dat, err := ioutil.ReadFile(fmt.Sprintf("%s/.kube/config", home))
	if err != nil {
		fmt.Errorf("%s", err)
	}
	err = yaml.Unmarshal(dat, &conf)
}

func ReplaceContext(name string) {
	home := os.Getenv("HOME")
	file := fmt.Sprintf("%s/.kube/config", home)
	dat, err := ioutil.ReadFile(file)
	if err != nil {
		log.Println(err)
	}
	var re = regexp.MustCompile(`current-context.*`)
	s := re.ReplaceAllString(string(dat), fmt.Sprintf("current-context: %s", name))
	if err = ioutil.WriteFile(file, []byte(s), 0666); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
