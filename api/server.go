package api

import (
	"encoding/base64"
	"strings"
	"time"

	"jarvis-trading-bot/analyzer"
	"jarvis-trading-bot/consts"
	"jarvis-trading-bot/structs"
	"jarvis-trading-bot/utils"
	"jarvis-trading-bot/utils/log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
)

type ApiServer struct {
	AppApi *fiber.App
}

func (a *ApiServer) Init() {
	a.AppApi = fiber.New()

	a.AppApi.Use(
		logger.New(), // add Logger middleware
	)

	a.AppApi.Post("/login", a.login)
	a.AppApi.Post("/create-admin", a.createAdminUser)
	a.AppApi.Post("/echo", a.echo)

	a.AppApi.Use(jwtware.New(jwtware.Config{
		SigningKey:   []byte(utils.GetStringConfig(consts.JwtSecretKey)),
		ErrorHandler: a.jwtError,
	}))

	a.InitValidateService()
	a.InitUserService()
	a.InitParametersService()
	a.InitBotParametersService()
	a.InitTradingStatusService()
	a.InitOnlyCalculateAnalyzerService()

	if err := a.AppApi.Listen(":3001"); err != nil {
		log.ErrorLogger.Fatal("Api:", err)
	} else {
		log.InfoLogger.Println("Server API has started")
	}
}

func (a *ApiServer) InitValidateService() {
	a.AppApi.Post("/validate", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})
}

func (a *ApiServer) InitParametersService() {
	a.AppApi.Get("/parameters", func(c *fiber.Ctx) error {
		idUser := a.getIntTokenField(c, "id")

		params := new(structs.Parameters)
		utils.DB.Where("user_id = ?", idUser).First(&params)

		if utils.DB.Error != nil {
			return c.Status(500).SendString(utils.DB.Error.Error())
		}

		if params == nil {
			return c.SendStatus(404)
		}

		return c.JSON(params)
	})

	a.AppApi.Post("/parameters", func(c *fiber.Ctx) error {
		params := new(structs.Parameters)
		if err := c.BodyParser(params); err != nil {
			return c.Status(500).SendString(err.Error())
		}

		params.UserId = a.getIntTokenField(c, "id")
		utils.DB.Save(params)
		if utils.DB.Error != nil {
			return c.Status(500).SendString(utils.DB.Error.Error())
		}

		return c.JSON(params)
	})

	a.AppApi.Put("/parameters", func(c *fiber.Ctx) error {
		params := new(structs.Parameters)
		if err := c.BodyParser(params); err != nil {
			return c.Status(500).SendString(err.Error())
		}

		params.UserId = a.getIntTokenField(c, "id")
		utils.DB.Save(params)
		if utils.DB.Error != nil {
			return c.Status(500).SendString(utils.DB.Error.Error())
		}

		return c.JSON(params)
	})

	a.AppApi.Delete("/parameters", func(c *fiber.Ctx) error {
		idUser := a.getIntTokenField(c, "id")

		utils.DB.Delete(&structs.Parameters{}, idUser)
		if utils.DB.Error != nil {
			return c.Status(500).SendString(utils.DB.Error.Error())
		}
		return c.SendStatus(200)
	})
}

func (a *ApiServer) InitBotParametersService() {
	a.AppApi.Get("/bot-parameters", func(c *fiber.Ctx) error {
		params := make([]structs.BotParameters, 0)
		userId := a.getIntTokenField(c, "id")

		utils.DB.Where("user_id = ?", userId).Find(&params)

		if utils.DB.Error != nil {
			return c.Status(500).SendString(utils.DB.Error.Error())
		}

		if params == nil {
			return c.SendStatus(404)
		}

		return c.JSON(params)
	})

	a.AppApi.Get("/bot-parameters/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		userId := a.getIntTokenField(c, "id")

		params := new(structs.BotParameters)
		utils.DB.Where("id = ? and user_id = ?", id, userId).First(&params)

		if utils.DB.Error != nil {
			return c.Status(500).SendString(utils.DB.Error.Error())
		}

		if params == nil {
			return c.SendStatus(404)
		}

		return c.JSON(params)
	})

	a.AppApi.Post("/bot-parameters", func(c *fiber.Ctx) error {
		params := new(structs.BotParameters)
		if err := c.BodyParser(params); err != nil {
			return c.Status(500).SendString(err.Error())
		}

		params.UserId = a.getIntTokenField(c, "id")
		result := utils.DB.Create(params)
		if result.Error != nil {
			log.ErrorLogger.Println(result.Error.Error())
			return c.Status(500).SendString(result.Error.Error())
		}
		return c.JSON(params)
	})

	a.AppApi.Put("/bot-parameters", func(c *fiber.Ctx) error {
		params := new(structs.BotParameters)
		if err := c.BodyParser(params); err != nil {
			return c.Status(500).SendString(err.Error())
		}

		params.UserId = a.getIntTokenField(c, "id")
		utils.DB.Save(params)
		if utils.DB.Error != nil {
			return c.Status(500).SendString(utils.DB.Error.Error())
		}

		return c.JSON(params)
	})

	a.AppApi.Delete("/bot-parameters/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		utils.DB.Delete(&structs.BotParameters{}, id)
		if utils.DB.Error != nil {
			return c.Status(500).SendString(utils.DB.Error.Error())
		}
		return c.SendStatus(200)
	})
}

