/*
Package cmd
Copyright © 2022 Signal Advisors <devteam@signaladvisors.com>
*/
package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	sigex "github.com/signaldevs/sigex/pkg"
	"github.com/spf13/cobra"
)

var (
	envFiles        []string
	envVars         map[string]string
	skipSecretsFlag bool
	debugFlag       bool
	versionFlag     bool
	version         string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sigex [flags] command",
	Short: "The sigex process runner",
	Long: `sigex is a process runner/executor with support for multiple .env file
configuration as well as automatic retrieval of secrets from 
supported secrets manager platforms.`,
	RunE: RootCmdRunE,
}

func RootCmdRunE(_ *cobra.Command, args []string) error {

	if versionFlag {
		fmt.Println(version)
		return nil
	}

	osHelper := sigex.GetOSHelper()

	if len(args) < 1 && !debugFlag {
		return fmt.Errorf("no command argument was provided")
	}

	// check to make sure the command binary actually exists

	var binary string
	var lpError error

	if !debugFlag {
		binary, lpError = osHelper.LookPath(args[0])
		if lpError != nil {
			return fmt.Errorf("invalid command: %e", lpError)
		}
	}

	// construct the complete environment to be
	// passed to the command process
	env := processEnv()

	// if in debug mode, stop here and just log
	// out the environment
	if debugFlag {
		logEnv(env)
		return nil
	}

	// execute the command with the processed environment
	// as separate lines
	execErr := osHelper.Exec(binary, args, env)
	if execErr != nil {
		return fmt.Errorf("unable to execute command: %e", execErr)
	}

	return nil
}

func RootCmdFlags(cmd *cobra.Command) {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.sigex.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	cmd.Flags().StringSliceVarP(&envFiles, "env-file", "f", []string{}, "specify one or more .env files to use")
	cmd.Flags().StringToStringVarP(&envVars, "env-var", "e", make(map[string]string), "specify one or more environment variables to use (ex: -e FOO=bar)")
	cmd.Flags().BoolVar(&skipSecretsFlag, "skip-secrets", false, "skip the automatic resolution of secret values")
	cmd.Flags().BoolVar(&debugFlag, "debug", false, "debug the resolved environment variables")
	cmd.Flags().BoolVar(&versionFlag, "version", false, "print the version and exit")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	RootCmdFlags(rootCmd)
	cobra.CheckErr(rootCmd.Execute())
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

	// Resolve tokenized secrets
	if !skipSecretsFlag {
		for key, element := range envMap {
			envMap[key] = sigex.ResolveSecret(element)
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

	envVarRegex := regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*=.+$`)

	for i := 0; i < len(lines); i++ {

		line := strings.TrimSpace(lines[i])

		// Skip if line is blank or is a comment
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Ensure the line follows the standard env file format: "VARIABLE=value"
		if !envVarRegex.MatchString(line) {
			continue
		}

		s := strings.Split(line, "=")

		if len(s) < 2 {
			continue
		}

		key := strings.TrimSpace(s[0])
		val := strings.TrimSpace(s[1])

		envMap[key] = val
	}
}

func getFileLines(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println(fmt.Errorf("error closing file: %v", err))
		}
	}(file)

	sc := bufio.NewScanner(file)
	lines := make([]string, 0)

	// Read through 'tokens' until an EOF is encountered.
	for sc.Scan() {
		// TODO: need to filter out comment lines, blank lines, etc
		lines = append(lines, sc.Text())
	}

	if err := sc.Err(); err != nil {
		log.Fatal(err)
	}

	return lines
}

func logEnv(env []string) {
	for i := 0; i < len(env); i++ {
		fmt.Println(env[i])
	}
}

func init() {

}
