// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	websocket "wszero/internal/handler/websocket"
	"wszero/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/wss",
				Handler: websocket.WssHandler(serverCtx),
			},
		},
		//rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)
}
