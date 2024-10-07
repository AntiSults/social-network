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

	followerInfo, err := middleware.GetUser(followerID)
	if err != nil {
		log.Printf("Failed to retrieve follower info: %v", err)
		return
	}
	dataJSON, err := json.Marshal(followerInfo)
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
	var data struct {
		Inviter   structs.User
		GroupName string
	}
	inviterInfo, err := middleware.GetUser(InviterID)
	if err != nil {
		log.Printf("Failed to retrieve inviter info: %v", err)
		return
	}
	groupName, err := sqlite.Db.GetGroupName(GroupID)
	if err != nil {
		log.Printf("Failed to retrieve group name info: %v", err)
		return
	}
	data.Inviter = *inviterInfo
	data.GroupName = groupName

	dataJSON, err := json.Marshal(&data)
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
