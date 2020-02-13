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

func (wac *Conn) CloseGroup(jid, userID string, close bool) (<-chan string, error) {
	fmt.Println("st")

		ts := time.Now().Unix()
		tag := fmt.Sprintf("%d.--%d", ts, wac.msgCount)
		fmt.Println("inside")

		//TODO: get proto or improve encoder to handle []interface{}
		//
		// p := buildParticipantNodes(participants)

		fmt.Println("inside")
		g := binary.Node{
			Description: "group",
			Attributes: map[string]string{
				"author": wac.session.Wid,
				"id":     tag,
				// "type":   "announce",
				"announce": "false",
				"jid": jid,
				// "announce": strconv.FormatBool(close),
			},
			// Content: nil,
		}
		// g.Attributes["subject"] = "test test3"

		fmt.Println("inside")

		n := binary.Node{
			Description: "action",
			Attributes: map[string]string{
				"type":  "set",
				"epoch": strconv.Itoa(wac.msgCount),
			},
			Content: []interface{}{g},
		}
		fmt.Println("inside last")

		return wac.writeBinary(n, group, ignore, tag)

		// if err != nil {
		// 	return "", err
		// }
		// fmt.Println(err)
		//
		// fmt.Println(ch)
		// var response map[string]interface{}
		//
		// select {
		// case r := <-ch:
		// 	if err := json.Unmarshal([]byte(r), &response); err != nil {
		// 		fmt.Println("error decoding response message: %v\n", err)
		// 		return "", fmt.Errorf("error decoding response message: %v\n", err)
		// 	}
		// case <-time.After(wac.msgTimeout):
		// 		fmt.Println("request timed out")
		// 	return "", fmt.Errorf("request timed out")
		// }
		// //
		// // if int(response["status"].(float64)) != 200 {
		// // 	return "", fmt.Errorf("request responded with %d", response["status"])
		// // }
		// fmt.Println(response)
		// return response["code"].(string), nil

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

func (wac *Conn) OpenGroup(jid, userID string, close bool) (<-chan string, error) {

		ts := time.Now().Unix()
		tag := fmt.Sprintf("%d.--%d", ts, wac.msgCount)
		fmt.Println("inside")

		//TODO: get proto or improve encoder to handle []interface{}
		//
		// p := buildParticipantNodes(participants)

		fmt.Println("inside")
		g := binary.Node{
			Description: "group",
			Attributes: map[string]string{
				"author": wac.session.Wid,
				"id":     tag,
				"type":   "announce",
				"jid": jid,
				"announce": strconv.FormatBool(close),
			},
			Content: nil,
		}
		fmt.Println("inside")

		n := binary.Node{
			Description: "action",
			Attributes: map[string]string{
				"type":  "set",
				"epoch": strconv.Itoa(wac.msgCount),
			},
			Content: []interface{}{g},
		}
		fmt.Println("inside")

		return wac.writeBinary(n, group, ignore, tag)
}


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