func (a *ApiServer) InitTradingStatusService() {
	a.AppApi.Get("/trading-status", func(c *fiber.Ctx) error {
		trades := make([]structs.TradingStatus, 0)
		utils.DB.Where("deleted_at is null").Find(&trades)

		if utils.DB.Error != nil {
			return c.Status(500).SendString(utils.DB.Error.Error())
		}

		if trades == nil {
			return c.SendStatus(404)
		}

		return c.JSON(trades)
	})
}

func (a *ApiServer) InitOnlyCalculateAnalyzerService() {
	a.AppApi.Get("/analysis/:idBot/candlesticks", func(c *fiber.Ctx) error {
		idBot := c.Params("idBot")
		start := c.Query("start")
		end := c.Query("end")

		a := new(analyzer.Analyzer)
		a.InitOnlyCalculate(idBot)
		k := a.CreateTradingStatusKeyByBotParamId(idBot)

		candles := a.GetCandlestickAnalysis(k, start, end)

		if len(candles) == 0 {
			return c.SendStatus(404)
		}
		return c.JSON(candles)
	})
}

func (a *ApiServer) InitUserService() {
	a.AppApi.Get("/user-data", func(c *fiber.Ctx) error {
		username := a.getTokenField(c, "username")
		log.InfoLogger.Println("Looking for user ", username)

		user := new(structs.User)
		utils.DB.Where("username = ?", username).First(&user)

		if utils.DB.Error != nil {
			return c.Status(500).SendString(utils.DB.Error.Error())
		}

		if user == nil {
			return c.SendStatus(404)
		}

		user.Password = "**********************"
		return c.JSON(user)
	})
}

func (a *ApiServer) login(c *fiber.Ctx) error {
	username := ""
	password := ""

	if c.Get("Authorization") != "" {
		authHeader := c.Get("Authorization")
		authParts := strings.Split(authHeader, " ")

		base64Text := make([]byte, base64.StdEncoding.DecodedLen(len(authParts[1])))
		base64.StdEncoding.Decode(base64Text, []byte(authParts[1]))
		userPass := string(base64Text[:])
		authParts = strings.Split(userPass, ":")

		username = authParts[0]
		password = authParts[1]
	} else {
		u := new(structs.User)
		if err := c.BodyParser(u); err != nil {
			return c.Status(500).SendString(err.Error())
		}

		username = u.Username
		password = utils.CreateHash(utils.GetStringConfig(consts.PasswordKey), u.Password)
	}

	user := new(structs.User)
	utils.DB.Where("username = ? AND password = ?", username, password).First(&user)
	if user == nil || len(user.Username) == 0 || len(user.Username) == 0 {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// Create the Claims
	claims := jwt.MapClaims{
		"id":       user.Id,
		"username": user.Username,
		"name":     user.Name,
		"email":    user.Email,
		"admin":    true,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(utils.GetStringConfig(consts.JwtSecretKey)))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t})
}

func (a *ApiServer) getTokenField(c *fiber.Ctx, field string) string {
	u := c.Locals("user").(*jwt.Token)
	claims := u.Claims.(jwt.MapClaims)
	return claims[field].(string)
}

func (a *ApiServer) getIntTokenField(c *fiber.Ctx, field string) uint {
	u := c.Locals("user").(*jwt.Token)
	claims := u.Claims.(jwt.MapClaims)
	v := claims[field].(float64)
	return uint(v)
}

func (a *ApiServer) jwtError(c *fiber.Ctx, err error) error {
	// Return status 401 and failed authentication error.
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Return status 401 and failed authentication error.
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": true,
		"msg":   err.Error(),
	})
}

func (a *ApiServer) createAdminUser(c *fiber.Ctx) error {
	adminUser := utils.GetStringConfig(consts.AdminUser)
	adminPass := utils.GetStringConfig(consts.AdminPass)
	adminEmail := utils.GetStringConfig(consts.AdminEmail)

	u := new(structs.User)
	utils.DB.Where("username = ?", adminUser).First(&u)
	if u == nil || len(u.Username) == 0 {
		adminPass = utils.CreateHash(utils.GetStringConfig(consts.PasswordKey), adminPass)

		user := new(structs.User)
		user.Name = "Admin User"
		user.Username = adminUser
		user.Password = adminPass
		user.Email = adminEmail
		user.IsAdmin = true

		utils.DB.Create(user)
		return c.Status(200).SendString("User Admin was created successfully")
	}
	return c.Status(200).SendString("User Admin was already created")
}

func (a *ApiServer) echo(c *fiber.Ctx) error {
	return c.Status(200).SendString("Hello! This is an Echo Test.")
}

func StartAndListen() {
	a := new(ApiServer)
	a.Init()
}
