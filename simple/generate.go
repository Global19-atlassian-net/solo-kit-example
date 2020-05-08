package main

import (
	"log"

	"github.com/solo-io/skv2/codegen"
	"github.com/solo-io/skv2/codegen/model"
	"github.com/solo-io/solo-kit/pkg/code-generator/sk_anyvendor"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// code generation for the simple example
func main() {

	cmd := &codegen.Command{
		// name of the operator/controller we are generating
		AppName: "simple-example",

		// we need to "vendor" protos in order to ensure compilation.
		// provide patterns here to match any project proto files
		AnyVendorConfig: sk_anyvendor.CreateDefaultMatchOptions([]string{
			"simple/api/*.proto",
		}),

		// define the API Groups this operator consumes
		Groups: []model.Group{
			{
				// kube GV info
				GroupVersion: schema.GroupVersion{
					Group:   "simple.skv2.solo.io",
					Version: "v1alpha1",
				},

				// go mod info
				Module: "github.com/solo-io/solo-kit-example",

				// the Resources "Kinds" defined for this GroupVersion
				Resources: []model.Resource{
					{
						// name of the kind
						Kind: "Circle",

						// the name of the Go Type that will be used for the CRD's Spec field
						Spec: model.Field{
							Type: model.Type{Name: "CircleSpec",
								GoPackage: "", /*Provide GoPackage to import the Type from another package*/
							},
						},

						// the name of the Go type that will be used for the CRD's Status field (optional)
						Status: &model.Field{
							Type: model.Type{Name: "CircleStatus"},
						},
					},

					// another example resource
					{
						Kind: "Square",
						Spec: model.Field{
							Type: model.Type{Name: "SquareSpec"},
						},
						Status: &model.Field{
							Type: model.Type{Name: "SquareStatus"},
						},
					},
				},

				// we usually define our Spec and Status types as protobufs. solo kit can compile
				// these for us.
				RenderProtos: true,
				// generate manifests for installing CRDs to kubernetes
				RenderManifests: true,
				// generate the Kube Go Structs themselves. disable when generating clients for external types
				RenderTypes: true,
				// render a strongly-typed API clients to work with the types.
				RenderClients: true,
				// render "Controller" code for subscribing to and reconciling the types.
				RenderController: true,

				// output all generated code to this directory
				ApiRoot: "simple/pkg/api",
			},
		},

		// generate CRD manifests in this directory
		ManifestRoot: "simple/install/helm/",
	}

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}

}
