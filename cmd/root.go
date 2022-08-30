/*
Copyright © 2022 Signal Advisors <devteam@signaladvisors.com>
*/
package cmd

import (
	"bufio"
	"context"
	"fmt"
	"hash/crc32"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"syscall"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/spf13/cobra"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

var (
	envFiles    []string
	envVars     map[string]string
	skipSecrets bool
	secretRegex *regexp.Regexp
)

type osHelperInterface interface {
	LookPath(string) (string, error)
	Exec(string, []string, []string) error
}

type osHelper struct{}

func (o osHelper) LookPath(path string) (string, error) {
	return exec.LookPath(path)
}

func (o osHelper) Exec(argv0 string, argv []string, envv []string) error {
	return syscall.Exec(argv0, argv, envv)
}

var osHelperInstance osHelperInterface

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sigex [flags] command",
	Short: "The sigex process runner",
	Long: `sigex is a process runner/executor with support for multiple .env file
configuration as well as automatic retrieval of secrets from 
supported secrets manager platforms.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: RootCmdRunE,
}

func RootCmdRunE(cmd *cobra.Command, args []string) error {
	secretRegex, _ = regexp.Compile(`^sigex-secret-(.*)\:\/\/(.*)$`)

	if len(args) < 1 {
		return fmt.Errorf("no command argument was provided")
	}

	binary, err := osHelperInstance.LookPath(args[0])
	if err != nil {
		return err
	}

	env := processEnv()

	execErr := osHelperInstance.Exec(binary, args, env)
	if execErr != nil {
		return execErr
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
	cmd.Flags().BoolVar(&skipSecrets, "skip-secrets", false, "skip the automatic resolution of secret values")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	RootCmdFlags(rootCmd)
	cobra.CheckErr(rootCmd.Execute())
}

func ResetVars() {
	envFiles = nil
	envVars = nil
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
	if !skipSecrets {
		for key, element := range envMap {
			if isSecretToken(element) {
				envMap[key] = resolveSecret(element)
			}
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
		// TODO: need to filter out comment lines, blank lines, etc
		lines = append(lines, sc.Text())
	}

	if err := sc.Err(); err != nil {
		log.Fatal(err)
	}

	return lines
}

func isSecretToken(token string) bool {
	matches := secretRegex.MatchString(token)
	return matches
}

func resolveSecret(token string) string {
	parts := secretRegex.FindStringSubmatch(token)
	if len(parts) < 3 {
		log.Fatalln("secret token in incorrect format: ", token)
	}
	// implement secret resolution (currently just returning the parsed token)

	secretPlatform := parts[1]

	var secret string

	if secretPlatform == "gcp" {
		secret = getGCPSecretVersion(parts[2])
	} else {
		log.Fatalln("unsupported secret platform: " + secretPlatform)
	}

	return secret
}

func getGCPSecretVersion(name string) string {
	// name := "projects/my-project/secrets/my-secret/versions/5"
	// name := "projects/my-project/secrets/my-secret/versions/latest"

	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Fatalln(fmt.Errorf("failed to create secretmanager client: %v", err))
	}
	defer client.Close()

	// Build the request.
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	// Call the API.
	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		log.Fatalln(fmt.Errorf("failed to access secret version: %v", err))
	}

	// Verify the data checksum.
	crc32c := crc32.MakeTable(crc32.Castagnoli)
	checksum := int64(crc32.Checksum(result.Payload.Data, crc32c))
	if checksum != *result.Payload.DataCrc32C {
		log.Fatalln(fmt.Errorf("data corruption detected in secret version"))
	}

	// WARNING: Do not print the secret in a production environment

	return string(result.Payload.Data)
}

func SetOSHelper(helper osHelperInterface) {
	osHelperInstance = helper
}

func init() {
	SetOSHelper(osHelper{})
	// osHelperInstance = osHelper{}
}
