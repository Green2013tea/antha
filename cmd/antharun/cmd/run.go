// run.go: Part of the Antha language
// Copyright (C) 2015 The Antha authors. All rights reserved.
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
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"

	"github.com/antha-lang/antha/antha/component/lib"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/cmd/antharun/frontend"
	"github.com/antha-lang/antha/cmd/antharun/param"
	"github.com/antha-lang/antha/cmd/antharun/pretty"
	"github.com/antha-lang/antha/cmd/antharun/spawn"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/target/auto"
	"github.com/antha-lang/antha/target/human"
	"github.com/antha-lang/antha/target/mixer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultPort = 50051
)

var runCmd = &cobra.Command{
	Use:          "antharun",
	Short:        "Run an antha workflow",
	RunE:         runWorkflow,
	SilenceUsage: true,
}

func makeMixerOpt() mixer.Opt {
	opt := mixer.Opt{}
	if i := viper.GetInt("maxPlates"); i != 0 {
		f := float64(i)
		opt.MaxPlates = &f
	}
	if i := viper.GetInt("maxWells"); i != 0 {
		f := float64(i)
		opt.MaxWells = &f
	}
	if i := viper.GetFloat64("residualVolumeWeight"); i != 0 {
		f := float64(i)
		opt.ResidualVolumeWeight = &f
	}
	opt.InputPlateType = GetStringSlice("inputPlateType")
	opt.OutputPlateType = GetStringSlice("outputPlateType")
	opt.TipType = GetStringSlice("tipType")
	opt.InputPlateFiles = GetStringSlice("inputPlates")
	return opt
}

func makeContext() (context.Context, error) {
	ctx := inject.NewContext(context.Background())
	for _, desc := range lib.GetComponents() {
		obj := desc.Constructor()
		runner, ok := obj.(inject.Runner)
		if !ok {
			return nil, fmt.Errorf("component %q has unexpected type %T", desc.Name, obj)
		}
		if err := inject.Add(ctx, inject.Name{Repo: desc.Name}, runner); err != nil {
			return nil, err
		}
	}
	return ctx, nil
}

type runOpt struct {
	MixerOpt       mixer.Opt
	Drivers        []string
	WorkflowFile   string
	ParametersFile string
}

func (a *runOpt) makeTarget(mixerOpt mixer.Opt) (*target.Target, error) {
	t := target.New()
	if len(a.Drivers) == 0 {
		if err := t.AddDevice(human.New(human.Opt{CanMix: true, CanIncubate: true})); err != nil {
			return nil, err
		}
	} else {
		if err := t.AddDevice(human.New(human.Opt{CanIncubate: true})); err != nil {
			return nil, err
		}
	}
	for _, uri := range a.Drivers {
		if d, err := auto.New(auto.Opt{
			Uri:  uri,
			Opts: []interface{}{mixerOpt},
		}); err != nil {
			return nil, err
		} else if err := t.AddDevice(d); err != nil {
			return nil, err
		}
	}
	return t, nil
}

func (a *runOpt) Run() error {
	wdata, err := ioutil.ReadFile(a.WorkflowFile)
	if err != nil {
		return err
	}

	pdata, err := ioutil.ReadFile(a.ParametersFile)
	if err != nil {
		return err
	}

	wdesc, params, err := param.TryExpand(wdata, pdata)
	if err != nil {
		return err
	}

	mixerOpt := mixer.DefaultOpt.Merge(params.Config).Merge(&a.MixerOpt)
	t, err := a.makeTarget(mixerOpt)
	if err != nil {
		return err
	}

	fe, err := frontend.New(frontend.Opt{})
	if err != nil {
		return err
	}
	defer fe.Shutdown()

	ctx, err := makeContext()
	if err != nil {
		return err
	}

	rout, err := execute.Run(ctx, execute.Opt{
		Target:   t,
		Workflow: wdesc,
		Params:   params,
	})
	if err != nil {
		return err
	}

	if err := pretty.Timeline(os.Stdout, rout); err != nil {
		return err
	}

	if err := pretty.Run(os.Stdout, rout, t); err != nil {
		return err
	}

	return nil
}

func runWorkflow(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	var drivers []string
	for idx, uri := range GetStringSlice("driver") {
		u, err := url.Parse(uri)
		if err != nil {
			return err
		}

		switch u.Scheme {
		case "go":
			p := u.Host + u.Path
			s, err := spawn.GoPackage(p, fmt.Sprintf("%d %s", idx, path.Base(u.Path)))
			if s != nil {
				defer s.Close()
			}
			if err != nil {
				return fmt.Errorf("cannot start package %s: %s", p, err)
			} else if err := s.Start(); err != nil {
				return fmt.Errorf("cannot start package %s: %s", p, err)
			} else if uri, err := s.Uri(); err != nil {
				return fmt.Errorf("cannot parse port for package %s: %s", p, err)
			} else {
				drivers = append(drivers, uri)
			}
		default:
			drivers = append(drivers, u.String())
		}
	}

	opt := &runOpt{
		MixerOpt:       makeMixerOpt(),
		Drivers:        drivers,
		WorkflowFile:   viper.GetString("workflow"),
		ParametersFile: viper.GetString("parameters"),
	}

	return opt.Run()
}

func init() {
	c := runCmd
	flags := c.Flags()

	//RootCmd.AddCommand(c)
	flags.String("parameters", "parameters.yml", "Parameters to workflow")
	flags.String("workflow", "workflow.json", "Workflow definition file")
	flags.StringSlice("driver", nil, "Uris of remote drivers ({tcp,go}://...)")
	flags.StringSlice("component", nil, "Uris of remote components ({tcp,go}://...)")
	flags.Int("maxPlates", 0, "Maximum number of plates")
	flags.Int("maxWells", 0, "Maximum number of wells on a plate")
	flags.Float64("residualVolumeWeight", 0.0, "Residual volume weight")
	flags.StringSlice("inputPlateType", nil, "Default input plate types (in order of preference)")
	flags.StringSlice("outputPlateType", nil, "Default output plate types (in order of preference)")
	flags.StringSlice("inputPlates", nil, "File containing input plates")
	flags.StringSlice("tipType", nil, "Names of permitted tip types")
}
