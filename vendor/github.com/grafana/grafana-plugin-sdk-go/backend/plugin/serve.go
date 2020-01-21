package plugin

import (
	plugin "github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

//ServeOpts options for serving plugins.
type ServeOpts struct {
	DiagnosticsServer DiagnosticsServer
	CoreServer        CoreServer
	TransformServer   TransformServer

	// GRPCServer factory method for creating GRPC server.
	// If nil, the default one will be used.
	GRPCServer func(options []grpc.ServerOption) *grpc.Server
}

// Serve starts serving the plugin over gRPC.
func Serve(opts ServeOpts) error {
	versionedPlugins := make(map[int]plugin.PluginSet)
	pSet := make(plugin.PluginSet)

	if opts.DiagnosticsServer != nil {
		pSet["diagnostics"] = &DiagnosticsGRPCPlugin{
			DiagnosticsServer: opts.DiagnosticsServer,
		}
	}

	if opts.CoreServer != nil {
		pSet["backend"] = &CoreGRPCPlugin{
			CoreServer: opts.CoreServer,
		}
	}

	if opts.TransformServer != nil {
		pSet["transform"] = &TransformGRPCPlugin{
			TransformServer: opts.TransformServer,
		}
	}

	versionedPlugins[ProtocolVersion] = pSet

	if opts.GRPCServer == nil {
		opts.GRPCServer = plugin.DefaultGRPCServer
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig:  Handshake,
		VersionedPlugins: versionedPlugins,
		GRPCServer:       opts.GRPCServer,
	})

	return nil
}