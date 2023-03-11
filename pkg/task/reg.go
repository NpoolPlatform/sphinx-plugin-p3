package task

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/sphinx-plugin-P3/pkg/coins"
	"github.com/NpoolPlatform/sphinx-plugin-P3/pkg/coins/getter"
	"github.com/NpoolPlatform/sphinx-plugin-P3/pkg/env"
	"github.com/NpoolPlatform/sphinx-plugin-P3/pkg/log"
)

const (
	// unit seconds
	getTransactionsTimeout = 60 * time.Second
	// unit seconds
	updateTransactionsTimeout = 10 * time.Second
)

type tworker struct {
	interval time.Duration
	// unit second
	// timeout int
	handle func(name string, interval time.Duration)
}

var (
	ErrTaskAlreadyRegister = errors.New("task already register")

	tworkers = make(map[string]tworker)
)

// interval unit second
// generate random number [3, 6)
func register(taskName string, interval time.Duration, handle func(name string, interval time.Duration)) error {
	if _, ok := tworkers[taskName]; ok {
		return ErrTaskAlreadyRegister
	}

	tworkers[taskName] = tworker{
		interval: interval,
		handle:   handle,
	}

	return nil
}

func fatalf(prefix, template string, args ...interface{}) {
	logger.Sugar().Fatalf(fmt.Sprintf("%s %v", prefix, template), args...)
}

func errorf(prefix, template string, args ...interface{}) {
	logger.Sugar().Errorf(fmt.Sprintf("%s %v", prefix, template), args...)
}

func warnf(prefix, template string, args ...interface{}) {
	logger.Sugar().Warnf(fmt.Sprintf("%s %v", prefix, template), args...)
}

func infof(prefix, template string, args ...interface{}) {
	logger.Sugar().Infof(fmt.Sprintf("%s %v", prefix, template), args...)
}

func setUpTokens() error {
	coinNetwork, _coinType, err := env.CoinInfo()
	if err != nil {
		return err
	}
	coinType := coins.CoinStr2CoinType(coinNetwork, _coinType)
	infos := getter.GetTokenInfos(coinType)
	fn, err := getter.GetTokenNetHandler(coinType)
	if err != nil {
		for _, info := range infos {
			info.DisableRegiste = false
		}
		return nil
	}

	_infos := make([]*coins.TokenInfo, len(infos))
	idx := 0
	for _, v := range infos {
		_infos[idx] = v
		idx++
	}

	return fn(_infos)
}

func Run() {
	go func() {
		err := setUpTokens()
		if err != nil {
			panic(err)
		}
	}()

	for name, tf := range tworkers {
		time.Sleep(time.Millisecond * time.Duration(500+rand.Int63n(200)))
		log.Infof("run task: %v ", name)
		go tf.handle(name, tf.interval)
	}
}
