package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"errors"

	"github.com/mcb/cf-service-resource/out"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <sources directory>\n", os.Args[0])
		os.Exit(1)
	}

	cloudFoundry := out.NewCloudFoundry()
	command := out.NewCommand(cloudFoundry)

	var request out.Request
	if err := json.NewDecoder(os.Stdin).Decode(&request); err != nil {
		fatal("reading request from stdin", err)
	}

	var err error
	if request.Params.Service == "" {
		err = errors.New("service")
	}
	if request.Params.Plan == "" {
		err = errors.New("plan")
	}
	if request.Params.InstanceName == "" {
		err = errors.New("instance_name")
	}
	if request.Params.ManifestPath == "" && request.Params.CurrentAppName == "" {
		err = errors.New("either manifest/current_app_name")
	}
	if err != nil {
		fatal("parameter required", err)
	}

	// read current app name from manifest
	if request.Params.ManifestPath != "" {
		// make it an absolute path
		request.Params.ManifestPath = filepath.Join(os.Args[1], request.Params.ManifestPath)

		manifestFiles, err := filepath.Glob(request.Params.ManifestPath)
		if err != nil {
			fatal("searching for manifest files", err)
		}

		if len(manifestFiles) != 1 {
			fatal("invalid manifest path", fmt.Errorf("found %d files instead of 1 at path: %s", len(manifestFiles), request.Params.ManifestPath))
		}

		manifest, err := out.NewManifest(manifestFiles[0])
		if err != nil {
			fatal("failed to load manifest file", err)
		}

		if manifest.Data["applications"] == nil {
			err := errors.New("applications required")
			fatal("invalid manifest file", err)
		}

		application, hasValue := manifest.Data["applications"].([]interface{})[0].(map[interface{}]interface{})
		if !hasValue {
			err := errors.New("structure")
			fatal("invalid manifest file", err)
		}
		if application["name"] == nil {
			err := errors.New("name required")
			fatal("invalid manifest file", err)
		}

		request.Params.CurrentAppName = application["name"].(string)
	}

	response, err := command.Run(request)
	if err != nil {
		fatal("running command", err)
	}

	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		fatal("writing response to stdout", err)
	}
}

func fatal(message string, err error) {
	fmt.Fprintf(os.Stderr, "error %s: %s\n", message, err)
	os.Exit(1)
}
