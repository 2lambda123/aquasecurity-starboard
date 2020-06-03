// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/aquasecurity/starboard/pkg/apis/aquasecurity/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// KubeHunterReportLister helps list KubeHunterReports.
// All objects returned here must be treated as read-only.
type KubeHunterReportLister interface {
	// List lists all KubeHunterReports in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.KubeHunterReport, err error)
	// Get retrieves the KubeHunterReport from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.KubeHunterReport, error)
	KubeHunterReportListerExpansion
}

// kubeHunterReportLister implements the KubeHunterReportLister interface.
type kubeHunterReportLister struct {
	indexer cache.Indexer
}

// NewKubeHunterReportLister returns a new KubeHunterReportLister.
func NewKubeHunterReportLister(indexer cache.Indexer) KubeHunterReportLister {
	return &kubeHunterReportLister{indexer: indexer}
}

// List lists all KubeHunterReports in the indexer.
func (s *kubeHunterReportLister) List(selector labels.Selector) (ret []*v1alpha1.KubeHunterReport, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.KubeHunterReport))
	})
	return ret, err
}

// Get retrieves the KubeHunterReport from the index for a given name.
func (s *kubeHunterReportLister) Get(name string) (*v1alpha1.KubeHunterReport, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("kubehunterreport"), name)
	}
	return obj.(*v1alpha1.KubeHunterReport), nil
}
