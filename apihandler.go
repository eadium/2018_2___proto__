package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"proto-game-server/api"
	m "proto-game-server/models"
	"proto-game-server/router"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

const (
	cookieSessionIdName    = "sessionId"
	sessionCtxParamName    = "session"
	leadersOffsetParamName = "offset"
	leadersCountParamName  = "count"
)

//посредник между сетью и логикой апи
type ApiHandler struct {
	apiService      *api.ApiService
	corsAllowedHost string
	staticRoot      string
}

//избавиться от хардкода коннекта к бд
func NewApiHandler(settings *ServerConfig) *ApiHandler {
	service, err := api.NewApiService(settings.DbConnector, settings.DbConnectionString)

	if err != nil {
		panic(err)
	}

	return &ApiHandler{
		apiService:      service,
		corsAllowedHost: settings.CorsAllowedHost,
		staticRoot:      settings.StaticRoot,
	}
}

func WriteResponse(response *api.ApiResponse, ctx router.IContext) {
	data, err := json.Marshal(response.Response)
	if err != nil {
		ctx.Logger().Error(err)
		return
	}

	ctx.ContentType("application/json")
	ctx.StatusCode(response.Code)
	ctx.Write(data)

	ctx.Logger().Debugf("%s", response)
}

//регистрация
//обязательно нужно реализовать
func (h *ApiHandler) AddUser(ctx router.IContext) {
	user := new(m.User)
	ctx.ReadJSON(user)

	//можно потом добавить валидацию, но не сейчас

	//передаем уюзера из тела запроса в хранилище юзеров на регистрацию
	WriteResponse(h.apiService.Users.Add(user), ctx)
}

func (h *ApiHandler) DeleteUser(ctx router.IContext) {
	user := new(m.User)
	ctx.ReadJSON(user)

	WriteResponse(h.apiService.Users.Remove(user), ctx)
}

func (h *ApiHandler) UpdateUser(ctx router.IContext) {
	_, ok := ctx.CtxParam(sessionCtxParamName)
	if !ok {
		WriteResponse(&api.ApiResponse{
			Code:     http.StatusNotFound,
			Response: "Session not found"},
			ctx)
		return
	}

	user := new(m.User)
	ctx.ReadJSON(user)

	WriteResponse(h.apiService.Users.Update(user), ctx)
}

func (h *ApiHandler) GetUser(ctx router.IContext) {
	user := new(m.User)
	ctx.ReadJSON(user)

	params := ctx.UrlParams()

	WriteResponse(h.apiService.Users.Get(params["slug"]), ctx)
}

func (h *ApiHandler) Profile(ctx router.IContext) {
	data, ok := ctx.CtxParam(sessionCtxParamName)
	if !ok {
		WriteResponse(&api.ApiResponse{
			Code:     http.StatusInternalServerError,
			Response: "Session not found"}, ctx)
		return
	}

	session := data.(*m.Session)
	println(session.Token)
	println(session.User.Nickname)
	WriteResponse(&api.ApiResponse{Code: http.StatusOK, Response: session.User}, ctx)
}

func (h *ApiHandler) GetLeaders(ctx router.IContext) {
	params := ctx.UrlParams()

	offset, offsetErr := strconv.Atoi(params[leadersOffsetParamName])
	limit, limitErr := strconv.Atoi(params[leadersCountParamName])

	if offsetErr != nil || limitErr != nil {
		WriteResponse(&api.ApiResponse{
			http.StatusBadRequest, ""}, ctx)
	}

	WriteResponse(h.apiService.Scores.Get(offset, limit), ctx)
}

func (h *ApiHandler) GetSession(ctx router.IContext) {
	session, _ := ctx.CtxParam(sessionCtxParamName)
	WriteResponse(&api.ApiResponse{Code: http.StatusOK, Response: session}, ctx)
}

func (h *ApiHandler) Test(ctx router.IContext) {
	ctx.StatusCode(http.StatusOK)
}

func (h *ApiHandler) Logout(ctx router.IContext) {
	sessionid, ok := ctx.CtxParam(sessionCtxParamName)
	if !ok {
		WriteResponse(&api.ApiResponse{
			Code:     http.StatusNotFound,
			Response: "Session not found"},
			ctx)
		return
	}

	session := new(m.Session)
	session.Token = sessionid.(string)

	WriteResponse(h.apiService.Sessions.Remove(session), ctx)
}

