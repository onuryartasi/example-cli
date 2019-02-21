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

var answers = struct {
	Name  string
	Color string
}{}
var conf config.KubeConfig

// the questions to ask
var simpleQs = []*survey.Question{
	{
		Name: "name",
		Prompt: &survey.Input{
			Message: "What is your name?",
		},
		Validate:  survey.Required,
		Transform: survey.Title,
	},
	{
		Name: "color",
		Prompt: &survey.Select{
			Message: "Choose a color:",
			Options: []string{"red", "blue", "green"},
		},
		Validate: survey.Required,
	},
}

func init() {
	rootCmd.AddCommand(simpleCmd)
	rootCmd.AddCommand(contextCmd)
}

var simpleCmd = &cobra.Command{
	Use:   "simple",
	Short: "Cobra test survey",
	Long:  `Cobra cli testing survey`,
	Run: func(cmd *cobra.Command, args []string) {
		// ask the question
		err := survey.Ask(simpleQs, &answers)

		if err != nil {
			fmt.Println(err.Error())
			return
		}
		// print the answers
		fmt.Printf("%s chose %s.\n", answers.Name, answers.Color)
	},
}

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

func GetContext() {
	user := os.Getenv("USER")
	dat, err := ioutil.ReadFile(fmt.Sprintf("/home/%s/.kube/config", user))
	if err != nil {
		fmt.Errorf("%s", err)
	}
	err = yaml.Unmarshal(dat, &conf)
}
