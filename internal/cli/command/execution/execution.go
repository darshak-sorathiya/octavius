package execution

import (
	"fmt"
	"octavius/internal/cli/client"
	"octavius/internal/cli/daemon"
	"octavius/internal/pkg/log"
	"strings"

	"github.com/spf13/cobra"
)

// NewCmd create a command for execution
func NewCmd(octaviusDaemon daemon.Client) *cobra.Command {
	return &cobra.Command{
		Use:     "execute",
		Short:   "Execute the existing job",
		Long:    "This command helps to execute the job which is already created in server",
		Example: fmt.Sprintf("octavius execute <job-name> arg1=argvalue1 arg2=argvalue2"),
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			jobName := args[0]
			jobData := map[string]string{}

			for i := 1; i < len(args); i++ {
				arg := strings.Split(args[i], "=")
				jobData[arg[0]] = arg[1]
			}
			client := &client.GrpcClient{}
			response, err := octaviusDaemon.ExecuteJob(jobName, jobData, client)
			if err != nil {
				log.Error(err, "error in executing job")
				return
			}
			log.Info(response.Status)
		},
	}
}
