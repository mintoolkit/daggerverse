# slim: Basic Dagger Module (mint "slim" command)

Dagger module to slim the target container image.

This is a fork of the original code from https://github.com/shykes/daggerverse/tree/main/slim

## Functions/Commands

* `slim` - slim the target image producing a slim version of the image as its output (see the `slim` minToolkit command for more details).

* `compare` - slim the target image and then create a temporary container with the original and minified images mounted to inspect the changes.


## `slim` Function/Command Flags

* `--container CONTAINER_IMAGE_NAME` - target container image (required parameter)
* `--mode` - the engine execution mode `docker` (containerized version using Docker service) or `native` (native Dagger execution mode) (`docker` by default, optional)
* `--probe-http true|false` - enable HTTP probing of the temporary container (enabled by default, optional)
* `--probe-http-exit-on-failure true|false` - fail minification if all http probes failed (optional)
* `--publish-exposed-ports true|false` - publish all exposed target container ports when minifying to allow other components to interact with the temporary container (optional)
* `--probe-http-ports` - a comma separated list of ports to http probe (a subset of all exposed ports) (optional)
* `--probe-exec SHELL_CMD_TO_RUN` - enable exec-based probing of the temporary container (optional)
* `--continue-after MODE_ENUM` - select the continue-after mode (optional)
* `--show-clogs true|false` - show temporary container logs (optional)
* `--slim-debug true|false` - enable debug output in minToolkit (optional)

## Optional (With*) Parameter Functions/Commands

* `with-include-path --val STRING` (`WithIncludePath`)
* `with-include-bin --val STRING` (`WithIncludeBin`)
* `with-include-exe --val STRING` (`WithIncludeExe`)
* `with-include-shell --val BOOL` (`WithIncludeShell`)
* `with-include-new --val BOOL` (`WithIncludeNew`)
* `with-include-zoneinfo --val BOOL` (`WithIncludeZoneinfo`)
* `with-preserve-path --val STRING` (`WithPreservePath`)
* `with-exclude-pattern --val STRING` (`WithExcludePattern`)
* `with-env --val STRING` (`WithEnv`)
* `with-sensor-ipc-mode --val STRING` (`WithSensorIpcMode`)
* `with-sensor-ipc-endpoint --val STRING` (`WithSensorIpcEndpoint`)
* `with-source-ptrace --val STRING` (`WithSourcePtrace`)
* `with-image-build-engine --val STRING` (`WithImageBuildEngine`)
* `with-image-build-arch --val STRING` (`WithImageBuildArch`)
* `with-exec-probe --val STRING` (`WithExecProbe`)
* `with-http-probe-cmd --val STRING` (`WithHttpProbeCmd`)
* `with-expose-port --val STRING` (`WithExposePort`)
* `with-publish-port --val STRING` (`WithPublishPort`)

## Demo Steps

### Remote Mode From Its External Location

Call the `slim` module function to minify the target image and expose its network port when the minified image is executed at end:

`dagger up -m github.com/mintoolkit/daggerverse/slim --port 8080:80 slim --container nginx:latest`

Minify the target image and save the minified image as a tar file:

`dagger download -m github.com/mintoolkit/daggerverse/slim --output ./nginx-slim.tar slim --container nginx:latest`

Load the minified image from the saved tar file:

`docker load -i ./nginx-slim.tar`

The `docker load` command will print `Loaded image ID: YOUR_IMAGE_HASH`

Tag the loaded image, so it's easier to use later:

`docker tag YOUR_IMAGE_HASH nginx-slim:latest`

Run the minified container image in your host environment:

`docker run -it --rm -p 8888:80 nginx-slim:latest`


Note that you can use a specific version of the module by specifying its commit (e.g., `github.com/mintoolkit/daggerverse/slim@05e2410ce0725ffd553d537dfdc9003f643a725a` instead of simply `github.com/mintoolkit/daggerverse/slim`)

### Local Mode From the Module Itself

Call the `slim` module function to minify the target image and expose its network port when the minified image is executed at end:

`dagger up --port 8080:80 slim --container nginx:latest`

Minify the target image and save the minified image as a tar file:

`dagger download --output ./nginx-slim.tar slim --container nginx:latest`

Load the minified image from the saved tar file:

`docker load -i ./nginx-slim.tar`

The `docker load` command will print `Loaded image ID: YOUR_IMAGE_HASH`

Tag the loaded image, so it's easier to use later:

`docker tag YOUR_IMAGE_HASH nginx-slim:latest`

Run the minified container image in your host environment:

`docker run -it --rm -p 8888:80 nginx-slim:latest`


# Notes

The examples repo has many minification examples for different application stacks and different base images: https://github.com/mintoolkit/examples

See the main repo to get more information about the available flags (look for the `slim` command flags): https://github.com/mintoolkit/mint
