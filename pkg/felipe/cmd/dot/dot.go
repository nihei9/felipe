package dot

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/awalterschulze/gographviz"
	"github.com/nihei9/felipe/graph"
	"github.com/nihei9/felipe/pkg/felipe/definitions"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var (
	flagSrcFile  string
	flagFaceFile string
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dot",
		Short: "dot generate .dot files.",
		Long:  "dot generate .dot files.",
		RunE:  run,
	}
	cmd.Flags().StringVarP(&flagSrcFile, "src_file", "s", "", "file path that defines components (default: stdin)")
	cmd.Flags().StringVarP(&flagFaceFile, "face", "f", "", "file path that defines faces for image generates from DOT")

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	def, err := readComponentsDefinition(flagSrcFile)
	if err != nil {
		return err
	}

	cs := graph.NewComponents()
	{
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
		err = cs.Complement()
		if err != nil {
			return err
		}
	}

	fs := []*graph.Face{}
	if flagFaceFile != "" {
		def, err := readFacesDefinition(flagFaceFile)
		if err != nil {
			return err
		}
		err = def.Validate()
		if err != nil {
			return err
		}

		if def.Kind != definitions.DefinitionKindFaces {
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

	err = writeDot(cs, cs, fs, os.Stdout)
	if err != nil {
		return err
	}

	return nil
}

func readComponentsDefinition(filePath string) (*definitions.ComponentsDefinition, error) {
	var r io.Reader
	if filePath != "" {
		f, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		r = f
	} else {
		r = os.Stdin
	}

	return definitions.ReadComponentsDefinition(r)
}

func readFacesDefinition(filePath string) (*definitions.FacesDefinition, error) {
	def := &definitions.FacesDefinition{}
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	src, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(src, def)
	if err != nil {
		return nil, err
	}

	return def, nil
}

func writeDot(group *graph.Components, cs *graph.Components, fs []*graph.Face, w io.Writer) error {
	dot, err := genDot(group, cs, fs)
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(w, dot)
	if err != nil {
		return err
	}

	return nil
}

func genDot(group *graph.Components, cs *graph.Components, fs []*graph.Face) (string, error) {
	ast, _ := gographviz.ParseString("digraph G {}")
	g := gographviz.NewGraph()
	err := gographviz.Analyse(ast, g)
	if err != nil {
		return "", err
	}
	g.AddAttr("G", "rankdir", "LR")
	g.AddAttr("G", "fontsize", "11.0")

	for _, c := range group.Components() {
		nAttrs, err := genAttributes(c, fs)
		nAttrs["penwidth"] = "0.75"
		if err != nil {
			return "", err
		}
		err = g.AddNode("G", fmt.Sprintf("\"%s\"", c.ID().String()), nAttrs)
		if err != nil {
			return "", err
		}
		for _, dcid := range c.Dependencies() {
			d, _ := cs.Get(dcid)
			nAttrs, err := genAttributes(d, fs)
			nAttrs["penwidth"] = "0.75"
			if err != nil {
				return "", err
			}
			err = g.AddNode("G", fmt.Sprintf("\"%s\"", d.ID().String()), nAttrs)
			if err != nil {
				return "", err
			}

			eAttrs := map[string]string{
				"arrowsize": "0.75",
				"penwidth":  "0.75",
			}
			err = g.AddEdge(fmt.Sprintf("\"%s\"", c.ID().String()), fmt.Sprintf("\"%s\"", d.ID().String()), true, eAttrs)
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
		ok, err := graph.Match(c, f.Condition)
		if err != nil {
			return nil, err
		}
		if !ok {
			continue
		}
		for k, v := range f.Attributes {
			if k == "label" {
				placeholders := []string{}
				capture := false
				var start int
				var end int
				for i, char := range v {
					switch char {
					case '{':
						if capture {
							return nil, fmt.Errorf("an embeded label cannot be nested")
						}
						capture = true
						start = i
					case '}':
						if !capture {
							return nil, fmt.Errorf("an embeded label is malformed")
						}
						capture = false
						end = i

						placeholder := v[start : end+1]
						placeholders = append(placeholders, placeholder)
					}
				}

				embeddedValues := []string{}
				for _, p := range placeholders {
					labelName := strings.TrimSpace(p[1 : len(p)-1])
					vs, ok := c.Labels()[labelName]
					if !ok {
						return nil, fmt.Errorf("undefined label `%s` cannot use in `name` directive", labelName)
					}
					if len(vs) != 1 {
						return nil, fmt.Errorf("a label used as the embeded label must have just one value; `%s` has %v values", labelName, len(vs))
					}
					embeddedValues = append(embeddedValues, vs[0])
				}

				label := v
				for i, p := range placeholders {
					label = strings.Replace(label, p, embeddedValues[i], 1)
				}
				attrs[k] = label
			} else {
				attrs[k] = v
			}
		}
	}

	return attrs, nil
}
