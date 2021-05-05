package configauditreport

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"context"
	"time"

	"github.com/aquasecurity/starboard/itest/helper"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	assertionTimeout   = 3 * time.Minute
	namespaceName      = corev1.NamespaceDefault
	baseDeploymentName = "wordpress"
)

type SharedBehaviorInputs struct {
	client.Client
	*helper.Helper
}

// SharedBehavior is a Ginkgo container that describes the behavior of a configuration checker.
func SharedBehavior(inputs *SharedBehaviorInputs) func() {
	return func() {
		Context("When unmanaged Pod is created", func() {

			var ctx context.Context
			var pod *corev1.Pod

			BeforeEach(func() {
				ctx = context.Background()
				pod = helper.NewPod().
					WithRandomName("unmanaged-nginx").
					WithNamespace(corev1.NamespaceDefault).
					WithContainer("nginx", "nginx:1.16").
					Build()

				err := inputs.Create(ctx, pod)
				Expect(err).ToNot(HaveOccurred())
			})

			It("Should create ConfigAuditReport", func() {
				Eventually(inputs.HasConfigAuditReportOwnedBy(pod), assertionTimeout).Should(BeTrue())
			})

			AfterEach(func() {
				err := inputs.Delete(ctx, pod)
				Expect(err).ToNot(HaveOccurred())
			})

		})

		Context("When Deployment is created", func() {

			var ctx context.Context
			var deploy *appsv1.Deployment

			BeforeEach(func() {
				ctx = context.Background()
				deploy = helper.NewDeployment().
					WithRandomName(baseDeploymentName).
					WithNamespace(namespaceName).
					WithContainer("wordpress", "wordpress:4.9").
					Build()

				err := inputs.Create(ctx, deploy)
				Expect(err).ToNot(HaveOccurred())
				Eventually(inputs.HasActiveReplicaSet(namespaceName, deploy.Name), assertionTimeout).Should(BeTrue())
			})

			It("Should create ConfigAuditReport", func() {
				rs, err := inputs.GetActiveReplicaSetForDeployment(namespaceName, deploy.Name)
				Expect(err).ToNot(HaveOccurred())
				Expect(rs).ToNot(BeNil())

				Eventually(inputs.HasConfigAuditReportOwnedBy(rs), assertionTimeout).Should(BeTrue())
			})

			AfterEach(func() {
				err := inputs.Delete(ctx, deploy)
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("When Deployment is rolling updated", func() {

			var ctx context.Context
			var deploy *appsv1.Deployment

			BeforeEach(func() {
				By("Creating Deployment")
				ctx = context.Background()
				deploy = helper.NewDeployment().
					WithRandomName(baseDeploymentName).
					WithNamespace(namespaceName).
					WithContainer("wordpress", "wordpress:4.9").
					Build()

				err := inputs.Create(ctx, deploy)
				Expect(err).ToNot(HaveOccurred())
				Eventually(inputs.HasActiveReplicaSet(namespaceName, deploy.Name), assertionTimeout).Should(BeTrue())
			})

			It("Should create ConfigAuditReport for new ReplicaSet", func() {
				By("Getting current active ReplicaSet")
				rs, err := inputs.GetActiveReplicaSetForDeployment(namespaceName, deploy.Name)
				Expect(err).ToNot(HaveOccurred())
				Expect(rs).ToNot(BeNil())

				By("Waiting for ConfigAuditReport")
				Eventually(inputs.HasConfigAuditReportOwnedBy(rs), assertionTimeout).Should(BeTrue())

				By("Updating Deployment image to wordpress:5")
				err = inputs.UpdateDeploymentImage(namespaceName, deploy.Name)
				Expect(err).ToNot(HaveOccurred())

				Eventually(inputs.HasActiveReplicaSet(namespaceName, deploy.Name), assertionTimeout).Should(BeTrue())

				By("Getting new active ReplicaSet")
				rs, err = inputs.GetActiveReplicaSetForDeployment(namespaceName, deploy.Name)
				Expect(err).ToNot(HaveOccurred())
				Expect(rs).ToNot(BeNil())

				By("Waiting for new ConfigAuditReport")
				Eventually(inputs.HasConfigAuditReportOwnedBy(rs), assertionTimeout).Should(BeTrue())
			})

			AfterEach(func() {
				err := inputs.Delete(ctx, deploy)
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("When CronJob is created", func() {

			var ctx context.Context
			var cronJob *batchv1beta1.CronJob

			BeforeEach(func() {
				ctx = context.Background()
				cronJob = &batchv1beta1.CronJob{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: namespaceName,
						Name:      "hello-" + rand.String(5),
					},
					Spec: batchv1beta1.CronJobSpec{
						Schedule: "*/1 * * * *",
						JobTemplate: batchv1beta1.JobTemplateSpec{
							Spec: batchv1.JobSpec{
								Template: corev1.PodTemplateSpec{
									Spec: corev1.PodSpec{
										RestartPolicy: corev1.RestartPolicyOnFailure,
										Containers: []corev1.Container{
											{
												Name:  "hello",
												Image: "busybox",
												Command: []string{
													"/bin/sh",
													"-c",
													"date; echo Hello from the Kubernetes cluster",
												},
											},
										},
									},
								},
							},
						},
					},
				}
				err := inputs.Create(ctx, cronJob)
				Expect(err).ToNot(HaveOccurred())
			})

			It("Should create ConfigAuditReport", func() {
				Eventually(inputs.HasConfigAuditReportOwnedBy(cronJob), assertionTimeout).Should(BeTrue())
			})

			AfterEach(func() {
				err := inputs.Delete(ctx, cronJob)
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("When ConfigAuditReport is deleted", func() {

			var ctx context.Context
			var deploy *appsv1.Deployment

			BeforeEach(func() {
				By("Creating Deployment")
				ctx = context.Background()
				deploy = helper.NewDeployment().
					WithRandomName(baseDeploymentName).
					WithNamespace(namespaceName).
					WithContainer("wordpress", "wordpress:4.9").
					Build()

				err := inputs.Create(ctx, deploy)
				Expect(err).ToNot(HaveOccurred())
				Eventually(inputs.HasActiveReplicaSet(namespaceName, deploy.Name), assertionTimeout).Should(BeTrue())
			})

			It("Should rescan Deployment when ConfigAuditReport is deleted", func() {
				By("Getting active ReplicaSet")
				rs, err := inputs.GetActiveReplicaSetForDeployment(namespaceName, deploy.Name)
				Expect(err).ToNot(HaveOccurred())
				Expect(rs).ToNot(BeNil())

				By("Waiting for ConfigAuditReport")
				Eventually(inputs.HasConfigAuditReportOwnedBy(rs), assertionTimeout).Should(BeTrue())
				By("Deleting ConfigAuditReport")
				err = inputs.DeleteConfigAuditReportOwnedBy(rs)
				Expect(err).ToNot(HaveOccurred())

				By("Waiting for new ConfigAuditReport")
				Eventually(inputs.HasConfigAuditReportOwnedBy(rs), assertionTimeout).Should(BeTrue())
			})

			AfterEach(func() {
				err := inputs.Delete(ctx, deploy)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	}
}
