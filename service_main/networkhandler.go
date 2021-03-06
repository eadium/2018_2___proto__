package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"proto-game-server/api"
	"proto-game-server/game"
	"proto-game-server/metrics"
	"proto-game-server/router"
	"strconv"
	"time"

	m "proto-game-server/models"

	ws "github.com/gorilla/websocket"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	cookieSessionIdName    = "sessionId"
	sessionCtxParamName    = "session"
	leadersOffsetParamName = "offset"
	leadersCountParamName  = "limit"
)

var upgrader = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

//посредник между сетью и логикой апи
type NetworkHandler struct {
	apiService      *api.ApiService
	game            *game.Game
	prof            *metrics.Profiler
	corsAllowedHost []string
	staticRoot      string
}

func NewNetworkHandler(settings *api.ServerConfig, logger router.ILogger) *NetworkHandler {
	service, err := api.NewApiService(settings)

	if err != nil {
		panic(err)
	}

	game := game.NewGame(logger)
	go game.Start()

	prof := metrics.NewProfiler()

	return &NetworkHandler{
		corsAllowedHost: settings.CorsAllowedHost,
		apiService:      service,
		staticRoot:      settings.StaticRoot,
		game:            game,
		prof:            prof,
	}
}

func (h *NetworkHandler) WriteResponse(response *api.ApiResponse, ctx router.IContext) {
	data, err := response.MarshalJSON()
	if err != nil {
		ctx.Logger().Error(err)
		return
	}
	h.prof.HitsStats.WithLabelValues(strconv.Itoa(response.Code)).Add(1)
	ctx.ContentType("application/json")
	ctx.StatusCode(response.Code)
	ctx.Write(data)

	ctx.Logger().Debugf("%s", response)
}

func (h *NetworkHandler) Authorize(ctx router.IContext) {
	user := new(m.User)
	ctx.ReadEasyJSON(user)

	//хранилище создают сессию и возвращает нам ид сессии, который записывам в куки
	serviceContext := context.Background()
	session, err := h.apiService.Sessions.Auth(serviceContext, user)
	if err != nil || session == nil {
		ctx.Logger().Debugf(err.Error())
		h.WriteResponse(&api.ApiResponse{
			Code: http.StatusBadRequest,
			Response: &m.Error{Code: http.StatusBadRequest,
				Message: "Wrong login or password"}},
			ctx)
		return
	}

	//записываем ид сессии в куки
	//при каждом запросе, требующем аутнетификацию, будет браться данная кука и искаться в хранилище
	err = ctx.SetCookie(&http.Cookie{Name: cookieSessionIdName, Value: session.Id})
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("FAILED TO WRITE SESSION TO COOKIE %v", session.Id))
		ctx.StatusCode(http.StatusBadRequest)
	} else {
		ctx.Logger().Notice(session.Id)
		ctx.StatusCode(http.StatusOK)
	}
}

//миддлварь для аутентификации
func (h *NetworkHandler) AuthMiddleware(next router.HandlerFunc) router.HandlerFunc {
	return func(ctx router.IContext) {
		//тут должно быть получение id сессии из кукисов
		//попытка найти сессию в хранилище сессий и вызов след обработчика если все норм
		sessionCookie, err := ctx.GetCookie(cookieSessionIdName)
		if err != nil {
			h.WriteResponse(&api.ApiResponse{
				Code:     http.StatusNotFound,
				Response: "Session not found"},
				ctx)
			return
		}

		//поиск сессии по ИД в хранилище
		sessionId := &m.SessionId{}
		sessionId.Id = sessionCookie.Value

		serviceContext := context.Background()
		session, err := h.apiService.Sessions.Check(serviceContext, sessionId)

		if err != nil {
			h.WriteResponse(&api.ApiResponse{
				Code:     http.StatusUnauthorized,
				Response: "You are not authorized"},
				ctx)
			return
		}

		if !session.IsAlive() {
			h.WriteResponse(&api.ApiResponse{
				Code:     http.StatusUnauthorized,
				Response: "Session timeout"},
				ctx)
			return
		}

		ctx.AddCtxParam(sessionCtxParamName, session)
		next(ctx)
	}
}

