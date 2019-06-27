package util

import (
  "go.uber.org/zap"
  "log"
  "sync"
)

var once sync.Once
var logger *zap.Logger

func GetLogger() *zap.Logger {
  once.Do(func() {
    l, err := zap.NewDevelopment()
    if err != nil {
      log.Fatalf("Error starting logger: %s", err.Error())
    }
    logger = l
  })
  return logger
}

func GetSugaredLogger() *zap.SugaredLogger {
  return GetLogger().Sugar()
}