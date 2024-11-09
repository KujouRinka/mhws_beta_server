package backend

type SystemPkt struct {
	ApiTimeout     int    `json:"api_timeout"`
	CustomProperty string `json:"custom_property"`
	JsonVer        string `json:"json_ver"`
	MMR            string `json:"mmr"`
	MTM            string `json:"mtm"`
	MTMs           string `json:"mtms"`
	NKM            string `json:"nkm"`
	Revision       string `json:"revision"`
	Selector       string `json:"selector"`
	Title          string `json:"title"`
	TMR            string `json:"tmr"`
	WLT            string `json:"wlt"`
	WorkingState   string `json:"working_state"`
}

type CustomProperty struct {
	ObtInfo *ObtInfo `json:"obt_info"`
	QA3     *QA3     `json:"qa3"`
}

type ObtInfo struct {
	Env       int   `json:"env"`
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
}

type QA3 struct {
	Api    string `json:"api"`
	Notify string `json:"notify"`
}

type AuthLogin struct {
	SessionId        string `msgpack:"SessionId"`
	UserId           string `msgpack:"UserId"`
	IsInCommunityBan bool   `msgpack:"IsInCommunityBan"`
}

type Hunter struct {
	HunterSaveList []HunterInfo `msgpack:"HunterSaveList"`
	UsingSaveSlot  int          `msgpack:"UsingSaveSlot"`
}

type HunterInfo struct {
	HunterId   string `msgpack:"HunterId"`
	HunterName string `msgpack:"HunterName"`
	OtomoName  string `msgpack:"OtomoName"`
	SaveSlot   int    `msgpack:"SaveSlot"`
}

type HunterSyncResponse struct {
	InvalidSaveSlotInfoList   interface{}          `msgpack:"InvalidSaveSlotInfoList"`
	InvalidClientHunterIdList interface{}          `msgpack:"InvalidClientHunterIdList"`
	SaveSlotInfoList          []HunterSaveResponse `msgpack:"SaveSlotInfoList"`
}

type HunterSaveResponse struct {
	HunterInfo
	ShortId string `msgpack:"ShortId"`
}

type AuthLoginResponse struct {
	SessionId        string `msgpack:"SessionId"`
	UserId           string `msgpack:"UserId"`
	IsInCommunityBan bool   `msgpack:"IsInCommunityBan"`
}

type HunterUpload struct {
	UploadUrl     string        `msgpack:"UploadUrl"`
	SignedHeaders []SignHeaders `msgpack:"SignedHeaders"`
}

type SignHeaders struct {
	HeaderKey    string   `msgpack:"HeaderKey"`
	HeaderValues []string `msgpack:"HeaderValues"`
}

type FollowTotalList struct {
	FollowList      []interface{} `msgpack:"FollowList"`
	LastOperationId string        `msgpack:"LastOperationId"`
}

type EmptyList struct {
	List []interface{} `msgpack:"List"`
}

type BlockList struct {
	IsConsistent  bool          `msgpack:"IsConsistent"`
	BlockedHunter []interface{} `msgpack:"BlockedHunter"`
	OperationId   int           `msgpack:"OperationId"`
}

type FriendList struct {
	FriendList []interface{} `msgpack:"FriendList"`
}

type LobbyAutoJoin struct {
	Endpoints []string `msgpack:"Endpoints"`
}