func (h *NetworkHandler) CorsSetup(ctx router.IContext) {
	origin := ctx.GetHeader("Origin")
	allowedHost := ScanSlice(h.corsAllowedHost, origin)
	ctx.Header("Access-Control-Allow-Origin", allowedHost)
	ctx.Header("Access-Control-Allow-Credentials", "true")
	ctx.Header("Access-Control-Allow-Headers", "Content-Type")
	ctx.Header("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, OPTIONS, PATCH")
}

func ScanSlice(s []string, key string) string {
	for _, a := range s {
		if a == key {
			return a
		}
	}
	return ""
}

func (h *NetworkHandler) CorsEnableMiddleware(next router.HandlerFunc) router.HandlerFunc {
	return func(ctx router.IContext) {
		h.CorsSetup(ctx)
		next(ctx)
	}
}

func (h *NetworkHandler) verifyDomain(ctx router.IContext) {
	message := "loaderio-3b73ee37ac50f8785f6e274aba668913"
	ctx.Write([]byte(message))
}

func (h *NetworkHandler) AddUser(ctx router.IContext) {

	user := new(m.User)
	ctx.ReadEasyJSON(user)

	//можно потом добавить валидацию, но не сейчас

	//передаем уюзера из тела запроса в хранилище юзеров на регистрацию
	h.WriteResponse(h.apiService.Users.Add(user), ctx)
}

func (h *NetworkHandler) DeleteUser(ctx router.IContext) {

	user := new(m.User)
	ctx.ReadEasyJSON(user)

	h.WriteResponse(h.apiService.Users.Remove(user), ctx)
}

func (h *NetworkHandler) UpdateUser(ctx router.IContext) {

	_, ok := ctx.CtxParam(sessionCtxParamName)
	if !ok {
		h.WriteResponse(&api.ApiResponse{
			Code:     http.StatusNotFound,
			Response: "Session not found"},
			ctx)
		return
	}

	user := new(m.User)
	ctx.ReadEasyJSON(user)

	h.WriteResponse(h.apiService.Users.Update(user), ctx)
}

func (h *NetworkHandler) GetUser(ctx router.IContext) {

	user := new(m.User)
	ctx.ReadEasyJSON(user)

	params := ctx.UrlParams()

	h.WriteResponse(h.apiService.Users.Get(params["slug"]), ctx)
}

func (h *NetworkHandler) Profile(ctx router.IContext) {
	data, ok := ctx.CtxParam(sessionCtxParamName)
	if !ok {
		h.WriteResponse(&api.ApiResponse{
			Code:     http.StatusInternalServerError,
			Response: "Session not found"}, ctx)
		return
	}

	session := data.(*m.Session)
	h.WriteResponse(&api.ApiResponse{Code: http.StatusOK, Response: session.User}, ctx)
}

func (h *NetworkHandler) GetLeaders(ctx router.IContext) {

	params := ctx.UrlParams()

	offset, offsetErr := strconv.Atoi(params[leadersOffsetParamName])
	limit, limitErr := strconv.Atoi(params[leadersCountParamName])

	if offsetErr != nil || limitErr != nil {
		h.WriteResponse(&api.ApiResponse{
			http.StatusBadRequest, ""}, ctx)
	}

	if offset < 0 {
		offset = 0
	}

	if limit < 1 {
		limit = 1
	}

	h.WriteResponse(h.apiService.Scores.Get(offset, limit), ctx)
}

func (h *NetworkHandler) GetSession(ctx router.IContext) {

	session, _ := ctx.CtxParam(sessionCtxParamName)
	h.WriteResponse(&api.ApiResponse{Code: http.StatusOK, Response: session}, ctx)
}

