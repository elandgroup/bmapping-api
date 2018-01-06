package main

import (
	"bmapping-api/models"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/asaskevich/govalidator"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	configutil "github.com/pangpanglabs/goutils/config"
	"github.com/pangpanglabs/goutils/echomiddleware"
	cache "github.com/patrickmn/go-cache"
	relaxmid "github.com/relax-space/go-kitt/echomiddleware"

	"bmapping-api/controllers"
	// "bmapping-api/models"
)

func main() {
	appEnv := flag.String("app-env", os.Getenv("APP_ENV"), "app env")
	connEnv := flag.String("BMAPPING_ENV", os.Getenv("BMAPPING_ENV"), "BMAPPING_ENV")
	flag.Parse()

	var c Config
	if err := configutil.Read(*appEnv, &c); err != nil {
		panic(err)
	}
	db, err := initDB(c.Database.Driver, *connEnv)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	localCache := cache.New(cache.NoExpiration, 0)
	configMap := map[string]interface{}{
		"|cache": localCache,
	}

	e := echo.New()
	//eland
	controllers.ElandStoreGroupApiController{}.Init(e.Group("/v3/eland/storegroups"))
	controllers.ElandStoreApiController{}.Init(e.Group("/v3/eland/stores"))
	//green
	controllers.GreenStoreApiController{}.Init(e.Group("/v3/green/stores"))

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	e.Use(middleware.RequestID())
	e.Use(echomiddleware.ContextLogger())
	e.Use(relaxmid.ContextConfig(configMap))
	e.Use(echomiddleware.ContextDB(c.Service, db, echomiddleware.KafkaConfig(c.Database.Logger.Kafka)))
	e.Use(echomiddleware.BehaviorLogger(c.Service, echomiddleware.KafkaConfig(c.BehaviorLog.Kafka)))
	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
		Skipper: func(c echo.Context) bool {
			ignore := []string{
				"/ping",
			}
			for _, i := range ignore {
				if strings.HasPrefix(c.Request().URL.Path, i) {
					return true
				}
			}

			return false
		},
	}))
	//e.Renderer = echotpl.New()
	e.Validator = &Validator{}

	e.Debug = c.Debug

	if err := e.Start(":" + c.HttpPort); err != nil {
		log.Println(err)
	}
}

func initDB(driver, connection string) (*xorm.Engine, error) {
	db, err := xorm.NewEngine(driver, connection)
	if err != nil {
		return nil, err
	}
	db.Sync(new(models.ElandStoreGroup),
		new(models.ElandStore),
		new(models.ElandMappingStoreIpay),
		new(models.GreenStoreGroup),
		new(models.GreenStore),
		new(models.GreenMappingStoreIpay),
		new(models.IpayType),
	)
	return db, nil
}

type Config struct {
	Database struct {
		Driver     string
		Connection string
		Logger     struct {
			Kafka echomiddleware.KafkaConfig
		}
	}
	BehaviorLog struct {
		Kafka echomiddleware.KafkaConfig
	}
	Trace struct {
		Zipkin echomiddleware.ZipkinConfig
	}

	Debug    bool
	Service  string
	HttpPort string
}

type Validator struct{}

func (v *Validator) Validate(i interface{}) error {
	_, err := govalidator.ValidateStruct(i)
	return err
}
