package apps

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/apps"
	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/types"
)

func MakeInstallOpenebs() *cobra.Command {
	var openebs = &cobra.Command{
		Use:   "openebs",
		Short: "Install openebs",
		Long:  `Install openebs`,
		Example: `arkade install openebs`,
		SilenceUsage: true,
	}

	openebs.Flags().StringP("namespace", "n", "openebs", "The namespace to install the chart")
	openebs.Flags().Bool("update-repo", true, "Update the helm repo")
	openebs.Flags().Bool("ndm", false, "Enable Node Disk Manager")
	openebs.Flags().StringArray("set", []string{}, 
		"Use custom flags or override existing flags \n(example --set serviceAccount.create=true)")

	openebs.PreRunE = func(cmd *cobra.Command, args []string) error {
		_, err := cmd.Flags().GetString("namespace")
		if err != nil {
			return err
		}
		_, err = cmd.Flags().GetBool("ndm")
		if err != nil {
			return err
		}
		_, err = cmd.Flags().GetStringArray("set")
		if err != nil {
			return err
		}

		_, err = cmd.Flags().GetBool("update-repo")
		if err != nil {
			return err
		}
		return nil
	}

	openebs.RunE = func(cmd *cobra.Command, args []string) error {
		kubeConfigPath, _ := cmd.Flags().GetString("kubeconfig")
		if err := config.SetKubeconfig(kubeConfigPath); err != nil {
			return err
		}

		namespace, _ := cmd.Flags().GetString("namespace")
		updateRepo, _ := cmd.Flags().GetBool("update-repo")
		customFlags, _ := cmd.Flags().GetStringArray("set")
		ndmEnabled, _ := cmd.Flags().GetBool("ndm")
		overrides := map[string]string{
			"ndm.enabled":              "false",
		}
		if ndmEnabled {
			overrides["ndm.enabled"] = "true"
		}
		if err := config.MergeFlags(overrides, customFlags); err != nil {
			return err
		}
		openebsOptions := types.DefaultInstallOptions().
			WithNamespace(namespace).
			WithHelmRepo("openebs/openebs").
			WithHelmURL("https://openebs.github.io/charts").
			WithOverrides(overrides).
			WithHelmUpdateRepo(updateRepo).
			WithInstallNamespace(true).
			WithKubeconfigPath(kubeConfigPath)

		_, err := apps.MakeInstallChart(openebsOptions)
		if err != nil {
			return err
		}

		fmt.Println(openebsInstallMsg)

		return nil
	}

	return openebs
}

const OpenebsInfoMsg = `Please wait for several minutes for openebs deployment to complete.
To test the usage of openEBS, simply deploy postgresql app with a storage class.
In this example we deploy postgresql using hostpath storage class.
arkade install postgresql --set primary.persistance.storageClass="openebs-hostpath"

Once postgresql deployed, verify the dynamic pv attached.
kubectl get pv,pvc (assuming the app deployed in teh default namespace)
`

const openebsInstallMsg = `=======================================================================
=                     openebs has been installed.                       =
=======================================================================` +
	"\n\n" + OpenebsInfoMsg + "\n\n" + pkg.SupportMessageShort