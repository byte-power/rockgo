package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/byte-power/rockgo/rock"
	"github.com/kataras/iris/v12"
)

func handleWorkspaces(app rock.Application) {
	app.NewService("workspaces", "/workspaces").
		Get(func(ctx iris.Context) {
			ctx.ResponseWriter().Write([]byte("workspaces: ..."))
		})
	app.NewService("workspace", "/workspace/{id:int}").
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
		Delete(handleWorkspaceDelete)
	app.NewService("delete_workspace", "/delete_workspace/{id:int}").
		Delete(handleWorkspaceDelete)
}

func handleWorkspaceDelete(ctx iris.Context) {
	ctx.ResponseWriter().WriteHeader(http.StatusOK)
}
