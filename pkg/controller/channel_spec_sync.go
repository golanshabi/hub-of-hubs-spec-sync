package controller

import (
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"k8s.io/apimachinery/pkg/api/equality"
	channelsv1 "open-cluster-management.io/multicloud-operators-channel/pkg/apis/apps/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func addChannelController(mgr ctrl.Manager, databaseConnectionPool *pgxpool.Pool) error {
	controlBuilder := ctrl.NewControllerManagedBy(mgr).For(&channelsv1.Channel{})
	controlBuilder = controlBuilder.WithEventFilter(generateNamespacePredicate())

	err := controlBuilder.Complete(&genericSpecToDBReconciler{
		client:                 mgr.GetClient(),
		databaseConnectionPool: databaseConnectionPool,
		log:                    ctrl.Log.WithName("channel-spec-syncer"),
		tableName:              "channels",
		finalizerName:          "hub-of-hubs.open-cluster-management.io/channel-cleanup",
		createInstance:         func() client.Object { return &channelsv1.Channel{} },
		cleanStatus:            cleanChannelStatus,
		areEqual:               areChannelsEqual,
	})
	if err != nil {
		return fmt.Errorf("failed to add channel controller to the manager: %w", err)
	}

	return nil
}

func cleanChannelStatus(instance client.Object) {
	channel, ok := instance.(*channelsv1.Channel)
	if !ok {
		panic("wrong instance passed to cleanSubscriptionStatus: not channelsv1.Channel")
	}

	channel.Status = channelsv1.ChannelStatus{}
}

func areChannelsEqual(instance1, instance2 client.Object) bool {
	annotationMatch := equality.Semantic.DeepEqual(instance1.GetAnnotations(), instance2.GetAnnotations())

	channel1, ok1 := instance1.(*channelsv1.Channel)
	channel2, ok2 := instance2.(*channelsv1.Channel)
	specMatch := ok1 && ok2 && equality.Semantic.DeepEqual(channel1.Spec, channel2.Spec)

	return annotationMatch && specMatch
}
