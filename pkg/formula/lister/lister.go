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

package lister

import (
	"github.com/spf13/cobra"
)

type Lister interface {
	SetRoot(rootCmd *cobra.Command)
	List() []*cobra.Command
}

type DefaultLister struct {
	rootCmd *cobra.Command
}

func NewLister() Lister {
	return &DefaultLister{}
}

func (l *DefaultLister) SetRoot(rootCmd *cobra.Command) {
	l.rootCmd = rootCmd
}

func (l *DefaultLister) List() []*cobra.Command {
	return l.rootCmd.Commands()
}
