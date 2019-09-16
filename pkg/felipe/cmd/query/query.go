package query

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/nihei9/felipe/graph"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var (
	flagSrcDir string
	flagLabel  string
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query",
		Short: "query generate a set of components specified by a query.",
		Long:  "query generate a set of components specified by a query.",
		RunE:  run,
	}
	cmd.Flags().StringVarP(&flagSrcDir, "src_dir", "s", "./", "directory to read definitions of components from")
	cmd.Flags().StringVarP(&flagLabel, "label", "l", "", "query label")

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	defFiles, err := listDefinitionFiles(flagSrcDir)
	if err != nil {
		return err
	}

	cs := graph.NewComponents()
	for _, defFile := range defFiles {
		def, err := readDefinition(defFile)
		if err != nil {
			return err
		}
		err = def.validateAndComplement()
		if err != nil {
			continue
		}
		if def.Kind != definitionKindComponents {
			continue
		}

		for _, cDef := range def.Components {
			c := graph.NewComponent(cDef.Name, cDef.Base, !cDef.Hide)
			for k, vs := range cDef.Labels {
				for _, v := range vs {
					c.Label(k, v)
				}
			}
			for _, dDef := range cDef.Dependencies {
				c.DependOn(graph.NewComponentID(dDef.Name))
			}
			cs.Add(c)
		}
	}
	err = cs.Complement()
	if err != nil {
		return err
	}

	var result *graph.Components
	if flagLabel != "" {
		l := strings.Split(flagLabel, "=")
		if len(l) != 2 {
			return fmt.Errorf("query label is malformed; got: %v", flagLabel)
		}
		condK := strings.TrimSpace(l[0])
		condV := strings.TrimSpace(l[1])

		cond := graph.NewCondition()
		cond.AddMatcher(graph.NewLabelsMatcher(map[string]string{condK: condV}))
		result, err = graph.Query(cs, cond, nil)
		if err != nil {
			return err
		}
	} else {
		cond := graph.NewCondition()
		cond.AddMatcher(graph.NewAnyMatcher())
		result, err = graph.Query(cs, cond, nil)
		if err != nil {
			return err
		}
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

func readDefinition(filePath string) (*definition, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	src, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	def := &definition{}
	err = yaml.Unmarshal(src, def)
	if err != nil {
		return nil, err
	}

	return def, nil
}

func writeResult(cs *graph.Components) error {
	components := []*component{}
	for _, c := range cs.Components() {
		deps := []*dependentComponent{}
		for _, d := range c.Dependencies() {
			deps = append(deps, &dependentComponent{
				Name: d.String(),
			})
		}

		components = append(components, &component{
			Name:         c.ID().String(),
			Hide:         false,
			RawLabels:    c.Labels(),
			Dependencies: deps,
		})
	}

	def := &definition{
		Version:    "1",
		Kind:       definitionKindComponents,
		Components: components,
	}

	data, err := yaml.Marshal(def)
	if err != nil {
		return err
	}

	fmt.Printf("%s", data)

	return nil
}
