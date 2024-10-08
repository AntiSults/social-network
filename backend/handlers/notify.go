package handlers

import (
	"encoding/json"
	"log"
	"social-network/db/sqlite"
	"social-network/middleware"
	"social-network/sockets"
	"social-network/structs"
)

func triggerFollowNotification(userID, followerID int) {
	var Data struct {
		User structs.User
	}
	followerInfo, err := middleware.GetUser(followerID)
	if err != nil {
		log.Printf("Failed to retrieve follower info: %v", err)
		return
	}
	Data.User = *followerInfo

	dataJSON, err := json.Marshal(&Data)
	if err != nil {
		log.Printf("Error marshalling follower info: %v", err)
		return
	}
	event := sockets.Event{
		Type:    sockets.EventFollowNotify,
		Payload: dataJSON,
	}
	manager := sockets.GetManager()

	client, ok := manager.ClientsByUserID[userID]
	if !ok {
		log.Printf("User with ID %d not connected", userID)
		return
	}
	if err := sockets.GetManager().HandleNotify(event, client); err != nil {
		log.Printf("Error triggering follow notification: %v", err)
	}
}

func triggerGroupInvite(userID, GroupID, InviterID int) {
	var Data struct {
		User      structs.User
		GroupName string
	}
	inviterInfo, err := middleware.GetUser(InviterID)
	if err != nil {
		log.Printf("Failed to retrieve inviter info: %v", err)
		return
	}
	groupName, _, err := sqlite.Db.GetGroupNameAndCreatorID(GroupID)
	if err != nil {
		log.Printf("Failed to retrieve group name info: %v", err)
		return
	}
	Data.User = *inviterInfo
	Data.GroupName = groupName

	dataJSON, err := json.Marshal(&Data)
	if err != nil {
		log.Printf("Error marshalling follower info: %v", err)
		return
	}
	event := sockets.Event{
		Type:    sockets.EventGroupInvite,
		Payload: dataJSON,
	}
	manager := sockets.GetManager()

	client, ok := manager.ClientsByUserID[userID]
	if !ok {
		log.Printf("User with ID %d not connected", userID)
		return
	}
	if err := sockets.GetManager().HandleNotify(event, client); err != nil {
		log.Printf("Error triggering follow notification: %v", err)
	}
}

func triggerGroupJoin(GroupID, reqUserID int) {
	var Data struct {
		User      structs.User
		GroupName string
	}
	reqUserInfo, err := middleware.GetUser(reqUserID)
	if err != nil {
		log.Printf("Failed to retrieve inviter info: %v", err)
		return
	}
	groupName, creatorID, err := sqlite.Db.GetGroupNameAndCreatorID(GroupID)
	if err != nil {
		log.Printf("Failed to retrieve group name info: %v", err)
		return
	}
	Data.User = *reqUserInfo
	Data.GroupName = groupName

	dataJSON, err := json.Marshal(&Data)
	if err != nil {
		log.Printf("Error marshalling follower info: %v", err)
		return
	}
	event := sockets.Event{
		Type:    sockets.EventGroupJoin,
		Payload: dataJSON,
	}
	manager := sockets.GetManager()

	client, ok := manager.ClientsByUserID[creatorID]
	if !ok {
		log.Printf("User with ID %d not connected", creatorID)
		return
	}
	if err := sockets.GetManager().HandleNotify(event, client); err != nil {
		log.Printf("Error triggering follow notification: %v", err)
	}
}
