// A Dagger module to Slim application container image size
package main

import (
	"context"
	"fmt"
	"runtime"
)

const (
	//todo: multi-arch engine image
	engineImageARM = "index.docker.io/mintoolkit/mint-arm"
	engineImageAMD = "index.docker.io/mintoolkit/mint"
	archAMD64      = "amd64"
	archARM64      = "arm64"

	outputImageTag = "slim-output:latest"
	outputImageTar = "output.tar"

	flagDebug = "--debug"
	trueValue = "true"
	cmdSlim   = "slim"

	modeDocker = "docker"
	modeNative = "native"

	flagShowClogs              = "--show-clogs"
	flagHttpProbe              = "--http-probe"
	flagHttpProbeCmd           = "--http-probe-cmd"
	flagHttpProbePorts         = "--http-probe-ports"
	flagHttpProbeExitOnFailure = "--http-probe-exit-on-failure"

	flagPublishPort         = "--publish-port"
	flagPublishExposedPorts = "--publish-exposed-ports"

	flagExecProbe = "--exec"

	flagIncludePath     = "--include-path"
	flagIncludeBin      = "--include-bin"
	flagIncludeExe      = "--include-exe"
	flagIncludeShell    = "--include-shell"
	flagIncludeNew      = "--include-new"
	flagIncludeZoneInfo = "--include-zoneinfo"
	flagPreservePath    = "--preserve-path"
	flagExcludePattern  = "--exclude-pattern"
	flagEnv             = "--env"
	flagExpose          = "--expose"
	flagContinueAfter   = "--continue-after"

	flagSensorIPCMode     = "--sensor-ipc-mode"
	flagSensorIPCEndpoint = "--sensor-ipc-endpoint"

	flagRTASourcePT      = "--rta-source-ptrace"
	flagImageBuildEngine = "--image-build-engine"
	flagImageBuildArch   = "--image-build-arch"
)

type Slim struct {
	includePaths      []string
	includeBins       []string
	includeExes       []string
	includeShell      *bool
	includeNew        *bool
	includeZoneinfo   *bool
	preservePaths     []string
	excludePatterns   []string
	envVars           []string
	sensorIPCMode     string
	sensorIPCEndpoint string
	rtaSourcePT       *bool
	imageBuildEngine  string
	imageBuildArch    string
	execProbes        []string
	httpProbeCmds     []string
	exposePorts       []string
	publishPorts      []string
}

