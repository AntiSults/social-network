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
		User []structs.User
	}
	followerInfo, err := middleware.GetUser(followerID)
	if err != nil {
		log.Printf("Failed to retrieve follower info: %v", err)
		return
	}
	Data.User = append(Data.User, *followerInfo)

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

	client, ok := manager.ClientsByUserID[userID]["notify"]
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
		User      []structs.User
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
	Data.User = append(Data.User, *inviterInfo)
	Data.GroupName = groupName

	dataJSON, err := json.Marshal(&Data)
	if err != nil {
		log.Printf("Error marshalling group invite info: %v", err)
		return
	}
	event := sockets.Event{
		Type:    sockets.EventGroupInvite,
		Payload: dataJSON,
	}
	manager := sockets.GetManager()

	client, ok := manager.ClientsByUserID[userID]["notify"]
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
		User      []structs.User
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
	Data.User = append(Data.User, *reqUserInfo)
	Data.GroupName = groupName

	dataJSON, err := json.Marshal(&Data)
	if err != nil {
		log.Printf("Error marshalling group join info: %v", err)
		return
	}
	event := sockets.Event{
		Type:    sockets.EventGroupJoin,
		Payload: dataJSON,
	}
	manager := sockets.GetManager()

	client, ok := manager.ClientsByUserID[creatorID]["notify"]
	if !ok {
		log.Printf("User with ID %d not connected", creatorID)
		return
	}
	if err := sockets.GetManager().HandleNotify(event, client); err != nil {
		log.Printf("Error triggering follow notification: %v", err)
	}
}

func triggerGroupEventNotify(groupEvent structs.Event) {
	var Data struct {
		GroupEvent []structs.Event
		GroupName  string
	}
	groupName, _, err := sqlite.Db.GetGroupNameAndCreatorID(groupEvent.GroupID)
	if err != nil {
		log.Printf("Failed to retrieve group name info: %v", err)
		return
	}
	Data.GroupEvent = append(Data.GroupEvent, groupEvent)
	Data.GroupName = groupName

	dataJSON, err := json.Marshal(&Data)
	if err != nil {
		log.Printf("Error marshalling Group Event info: %v", err)
		return
	}
	event := sockets.Event{
		Type:    sockets.EventNewGroupEvent,
		Payload: dataJSON,
	}
	manager := sockets.GetManager()

	client, ok := manager.ClientsByUserID[groupEvent.UserID]["notify"]
	if !ok {
		log.Printf("User with ID %d is not online", groupEvent.UserID)
		return
	}
	if err := sockets.GetManager().HandleNotify(event, client); err != nil {
		log.Printf("Error triggering follow notification: %v", err)
	}
}
