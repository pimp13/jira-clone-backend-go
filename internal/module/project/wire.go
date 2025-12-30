//go:build wireinject
// +build wireinject

package project

import "github.com/google/wire"

var Provider = wire.NewSet(
	NewProjectService,
	NewProjectController,
	NewProjectMiddleware,
)