func (h *NetworkHandler) Test(ctx router.IContext) {
	// h.WriteResponse(&{Code: http.StatusOK, Message: "OK"}, ctx)
	ctx.StatusCode(http.StatusOK)
}

func (h *NetworkHandler) Logout(ctx router.IContext) {

	sessionid, ok := ctx.CtxParam(sessionCtxParamName)
	if !ok {
		h.WriteResponse(&api.ApiResponse{
			Code:     http.StatusNotFound,
			Response: "Session not found"},
			ctx)
		return
	}

	session := sessionid.(*m.Session)
	serviceContext := context.Background()
	_, err := h.apiService.Sessions.LogOut(serviceContext, session)

	if err != nil {
		h.WriteResponse(&api.ApiResponse{
			Code:     http.StatusNotFound,
			Response: "Session not found"},
			ctx)
		return
	}

	h.WriteResponse(&api.ApiResponse{
		Code:     http.StatusGone,
		Response: "session terminated"},
		ctx)
}

func (h *NetworkHandler) GetStatic(ctx router.IContext) {
	params := ctx.UrlParams()
	file := fmt.Sprintf("%v/%v", h.staticRoot, params["file"])

	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		ctx.Logger().Error(err)
		ctx.StatusCode(http.StatusNotFound)
		return
	}

	ctx.StatusCode(http.StatusOK)
	ctx.ContentType("image/png")
	ctx.Write(bytes)
}

func (h *NetworkHandler) Upload(ctx router.IContext) {

	r := ctx.Request()

	// the FormFile function takes in the POST input id file
	file, header, err := r.FormFile("file")
	if err != nil {
		h.WriteResponse(&api.ApiResponse{Code: http.StatusBadRequest, Response: err}, ctx)
		return
	}

	defer file.Close()

	fileName := fmt.Sprintf("%v-%v", time.Now(), header.Filename)
	out, err := os.Create(fileName)
	if err != nil {
		h.WriteResponse(&api.ApiResponse{Code: http.StatusInternalServerError, Response: err}, ctx)
		return
	}

	defer out.Close()

	// write the content from POST to the file
	_, err = io.Copy(out, file)
	if err != nil {
		file := &m.File{fileName}
		h.WriteResponse(&api.ApiResponse{Code: http.StatusBadRequest, Response: file}, ctx)
		return
	}

	h.WriteResponse(&api.ApiResponse{Code: http.StatusInternalServerError, Response: err}, ctx)
}

func (h *NetworkHandler) Panic(ctx router.IContext) {
	panic("panic")
}

//function for testing cookie adding
func (h *NetworkHandler) AddCookie(ctx router.IContext) {
	//записываем ид сессии в куки
	//при каждом запросе, требующем аутнетификацию, будет браться данная кука и искаться в хранилище
	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := &http.Cookie{
		Name:    "csrftoken",
		Value:   "abcd",
		Expires: expiration,
		Path:    "/"}

	err := ctx.SetCookie(cookie)
	if err != nil {
		ctx.Logger().Critical(err)
	}

	ctx.StatusCode(http.StatusOK)
	ctx.Write([]byte("COOKIE"))
}

func (h *NetworkHandler) ConnectPlayer(ctx router.IContext) {

	w := ctx.Writer()
	r := ctx.Request()

	sessionData, _ := ctx.CtxParam(sessionCtxParamName)
	session := sessionData.(*m.Session)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.WriteResponse(&api.ApiResponse{Code: http.StatusBadRequest, Response: err}, ctx)
		return
	}

	player := game.NewPlayer(session, conn)
	h.game.AddPlayer(player)
}

func (h *NetworkHandler) Metrics(ctx router.IContext) {
	r := ctx.Request()
	w := ctx.Writer()

	hand := promhttp.Handler()
	hand.ServeHTTP(w, r)
}
