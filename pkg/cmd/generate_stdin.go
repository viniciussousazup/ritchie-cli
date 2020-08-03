/*
 * Copyright 2020 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/formula/lister"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

type GenerateStdinCmd struct {
	lister lister.Lister
}

func NewGenerateStdinCmd(lister lister.Lister) *cobra.Command {
	lr := GenerateStdinCmd{
		lister: lister,
	}
	cmd := &cobra.Command{
		Use:     "stdin",
		Short:   "Generate a stdin for a formula",
		Example: "rit generate stdin",
		RunE:    lr.runFunc(),
	}
	return cmd
}

func (gsc GenerateStdinCmd) runFunc() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		selectFormulaCmd(gsc.lister.List())

		prompt.Success("Hello World.")
		return nil
	}
}

func selectFormulaCmd(commands []*cobra.Command) error {
	var cmds = map[string][]*cobra.Command{}
	var cmdsKeys = []string{}
	for _, c := range commands {
		if !c.Hidden {
			k := c.Name()
			cmds[k] = c.Commands()
			cmdsKeys = append(cmdsKeys, k)
		}

	}
	selected, err := prompt.NewSurveyList().List("Select command", cmdsKeys)
	if err != nil {
		return err
	}
	cmdS, _ := cmds[selected]
	if len(cmds) != 0 {
		return selectFormulaCmd(cmdS)
	}
	return nil
}
