/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package e2e

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/runtime"

	"sigs.k8s.io/cluster-api/test/framework"
	"sigs.k8s.io/cluster-api/test/framework/management/kind"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/api/v1alpha3"
	frameworkx "sigs.k8s.io/cluster-api-provider-vsphere/test/e2e/framework"
	kindx "sigs.k8s.io/cluster-api-provider-vsphere/test/e2e/kind"
)

var (
	mgmt *kind.Cluster
	ctx  = context.Background()
)

func TestCAPV(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CAPV e2e Suite")
}

var _ = BeforeSuite(func() {
	By("cleaning up previous kind cluster")
	kindx.TeardownIfExists(ctx, frameworkx.Flags.ManagementClusterName)

	By("initializing the vSphere session", initVSphereSession)

	By("initializing the runtime.Scheme")
	scheme := runtime.NewScheme()
	Expect(infrav1.AddToScheme(scheme)).To(Succeed())

	mgmt = frameworkx.InitManagementCluster(ctx, &frameworkx.InitManagementClusterInput{
		InfraNamespace: "capv-system",
		InfraComponentGenerators: []framework.ComponentGenerator{
			providerGenerator{},
			credentialsGenerator{},
		},
		Scheme: scheme,
	})
})

var _ = AfterSuite(func() {
	By("tearing down the management cluster")
	Expect(mgmt.Teardown(ctx)).To(Succeed())
})
