package utils

import (
	"github.com/aquasecurity/starboard/pkg/ext"
	"github.com/gorhill/cronexpr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"time"
)

// NextCronDuration check if next cron activation time has exceeded if so return true
// if activation time has not reached return false and remaining time
// in case it failed to parse cron expression return error
func NextCronDuration(cronString string, creationTime time.Time, clock ext.Clock) (time.Duration, error) {
	expr, err := cronexpr.Parse(cronString)
	if err != nil {
		return time.Duration(0), err
	}
	return timeToExpiration(expr.Next(creationTime), clock), nil
}

// DurationExceeded  check if duration is now meaning zero
func DurationExceeded(duration time.Duration) bool {
	return duration.Nanoseconds() <= 0
}

// timeToExpiration  return the duration between time to expiration
func timeToExpiration(expiresAt time.Time, clock ext.Clock) time.Duration {
	return expiresAt.Sub(clock.Now())
}

// IsTTLExpired check whether current time has exceeded creation time + ttl duration
func IsTTLExpired(ttl time.Duration, creationTime time.Time, clock ext.Clock) (bool, time.Duration) {
	durationToTTLExpiration := timeToExpiration(creationTime.Add(ttl), clock)
	return DurationExceeded(durationToTTLExpiration), durationToTTLExpiration
}

// TransformObjMetadata removes managed fields and the "kubectl.kubernetes.io/last-applied-configuration"
// annotation from Kubernetes objects. This ensures that these fields are not stored in the client-go cache.
func TransformObjMetadata(obj interface{}) (interface{}, error) {
	obj, err := cache.TransformStripManagedFields()(obj)
	if err != nil {
		return obj, err
	}

	if metaObj, ok := obj.(metav1.ObjectMetaAccessor); ok {
		annotations := metaObj.GetObjectMeta().GetAnnotations()
		if annotations != nil {
			delete(annotations, "kubectl.kubernetes.io/last-applied-configuration")
			metaObj.GetObjectMeta().SetAnnotations(annotations)
		}
	}
	return obj, nil
}
