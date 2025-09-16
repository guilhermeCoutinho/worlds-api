package services

import (
	"context"

	"github.com/sirupsen/logrus"
)

func SafeGo(ctx context.Context, fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logrus.WithField("method", "SafeGo").Errorf("Recovered from panic: %v", r)
			}
		}()
		fn()
	}()
}
