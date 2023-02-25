package websocket

import (
	"github.com/gorilla/websocket"
	"net/http"
	wlogic "wszero/internal/logic/websocket"
	"wszero/internal/svc"
)

func WssHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//WhoAmI := r.Header.Get("WhoAmI")
		conn, err := (&websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { // CheckOrigin解决跨域问题
				return true
			}}).Upgrade(w, r, nil) // 升级成ws协议
		if err != nil {
			return
		}
		l := wlogic.NewWssLogic(r.Context(), svcCtx, conn)
		go l.Read("WhoAmI")
		go l.Write("WhoAmI")

	}
}
