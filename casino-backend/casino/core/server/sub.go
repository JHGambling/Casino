package server

import (
	"jhgambling/backend/core/utils"
	"jhgambling/protocol"
	"jhgambling/protocol/models"
)

type DBSubscription struct {
	TableID    string      `json:"tableID"`
	ResourceID interface{} `json:"resourceID"`
}

type SubscriptionManager struct {
	gateway               *Gateway
	ChangedRecordsChannel chan protocol.SubChangedRecord
}

func NewSubscriptionsManager(gateway *Gateway) *SubscriptionManager {
	return &SubscriptionManager{
		gateway:               gateway,
		ChangedRecordsChannel: make(chan protocol.SubChangedRecord, 256),
	}
}

func (sub *SubscriptionManager) Update() {
	select {
	case changedRecord := <-sub.ChangedRecordsChannel:
		sub.handleChangedRecord(changedRecord)
	default:
		break
	}
}

func (sub *SubscriptionManager) handleChangedRecord(rec protocol.SubChangedRecord) {
	utils.Log("debug", "casino::server", "[sub] op:'", rec.Operation, "' table:'", rec.TableID, "' resource:'", rec.ResourceID, "'")

	for _, client := range sub.gateway.Clients {
		for _, subscription := range client.Subscriptions {
			isSubscribed := sub.isSubscribed(subscription, rec)
			if isSubscribed {
				// Client is subscribed to this record change, but we
				// have to check if the user is allowed to view this record at all
				if sub.canViewRecord(client.authenticatedAs, rec) {
					client.SendSubscriptionUpdatePacket(rec)
				}
			}
		}
	}
}

// Returns whether or not a SubChangedRecord matches the DBSubscription
func (sub *SubscriptionManager) isSubscribed(subscription DBSubscription, record protocol.SubChangedRecord) bool {
	if subscription.TableID != record.TableID {
		// Not subscribed to this table or any of its
		// records so we can just return false
		return false
	}

	if isZero(subscription.ResourceID) {
		// No specific record is being watched, so we are interested in the
		// entire table -> the record matches all requirements
		return true
	} else {
		// We are looking for a specific record
		return record.ResourceID == subscription.ResourceID
	}
}

func (sub *SubscriptionManager) canViewRecord(userID uint, record protocol.SubChangedRecord) bool {
	user, err := sub.gateway.ctx.Database.GetUserTable().FindByID(userID)
	if err != nil {
		utils.Log("warn", "casino::server", "[sub] canViewRecord() failed to find user with ID:", userID)
		return false
	}

	table, err := sub.gateway.ctx.Database.GetTable(record.TableID)
	if err != nil {
		utils.Log("warn", "casino::server", "[sub] canViewRecord() failed to find table with ID:", record.TableID)
		return false
	}

	userModel, ok := user.(*models.UserModel)
	if !ok {
		utils.Log("warn", "casino::server", "[sub] canViewRecord() failed to convert user")
		return false
	}

	return table.CanViewChangedRecord(*userModel, record)
}

func isZero(val interface{}) bool {
	switch v := val.(type) {
	case int:
		return v == 0
	case int64:
		return v == 0
	case float64:
		return v == 0
	case uint:
		return v == 0
	default:
		return false
	}
}
