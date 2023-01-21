package mesh

import (
	"fmt"
	"strings"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	removeCmd = &cobra.Command{
		Use:   "remove",
		Short: "remove a service mesh in the kubernetes cluster",
		Long:  `remove service mesh in the connected kubernetes cluster`,
		Example: `
// Remove a service mesh
mesheryctl mesh remove [mesh adapter name]
		`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			utils.Log.Info("Verifying prerequisites...")
			mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
			if err != nil {
				return errors.Wrap(err, "error processing config")
			}

			if len(args) < 1 {
				meshName, err = validateMesh(mctlCfg, "")
			} else {
				meshName, err = validateMesh(mctlCfg, args[0])
			}
			if err != nil {
				return errors.Wrap(err, "error validating request")
			}

			if err = validateAdapter(mctlCfg, meshName); err != nil {
				return errors.Wrap(err, "adapter not valid")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			if len(args) > 1 {
				return errors.New(utils.MeshError("'mesheryctl mesh remove' should not have more than one argument, it can remove only one adapter at a time.\n"))
			}

			s := utils.CreateDefaultSpinner(fmt.Sprintf("Removing %s", meshName), fmt.Sprintf("\n%s service mesh removed successfully", meshName))
			mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
			if err != nil {
				return errors.Wrap(err, "error processing config")
			}

			s.Start()
			_, err = sendOperationRequest(mctlCfg, strings.ToLower(meshName), true)
			if err != nil {
				return errors.Wrap(err, "error installing service mesh")
			}
			s.Stop()

			return nil
		},
	}
)

func init() {
	removeCmd.Flags().StringVarP(
		&namespace, "namespace", "n", "default",
		"Kubernetes namespace where the mesh is deployed",
	)
}
