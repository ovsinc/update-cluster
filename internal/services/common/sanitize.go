package common

import (
	"github.com/microcosm-cc/bluemonday"
)

var Policy = bluemonday.UGCPolicy()
