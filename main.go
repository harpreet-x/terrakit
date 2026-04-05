// Copyright (c) 2026 TerraKit. Licensed under BSL 1.1.

package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/harpreet-x/terrakit/internal/provider"
)

// version is overridden at release time via -ldflags.
var version string = "dev"

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "run the provider with debugger support (e.g. delve)")
	flag.Parse()

	opts := providerserver.ServeOpts{
		// Address must match the source address used in required_providers.
		Address: "registry.terraform.io/harpreet-x/terrakit",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.New(version), opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}
