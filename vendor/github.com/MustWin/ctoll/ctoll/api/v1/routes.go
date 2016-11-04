package v1

import "github.com/gorilla/mux"

const (
	RouteNameBase            = "base"
	RouteNameOrgs            = "orgs"
	RouteNameBillingPlans    = "billingplans"
	RouteNameBillingPlan     = "billingplan"
	RouteNameMeterEvents     = "meterevents"
	RouteNameMeterSessions   = "metersessions"
	RouteNameMeterImageNames = "imagenames"
	RouteNameMeterImageTags  = "imagetags"
	RouteNameMeterLabels     = "labels"
	RouteNameAPIKeys         = "apikeys"
	RouteNameAPIKey          = "apikey"
	RouteNameBillingModel    = "billingmodel"
	RouteNameDistribution    = "distribution"
	RouteNameCostSavings     = "costsavings"
)

func Router() *mux.Router {
	return RouterWithPrefix("")
}

func RouterWithPrefix(prefix string) *mux.Router {
	rootRouter := mux.NewRouter()
	router := rootRouter
	if prefix != "" {
		router = router.PathPrefix(prefix).Subrouter()
	}

	router.StrictSlash(true)
	for _, descriptor := range routeDescriptors {
		router.Path(descriptor.Path).Name(descriptor.Name)
	}

	return rootRouter
}
