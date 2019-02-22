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
	"bytes"
	"fmt"
	"log"

	types "github.com/onuryartasi/example-cli/types"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
	"gopkg.in/AlecAivazis/survey.v1"
	"gopkg.in/yaml.v2"
)

var data types.KubeConfig
var Host, User, Password string

// sshCmd represents the ssh command
var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "Get kubernetes config from remote server with ssh",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		selectedContext := []string{}
		answer := struct {
			Host     string
			User     string
			Password string
		}{}
		qs := []*survey.Question{
			{
				Name: "Host",
				Prompt: &survey.Input{
					Message: "Server: ",
				},
				Validate: survey.Required,
			},
			{
				Name: "User",
				Prompt: &survey.Input{
					Message: "Username: ",
				},
				Validate: survey.Required,
			},
			{
				Name: "Password",
				Prompt: &survey.Password{
					Message: "Password: ",
				},
				Validate: survey.Required,
			},
		}
		err := survey.Ask(qs, &answer)
		if err != nil {
			log.Fatal(err)
		}
		config := ConnectServer(answer.Host, answer.User, answer.Password)
		err = yaml.Unmarshal([]byte(config), &data)
		if err != nil {
			log.Fatalf("%s", err)
		}
		var names = []string{}
		for _, value := range data.Contexts {
			names = append(names, value.Name)
		}

		qs2 := []*survey.Question{
			{Name: "context",
				Prompt: &survey.MultiSelect{
					Message: "Choose a context:",
					Options: names,
				},
				Validate: survey.Required,
			}}
		err = survey.Ask(qs2, &selectedContext)
		fmt.Println(selectedContext)

	},
}

func init() {
	rootCmd.AddCommand(sshCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sshCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sshCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func ConnectServer(host, username, password string) string {
	hostKey := ssh.InsecureIgnoreHostKey()
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
	}
	config.HostKeyCallback = hostKey
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", host), config)
	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}

	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}
	defer session.Close()

	// Once a Session is created, you can execute a single command on
	// the remote side using the Run method.
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(fmt.Sprintf("echo \"%s\" | sudo -S cat /etc/kubernetes/admin.conf", password)); err != nil {
		log.Fatalf("Failed to run: %s", err)
	}

	return b.String()
}
