package backend

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