func (h *ApiHandler) GetStatic(ctx router.IContext) {
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

func (h *ApiHandler) Upload(ctx router.IContext) {
	r := ctx.Request()

	// the FormFile function takes in the POST input id file
	file, header, err := r.FormFile("file")
	if err != nil {
		WriteResponse(&api.ApiResponse{Code: http.StatusBadRequest, Response: err}, ctx)
		return
	}

	defer file.Close()

	fileName := fmt.Sprintf("%v-%v", time.Now(), header.Filename)
	out, err := os.Create(fileName)
	if err != nil {
		WriteResponse(&api.ApiResponse{Code: http.StatusInternalServerError, Response: err}, ctx)
		return
	}

	defer out.Close()

	// write the content from POST to the file
	_, err = io.Copy(out, file)
	if err != nil {
		file := &m.File{fileName}
		WriteResponse(&api.ApiResponse{Code: http.StatusBadRequest, Response: file}, ctx)
		return
	}

	WriteResponse(&api.ApiResponse{Code: http.StatusInternalServerError, Response: err}, ctx)
}

//миддлварь для аутентификации
func (h *ApiHandler) AuthMiddleware(next router.HandlerFunc) router.HandlerFunc {
	return func(ctx router.IContext) {
		//тут должно быть получение id сессии из кукисов
		//попытка найти сессию в хранилище сессий и вызов след обработчика если все норм
		sessionCookie, err := ctx.GetCookie(cookieSessionIdName)
		if err != nil {
			WriteResponse(&api.ApiResponse{
				Code:     http.StatusNotFound,
				Response: "Session not found"},
				ctx)
			return
		}
		println(sessionCookie.Value)

		//поиск сессии по ИД в хранилище
		session, sessionExists := h.apiService.Sessions.GetById(sessionCookie.Value)
		if !sessionExists {
			WriteResponse(&api.ApiResponse{
				Code:     http.StatusUnauthorized,
				Response: "You are not authorized"},
				ctx)
			return
		}

		if session.TTL <= time.Now().Unix() {
			WriteResponse(&api.ApiResponse{
				Code:     http.StatusUnauthorized,
				Response: "Session timeout"},
				ctx)
			return
		}
		println(session.Token)
		ctx.AddCtxParam(sessionCtxParamName, session)
		next(ctx)
	}
}

//настройка cors'a
func (h *ApiHandler) CorsSetup(ctx router.IContext) {
	ctx.Header("Access-Control-Allow-Origin", h.corsAllowedHost)
	ctx.Header("Access-Control-Allow-Credentials", "true")
	ctx.Header("Access-Control-Allow-Headers", "Content-Type")
	ctx.Header("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, OPTIONS, PATCH")
	ctx.Header("Access-Control-Request-Methods", "GET, PUT, POST, DELETE, OPTIONS, PATCH")

}

func (h *ApiHandler) CorsEnableMiddleware(next router.HandlerFunc) router.HandlerFunc {
	return func(ctx router.IContext) {
		h.CorsSetup(ctx)
		next(ctx)
	}
}

func (h *ApiHandler) Authorize(ctx router.IContext) {
	user := new(m.User)
	ctx.ReadJSON(user)

	//хранилище создают сессию и возвращает нам ид сессии, который записывам в куки
	sessionId, ok := h.apiService.Sessions.Create(user)
	if !ok {
		ctx.Logger().Debugf("unauthorized request %s\n", ctx.RequestURI())
		WriteResponse(&api.ApiResponse{
			Code: http.StatusBadRequest,
			Response: &m.Error{Code: http.StatusBadRequest,
				Message: "Wrong login or password"}},
			ctx)
		return
	}

	//записываем ид сессии в куки
	//при каждом запросе, требующем аутнетификацию, будет браться данная кука и искаться в хранилище
	err := ctx.SetCookie(&http.Cookie{Name: cookieSessionIdName, Value: sessionId})
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("FAILED TO WRITE SESSION TO COOKIE %v", sessionId))
		ctx.StatusCode(http.StatusBadRequest)
	} else {
		ctx.Logger().Notice(sessionId)
		ctx.StatusCode(http.StatusOK)
	}
}

//function for testing cookie adding
func (h *ApiHandler) AddCookie(ctx router.IContext) {
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

func (h *ApiHandler) verifyDomain(ctx router.IContext) {
	message := "loaderio-3b73ee37ac50f8785f6e274aba668913"
	ctx.Write([]byte(message))
}
