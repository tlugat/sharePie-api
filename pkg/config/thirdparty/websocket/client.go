package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"sharePie-api/internal/models"
	"sharePie-api/internal/types"
	"strconv"
	"time"
)

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	room string
	send chan []byte
	user models.User
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Println("error unmarshaling message:", err)
			return
		}
		c.handleMessage(msg)
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.conn.WriteMessage(websocket.TextMessage, message)
		case <-ticker.C:
			c.conn.WriteMessage(websocket.PingMessage, []byte{})
		}
	}
}

func (c *Client) handleMessage(message Message) {
	switch message.Type {
	case "updateEvent":
		c.handleUpdateEvent(message.Payload)
	case "createExpense":
		c.handleCreateExpense(message.Payload)
	case "updateExpense":
		c.handleUpdateExpense(message.Payload)
	case "deleteExpense":
		c.handleDeleteExpense(message.Payload)
	default:
		log.Println("unknown event type:", message.Type)
	}
}

func (c *Client) handleUpdateEvent(payload json.RawMessage) {
	var input types.UpdateEventInput
	if err := json.Unmarshal(payload, &input); err != nil {
		fmt.Println("Failed to unmarshal data:", err)
		return
	}
	eventId, err := strconv.ParseUint(c.room, 10, 32)
	if err != nil {
		fmt.Println("Invalid eventId:", err)
		return
	}

	updatedEvent, err := c.hub.eventService.Update(uint(eventId), input)
	if err != nil {
		fmt.Println("Failed to update event:", err)
		return
	}

	updatedEventJSON, err := json.Marshal(updatedEvent)
	if err != nil {
		fmt.Println("Failed to marshal updated event:", err)
		return
	}

	message, err := json.Marshal(Message{
		Type:    "event",
		Payload: updatedEventJSON,
	})
	if err != nil {
		fmt.Println("Failed to marshal updated event:", err)
		return
	}

	c.hub.rooms[c.room].broadcast <- message

	users, err := c.hub.eventService.GetUsers(uint(eventId))
	if err != nil {
		fmt.Println("Failed to get event users:", err)
		return
	}

	usersJSON, err := json.Marshal(users)
	if err != nil {
		fmt.Println("Failed to marshal users:", err)
		return
	}

	usersMessage, err := json.Marshal(Message{
		Type:    "users",
		Payload: usersJSON,
	})
	if err != nil {
		fmt.Println("Failed to marshal updated event:", err)
		return
	}

	c.hub.rooms[c.room].broadcast <- usersMessage
}

func (c *Client) handleCreateExpense(payload json.RawMessage) {
	var input types.CreateExpenseInput
	if err := json.Unmarshal(payload, &input); err != nil {
		fmt.Println("Failed to unmarshal data:", err)
		return
	}

	_, err := c.hub.expenseService.Create(input, c.user)
	if err != nil {
		fmt.Println("Failed to create expense:", err)
		return
	}

	c.refreshExpenses()
}

func (c *Client) handleUpdateExpense(payload json.RawMessage) {
	var input types.UpdateExpenseInput
	if err := json.Unmarshal(payload, &input); err != nil {
		fmt.Println("Failed to unmarshal data:", err)
		return
	}

	_, err := c.hub.expenseService.Update(input.ID, input)
	if err != nil {
		fmt.Println("Failed to update expense:", err)
		return
	}

	c.refreshExpenses()
}

func (c *Client) handleDeleteExpense(payload json.RawMessage) {
	var input DeleteExpenseInput
	if err := json.Unmarshal(payload, &input); err != nil {
		fmt.Println("Failed to unmarshal data:", err)
		return
	}

	err := c.hub.expenseService.Delete(input.ID)
	if err != nil {
		fmt.Println("Failed to delete expense:", err)
		return
	}

	c.refreshExpenses()
}

func (c *Client) refreshExpenses() {
	eventId, err := strconv.ParseUint(c.room, 10, 32)
	if err != nil {
		fmt.Println("Invalid eventId:", err)
		return
	}

	expenses, err := c.hub.expenseService.FindByEventId(uint(eventId))
	if err != nil {
		fmt.Println("Failed to get expenses:", err)
		return
	}

	expensesJSON, err := json.Marshal(expenses)
	if err != nil {
		fmt.Println("Failed to marshal expenses:", err)
		return
	}

	message, err := json.Marshal(Message{
		Type:    "expenses",
		Payload: expensesJSON,
	})
	if err != nil {
		fmt.Println("Failed to marshal expenses message:", err)
		return
	}

	event, err := c.hub.eventService.FindOne(uint(eventId))
	if err != nil {
		fmt.Println("Failed to get event:", err)
		return
	}

	c.hub.rooms[c.room].broadcast <- message
	c.refreshBalances(event)
	c.refreshTransactions(event)
	c.refreshUsers(event)
}

func (c *Client) refreshBalances(event models.Event) {
	balances, err := c.hub.eventService.GetBalances(event)
	if err != nil {
		fmt.Println("Failed to get balances:", err)
		return
	}

	balancesJSON, err := json.Marshal(balances)
	if err != nil {
		fmt.Println("Failed to marshal balances:", err)
		return
	}

	message, err := json.Marshal(Message{
		Type:    "balances",
		Payload: balancesJSON,
	})
	if err != nil {
		fmt.Println("Failed to marshal balances message:", err)
		return
	}

	c.hub.rooms[c.room].broadcast <- message
}

func (c *Client) refreshTransactions(event models.Event) {
	transactions, err := c.hub.eventService.GetTransactions(event)
	if err != nil {
		fmt.Println("Failed to create transactions:", err)
		return
	}

	transactionsJSON, err := json.Marshal(transactions)
	if err != nil {
		fmt.Println("Failed to marshal transactions:", err)
		return
	}

	message, err := json.Marshal(Message{
		Type:    "transactions",
		Payload: transactionsJSON,
	})
	if err != nil {
		fmt.Println("Failed to marshal transactions message:", err)
		return
	}

	c.hub.rooms[c.room].broadcast <- message
}

func (c *Client) refreshUsers(event models.Event) {
	users, err := c.hub.eventService.GetUsers(event.ID)
	if err != nil {
		fmt.Println("Failed to get users:", err)
		return
	}

	usersJSON, err := json.Marshal(users)
	if err != nil {
		fmt.Println("Failed to marshal users:", err)
		return
	}

	message, err := json.Marshal(Message{
		Type:    "users",
		Payload: usersJSON,
	})
	if err != nil {
		fmt.Println("Failed to marshal users message:", err)
		return
	}

	c.hub.rooms[c.room].broadcast <- message
}
