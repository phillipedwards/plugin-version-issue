package main

import (
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi-vault/sdk/v5/go/vault"
	"github.com/pulumi/pulumi-vault/sdk/v5/go/vault/generic"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		cfg := config.New(ctx, "vault")

		vp, err := vault.NewProvider(ctx, "vault", &vault.ProviderArgs{
			Address: pulumi.String(cfg.Require("address")),
			Token:   cfg.RequireSecret("token"),
		})

		if err != nil {
			return err
		}

		j := `
		{
			"key": "value"
		}
		`
		_, err = generic.NewSecret(ctx, "secret", &generic.SecretArgs{
			Path:     pulumi.String("secret/test"),
			DataJson: pulumi.String(j),
		}, pulumi.Provider(vp))

		if err != nil {
			return err
		}

		rp, err := random.NewProvider(ctx, "random", &random.ProviderArgs{})
		if err != nil {
			return nil
		}

		_, err = random.NewRandomId(ctx, "id", &random.RandomIdArgs{
			ByteLength: pulumi.Int(4),
		},
			pulumi.Provider(rp),
		)

		if err != nil {
			return err
		}

		return nil
	})
}
