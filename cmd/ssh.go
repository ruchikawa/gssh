/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

// ServerInfo means Server information
// TODO: in the future, need re-factoring this struture to use map.
type ServerInfo struct {
	HOST string
	PORT string
	USER string
	// Passは流石に標準入力にしましょう
	PASS string
}

// sshCmd represents the ssh command
var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		si := ServerInfo{}
		si.USER, _ = cmd.Flags().GetString("USER")
		si.HOST, _ = cmd.Flags().GetString("HOST")
		si.PORT, _ = cmd.Flags().GetString("PORT")
		si.PASS, _ = cmd.Flags().GetString("PASS")
		sshToServer(&si)
	},
}

/*
Run: func(cmd *cobra.Command, args []string) {
	si := ServerInfo{}

	si.USER, _ = sshCmd.Flags().GetString("USER")
	si.HOST, _ = sshCmd.Flags().GetString("HOST")
	si.PORT, _ = sshCmd.Flags().GetString("PORT")
	si.PASS, _ = sshCmd.Flags().GetString("PASS")
	fmt.Println("%s", si.USER)

},*/
func init() {
	rootCmd.AddCommand(sshCmd)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:

	sshCmd.PersistentFlags().String("USER", "u", "A help for foo")
	sshCmd.PersistentFlags().String("HOST", "h", "A help for foo")
	sshCmd.PersistentFlags().String("PORT", "p", "A help for foo")
	sshCmd.PersistentFlags().String("PASS", "pass", "A help for foo")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	sshCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// ここでcobra.commandのポインタを参照するのが正しいのでは？
func sshToServer(si *ServerInfo) {

	// Create sshClientConfig
	sshConfig := &ssh.ClientConfig{
		User: si.USER,
		Auth: []ssh.AuthMethod{
			ssh.Password(si.PASS),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// SSH connect.
	client, err := ssh.Dial("tcp", si.HOST+":"+si.PORT, sshConfig)

	// Create Session
	session, err := client.NewSession()
	defer session.Close()

	// キー入力を接続先が認識できる形式に変換する(ここがキモ)
	fd := int(os.Stdin.Fd())
	state, err := terminal.MakeRaw(fd)
	if err != nil {
		fmt.Println(err)
	}
	defer terminal.Restore(fd, state)

	// ターミナルサイズの取得
	w, h, err := terminal.GetSize(fd)
	if err != nil {
		fmt.Println(err)
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	err = session.RequestPty("xterm", h, w, modes)
	if err != nil {
		fmt.Println(err)
	}

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	err = session.Shell()
	if err != nil {
		fmt.Println(err)
	}

	err = session.Wait()
	if err != nil {
		fmt.Println(err)
	}
}
