package buildrun

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	buildv1alpha1 "github.com/shipwright-io/build/pkg/apis/build/v1alpha1"

	"github.com/shipwright-io/cli/pkg/shp/cmd/runner"
	"github.com/shipwright-io/cli/pkg/shp/params"
	"github.com/shipwright-io/cli/pkg/shp/resource"
)

// ListCommand contains data input from user for list sub-command
type ListCommand struct {
	cmd *cobra.Command

	noHeader bool
}

func listCmd() runner.SubCommand {
	listCmd := &ListCommand{
		cmd: &cobra.Command{
			Use:   "list [flags]",
			Short: "List Builds",
		},
	}

	listCmd.cmd.Flags().BoolVar(&listCmd.noHeader, "no-header", false, "Do not show columns header in list output")

	return listCmd
}

// Cmd returns cobra command object
func (c *ListCommand) Cmd() *cobra.Command {
	return c.cmd
}

// Complete fills in data provided by user
func (c *ListCommand) Complete(params *params.Params, args []string) error {
	return nil
}

// Validate validates data input by user
func (c *ListCommand) Validate() error {
	return nil
}

// Run executes list sub-command logic
func (c *ListCommand) Run(params *params.Params, io *genericclioptions.IOStreams) error {
	// TODO: Support multiple output formats here, not only tabwriter
	//       find out more in kubectl libraries and use them

	writer := tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0)
	columnNames := "NAME\tSTATUS"
	columnTemplate := "%s\t%s\n"

	brr := resource.GetBuildRunResource(params)

	var brs buildv1alpha1.BuildRunList
	if err := brr.List(c.cmd.Context(), &brs); err != nil {
		return err
	}

	if !c.noHeader {
		fmt.Fprintln(writer, columnNames)
	}

	for _, br := range brs.Items {
		name := br.Name
		status := string(metav1.ConditionUnknown)
		for _, condition := range br.Status.Conditions {
			if condition.Type == buildv1alpha1.Succeeded {
				status = condition.Reason
				break
			}
		}

		fmt.Fprintf(writer, columnTemplate, name, status)
	}

	writer.Flush()

	return nil
}
