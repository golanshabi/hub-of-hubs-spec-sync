// Copyright (c) 2020 Red Hat, Inc.
// Copyright Contributors to the Open Cluster Management project

package controller

import (
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	configv1 "github.com/stolostron/hub-of-hubs-data-types/apis/config/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	ctrl "sigs.k8s.io/controller-runtime"
	client "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

const (
	hohSystemNamespace = "hoh-system"
)

func addHubOfHubsConfigController(mgr ctrl.Manager, databaseConnectionPool *pgxpool.Pool) error {
	err := ctrl.NewControllerManagedBy(mgr).
		For(&configv1.Config{}).
		WithEventFilter(predicate.NewPredicateFuncs(func(object client.Object) bool {
			return object.GetNamespace() == hohSystemNamespace
		})).
		Complete(&genericSpecToDBReconciler{
			client:                 mgr.GetClient(),
			databaseConnectionPool: databaseConnectionPool,
			log:                    ctrl.Log.WithName("hoh-config-spec-syncer"),
			tableName:              "configs",
			finalizerName:          "hub-of-hubs.open-cluster-management.io/hoh-config-cleanup",
			createInstance:         func() client.Object { return &configv1.Config{} },
			cleanStatus:            cleanConfigStatus,
			areEqual:               areConfigsEqual,
		})
	if err != nil {
		return fmt.Errorf("failed to add HubOfHubsConfig Controller to the manager: %w", err)
	}

	return nil
}

func cleanConfigStatus(instance client.Object) {
	config, ok := instance.(*configv1.Config)

	if !ok {
		panic("wrong instance passed to cleanConfigStatus: not configv1.Config")
	}

	config.Status = configv1.ConfigStatus{}
}

func areConfigsEqual(instance1, instance2 client.Object) bool {
	annotationMatch := equality.Semantic.DeepEqual(instance1.GetAnnotations(), instance2.GetAnnotations())

	config1, ok1 := instance1.(*configv1.Config)
	config2, ok2 := instance2.(*configv1.Config)
	specMatch := ok1 && ok2 && equality.Semantic.DeepEqual(config1.Spec, config2.Spec)

	return annotationMatch && specMatch
}
