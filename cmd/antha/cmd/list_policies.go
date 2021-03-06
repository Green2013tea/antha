// list_policies.go: Part of the Antha language
// Copyright (C) 2016 The Antha authors. All rights reserved.
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
//
// For more information relating to the software or licensing issues please
// contact license@antha-lang.org or write to the Antha team c/o
// Synthace Ltd. The London Bioscience Innovation Centre
// 2 Royal College St, London NW1 0NH UK

package cmd

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"github.com/ghodss/yaml"
	"github.com/mgutz/ansi"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listPoliciesCmd = &cobra.Command{
	Use:   "policies",
	Short: "List available antha liquid handling policies",
	RunE:  listPolicies,
}

type simplePolicy struct {
	Name       string
	Properties map[string]interface{}
}

type simplePolicies []simplePolicy

func (a simplePolicies) Less(i, j int) bool {
	return a[i].Name < a[j].Name
}

func (a simplePolicies) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a simplePolicies) Len() int {
	return len(a)
}

func listPolicies(cmd *cobra.Command, args []string) error {
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	red := func(x string) string {
		return ansi.Color(x, "red")
	}

	var ps simplePolicies
	for name, p := range liquidhandling.MakePolicies() {

		ps = append(ps, simplePolicy{
			Name:       name,
			Properties: p,
		})
	}

	sort.Sort(ps)

	output := viper.GetString("output")
	switch output {
	case jsonOutput:
		bs, err := json.MarshalIndent(ps, "", "  ")
		if err != nil {
			return err
		}
		_, err = fmt.Println(string(bs))
		return err
	case yamlOutput:
		bs, err := yaml.Marshal(ps)
		if err != nil {
			return err
		}
		_, err = fmt.Print(string(bs))
		return err
	case textOutput:
		var lines []string
		lines = append(lines, red("PolicyName")+" Properties")

		for _, p := range ps {
			var kvs []string
			for k, v := range p.Properties {
				kvs = append(kvs, fmt.Sprintf("%s: %v", k, v))
			}
			sort.Strings(kvs)

			lines = append(lines, red(p.Name)+" "+fmt.Sprint(kvs))
		}

		_, err := fmt.Println(strings.Join(lines, "\n"))
		return err
	case csvOutput:
		var lines []string
		lines = append(lines, "PolicyName,Properties")

		for _, p := range ps {
			var kvs []string
			for k, v := range p.Properties {
				kvs = append(kvs, fmt.Sprintf("%s: %v", k, v))
			}
			sort.Strings(kvs)

			lines = append(lines, p.Name+","+strings.Join(kvs, ","))
		}
		_, err := fmt.Println(strings.Join(lines, "\n"))
		return err
	default:
		return fmt.Errorf("unknown output format %q", output)
	}
}

func init() {
	c := listPoliciesCmd
	listCmd.AddCommand(c)
}
