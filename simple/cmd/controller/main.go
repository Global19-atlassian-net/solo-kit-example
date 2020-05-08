package main

import (
	"context"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	log2 "log"
	"math"
	"sigs.k8s.io/controller-runtime/pkg/log"
	zaputil "sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/solo-io/skv2/pkg/reconcile"
	"github.com/solo-io/solo-kit-example/simple/pkg/api/simple.skv2.solo.io/v1alpha1"
	"github.com/solo-io/solo-kit-example/simple/pkg/api/simple.skv2.solo.io/v1alpha1/controller"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func main() {
	setupLogging()
	if err := run(); err != nil {
		log2.Fatal(err)
	}
}

func run() error {
	ctx := context.Background() // initialize a context

	kubeCfg, err := config.GetConfig() // initialze a KubeConfig
	if err != nil {
		return err
	}

	// initialize a "Manager", which is a generic handle to interacting with Kubernetes.
	// managers drive the internal components of Solo-Kit V2
	mgr, err := manager.New(kubeCfg, manager.Options{})
	if err != nil {
		return err
	}

	// register our types with the manager's Scheme.
	// this is required to enable using the Clients and Controllers with our types.
	if err := v1alpha1.AddToScheme(mgr.GetScheme()); err != nil {
		return err
	}

	// initialize circle reconcile loop
	{

		// create Reconcile Loop for reconciling Circles.
		// a single Reconcile Loop can be used to register multiple reconcilers for the given type.
		loop := controller.NewCircleReconcileLoop("circles", mgr)

		// run our circle reconciler
		err := loop.RunCircleReconciler(&circleReconciler{
			circles: v1alpha1.NewCircleClient(mgr.GetClient()), // construct a client for working with Circles
		})
		if err != nil {
			return err
		}
	}

	// initialize square reconcile loop
	{
		loop := controller.NewSquareReconcileLoop("squares", mgr)

		err := loop.RunSquareReconciler(&squareReconciler{
			squares: v1alpha1.NewSquareClient(mgr.GetClient()),
		})
		if err != nil {
			return err
		}
	}

	// finally, start the manager (blocking)
	return mgr.Start(ctx.Done())
}

type squareReconciler struct {
	ctx     context.Context
	squares v1alpha1.SquareClient
}

func (c *squareReconciler) ReconcileSquare(obj *v1alpha1.Square) (reconcile.Result, error) {
	// perform business logic
	area := float32(math.Pow(float64(obj.Spec.Width), 2))

	if obj.Status.Area != area {
		// report the result to the status
		obj.Status.Area = area

		err := c.squares.UpdateSquareStatus(c.ctx, obj)

		return reconcile.Result{}, err
	} else {

		// nothing to do, the object is up-to-date
		return reconcile.Result{}, nil
	}
}

type circleReconciler struct {
	ctx     context.Context
	circles v1alpha1.CircleClient
}

func (c *circleReconciler) ReconcileCircle(obj *v1alpha1.Circle) (reconcile.Result, error) {
	// perform business logic
	area := float32(math.Pow(float64(obj.Spec.Radius), 2)) * math.Pi

	if obj.Status.Area != area {
		// report the result to the status
		obj.Status.Area = area

		err := c.circles.UpdateCircleStatus(c.ctx, obj)

		return reconcile.Result{}, err
	} else {

		// nothing to do, the object is up-to-date
		return reconcile.Result{}, nil
	}
}

func setupLogging() {
	logconfig := zap.NewDevelopmentConfig()

		logconfig.Level.SetLevel(zap.DebugLevel)

	logger := zaputil.NewRaw(
		zaputil.UseDevMode(true),
	)
	zap.ReplaceGlobals(logger)
	log.SetLogger(zapr.NewLogger(logger))
}