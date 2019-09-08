package dot

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/awalterschulze/gographviz"
	"github.com/nihei9/felipe/graph"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var (
	flagSrcDir     string
	flagOutputFile string
	flagLabel      string
	flagFaceFile   string
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dot",
		Short: "dot generate .dot files.",
		Long:  "dot generate .dot files.",
		RunE:  run,
	}
	cmd.Flags().StringVarP(&flagSrcDir, "src_dir", "s", "./", "directory to read definitions of components from")
	cmd.Flags().StringVarP(&flagOutputFile, "output_file", "o", "", "file path to write DOT to (default: stdout)")
	cmd.Flags().StringVarP(&flagLabel, "label", "l", "", "query label")
	cmd.Flags().StringVarP(&flagFaceFile, "face", "f", "", "file path that defines faces for image generates from DOT")

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
		err = def.validate()
		if err != nil {
			return err
		}

		if def.Kind != definitionKindComponents {
			continue
		}

		for _, cDef := range def.Components {
			c := graph.NewComponent(cDef.Name)
			for k, v := range cDef.Labels {
				c.Label(k, v)
			}
			for _, dDef := range cDef.Dependencies {
				d := graph.NewComponent(dDef.Name)
				c.DependOn(d)
			}
			cs.AddComponent(c)
		}
	}

	fs := []*graph.Face{}
	if flagFaceFile != "" {
		def, err := readDefinition(flagFaceFile)
		if err != nil {
			return err
		}
		err = def.validate()
		if err != nil {
			return err
		}

		if def.Kind != definitionKindFaces {
			return fmt.Errorf("kind of specified face file is not `faces`; got: %v", def.Kind)
		}

		for _, fDef := range def.Faces {
			f := graph.NewFace()
			m := graph.NewLabelsMatcher(fDef.Targets.MatchLabels)
			f.AddTarget(m)
			f.AddAttributes(fDef.Attributes)
			fs = append(fs, f)
		}
	}

	group, err := cs.Query(flagLabel)
	if err != nil {
		return err
	}
	writeDot(group, fs, flagOutputFile)

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

func writeDot(group []*graph.Component, fs []*graph.Face, filePath string) error {
	dot, err := genDot(group, fs)
	if err != nil {
		return err
	}

	if flagOutputFile != "" {
		f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			return err
		}
		defer f.Close()

		f.Write([]byte(dot))
	} else {
		fmt.Printf(dot)
	}

	return nil
}

func genDot(group []*graph.Component, fs []*graph.Face) (string, error) {
	ast, _ := gographviz.ParseString("digraph G {}")
	g := gographviz.NewGraph()
	err := gographviz.Analyse(ast, g)
	if err != nil {
		return "", err
	}
	g.AddAttr("G", "rankdir", "LR")
	g.AddAttr("G", "fontsize", "11.0")

	for _, c := range group {
		attrs, err := genAttributes(c, fs)
		if err != nil {
			return "", err
		}
		err = g.AddNode("G", fmt.Sprintf("\"%s\"", c.ID.String()), attrs)
		if err != nil {
			return "", err
		}
		for _, d := range c.Dependencies {
			attrs, err := genAttributes(d, fs)
			if err != nil {
				return "", err
			}
			err = g.AddNode("G", fmt.Sprintf("\"%s\"", d.ID.String()), attrs)
			if err != nil {
				return "", err
			}

			err = g.AddEdge(fmt.Sprintf("\"%s\"", c.ID.String()), fmt.Sprintf("\"%s\"", d.ID.String()), true, nil)
			if err != nil {
				return "", err
			}
		}
	}

	return g.String(), nil
}

func genAttributes(c *graph.Component, fs []*graph.Face) (map[string]string, error) {
	attrs := map[string]string{}
	for _, f := range fs {
		match := false
		for _, m := range f.Matchers {
			ok, err := m.Match(c)
			if err != nil {
				return nil, err
			}
			if !ok {
				continue
			}
			match = true
			break
		}
		if !match {
			continue
		}
		for k, v := range f.Attributes {
			attrs[k] = v
		}
	}

	return attrs, nil
}
