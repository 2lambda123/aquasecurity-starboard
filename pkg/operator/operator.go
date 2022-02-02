package operator

import (
	"context"
	"fmt"
	"github.com/aquasecurity/starboard/pkg/apis/aquasecurity/v1alpha1"
	"github.com/aquasecurity/starboard/pkg/configauditreport"
	"github.com/aquasecurity/starboard/pkg/ext"
	"github.com/aquasecurity/starboard/pkg/kube"
	"github.com/aquasecurity/starboard/pkg/kubebench"
	"github.com/aquasecurity/starboard/pkg/nsa"
	"github.com/aquasecurity/starboard/pkg/operator/controller"
	"github.com/aquasecurity/starboard/pkg/operator/etc"
	"github.com/aquasecurity/starboard/pkg/plugin"
	"github.com/aquasecurity/starboard/pkg/starboard"
	"github.com/aquasecurity/starboard/pkg/vulnerabilityreport"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

var (
	setupLog = log.Log.WithName("operator")
)

// Start starts all registered reconcilers and blocks until the context is cancelled.
// Returns an error if there is an error starting any reconciler.
func Start(ctx context.Context, buildInfo starboard.BuildInfo, operatorConfig etc.Config) error {
	installMode, operatorNamespace, targetNamespaces, err := operatorConfig.ResolveInstallMode()
	if err != nil {
		return fmt.Errorf("resolving install mode: %w", err)
	}
	setupLog.Info("Resolved install mode", "install mode", installMode,
		"operator namespace", operatorNamespace,
		"target namespaces", targetNamespaces)

	// Set the default manager options.
	options := manager.Options{
		Scheme:                 starboard.NewScheme(),
		MetricsBindAddress:     operatorConfig.MetricsBindAddress,
		HealthProbeBindAddress: operatorConfig.HealthProbeBindAddress,
	}

	if operatorConfig.LeaderElectionEnabled {
		options.LeaderElection = operatorConfig.LeaderElectionEnabled
		options.LeaderElectionID = operatorConfig.LeaderElectionID
		options.LeaderElectionNamespace = operatorNamespace
	}

	switch installMode {
	case etc.OwnNamespace:
		// Add support for OwnNamespace set in OPERATOR_NAMESPACE (e.g. `starboard-operator`)
		// and OPERATOR_TARGET_NAMESPACES (e.g. `starboard-operator`).
		setupLog.Info("Constructing client cache", "namespace", targetNamespaces[0])
		options.Namespace = targetNamespaces[0]
	case etc.SingleNamespace:
		// Add support for SingleNamespace set in OPERATOR_NAMESPACE (e.g. `starboard-operator`)
		// and OPERATOR_TARGET_NAMESPACES (e.g. `default`).
		cachedNamespaces := append(targetNamespaces, operatorNamespace)
		if operatorConfig.CISKubernetesBenchmarkEnabled {
			// Cache cluster-scoped resources such as Nodes
			cachedNamespaces = append(cachedNamespaces, "")
		}
		setupLog.Info("Constructing client cache", "namespaces", cachedNamespaces)
		options.NewCache = cache.MultiNamespacedCacheBuilder(cachedNamespaces)
	case etc.MultiNamespace:
		// Add support for MultiNamespace set in OPERATOR_NAMESPACE (e.g. `starboard-operator`)
		// and OPERATOR_TARGET_NAMESPACES (e.g. `default,kube-system`).
		// Note that you may face performance issues when using this mode with a high number of namespaces.
		// More: https://godoc.org/github.com/kubernetes-sigs/controller-runtime/pkg/cache#MultiNamespacedCacheBuilder
		cachedNamespaces := append(targetNamespaces, operatorNamespace)
		if operatorConfig.CISKubernetesBenchmarkEnabled {
			// Cache cluster-scoped resources such as Nodes
			cachedNamespaces = append(cachedNamespaces, "")
		}
		setupLog.Info("Constructing client cache", "namespaces", cachedNamespaces)
		options.NewCache = cache.MultiNamespacedCacheBuilder(cachedNamespaces)
	case etc.AllNamespaces:
		// Add support for AllNamespaces set in OPERATOR_NAMESPACE (e.g. `operators`)
		// and OPERATOR_TARGET_NAMESPACES left blank.
		setupLog.Info("Watching all namespaces")
	default:
		return fmt.Errorf("unrecognized install mode: %v", installMode)
	}

	kubeConfig, err := ctrl.GetConfig()
	if err != nil {
		return fmt.Errorf("getting kube client config: %w", err)
	}

	// The only reason we're using kubernetes.Clientset is that we need it to read Pod logs,
	// which is not supported by the client returned by the ctrl.Manager.
	kubeClientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return fmt.Errorf("constructing kube client: %w", err)
	}

	mgr, err := ctrl.NewManager(kubeConfig, options)
	if err != nil {
		return fmt.Errorf("constructing controllers manager: %w", err)
	}

	err = mgr.AddReadyzCheck("ping", healthz.Ping)
	if err != nil {
		return err
	}

	err = mgr.AddHealthzCheck("ping", healthz.Ping)
	if err != nil {
		return err
	}

	configManager := starboard.NewConfigManager(kubeClientset, operatorNamespace)
	err = configManager.EnsureDefault(context.Background())
	if err != nil {
		return err
	}

	starboardConfig, err := configManager.Read(context.Background())
	if err != nil {
		return err
	}

	objectResolver := kube.ObjectResolver{Client: mgr.GetClient()}
	limitChecker := controller.NewLimitChecker(operatorConfig, mgr.GetClient())
	logsReader := kube.NewLogsReader(kubeClientset)
	secretsReader := kube.NewSecretsReader(mgr.GetClient())

	if operatorConfig.VulnerabilityScannerEnabled {
		plugin, pluginContext, err := plugin.NewResolver().
			WithBuildInfo(buildInfo).
			WithNamespace(operatorNamespace).
			WithServiceAccountName(operatorConfig.ServiceAccount).
			WithConfig(starboardConfig).
			WithClient(mgr.GetClient()).
			GetVulnerabilityPlugin()
		if err != nil {
			return err
		}

		err = plugin.Init(pluginContext)
		if err != nil {
			return fmt.Errorf("initializing %s plugin: %w", pluginContext.GetName(), err)
		}

		if err = (&controller.VulnerabilityReportReconciler{
			Logger:         ctrl.Log.WithName("reconciler").WithName("vulnerabilityreport"),
			Config:         operatorConfig,
			ConfigData:     starboardConfig,
			Client:         mgr.GetClient(),
			ObjectResolver: objectResolver,
			LimitChecker:   limitChecker,
			LogsReader:     logsReader,
			SecretsReader:  secretsReader,
			Plugin:         plugin,
			PluginContext:  pluginContext,
			ReadWriter:     vulnerabilityreport.NewReadWriter(mgr.GetClient()),
		}).SetupWithManager(mgr); err != nil {
			return fmt.Errorf("unable to setup vulnerabilityreport reconciler: %w", err)
		}

		if operatorConfig.VulnerabilityScannerReportTTL != nil {
			if err = (&controller.TTLReportReconciler{
				Logger: ctrl.Log.WithName("reconciler").WithName("ttlreport"),
				Config: operatorConfig,
				Client: mgr.GetClient(),
			}).SetupWithManager(mgr); err != nil {
				return fmt.Errorf("unable to setup TTLreport reconciler: %w", err)
			}
		}
	}

	if operatorConfig.ConfigAuditScannerEnabled {
		plugin, pluginContext, err := createConfPlugin(buildInfo, operatorNamespace, operatorConfig, starboardConfig, mgr)
		if err != nil {
			return err
		}
		rw := configauditreport.NewReadWriter(mgr.GetClient())
		confAuditController := createConfAuditController(operatorConfig, starboardConfig, mgr, objectResolver, limitChecker, logsReader, plugin, controller.ConfigAuditJobName, rw.WriteClusterReport, rw.FindClusterReportByOwner, pluginContext)
		if err = confAuditController.SetupWithManager(mgr); err != nil {
			return fmt.Errorf("unable to setup configauditreport reconciler: %w", err)
		}
		if err = (&controller.PluginsConfigReconciler{
			Logger:        ctrl.Log.WithName("reconciler").WithName("pluginsconfig"),
			Config:        operatorConfig,
			Client:        mgr.GetClient(),
			Plugin:        plugin,
			PluginContext: pluginContext,
		}).SetupWithManager(mgr); err != nil {
			return fmt.Errorf("unable to setup %T: %w", controller.PluginsConfigReconciler{}, err)
		}
	}

	if operatorConfig.CISKubernetesBenchmarkEnabled {
		rw := kubebench.NewReadWriter(mgr.GetClient())
		cisBenchmarkController := createCisBenchmarkController(operatorConfig, starboardConfig, mgr, logsReader, controller.KubeBenchJobName, rw.Write, rw.FindByOwner, limitChecker)
		if err = cisBenchmarkController.SetupWithManager(mgr); err != nil {
			return fmt.Errorf("unable to setup ciskubebenchreport reconciler: %w", err)
		}
	}

	if operatorConfig.NsaEnabled {
		// create conf audit controller
		plugin, pluginContext, err := createConfPlugin(buildInfo, operatorNamespace, operatorConfig, starboardConfig, mgr)
		if err != nil {
			return err
		}
		writer := nsa.NewReadWriter(mgr.GetClient())
		confAuditController := createConfAuditController(operatorConfig,
			starboardConfig,
			mgr,
			objectResolver,
			limitChecker,
			logsReader,
			plugin,
			controller.NsaConfigAuditJobName,
			writer.WriteConfig,
			writer.FindByOwner,
			pluginContext)
		if err = confAuditController.SetupWithManager(mgr); err != nil {
			return fmt.Errorf("unable to setup configauditreport reconciler: %w", err)
		}
		// create cis-benchmark controller
		cisBenchmarkController := createCisBenchmarkController(operatorConfig,
			starboardConfig,
			mgr,
			logsReader,
			controller.NsaKubeBenchJobName,
			writer.WriteInfra,
			writer.FindByOwner,
			limitChecker)

		if err := controller.NewNsaReportReconciler(
			operatorConfig,
			cisBenchmarkController,
			confAuditController,
			ctrl.Log.WithName("reconciler").WithName("nsareport"),
			starboardConfig, mgr.GetClient(), writer).SetupWithManager(mgr); err != nil {
			return fmt.Errorf("unable to setup nsa reconciler: %w", err)
		}
	}

	setupLog.Info("Starting controllers manager")
	if err := mgr.Start(ctx); err != nil {
		return fmt.Errorf("starting controllers manager: %w", err)
	}

	return nil
}

