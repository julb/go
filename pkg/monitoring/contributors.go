package monitoring

import "github.com/julb/go/pkg/build"

// Generic contributor type.
type Contributor map[string]interface{}

var (
	healthContributor  = &Contributor{}
	infoContributor    = &Contributor{}
	runtimeContributor = &Contributor{}
)

func init() {
	// health contributor should contain status of system by default
	GetHealthContributor().WithFunc("status", func() interface{} {
		return GetSystemStatus()
	})

	// info contributor should contain build info by default
	GetInfoContributor().With("build", build.Info)
}

// Get health contributor
func GetHealthContributor() *Contributor {
	return healthContributor
}

// Get info contributor
func GetInfoContributor() *Contributor {
	return infoContributor
}

// Get runtime contributor
func GetRuntimeContributor() *Contributor {
	return runtimeContributor
}

// Add value to the health contributor
func (contributor *Contributor) With(key string, value interface{}) *Contributor {
	(*contributor)[key] = value
	return contributor
}

// Add map value to the health contributor.
func (contributor *Contributor) WithMap(valuesMap map[string]interface{}) *Contributor {
	for key, value := range valuesMap {
		contributor.With(key, value)
	}
	return contributor
}

// Add func to the health contributor
func (contributor *Contributor) WithFunc(key string, fn func() interface{}) *Contributor {
	return contributor.With(key, fn)
}

// Add func to the health contributor
func (contributor *Contributor) ToMap() map[string]interface{} {
	var evaluate = make(map[string]interface{})
	for key, value := range *contributor {
		switch typedValue := value.(type) {
		case func() interface{}:
			evaluate[key] = typedValue()
		default:
			evaluate[key] = typedValue
		}
	}
	return evaluate
}
