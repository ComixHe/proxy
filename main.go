package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/ComixHe/proxy/proxyService"
	"github.com/ComixHe/proxy/proxyService/interceptor"
	"github.com/ComixHe/proxy/utils"
	"github.com/godbus/dbus/v5"
)

func main() {
	log := utils.GetLogger()
	inInter, err := interceptor.New(interceptor.WithLogger(log))
	if err != nil {
		panic(err)
	}

	conn, err := dbus.ConnectSystemBus(dbus.WithIncomingInterceptor(inInter.GetBusName))
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancelCause(context.Background())

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

		sig := <-sigCh
		cancel(&utils.ErrCancelBySignal{Signal: sig})
	}()

	server, err := proxyService.New(
		proxyService.WithLogger(log),
		proxyService.WithDBusConnection(conn),
		proxyService.WithContext(ctx))

	if err != nil {
		panic(err)
	}

	err = server.RunDBusService(ctx)
	if err == nil {
		return
	}

	log.Debugw(
		"Deepin-Proxy exited with error.",
		"error", err,
	)

	var cancelBySignal *utils.ErrCancelBySignal
	if errors.As(err, &cancelBySignal) {
		log.Infow("Signal received, exiting...",
			"signal", cancelBySignal.Signal,
		)
		err = nil
	}

	err = conn.Close()
	if err != nil {
		log.Warnw("error:", err)
	}
	// content,err := config.LoadConfig()

	// globalConf, err := config.New(config.WithLogger(log), config.WithContent(content))
	// if err != nil{
	// 	log.Warn(err)
	// 	return
	// }
	// globalConf.Configs = make(map[int]*config.UserConfig)
	// userConf := config.UserConfig{}
	// err = userConf.Check()
	// if err != nil {
	// 	log.Warn(err)
	// 	return
	// }
	// globalConf.Configs[1000] = &userConf
	// conf, err := json.MarshalIndent(globalConf, "", "\t")
	// if err != nil{
	// 	log.Warn(err)
	// 	return
	// }

	// err = globalConf.WriteToFile()
	// if err != nil{
	// 	log.Warn(err)
	// 	return
	// }
	// log.Debug(string(conf))

}
