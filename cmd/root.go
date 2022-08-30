/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
)

var envFiles []string
var envVars map[string]string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sigex [flags] command",
	Short: "The sigex process runner",
	Long: `sigex runs processes and leverages multiple env files for 
ultimate flexibility.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Fatal("Not enough arguments provided")
		}

		binary, err := exec.LookPath(args[0])
		if err != nil {
			panic(err)
		}

		env := processEnv()

		execErr := syscall.Exec(binary, args, env)
		if execErr != nil {
			log.Fatal(execErr)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// Resolve all env vars using existing environment plus any additional env files
// in use
func processEnv() []string {
	envMap := make(map[string]string)

	// Add existing environment to map
	envLinesToMap(envMap, os.Environ())

	// Load other env files provided by flags
	if len(envFiles) > 0 {
		for i := 0; i < len(envFiles); i++ {
			envLinesToMap(envMap, getFileLines(envFiles[i]))
		}
	}

	// Load env vars supplied on the command line
	if len(envVars) > 0 {
		for key, element := range envVars {
			envMap[key] = element
		}
	}

	// Convert the map to lines
	lines := make([]string, 0)

	for key, element := range envMap {
		lines = append(lines, strings.Join([]string{key, element}, "="))
	}

	return lines
}

func envLinesToMap(envMap map[string]string, lines []string) {

	for i := 0; i < len(lines); i++ {
		s := strings.Split(lines[i], "=")
		if len(s) < 2 {
			continue
		}
		envMap[strings.Trim(s[0], " ")] = strings.Trim(s[1], " ")
	}
}

func getFileLines(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	sc := bufio.NewScanner(file)
	lines := make([]string, 0)

	// Read through 'tokens' until an EOF is encountered.
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}

	if err := sc.Err(); err != nil {
		log.Fatal(err)
	}

	return lines
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.sigex.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringSliceVarP(&envFiles, "env-file", "f", []string{}, "Specify one or more .env files to use")
	rootCmd.Flags().StringToStringVarP(&envVars, "env-var", "e", make(map[string]string), "Specify one or more environment variables to use (ex: -e FOO=bar)")
}
