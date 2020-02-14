package whatsapp

import (
	"encoding/json"
	"fmt"
	"time"
	"github.com/catalinaM/go-whatsapp/binary"
	"strconv"
)

func (wac *Conn) GetGroupMetaData(jid string) (<-chan string, error) {
	data := []interface{}{"query", "GroupMetadata", jid}
	return wac.writeJson(data)
}

func (wac *Conn) GroupAnnouceFlag(jid string, cid string, close bool) (<-chan string, error) {
		ts := time.Now().Unix()
		tag := fmt.Sprintf("%d.--%d", ts, wac.msgCount)

		a := binary.Node{
			Description: "announcement",
			Attributes: map[string]string{
				"value": "true",
			},
			Content: nil,
		}

		fmt.Println(wac.session.Wid)

 		g := binary.Node{
			Description: "group",
			Attributes: map[string]string{
				"id": tag,
				"jid": jid,
				"type":   "prop",
				"author": cid,
			},
			Content: []interface{}{a},
		}

		n := binary.Node{
			Description: "action",
			Attributes: map[string]string{
				"type":  "set",
				"epoch": strconv.Itoa(wac.msgCount),
			},
			Content: []interface{}{g},
		}

		ch, err := wac.writeBinary(n, group, ignore, tag)
		var response map[string]interface{}
		fmt.Println(n)
fmt.Println(err)

						select {
						case r := <-ch:
							if err := json.Unmarshal([]byte(r), &response); err != nil {
							 fmt.Println("error decoding response message: %v\n", err)
							}
						case <-time.After(wac.msgTimeout):
							 fmt.Println("rrr request timed out")
						}
						//
						// if int(response["status"].(float64)) != 200 {
						// 	fmt.Errorf("request responded with %d", response["status"])
						// }
		return ch,err

}
//         // data := []interface{}{"Chat", map[string]interface{}{
//         //         "id": jid,
//         //         "cmd": "action",
//         //         "data": []interface{}{
//         //                 "restrict", userID, close,
//         //         },
//         //         }}
// 				//
// 				// fmt.Println(data)
// 				// fmt.Println("LIB-----")
//         // return wac.setGroup("restrict", jid, "", nil)
//
//         //         //      return wac.writeJson(data2)
// 				request := []interface{}{"group", "inviteCode", jid}
// 				ch, err := wac.writeJson(request)
// 				if err != nil {
// 					return "", err
// 				}
// 				fmt.Println("inside")
//
// 				fmt.Println(ch)
// 				var response map[string]interface{}
//
// 				select {
// 				case r := <-ch:
// 					if err := json.Unmarshal([]byte(r), &response); err != nil {
// 						return "", fmt.Errorf("error decoding response message: %v\n", err)
// 					}
// 				case <-time.After(wac.msgTimeout):
// 					return "", fmt.Errorf("request timed out")
// 				}
//
// 				if int(response["status"].(float64)) != 200 {
// 					return "", fmt.Errorf("request responded with %d", response["status"])
// 				}
// fmt.Println(response)
// 				return response["code"].(string), nil
// }

func (wac *Conn) CreateGroup(subject string, participants []string) (<-chan string, error) {
	return wac.setGroup("create", "", subject, participants)
}

func (wac *Conn) UpdateGroupSubject(subject string, jid string) (<-chan string, error) {
	return wac.setGroup("subject", jid, subject, nil)
}

func (wac *Conn) SetAdmin(jid string, participants []string) (<-chan string, error) {
	return wac.setGroup("promote", jid, "", participants)
}

func (wac *Conn) RemoveAdmin(jid string, participants []string) (<-chan string, error) {
	return wac.setGroup("demote", jid, "", participants)
}

func (wac *Conn) AddMember(jid string, participants []string) (<-chan string, error) {
	return wac.setGroup("add", jid, "", participants)
}

func (wac *Conn) RemoveMember(jid string, participants []string) (<-chan string, error) {
	return wac.setGroup("remove", jid, "", participants)
}

func (wac *Conn) LeaveGroup(jid string) (<-chan string, error) {
	return wac.setGroup("leave", jid, "", nil)
}

func (wac *Conn) GroupInviteLink(jid string) (string, error) {
	request := []interface{}{"query", "inviteCode", jid}
	ch, err := wac.writeJson(request)
	if err != nil {
		return "", err
	}

	var response map[string]interface{}

	select {
	case r := <-ch:
		if err := json.Unmarshal([]byte(r), &response); err != nil {
			return "", fmt.Errorf("error decoding response message: %v\n", err)
		}
	case <-time.After(wac.msgTimeout):
		return "", fmt.Errorf("request timed out")
	}

	if int(response["status"].(float64)) != 200 {
		return "", fmt.Errorf("request responded with %d", response["status"])
	}

	return response["code"].(string), nil
}

func (wac *Conn) GroupAcceptInviteCode(code string) (jid string, err error) {
	request := []interface{}{"action", "invite", code}
	ch, err := wac.writeJson(request)
	if err != nil {
		return "", err
	}

	var response map[string]interface{}

	select {
	case r := <-ch:
		if err := json.Unmarshal([]byte(r), &response); err != nil {
			return "", fmt.Errorf("error decoding response message: %v\n", err)
		}
	case <-time.After(wac.msgTimeout):
		return "", fmt.Errorf("request timed out")
	}

	if int(response["status"].(float64)) != 200 {
		return "", fmt.Errorf("request responded with %d", response["status"])
	}

	return response["gid"].(string), nil
}
