package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/byte-power/rockgo/rock"
	"github.com/byte-power/rockgo/util"
	"github.com/kataras/iris/v12"
)

func handleArticles(app rock.Application) {
	app.Serve("arts", "/articles").
		Get(func(ctx iris.Context) {
			ctx.ResponseWriter().Write([]byte("articles: ..."))
		})
	
	app.Serve("art", "/art/{id:int}").
		Get(func(ctx iris.Context) {
			id, _ := strconv.Atoi(ctx.Params().Get("id"))
			if id <= 0 {
				ctx.ResponseWriter().WriteHeader(http.StatusBadRequest)
				return
			}
			ctx.Text(fmt.Sprint("w ", id))
		}).
		Put(func(ctx iris.Context) {
			rock.Logger("WS").Info("Put", "id", ctx.Params().Get("id"))
		}).
		Delete(func(ctx iris.Context) {
			id, _ := ctx.Params().GetInt("id")
			handleDeleteArticle(ctx, id)
		})
	// to be compatable with old API
	app.Serve("delete_art", "/delete_art").
		Post(func(ctx iris.Context) {
			id := util.AnyToInt64(ctx.Request().FormValue("id"))
			handleDeleteArticle(ctx, int(id))
		})
}

// view or handler
func handleDeleteArticle(ctx iris.Context, id int) {
	err := deleteArticle(id)
	if err == nil {
		ctx.StatusCode(http.StatusOK)
	} else {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(util.StrMap{"error": err.Error()})
	}
}

// controller
func deleteArticle(id int) error {
	print(id)
	return nil
}
