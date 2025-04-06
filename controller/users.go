package controller

import (
	"net/http"
	"strconv"

	"go_effective/httputil"
	"go_effective/model"

	"github.com/gin-gonic/gin"
)

// ShowUser godoc
//
//	@Summary		Возвращает пользователя
//	@Description	Возвращает пользователя по ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	model.User
//	@Failure		400	{object}	httputil.HTTPError
//	@Failure		404	{object}	httputil.HTTPError
//	@Failure		500	{object}	httputil.HTTPError
//	@Router			/api/v1/users/{id} [get]
func (c *Controller) ShowUser(ctx *gin.Context) {
	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}
	user, err := model.UserOne(c.dbconn, aid)
	if err != nil {
		httputil.NewError(ctx, http.StatusNotFound, err)
		return
	}
	ctx.JSON(http.StatusOK, user)
}

// ListUsers godoc
//
//		@Summary		Список пользователей
//		@Description	Возвращает список пользователей фильтром по name и пагинацией offset limit
//		@Tags			users
//		@Accept			json
//		@Produce		json
//		@Param			name	query		string	false	"name search by name"	Format(email)
//	 @Param			offset	query		int	false	"offset"
//	 @Param			limit	query		int	false	"limit"
//		@Success		200	{array}		model.User
//		@Failure		400	{object}	httputil.HTTPError
//		@Failure		404	{object}	httputil.HTTPError
//		@Failure		500	{object}	httputil.HTTPError
//		@Router			/api/v1/users [get]
func (c *Controller) ListUsers(ctx *gin.Context) {
	name := ctx.Request.URL.Query().Get("name")
	offsetStr := ctx.Request.URL.Query().Get("offset")
	limitStr := ctx.Request.URL.Query().Get("limit")
	limit := 10 // Значение по умолчанию
	offset := 0 // Значение по умолчанию

	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			limit = parsedLimit
		}
	}

	if offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil {
			offset = parsedOffset
		}
	}
	users, err := model.GetUsersByName(c.dbconn, name, limit, offset)
	if err != nil {
		httputil.NewError(ctx, http.StatusNotFound, err)
		return
	}
	ctx.JSON(http.StatusOK, users)
}

// AddUser godoc
//
//	@Summary		Добавление пользователя
//	@Description	Добавляет пользователя в бд по JSON
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			user	body		model.AddUser	true	"Add user"
//	@Success		200		{object}	model.User
//	@Failure		400		{object}	httputil.HTTPError
//	@Failure		404		{object}	httputil.HTTPError
//	@Failure		500		{object}	httputil.HTTPError
//	@Router			/api/v1/users [post]
func (c *Controller) AddUser(ctx *gin.Context) {
	var addUser model.AddUser
	if err := ctx.ShouldBindJSON(&addUser); err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}
	if err := addUser.Validation(); err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	user := model.User{
		Name:       addUser.Name,
		Surname:    addUser.Surname,
		Patronimic: addUser.Patronimic,
	}
	user.AddData()
	lastID, err := user.Insert(c.dbconn)
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}
	user.ID = lastID
	ctx.JSON(http.StatusOK, user)
}

// UpdateUser godoc
//
//	@Summary		Изменение информации о пользователе
//	@Description	Изменяет информацию о пользователе в соответствии с JSON
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"User ID"
//	@Param			user	body		model.UpdateUser	true	"Update user"
//	@Success		200		{object}	model.User
//	@Failure		400		{object}	httputil.HTTPError
//	@Failure		404		{object}	httputil.HTTPError
//	@Failure		500		{object}	httputil.HTTPError
//	@Router			/api/v1/users/{id} [patch]
func (c *Controller) UpdateUser(ctx *gin.Context) {
	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}
	var updateUser model.UpdateUser
	if err := ctx.ShouldBindJSON(&updateUser); err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}
	user := model.User{
		ID:         aid,
		Name:       updateUser.Name,
		Surname:    updateUser.Surname,
		Patronimic: updateUser.Patronimic,
	}
	editedUser, err := user.Update(c.dbconn)
	if err != nil {
		httputil.NewError(ctx, http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusOK, editedUser)
}

// DeleteUser godoc
//
//	@Summary		Удаление пользователя
//	@Description	Удаляет пользователя из бд по id
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"user ID"	Format(int64)
//	@Success		204	{object}	model.User
//	@Failure		400	{object}	httputil.HTTPError
//	@Failure		404	{object}	httputil.HTTPError
//	@Failure		500	{object}	httputil.HTTPError
//	@Router			/api/v1/users/{id} [delete]
func (c *Controller) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}
	err = model.Delete(c.dbconn, aid)
	if err != nil {
		httputil.NewError(ctx, http.StatusNotFound, err)
		return
	}
	ctx.JSON(http.StatusNoContent, gin.H{})
}
