package main

import (
	"fmt"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Deploy the bitnami/wordpress chart.

		wordpress, err := helm.NewRelease(ctx, "wpdev", &helm.ReleaseArgs{
			Version: pulumi.String("13.0.6"),
			Chart:   pulumi.String("wordpress"),
			Values:    pulumi.Map{"service": pulumi.StringMap{"type": pulumi.String("ClusterIP")}},
			RepositoryOpts: &helm.RepositoryOptsArgs{
				Repo: pulumi.String("https://charts.bitnami.com/bitnami"),
			},
		})

		service := pulumi.All(wordpress.Status.Namespace(), wordpress.Status.Name()).ApplyT(func(r interface{})(interface{}, error){

			arr := r.([]interface{})
			namespace := arr[0].(*string)
			name := arr[1].(*string)
			svc, err := corev1.GetService(ctx, "svc", pulumi.ID(fmt.Sprintf("%s/%s-wordpress", *namespace, *name)), nil)
			if err != nil {
				return "", nil
			}

			retval := []pulumi.StringPtrOutput {
				svc.Metadata.Name(),
				svc.Spec.ClusterIP(),

			}

			return retval, nil


		})

		//myservice := service.(*corev1.Service)
		ctx.Export("frontendIp", service)

		if err != nil {
			return err
		}

		return nil
	})
}