package starboard

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aquasecurity/starboard/pkg/apis/aquasecurity/v1alpha1"
	"github.com/google/go-containerregistry/pkg/name"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	coordinationv1 "k8s.io/api/coordination/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
)

func NewScheme() *runtime.Scheme {
	scheme := runtime.NewScheme()
	_ = corev1.AddToScheme(scheme)
	_ = appsv1.AddToScheme(scheme)
	_ = batchv1.AddToScheme(scheme)
	_ = batchv1beta1.AddToScheme(scheme)
	_ = rbacv1.AddToScheme(scheme)
	_ = v1alpha1.AddToScheme(scheme)
	_ = coordinationv1.AddToScheme(scheme)
	_ = apiextensionsv1.AddToScheme(scheme)
	return scheme
}

// BuildInfo holds build info such as Git revision, Git SHA-1,
// build datetime, and the name of the executable binary.
type BuildInfo struct {
	Version    string
	Commit     string
	Date       string
	Executable string
}

// Scanner represents unique, human readable identifier of a security scanner.
type Scanner string

const (
	Trivy    Scanner = "Trivy"
	Aqua     Scanner = "Aqua"
	Polaris  Scanner = "Polaris"
	Conftest Scanner = "Conftest"
)

const (
	keyVulnerabilityReportsScanner = "vulnerabilityReports.scanner"
	keyConfigAuditReportsScanner   = "configAuditReports.scanner"
)

// ConfigData holds Starboard configuration settings as a set
// of key-value pairs.
type ConfigData map[string]string

// ConfigManager defines methods for managing ConfigData.
type ConfigManager interface {
	EnsureDefault(ctx context.Context) error
	Read(ctx context.Context) (ConfigData, error)
	Delete(ctx context.Context) error
}

// GetDefaultConfig returns the default configuration settings.
func GetDefaultConfig() ConfigData {
	return map[string]string{
		keyVulnerabilityReportsScanner: string(Trivy),
		keyConfigAuditReportsScanner:   string(Polaris),

		"kube-bench.imageRef":  "docker.io/aquasec/kube-bench:v0.6.5",
		"kube-hunter.imageRef": "docker.io/aquasec/kube-hunter:0.6.3",
		"kube-hunter.quick":    "false",
	}
}

func (c ConfigData) GetScanJobTolerations() ([]corev1.Toleration, error) {
	scanJobTolerations := []corev1.Toleration{}
	if c["scanJob.tolerations"] == "" {
		return scanJobTolerations, nil
	}
	err := json.Unmarshal([]byte(c["scanJob.tolerations"]), &scanJobTolerations)

	return scanJobTolerations, err
}

func (c ConfigData) GetVulnerabilityReportsScanner() (Scanner, error) {
	var ok bool
	var value string
	if value, ok = c[keyVulnerabilityReportsScanner]; !ok {
		return "", fmt.Errorf("property %s not set", keyVulnerabilityReportsScanner)
	}

	switch Scanner(value) {
	case Trivy:
		return Trivy, nil
	case Aqua:
		return Aqua, nil
	}

	return "", fmt.Errorf("invalid value (%s) of %s; allowed values (%s, %s)",
		value, keyVulnerabilityReportsScanner, Trivy, Aqua)
}

func (c ConfigData) GetConfigAuditReportsScanner() (Scanner, error) {
	var ok bool
	var value string
	if value, ok = c[keyConfigAuditReportsScanner]; !ok {
		return "", fmt.Errorf("property %s not set", keyConfigAuditReportsScanner)
	}

	switch Scanner(value) {
	case Polaris:
		return Polaris, nil
	case Conftest:
		return Conftest, nil
	}
	return "", fmt.Errorf("invalid value (%s) of %s; allowed values (%s, %s)",
		value, keyConfigAuditReportsScanner, Polaris, Conftest)
}

func (c ConfigData) GetScanJobAnnotations() (map[string]string, error) {
	scanJobAnnotationsStr, found := c[AnnotationScanJobAnnotations]
	if !found || strings.TrimSpace(scanJobAnnotationsStr) == "" {
		return map[string]string{}, nil
	}

	scanJobAnnotationsMap := map[string]string{}
	for _, annotation := range strings.Split(scanJobAnnotationsStr, ",") {
		sepByEqual := strings.Split(annotation, "=")
		if len(sepByEqual) != 2 {
			return map[string]string{}, fmt.Errorf("custom annotations found to be wrongfully provided: %s", scanJobAnnotationsStr)
		}
		key, value := sepByEqual[0], sepByEqual[1]
		scanJobAnnotationsMap[key] = value
	}

	return scanJobAnnotationsMap, nil
}

