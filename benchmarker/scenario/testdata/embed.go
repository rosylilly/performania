package testdata

import "embed"

//go:embed icons/*
var IconFiles embed.FS

//go:embed covers/*
var CoverFiles embed.FS
