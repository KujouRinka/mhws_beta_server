package backend

import (
	"encoding/json"
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/vmihailenco/msgpack/v5"
)

const apiHost = "mhws.io"

var apis = []func(r *gin.Engine){
	registerSystemJson,
	registerListPartyQos,
	registerV1Api,
	registerAuth,
	registerInGame,
	registerWssHandler,
	registerOthers,
}

func RegisterHandler() *gin.Engine {
	r := gin.Default()
	for _, api := range apis {
		api(r)
	}
	return r
}

func registerSystemJson(r *gin.Engine) {
	r.GET("/systems/EAR-B-WW/00001/system.json", func(c *gin.Context) {
		m, err := filenameToMap("system.json")
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
		}
		c.JSON(200, m)
	})
	r.GET("/consents/EAR-B-WW/analysis/1/zh-hans.json", func(c *gin.Context) {
		m, err := filenameToMap("system.json")
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
		}
		c.JSON(200, m)
	})
}

func registerListPartyQos(r *gin.Engine) {
	r.POST("/MultiplayerServer/ListPartyQosServers", func(c *gin.Context) {
		m, err := filenameToMap("list_party_qos_servers.json")
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
		}
		c.JSON(200, m)
	})
}

func registerV1Api(r *gin.Engine) {
	g := r.Group("/v1")
	g.POST("/steam-steam/sign/EAR-B-WW", func(c *gin.Context) {
		m, err := filenameToMap("steam_sign_ear-b-ww.json")
		if err != nil {
			c.JSON(404, gin.H{"error": err.Error()})
		}
		c.JSON(200, m)
	})
	g.GET("/consent/restrictions/:country_code", func(c *gin.Context) {
		m, err := filenameToMap("restrictions.json")
		if err != nil {
			c.JSON(404, gin.H{"error": err.Error()})
		}
		c.JSON(200, m)
	})
	g.GET("/consent/countries/:country_code", func(c *gin.Context) {
		m, err := filenameToMap("countries.json")
		if err != nil {
			c.JSON(404, gin.H{"error": err.Error()})
		}
		c.JSON(200, m)
	})
	g.GET("/consent/documents/EAR-B-WW/:restriction/:lang/:tail", func(c *gin.Context) {
		m, err := filenameToMap("over.json")
		if err != nil {
			c.JSON(404, gin.H{"error": err.Error()})
		}
		c.JSON(200, m)
	})
	g.POST("/projects/*junk", func(c *gin.Context) {
		m, err := filenameToMap("projects.json")
		if err != nil {
			c.JSON(404, gin.H{"error": err.Error()})
		}
		c.JSON(200, m)
	})
}

// Character creation
var userId = newUUID()

func registerAuth(r *gin.Engine) {
	g := r.Group("/auth")

	g.POST("/login", func(c *gin.Context) {
		data := AuthLogin{
			SessionId:        newUUID(),
			UserId:           userId,
			IsInCommunityBan: false,
		}
		setMsgPack(c, data)
	})
	g.POST("/ticket", func(c *gin.Context) {
		data := "\x81\xa6" + "Ticket" + "\xd9\x24" + newUUID()

		setRawHeader(c, "x-session-nonce", newUUID())
		c.Data(200, "application/octct-stream", []byte(data))
	})
}

// Hunter profile
var hunterId string

