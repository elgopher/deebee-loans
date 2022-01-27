// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package database

import "github.com/elgopher/yala/logger"

var log logger.Global

func SetLoggerAdapter(adapter logger.Adapter) {
	log.SetAdapter(adapter)
}
