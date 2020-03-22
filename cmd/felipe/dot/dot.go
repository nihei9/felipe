package dot

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/awalterschulze/gographviz"
	"github.com/nihei9/felipe/component"
	"github.com/nihei9/felipe/definitions"
	"github.com/nihei9/felipe/query"
	"github.com/spf13/cobra"
)

type Face struct {
	Filter     query.Filter
	Attributes map[string]string
}

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

	cs := component.NewComponents()
	for _, cDef := range def.Components {
		c := definitions.MakeComponentEntity(cDef)
		cs.Add(c)
	}
	err = cs.Complement()
	if err != nil {
		return err
	}

	fs := []*Face{}
	if flagFaceFile != "" {
		def, err := readFacesDefinition(flagFaceFile)
		if err != nil {
			return err
		}

		for _, fDef := range def.Faces {
			f := &Face{
				Filter: query.LabelsFilter{
					Labels: fDef.Targets.MatchLabels,
				},
				Attributes: fDef.Attributes,
			}
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
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return definitions.ReadFacesDefinition(f)
}

func writeDot(group *component.Components, cs *component.Components, fs []*Face, w io.Writer) error {
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

func genDot(group *component.Components, cs *component.Components, fs []*Face) (string, error) {
	ast, _ := gographviz.ParseString("digraph G {}")
	g := gographviz.NewGraph()
	err := gographviz.Analyse(ast, g)
	if err != nil {
		return "", err
	}
	g.AddAttr("G", "rankdir", "LR")
	g.AddAttr("G", "fontsize", "11.0")

	for _, id := range group.GetIDs() {
		c, _ := group.Get(id)
		nAttrs, err := genAttributes(c, fs)
		if err != nil {
			return "", err
		}
		nAttrs["penwidth"] = "0.75"
		err = g.AddNode("G", fmt.Sprintf("\"%s\"", c.ID.String()), nAttrs)
		if err != nil {
			return "", err
		}
		for dcid, _ := range c.Dependencies {
			d, ok := cs.Get(dcid)
			if !ok {
				continue
			}
			nAttrs, err := genAttributes(d, fs)
			nAttrs["penwidth"] = "0.75"
			if err != nil {
				return "", err
			}
			err = g.AddNode("G", fmt.Sprintf("\"%s\"", d.ID.String()), nAttrs)
			if err != nil {
				return "", err
			}

			eAttrs := map[string]string{
				"arrowsize": "0.75",
				"penwidth":  "0.75",
			}
			err = g.AddEdge(fmt.Sprintf("\"%s\"", c.ID.String()), fmt.Sprintf("\"%s\"", d.ID.String()), true, eAttrs)
			if err != nil {
				return "", err
			}
		}
	}

	return g.String(), nil
}

func genAttributes(c *component.Component, fs []*Face) (map[string]string, error) {
	attrs := map[string]string{}
	for _, f := range fs {
		pass, err := f.Filter.Pass(c)
		if err != nil {
			return nil, err
		}
		if !pass {
			continue
		}
		for k, v := range f.Attributes {
			if k == "label" {
				label, err := constructLabel(c, v)
				if err != nil {
					return nil, err
				}
				attrs[k] = label
			} else {
				attrs[k] = v
			}
		}
	}

	return attrs, nil
}

func constructLabel(c *component.Component, template string) (string, error) {
	placeholders := []string{}
	capture := false
	var start int
	var end int
	for i, char := range template {
		switch char {
		case '{':
			if capture {
				return "", fmt.Errorf("an embeded label cannot be nested")
			}
			capture = true
			start = i
		case '}':
			if !capture {
				return "", fmt.Errorf("an embeded label is malformed")
			}
			capture = false
			end = i

			placeholder := template[start : end+1]
			placeholders = append(placeholders, placeholder)
		}
	}

	embeddedValues := []string{}
	for _, p := range placeholders {
		labelK := strings.TrimSpace(p[1 : len(p)-1])
		labelV, ok := c.Labels[labelK]
		if !ok {
			return "", fmt.Errorf("ID cannot include the undefined label `%s`", labelK)
		}
		embeddedValues = append(embeddedValues, labelV)
	}

	label := template
	for i, p := range placeholders {
		label = strings.Replace(label, p, embeddedValues[i], 1)
	}

	return label, nil
}