func createCisBenchmarkController(operatorConfig etc.Config,
	starboardConfig starboard.ConfigData,
	mgr ctrl.Manager,
	logsReader kube.LogsReader,
	jobNameFunc func(node *corev1.Node, f func(node *corev1.Node) string) string,
	writeFunc func(ctx context.Context, report v1alpha1.CISKubeBenchReport) error,
	findOwnerFunc func(ctx context.Context, node kube.ObjectRef) (interface{}, error),
	limitChecker controller.LimitChecker,
) *controller.CISKubeBenchReportReconciler {
	return &controller.CISKubeBenchReportReconciler{
		Logger:        ctrl.Log.WithName("reconciler").WithName("ciskubebenchreport"),
		Config:        operatorConfig,
		ConfigData:    starboardConfig,
		Client:        mgr.GetClient(),
		LogsReader:    logsReader,
		LimitChecker:  limitChecker,
		ReadWriter:    kubebench.NewReadWriter(mgr.GetClient()),
		Plugin:        kubebench.NewKubeBenchPlugin(ext.NewSystemClock(), starboardConfig),
		JobNameFunc:   jobNameFunc,
		WriteFunc:     writeFunc,
		FindOwnerFunc: findOwnerFunc,
	}
}

func createConfAuditController(operatorConfig etc.Config,
	starboardConfig starboard.ConfigData,
	mgr ctrl.Manager,
	objectResolver kube.ObjectResolver,
	limitChecker controller.LimitChecker,
	logsReader kube.LogsReader,
	plugin configauditreport.Plugin,
	jobNameFunc func(obj client.Object) string,
	writeFunc func(ctx context.Context, report v1alpha1.ClusterConfigAuditReport) error,
	findOwnerFunc func(ctx context.Context, node kube.ObjectRef) (interface{}, error),
	pluginContext starboard.PluginContext) *controller.ConfigAuditReportReconciler {
	pluginController := &controller.ConfigAuditReportReconciler{
		Logger:         ctrl.Log.WithName("reconciler").WithName("configauditreport"),
		Config:         operatorConfig,
		ConfigData:     starboardConfig,
		Client:         mgr.GetClient(),
		ObjectResolver: objectResolver,
		LimitChecker:   limitChecker,
		LogsReader:     logsReader,
		Plugin:         plugin,
		PluginContext:  pluginContext,
		ReadWriter:     configauditreport.NewReadWriter(mgr.GetClient()),
		JobNameFunc:    jobNameFunc,
		FindOwnerFunc:  findOwnerFunc,
		WriteFunc:      writeFunc,
	}
	return pluginController
}
func createConfPlugin(buildInfo starboard.BuildInfo, operatorNamespace string, operatorConfig etc.Config, starboardConfig starboard.ConfigData, manager ctrl.Manager) (configauditreport.Plugin, starboard.PluginContext, error) {
	plugin, pluginContext, err := plugin.NewResolver().
		WithBuildInfo(buildInfo).
		WithNamespace(operatorNamespace).
		WithServiceAccountName(operatorConfig.ServiceAccount).
		WithConfig(starboardConfig).
		WithClient(manager.GetClient()).
		GetConfigAuditPlugin()
	err = plugin.Init(pluginContext)
	if err != nil {
		return plugin, pluginContext, fmt.Errorf("initializing %s plugin: %w", pluginContext.GetName(), err)
	}
	return plugin, pluginContext, err
}