func registerOthers(r *gin.Engine) {
	r.POST("/delivery_data/get", func(c *gin.Context) {
		setRawHeader(c, "x-session-nonce", uuid.New().String())
		c.File("asserts/delivery_data_get.bin")
		c.Header("Content-Type", "application/octet-stream")
	})

	hunterG := r.Group("/hunter")
	hunterG.POST("/sync", func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
		}

		var hunter Hunter
		if err := msgpack.Unmarshal(body, &hunter); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		hunterId = hunter.HunterSaveList[0].HunterId
		if hunterId == "" {
			hunterId = newUUID()
		}

		data := HunterSyncResponse{
			InvalidSaveSlotInfoList:   nil,
			InvalidClientHunterIdList: nil,
			SaveSlotInfoList: []HunterSaveResponse{
				{
					HunterInfo: HunterInfo{
						HunterId:   hunterId,
						HunterName: hunter.HunterSaveList[0].HunterName,
						OtomoName:  hunter.HunterSaveList[0].OtomoName,
						SaveSlot:   0,
					},
					ShortId: "1A2B3C4D",
				},
			},
		}

		byteData, err := msgpack.Marshal(data)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// c.Header("x-session-nonce", newUUID())
		c.Writer.Header()["x-session-nonce"] = []string{newUUID()}
		c.Data(200, "application/octet-stream", byteData)
	})
	hunterG.POST("/character_creation/upload", func(c *gin.Context) {
		data := HunterUpload{
			UploadUrl: "https://" + apiHost + "/character-creation/b9/" + userId,
			SignedHeaders: []SignHeaders{
				{
					HeaderKey:    "Host",
					HeaderValues: []string{apiHost},
				}, {
					HeaderKey:    "Content-Length",
					HeaderValues: []string{"3"},
				},
			},
		}
		setMsgPack(c, data)
	})
	hunterG.POST("/profile/update", func(c *gin.Context) {
		data := HunterUpload{
			UploadUrl: "https://" + apiHost + "/hunter-profile/dd/" + hunterId,
			SignedHeaders: []SignHeaders{
				{
					HeaderKey:    "Host",
					HeaderValues: []string{apiHost},
				}, {
					HeaderKey:    "Content-Length",
					HeaderValues: []string{"14113"},
				},
			},
		}
		setMsgPack(c, data)
	})
	hunterG.POST("/update/rank", func(c *gin.Context) {
		data := struct{}{}
		setMsgPack(c, data)
	})

	r.POST("/obt/play", func(c *gin.Context) {
		setRawHeader(c, "x-session-nonce", uuid.New().String())
		c.Header("Content-Type", "application/octet-stream")
		c.Data(200, "application/octet-stream", []byte{0x80})
	})
	r.PUT("/character-creation/*junk", func(c *gin.Context) {
		c.Data(200, "", []byte{})
	})
	r.PUT("/hunter-profile/*junk", func(c *gin.Context) {
		c.Data(200, "", []byte{})
	})
}

func registerInGame(r *gin.Engine) {
	r.POST("/follow/total_list", func(c *gin.Context) {
		data := FollowTotalList{
			FollowList:      []interface{}{},
			LastOperationId: "",
		}
		setMsgPack(c, data)
	})
	r.POST("/offline/notification_list", func(c *gin.Context) {
		data := EmptyList{List: []interface{}{}}
		setMsgPack(c, data)
	})
	r.POST("/community/invitation/received_list", func(c *gin.Context) {
		data := EmptyList{List: []interface{}{}}
		setMsgPack(c, data)
	})
	r.POST("/block/list", func(c *gin.Context) {
		data := BlockList{
			IsConsistent:  true,
			BlockedHunter: []interface{}{},
			OperationId:   0,
		}
		setMsgPack(c, data)
	})
	r.POST("/friend/list", func(c *gin.Context) {
		data := FriendList{FriendList: []interface{}{}}
		setMsgPack(c, data)
	})
	r.POST("/lobby/auto_join", func(c *gin.Context) {
		data := LobbyAutoJoin{Endpoints: []string{apiHost + ":443"}}
		setMsgPack(c, data)
	})
}

func registerWssHandler(r *gin.Engine) {
	r.GET("/ws", func(c *gin.Context) {
		upgrader := websocket.Upgrader{}
		c.Header("Sec-WebSocket-Protocol", "access_token")
		ws, err := upgrader.Upgrade(c.Writer, c.Request, c.Writer.Header())
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
		}
		defer ws.Close()

		for {
			_, _, err = ws.ReadMessage()
			if err != nil {
				return
			}
		}
	})

	r.GET("/socket", func(c *gin.Context) {
		upgrader := websocket.Upgrader{}
		c.Header("Sec-WebSocket-Protocol", "access_token")
		ws, err := upgrader.Upgrade(c.Writer, c.Request, c.Writer.Header())
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
		}
		defer ws.Close()

		c.Header("Sec-WebSocket-Protocol", "access_token")
		message1 := "\x81\x01\x00\x00" + hunterId + userId
		err = ws.WriteMessage(websocket.BinaryMessage, []byte(message1))
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
		}
		message2 := "\x85\x00\x02\x01\x01\x63\x00\x00\x00" + "FAKENAME"
		err = ws.WriteMessage(websocket.BinaryMessage, []byte(message2))

		for i := 0; i < 8; i++ {
			_, _, err = ws.ReadMessage()
			if err != nil {
				break
			}
		}
	})
}

// --------------------------

func filenameToMap(filename string) (map[string]interface{}, error) {
	data, err := os.ReadFile("asserts/" + filename)
	if err != nil {
		return nil, err
	}

	var jsonData map[string]interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return nil, err
	}

	return jsonData, nil
}

func newUUID() string {
	return uuid.New().String()
}

func setRawHeader(c *gin.Context, key, value string) {
	c.Writer.Header()[key] = []string{value}
}

func setMsgPack(c *gin.Context, v interface{}) {
	byteData, err := msgpack.Marshal(v)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	setRawHeader(c, "x-session-nonce", uuid.New().String())
	c.Header("Content-Type", "application/octet-stream")
	c.Data(200, "application/octet-stream", byteData)
}
