package server

import "regexp"

var reHex = regexp.MustCompile(`^[0-9A-Fa-f]+$`)
var reBase64 = regexp.MustCompile(`^[A-Za-z0-9+/]+={0,2}$`)

var BufferSize = 1024
