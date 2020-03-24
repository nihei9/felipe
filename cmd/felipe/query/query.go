package query

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/nihei9/felipe/component"
	"github.com/nihei9/felipe/definitions"
	"github.com/nihei9/felipe/query"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var (
	flagFilter          string
	flagComplementation string
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query",
		Short: "query generate a set of components specified by a query.",
		Long:  "query generate a set of components specified by a query.",
		Args:  cobra.ExactArgs(1),
		RunE:  run,
	}
	cmd.Flags().StringVarP(&flagFilter, "filter", "f", "", "filter used in the query")
	cmd.Flags().StringVarP(&flagComplementation, "complementation", "c", "", "complementation used in the query")

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	defFiles, err := listDefinitionFiles(args[0])
	if err != nil {
		return err
	}

	cs := component.NewComponents()
	for _, defFile := range defFiles {
		def, err := readComponentsDefinition(defFile)
		if err != nil {
			return err
		}

		for _, cDef := range def.Components {
			c := definitions.MakeComponentEntity(cDef)
			cs.Add(c)
		}
	}
	err = cs.Complement()
	if err != nil {
		return err
	}

	var filter query.Filter
	if flagFilter != "" {
		f := strings.Split(flagFilter, "=")
		if len(f) != 2 {
			return fmt.Errorf("filter is malformed; got: %v", flagFilter)
		}
		k := strings.TrimSpace(f[0])
		v := strings.TrimSpace(f[1])

		filter = query.LabelsFilter{
			Labels: map[string]string{k: v},
		}
	} else {
		filter = query.AllPassFilter{}
	}

	var complementer query.Complementer
	if flagComplementation != "" {
		f := strings.Split(flagComplementation, "=")
		if len(f) != 2 {
			return fmt.Errorf("complementation is malformed; got: %v", flagComplementation)
		}
		k := strings.TrimSpace(f[0])
		v := strings.TrimSpace(f[1])

		switch k {
		case "dep":
			depth, err := strconv.Atoi(v)
			if err != nil {
				return err
			}
			complementer = query.DependenciesComplementer{
				AllComponents: cs,
				Depth:         depth,
			}
		case "rdep":
			depth, err := strconv.Atoi(v)
			if err != nil {
				return err
			}
			complementer = query.ReverseDependenciesComplementer{
				AllComponents: cs,
				Depth:         depth,
			}
		default:
			return fmt.Errorf("invalid complementation; got: %v", k)
		}
	} else {
		complementer = query.DependenciesComplementer{
			AllComponents: cs,
			Depth:         -1,
		}
	}

	result, err := query.Query{
		Components:   cs,
		Filter:       filter,
		Complementer: complementer,
	}.Do()
	if err != nil {
		return err
	}

	err = writeResult(result)
	if err != nil {
		return err
	}

	return nil
}

func listDefinitionFiles(srcDir string) ([]string, error) {
	return filepath.Glob(filepath.Join(srcDir, "*.yaml"))
}

func readComponentsDefinition(filePath string) (*definitions.ComponentsDefinition, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return definitions.ReadComponentsDefinition(f)
}

func writeResult(cs *component.Components) error {
	components := []*definitions.Component{}
	for _, id := range cs.GetIDs() {
		c, _ := cs.Get(id)
		deps := []*definitions.DependentComponent{}
		for d, _ := range c.Dependencies {
			deps = append(deps, &definitions.DependentComponent{
				ID: d.String(),
			})
		}

		components = append(components, &definitions.Component{
			ID:           c.ID.String(),
			Hide:         false,
			Labels:       c.Labels,
			Dependencies: deps,
		})
	}

	def := &definitions.ComponentsDefinition{
		Version:    "1",
		Kind:       definitions.DefinitionKindComponents,
		Components: components,
	}

	data, err := yaml.Marshal(def)
	if err != nil {
		return err
	}

	fmt.Printf("%s", data)

	return nil
}