func (c ConfigData) GetScanJobTemplateLabels() (labels.Set, error) {
	scanJobTemplateLabelsStr, found := c[AnnotationScanJobTemplateLabels]
	if !found || strings.TrimSpace(scanJobTemplateLabelsStr) == "" {
		return labels.Set{}, nil
	}

	scanJobTemplateLabelsMap := map[string]string{}
	for _, annotation := range strings.Split(scanJobTemplateLabelsStr, ",") {
		sepByEqual := strings.Split(annotation, "=")
		if len(sepByEqual) != 2 {
			return labels.Set{}, fmt.Errorf("custom template labels found to be wrongfully provided: %s", scanJobTemplateLabelsStr)
		}
		key, value := sepByEqual[0], sepByEqual[1]
		scanJobTemplateLabelsMap[key] = value
	}

	return scanJobTemplateLabelsMap, nil
}

func (c ConfigData) GetKubeBenchImageRef() (string, error) {
	return c.GetRequiredData("kube-bench.imageRef")
}

func (c ConfigData) GetKubeHunterImageRef() (string, error) {
	return c.GetRequiredData("kube-hunter.imageRef")
}

func (c ConfigData) GetKubeHunterQuick() (bool, error) {
	val, ok := c["kube-hunter.quick"]
	if !ok {
		return false, nil
	}
	if val != "false" && val != "true" {
		return false, fmt.Errorf("property kube-hunter.quick must be either \"false\" or \"true\", got %q", val)
	}
	return val == "true", nil
}

func (c ConfigData) GetRequiredData(key string) (string, error) {
	var ok bool
	var value string
	if value, ok = c[key]; !ok {
		return "", fmt.Errorf("property %s not set", key)
	}
	return value, nil
}

// GetVersionFromImageRef returns the image identifier for the specified image reference.
func GetVersionFromImageRef(imageRef string) (string, error) {
	ref, err := name.ParseReference(imageRef)
	if err != nil {
		return "", fmt.Errorf("parsing reference: %w", err)
	}

	var version string
	switch t := ref.(type) {
	case name.Tag:
		version = t.TagStr()
	case name.Digest:
		version = t.DigestStr()
	}

	return version, nil
}

// NewConfigManager constructs a new ConfigManager that is using kubernetes.Interface
// to manage ConfigData backed by the ConfigMap stored in the specified namespace.
func NewConfigManager(client kubernetes.Interface, namespace string) ConfigManager {
	return &configManager{
		client:    client,
		namespace: namespace,
	}
}

type configManager struct {
	client    kubernetes.Interface
	namespace string
}

func (c *configManager) EnsureDefault(ctx context.Context) error {
	_, err := c.client.CoreV1().ConfigMaps(c.namespace).Get(ctx, ConfigMapName, metav1.GetOptions{})
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return fmt.Errorf("getting config: %w", err)
		}
		_, err = c.client.CoreV1().ConfigMaps(c.namespace).Create(ctx, &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: c.namespace,
				Name:      ConfigMapName,
				Labels: labels.Set{
					LabelK8SAppManagedBy: "starboard",
				},
			},
			Data: GetDefaultConfig(),
		}, metav1.CreateOptions{})

		if err != nil {
			return fmt.Errorf("creating config: %w", err)
		}
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: c.namespace,
			Name:      SecretName,
			Labels: labels.Set{
				LabelK8SAppManagedBy: "starboard",
			},
		},
	}
	_, err = c.client.CoreV1().Secrets(c.namespace).Create(ctx, secret, metav1.CreateOptions{})
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return err
	}

	return nil
}

func (c *configManager) Read(ctx context.Context) (ConfigData, error) {
	cm, err := c.client.CoreV1().ConfigMaps(c.namespace).Get(ctx, ConfigMapName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	secret, err := c.client.CoreV1().Secrets(c.namespace).Get(ctx, SecretName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	var data = make(map[string]string)

	for k, v := range cm.Data {
		data[k] = v
	}

	for k, v := range secret.Data {
		data[k] = string(v)
	}

	return data, nil
}

func (c *configManager) Delete(ctx context.Context) error {
	err := c.client.CoreV1().ConfigMaps(c.namespace).Delete(ctx, ConfigMapName, metav1.DeleteOptions{})
	if err != nil && !apierrors.IsNotFound(err) {
		return err
	}
	err = c.client.CoreV1().ConfigMaps(c.namespace).Delete(ctx, GetPluginConfigMapName(string(Polaris)), metav1.DeleteOptions{})
	if err != nil && !apierrors.IsNotFound(err) {
		return err
	}
	err = c.client.CoreV1().Secrets(c.namespace).Delete(ctx, SecretName, metav1.DeleteOptions{})
	if err != nil && !apierrors.IsNotFound(err) {
		return err
	}
	return nil
}

// LinuxNodeAffinity constructs a new Affinity resource with linux supported nodes.
func LinuxNodeAffinity() *corev1.Affinity {
	return &corev1.Affinity{
		NodeAffinity: &corev1.NodeAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
				NodeSelectorTerms: []corev1.NodeSelectorTerm{
					{
						MatchExpressions: []corev1.NodeSelectorRequirement{
							{
								Key:      "kubernetes.io/os",
								Operator: corev1.NodeSelectorOpIn,
								Values:   []string{"linux"},
							},
						},
					},
				}}}}
}
