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

package linux

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMapYumPackageVersion(t *testing.T) {
	testCases := []struct {
		pkg               string
		kubernetesVersion string
		expected          string
	}{
		{kubeadm, "1.14.0", "kubeadm-1.14.0-0"},
		{kubectl, "1.14.0", "kubectl-1.14.0-0"},
		{kubelet, "1.14.0", "kubelet-1.14.0-0"},
		{kubernetescni, "1.14.0", "kubernetes-cni-0.7.5-0"},
		{kubernetescni, "1.13.0", "kubernetes-cni-0.6.0-0"},
	}
	for _, tc := range testCases {
		got := mapYumPackageVersion(tc.pkg, tc.kubernetesVersion)
		require.Equal(t, tc.expected, got)
	}
}

func TestParseRpmPackageOutput(t *testing.T) {
	testCases := []struct {
		pkg     string
		name    string
		version string
		release string
		arch    string
		err     bool
	}{
		{"kubernetes-cni-0.7.5-0.x86_64", "kubernetes-cni", "0.7.5", "0", "x86_64", false},
		{"kubeadm-1.14.0-0.x86_64", "kubeadm", "1.14.0", "0", "x86_64", false},
		{"kubeadm", "", "", "", "", true},
		{"util-linux-2.23.2-59.el7.x86_64", "util-linux", "2.23.2", "59.el7", "x86_64", false},
		{"systemd-219-62.el7_6.5.x86_64", "systemd", "219", "62.el7_6.5", "x86_64", false},
	}
	for _, tc := range testCases {
		name, ver, rel, arch, err := parseRpmPackageOutput(tc.pkg)
		if tc.err {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}
		require.Equal(t, tc.name, name)
		require.Equal(t, tc.version, ver)
		require.Equal(t, tc.release, rel)
		require.Equal(t, tc.arch, arch)
	}
}
