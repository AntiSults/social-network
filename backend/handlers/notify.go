package handlers

import (
	"encoding/json"
	"log"
	"social-network/middleware"
	"social-network/sockets"
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
		Type:    sockets.EventNotify,
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
