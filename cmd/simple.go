package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
)

var answers = struct {
	Name  string
	Color string
}{}

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