func (s *Slim) Slim(
	ctx context.Context,
	container *Container,
	// Execution mode to use
	// +optional
	// +default="docker"
	mode string,
	// Enable running HTTP probes against the temporary container (test to false to disable)
	// +optional
	// +default=true
	probeHttp bool,
	// Probe HTTP - exit when all HTTP probe commands fail
	// +optional
	// +default=true
	probeHttpExitOnFailure bool,
	// Probe HTTP - comma separated subset of ports to probe
	// +optional
	probeHttpPorts string,
	// Map all exposed ports to the same host ports analyzing image at runtime
	// +optional
	// +default=true
	publishExposedPorts bool,
	// Select when to start processing the collected telemetry - enter | signal | probe | exec | timeout-number-in-seconds | container.probe (can combine probe and exec like this: probe&exe)
	// +optional
	continueAfter string,
	// Show container logs from the container used to perform dynamic inspection
	// +optional
	// +default=false
	showClogs bool,
	// Show debugging information
	// +optional
	// +default=false
	slimDebug bool,
) (*Container, error) {
	switch mode {
	case modeDocker, modeNative:
	default:
		mode = modeDocker
	}

	if mode != modeDocker {
		return nil, fmt.Errorf("unsupported mode - %s", mode)
	}

	// Start an ephemeral dockerd
	dockerd := dag.Docker().Engine()
	docker := dag.Docker().Cli(DockerCliOpts{
		Engine: dockerd,
	})

	// Load the input container into the dockerd
	imgRef, err := docker.Import(container).Ref(ctx)
	if err != nil {
		return container, err
	}

	var cargs []string
	if slimDebug {
		cargs = append(cargs, flagDebug)
	}

	cargs = append(cargs, cmdSlim)
	cargs = append(cargs, "--tag")
	cargs = append(cargs, outputImageTag)
	cargs = append(cargs, "--target")
	cargs = append(cargs, imgRef)

	if showClogs {
		cargs = append(cargs, flagShowClogs)
	}

	//pick up 'false' values too
	cargs = append(cargs, flagHttpProbe, fmt.Sprintf("%v", probeHttp))
	cargs = append(cargs, flagHttpProbeExitOnFailure, fmt.Sprintf("%v", probeHttpExitOnFailure))
	cargs = append(cargs, flagPublishExposedPorts, fmt.Sprintf("%v", publishExposedPorts))

	if probeHttpPorts != "" {
		cargs = append(cargs, flagHttpProbePorts, probeHttpPorts)
	}

	for _, val := range s.exposePorts {
		cargs = append(cargs, flagExpose, val)
	}

	for _, val := range s.publishPorts {
		cargs = append(cargs, flagPublishPort, val)
	}

	for _, val := range s.httpProbeCmds {
		cargs = append(cargs, flagHttpProbeCmd, val)
	}

	if len(s.execProbes) > 0 {
		//todo: support multiple exec probes (using the first one for now)
		cargs = append(cargs, flagExecProbe, s.execProbes[0])
	}

	for _, val := range s.includePaths {
		cargs = append(cargs, flagIncludePath, val)
	}

	for _, val := range s.includeBins {
		cargs = append(cargs, flagIncludeBin, val)
	}

	for _, val := range s.includeExes {
		cargs = append(cargs, flagIncludeExe, val)
	}

	for _, val := range s.preservePaths {
		cargs = append(cargs, flagPreservePath, val)
	}

	for _, val := range s.excludePatterns {
		cargs = append(cargs, flagExcludePattern, val)
	}

	for _, val := range s.envVars {
		cargs = append(cargs, flagEnv, val)
	}

	if continueAfter != "" {
		cargs = append(cargs, flagContinueAfter, continueAfter)
	}

	if s.sensorIPCMode != "" {
		cargs = append(cargs, flagSensorIPCMode, s.sensorIPCMode)
	}

	if s.sensorIPCEndpoint != "" {
		cargs = append(cargs, flagSensorIPCEndpoint, s.sensorIPCEndpoint)
	}

	if s.imageBuildArch != "" {
		cargs = append(cargs, flagImageBuildArch, s.imageBuildArch)
	}

	if s.imageBuildEngine != "" {
		cargs = append(cargs, flagImageBuildEngine, s.imageBuildEngine)
	}

	if s.rtaSourcePT != nil {
		cargs = append(cargs, flagRTASourcePT, fmt.Sprintf("%v", *s.rtaSourcePT))
	}

	if s.includeZoneinfo != nil {
		cargs = append(cargs, flagIncludeZoneInfo, fmt.Sprintf("%v", *s.includeZoneinfo))
	}

	if s.includeNew != nil {
		cargs = append(cargs, flagIncludeNew, fmt.Sprintf("%v", *s.includeNew))
	}

	if s.includeShell != nil {
		cargs = append(cargs, flagIncludeShell, fmt.Sprintf("%v", *s.includeShell))
	}

	//reuse the param to show the constructed command line:
	if slimDebug {
		fmt.Printf("Slim(Toolkit) params: %#v\n", cargs)
	}

	// Setup the slim container, attached to the dockerd
	slim := dag.
		Container().
		From(engineImage()).
		WithServiceBinding("dockerd", dockerd).
		WithEnvVariable("DOCKER_HOST", "tcp://dockerd:2375").
		WithExec(cargs)

	// Force execution of the slim command
	slim, err = slim.Sync(ctx)
	if err != nil {
		return container, err
	}

	// Extract the resulting image back into a container
	return docker.Image(DockerCliImageOpts{
		Repository: "slim-output",
		Tag:        "latest",
	}).Export(), nil
}

