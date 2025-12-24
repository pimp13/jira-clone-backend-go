//go:build wireinject
// +build wireinject

package workspace

import "github.com/google/wire"

var Provider = wire.NewSet(
	NewWorkspaceService,
	NewWorkspaceController,
	NewWorkspaceMiddleware,
)
