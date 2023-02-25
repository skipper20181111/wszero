package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"log"
	"strconv"
	"strings"
	"time"
	"wszero/internal/svc"
	"wszero/internal/types"
)

type WssLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	conn   *websocket.Conn
	WhoAmI string
}

func NewWssLogic(ctx context.Context, svcCtx *svc.ServiceContext, conn *websocket.Conn) *WssLogic {

	//svcCtx.RedisClient = redis.NewClient(&redis.Options{
	//	Addr:     svcCtx.Config.Cache[0].Host,
	//	Password: svcCtx.Config.Cache[0].Pass, // no password set
	//	DB:       0,                           // use default DB
	//})
	return &WssLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		conn:   conn,
	}
}

func (l *WssLogic) Read(WhoAmI string) {
	_, p, err := l.conn.ReadMessage()
	fmt.Println(string(p), "read")
	if err != nil {
		log.Println(err)
		return
	}
	l.WhoAmI = string(p)
	WhoAmI = l.WhoAmI
	defer func() { // 避免忘记关闭，所以要加上close
		fmt.Println("read is down")
		l.conn.Close()
	}()
	HistoryMsg := l.svcCtx.RedisClient.ZRevRangeByScore(context.Background(), WhoAmI+"get", &redis.ZRangeBy{Min: "3", Max: strconv.FormatInt(time.Now().UnixNano()+999999999999, 10), Offset: 0, Count: 10})
	for i := len(HistoryMsg.Val()); i > 0; i-- {
		Msgformat := &types.MsgFormat{}
		json.Unmarshal([]byte(HistoryMsg.Val()[i-1]), Msgformat)
		l.conn.WriteMessage(websocket.TextMessage, []byte(Msgformat.Sender+":"+Msgformat.MsgString))
	}
	for {
		fmt.Println("进入主循环并等待")
		_, p, err := l.conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println("写入数据")
		fmt.Println(l.WhoAmI, string(p))
		l.RedisPublic(p)
	}

}
func (l *WssLogic) Write(WhoAmI string) {
	//_, p, err := l.conn.ReadMessage()
	//fmt.Println(string(p), "write")
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//l.WhoAmI = string(p)
	for {
		if l.WhoAmI == "" {
			time.Sleep(time.Second)
		} else {
			fmt.Println(l.WhoAmI, "初始化完成")
			break
		}
	}
	defer func() { // 避免忘记关闭，所以要加上close
		fmt.Println("write is down")
		l.conn.Close()
	}()
	subscribe := l.svcCtx.RedisClient.Subscribe(context.Background(), l.WhoAmI)
	for {
		fmt.Println("进入writ的主循环", subscribe)
		message, err := subscribe.ReceiveMessage(context.Background())
		if message == nil {
			time.Sleep(time.Second)
			fmt.Println("无数据")
			continue
		}
		fmt.Println("得到数据")
		Msgformat := &types.MsgFormat{}
		json.Unmarshal([]byte(message.Payload), Msgformat)
		err = l.conn.WriteMessage(websocket.TextMessage, []byte(Msgformat.Sender+":"+Msgformat.MsgString))
		if err != nil {
			return
		}
	}
}
func (l *WssLogic) RedisPublic(p []byte) {
	defer func() {
		if e := recover(); e != nil {
			return
		}
	}()
	// 写入订阅redis
	split := strings.Split(string(p), "!@##$")
	Receiver := split[0]
	Msg := split[1]
	Msgformat := &types.MsgFormat{Sender: l.WhoAmI, Receiver: Receiver, MsgString: Msg, Timestamp: time.Now().Format("2006-01-02 15:04:05")}
	marshal, _ := json.Marshal(Msgformat)
	publish := l.svcCtx.RedisClient.Publish(context.Background(), Receiver, string(marshal))
	fmt.Println("publish ：", publish.Val(), string(marshal))
	// 写入sender和receiver

	l.svcCtx.RedisClient.ZAdd(context.Background(), l.WhoAmI+"send", &redis.Z{Score: float64(time.Now().UnixNano()), Member: string(marshal)})
	l.svcCtx.RedisClient.ZAdd(context.Background(), Receiver+"get", &redis.Z{Score: float64(time.Now().UnixNano()), Member: string(marshal)})
}
