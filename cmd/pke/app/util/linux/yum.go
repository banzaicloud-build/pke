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
	"fmt"
	"io"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/banzaicloud/pke/cmd/pke/app/util/runner"
	"github.com/pkg/errors"
)

const (
	cmdYum                    = "/bin/yum"
	cmdRpm                    = "/bin/rpm"
	kubeadm                   = "kubeadm"
	kubectl                   = "kubectl"
	kubelet                   = "kubelet"
	kubernetescni             = "kubernetes-cni"
	disableExcludesKubernetes = "--disableexcludes=kubernetes"

	dotS  = "."
	dashS = "-"
)

var (
	errorUnableToParseRPMOutput = errors.New("Unable to parse rpm output")
)

func YumInstall(out io.Writer, packages []string) error {
	err := runner.Cmd(out, cmdYum, append([]string{"install", "-y"}, packages...)...).CombinedOutputAsync()
	if err != nil {
		return err
	}

	for _, pkg := range packages {
		if pkg[:1] == "-" {
			continue
		}

		name, ver, rel, arch, err := rpmQuery(out, pkg)
		if err != nil {
			return err
		}
		if name == pkg ||
			name+"-"+ver == pkg ||
			name+"-"+ver+"-"+rel == pkg ||
			name+"-"+ver+"-"+rel+"."+arch == pkg {
			continue
		}
		return errors.New(fmt.Sprintf("expected packgae version after installation: %q, got: %q", pkg, name+"-"+ver+"-"+rel+"."+arch))
	}

	return nil
}

func rpmQuery(out io.Writer, pkg string) (name, version, release, arch string, err error) {
	b, err := runner.Cmd(out, cmdRpm, []string{"-q", pkg}...).Output()
	if err != nil {
		return
	}

	return parseRpmPackageOutput(string(b))
}

func parseRpmPackageOutput(pkg string) (name, version, release, arch string, err error) {
	idx := strings.LastIndex(pkg, dotS)
	if idx < 0 {
		err = errorUnableToParseRPMOutput
		return
	}
	arch = pkg[idx+1:]

	pkg = pkg[:idx]
	idx = strings.LastIndex(pkg, dashS)
	if idx < 0 {
		err = errorUnableToParseRPMOutput
		return
	}
	release = pkg[idx+1:]

	pkg = pkg[:idx]
	idx = strings.LastIndex(pkg, dashS)
	if idx < 0 {
		err = errorUnableToParseRPMOutput
		return
	}
	version = pkg[idx+1:]
	name = pkg[:idx]

	return
}

var _ KubernetesPackages = (*YumInstaller)(nil)

type YumInstaller struct{}

func NewYumInstaller() *YumInstaller {
	return &YumInstaller{}
}

func (y *YumInstaller) InstallKubernetesPackages(out io.Writer, kubernetesVersion string) error {
	p := []string{
		mapYumPackageVersion(kubelet, kubernetesVersion),
		mapYumPackageVersion(kubeadm, kubernetesVersion),
		mapYumPackageVersion(kubectl, kubernetesVersion),
		mapYumPackageVersion(kubernetescni, kubernetesVersion),
		disableExcludesKubernetes,
	}

	return YumInstall(out, p)
}

func (y *YumInstaller) InstallKubeadmPackage(out io.Writer, kubernetesVersion string) error {
	pkg := []string{
		mapYumPackageVersion(kubeadm, kubernetesVersion),
		mapYumPackageVersion(kubelet, kubernetesVersion),       // kubeadm dependency
		mapYumPackageVersion(kubernetescni, kubernetesVersion), // kubeadm dependency
		"--disableexcludes=kubernetes",
	}
	return YumInstall(out, pkg)
}

func mapYumPackageVersion(pkg, kubernetesVersion string) string {
	switch pkg {
	case kubeadm:
		return "kubeadm-" + kubernetesVersion + "-0"

	case kubectl:
		return "kubectl-" + kubernetesVersion + "-0"

	case kubelet:
		return "kubelet-" + kubernetesVersion + "-0"

	case kubernetescni:
		ver, _ := semver.NewVersion(kubernetesVersion)
		if !ver.LessThan(semver.MustParse("1.14.0")) {
			return "kubernetes-cni-0.7.5-0"
		}
		return "kubernetes-cni-0.6.0-0"

	default:
		return ""
	}
}
