// Copyright © 2019 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"io"

	"github.com/banzaicloud/pke/cmd/pke/app/phases"
	"github.com/banzaicloud/pke/cmd/pke/app/phases/kubeadm/upgrade/controlplane"
	"github.com/banzaicloud/pke/cmd/pke/app/phases/kubeadm/upgrade/node"
	"github.com/banzaicloud/pke/cmd/pke/app/phases/kubeadm/version"
	"github.com/spf13/cobra"
)

func NewCmdUpgrade(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade a single Banzai Cloud Pipeline Kubernetes Engine (PKE) machine",
		Args:  cobra.NoArgs,
		RunE:  phases.RunEAllSubcommands,
	}

	cmd.AddCommand(upgradeMaster(out))
	cmd.AddCommand(upgradeWorker(out))

	return cmd
}

func upgradeMaster(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "master",
		Short: "Upgrade a single Banzai Cloud Pipeline Kubernetes Engine (PKE) master machine",
		Args:  cobra.NoArgs,
		RunE:  phases.RunEAllSubcommands,
	}

	cmd.AddCommand(version.NewCommand(out))
	cmd.AddCommand(controlplane.NewCommand(out))

	phases.MakeRunnable(cmd)

	return cmd
}

func upgradeWorker(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "worker",
		Short: "Upgrade a single Banzai Cloud Pipeline Kubernetes Engine (PKE) worker machine",
		Args:  cobra.NoArgs,
		RunE:  phases.RunEAllSubcommands,
	}

	cmd.AddCommand(version.NewCommand(out))
	cmd.AddCommand(node.NewCommand(out))

	phases.MakeRunnable(cmd)

	return cmd
}
