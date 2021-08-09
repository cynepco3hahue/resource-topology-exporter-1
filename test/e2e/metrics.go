package e2e

import (
	"context"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/test/e2e/framework"

	"github.com/k8stopologyawareschedwg/resource-topology-exporter/test/e2e/utils"
)

var _ = ginkgo.Describe("[RTE] metrics", func() {
	var (
		initialized         bool
	)

	f := framework.NewDefaultFramework("metrics")

	ginkgo.BeforeEach(func() {
		var err error

		if !initialized {
			// if we'll get an error here the system is not ready to be tested
			_, err = f.ClientSet.CoreV1().Nodes().Get(context.TODO(), getNodeName(), metav1.GetOptions{})
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			initialized = true
		}
	})

	ginkgo.Context("With prometheus endpoint configured", func() {
		ginkgo.It("should have some metrics exported", func() {
			pods, err := utils.GetByRegex(f.ClientSet, defaultNamespace, "resource-topology-exporter-*")
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(len(pods)).Should(gomega.BeEquivalentTo(1))

			stdout, stderr, err := f.ExecWithOptions(framework.ExecOptions{
				Command:            []string{"curl", "http://127.0.0.1:2112/metrics"},
				Namespace:          getNamespaceName(),
				PodName:            pods[0].Name,
				ContainerName:      rteContainerName,
				Stdin:              nil,
				CaptureStdout:      true,
				CaptureStderr:      true,
				PreserveWhitespace: false,
			})
			gomega.Expect(err).NotTo(gomega.HaveOccurred(), "%s", stderr)
			gomega.Expect(stdout).To(gomega.ContainSubstring("topology_updater_api_call_failures_total"))
			gomega.Expect(stdout).To(gomega.ContainSubstring("topology_updater_hw_topology_update_operation_measurement"))
		})
	})
})