func (s *Slim) Compare(
	ctx context.Context,
	container *Container,
	// Execution mode to use
	// +optional
	// +default="docker"
	mode string,
	// Run HTTP probes against the temporary container
	// +optional
	// +default=true
	probeHttp bool,
	// Probe HTTP - exit on failure - TBD - add real desc
	// +optional
	// +default=true
	probeHttpExitOnFailure bool,
	// Probe HTTP - comma separated subset of ports to probe
	// +optional
	probeHttpPorts string,
	// Probe HTTP - publish exposed ports - TBD - add real desc
	// +optional
	// +default=true
	publishExposedPorts bool,
	// Continue after mode - TBD - add real desc
	// +optional
	continueAfter string,
	// Show temporary container logs - TBD - add real desc
	// +optional
	// +default=false
	showClogs bool,
	// Show debug messages - TBD - add real desc
	// +optional
	// +default=false
	slimDebug bool,
) (*Container, error) {
	slimmed, err := s.Slim(ctx,
		container,
		mode,
		probeHttp,
		probeHttpExitOnFailure,
		probeHttpPorts,
		publishExposedPorts,
		continueAfter,
		showClogs,
		slimDebug)
	if err != nil {
		return nil, err
	}

	debug := dag.
		Container().
		From("alpine").
		WithMountedDirectory("before", slimmed.Rootfs()).
		WithMountedDirectory("after", container.Rootfs())
	return debug, nil
}

// MORE OPTIONAL PARAMS

func (s *Slim) WithIncludePath(val string) *Slim {
	s.includePaths = append(s.includePaths, val)
	return s
}

func (s *Slim) WithIncludeBin(val string) *Slim {
	s.includeBins = append(s.includeBins, val)
	return s
}

func (s *Slim) WithIncludeExe(val string) *Slim {
	s.includeExes = append(s.includeExes, val)
	return s
}

func (s *Slim) WithIncludeShell(val bool) *Slim {
	s.includeShell = &val
	return s
}

func (s *Slim) WithIncludeNew(val bool) *Slim {
	s.includeNew = &val
	return s
}

func (s *Slim) WithIncludeZoneinfo(val bool) *Slim {
	s.includeZoneinfo = &val
	return s
}

func (s *Slim) WithPreservePath(val string) *Slim {
	s.preservePaths = append(s.preservePaths, val)
	return s
}

func (s *Slim) WithExcludePattern(val string) *Slim {
	s.excludePatterns = append(s.excludePatterns, val)
	return s
}

func (s *Slim) WithEnv(val string) *Slim {
	s.envVars = append(s.envVars, val)
	return s
}

func (s *Slim) WithSensorIpcMode(val string) *Slim {
	s.sensorIPCMode = val
	return s
}

func (s *Slim) WithSensorIpcEndpoint(val string) *Slim {
	s.sensorIPCEndpoint = val
	return s
}

func (s *Slim) WithSourcePtrace(val bool) *Slim {
	s.rtaSourcePT = &val
	return s
}

func (s *Slim) WithImageBuildEngine(val string) *Slim {
	s.imageBuildEngine = val
	return s
}

func (s *Slim) WithImageBuildArch(val string) *Slim {
	s.imageBuildArch = val
	return s
}

func (s *Slim) WithExecProbe(val string) *Slim {
	s.execProbes = append(s.execProbes, val)
	return s
}

func (s *Slim) WithHttpProbeCmd(val string) *Slim {
	s.httpProbeCmds = append(s.httpProbeCmds, val)
	return s
}

func (s *Slim) WithExposePort(val string) *Slim {
	s.exposePorts = append(s.exposePorts, val)
	return s
}

func (s *Slim) WithPublishPort(val string) *Slim {
	s.publishPorts = append(s.publishPorts, val)
	return s
}

// SUPPORTING FUNCTIONS:

func engineImage() string {
	switch runtime.GOARCH {
	case archAMD64:
		return engineImageAMD
	case archARM64:
		return engineImageARM
	default:
		return "" //let it error :)
	}
}
