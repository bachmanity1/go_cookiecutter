//go:generate swagger generate spec
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	ct "pandita/api/controller"
	mw "pandita/api/middleware"
	conf "pandita/conf"
	_Repo "pandita/repository"
	"pandita/util"
	"runtime"
	"syscall"
	"time"

	"github.com/juju/errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	banner = "\n" +
		"                                   ___         ___               \n" +
		"                                  (   )  .-.  (   )              \n" +
		"   .-..     .---.   ___ .-.     .-.| |  ( __)  | |_       .---.  \n" +
		"  /    \\   / .-, \\ (   )   \\   /   \\ |  (''\") (   __)    / .-, \\ \n" +
		" ' .-,  ; (__) ; |  |  .-. .  |  .-. |   | |   | |      (__) ; | \n" +
		" | |  . |   .'`  |  | |  | |  | |  | |   | |   | | ___    .'`  | \n" +
		" | |  | |  / .'| |  | |  | |  | |  | |   | |   | |(   )  / .'| | \n" +
		" | |  | | | /  | |  | |  | |  | |  | |   | |   | | | |  | /  | | \n" +
		" | |  ' | ; |  ; |  | |  | |  | '  | |   | |   | ' | |  ; |  ; | \n" +
		" | `-'  ' ' `-'  |  | |  | |  ' `-'  /   | |   ' `-' ;  ' `-'  | \n" +
		" | \\__.'  `.__.'_. (___)(___)  `.__,'   (___)   `.__.   `.__.'_. \n" +
		" | |                                                             \n" +
		"(___)                                                            \n" +
		"%s\n" +
		" => Starting listen %s\n"
)

var (
	// BuildDate for Program BuildDate
	BuildDate string
	// Version for Program Version
	Version string
	svrInfo = fmt.Sprintf("pandita %s(%s)", Version, BuildDate)
	mlog    *util.MLogger
)

func init() {
	// use all cpu
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Printf("Init %s\n", svrInfo)
}

func main() {

	Pandita := conf.Pandita
	if Pandita.GetBool("version") {
		fmt.Printf("%s\n", svrInfo)
		os.Exit(0)
	}
	Pandita.SetProfile()
	mlog, _ = util.InitLog("main", Pandita.GetString("loglevel"))

	e := echoInit(Pandita)
	sc := sigInit(e)

	// Prepare Server
	db := _Repo.InitDB(Pandita)
	if err := ct.InitHandler(Pandita, e, db); err != nil {
		mlog.Errorw("InitHandler", "err", errors.Details(err))
		os.Exit(1)
	}

	if !prepareServer(Pandita, sc) {
		os.Exit(1)
	}

	startServer(Pandita, e)
}

func echoInit(pandita *conf.ViperConfig) (e *echo.Echo) {

	// Echo instance
	e = echo.New()

	// Middleware
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())

	// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.POST, echo.GET, echo.PUT, echo.DELETE},
	}))
	// Ping Check
	e.GET("/healthCheck", func(c echo.Context) error { return c.String(http.StatusAlreadyReported, "pandita API Alive!\n") })
	e.POST("/healthCheck", func(c echo.Context) error { return c.String(http.StatusAlreadyReported, "pandita API Alive!\n") })

	e.Use(mw.ZapLogger(mlog))
	e.HideBanner = true

	return e
}

func sigInit(e *echo.Echo) chan os.Signal {

	// Signal
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	go func() {
		sig := <-sc
		e.Logger.Error("Got signal", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := e.Shutdown(ctx); err != nil {
			e.Logger.Error(err)
		}
		signal.Stop(sc)
		close(sc)
	}()

	return sc
}

func startServer(pandita *conf.ViperConfig, e *echo.Echo) {
	// Start Server
	apiServer := fmt.Sprintf("0.0.0.0:%d", pandita.GetInt("port"))
	mlog.Infow("Starting server", "info", svrInfo, "listen", apiServer)
	fmt.Printf(banner, svrInfo, apiServer)
	if pandita.GetBool("debug_route") {
		data, err := json.MarshalIndent(e.Routes(), "", "  ")
		if err != nil {
			mlog.Errorw("Bad routes", "err", err)
			os.Exit(1)
		}
		fmt.Printf("%s\n", data)
	}

	if err := e.Start(apiServer); err != nil {
		mlog.Errorw("End server", "err", err)
	}
}
