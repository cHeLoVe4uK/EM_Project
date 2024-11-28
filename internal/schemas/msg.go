package schemas

type RequestDeleteMsg struct {
	MsgID string `json:"msg_id"`
}

type RequestUpdateMsg struct {
	MsgID string `json:"msg_id"`
	Text  string `json:"text"`
}
