package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"

	uuid "github.com/google/uuid"
)

const (
	BOSS                 = "goofy"
	CREATE_INSTRUCTION   = "CreateCoin"
	TRANSFER_INSTRUCTION = "TransferCoin:"
)

type User struct {
	UUID       string
	name       string
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

type Node struct {
	payload   string
	signature []byte
	prev      *Node
	ownerId   string
}

//global variables

var usersStorage = make(map[string]*User) //for storing users
var ledger *Node
var GOOFYS_UUID string

func generateKeys() (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	publicKey := &privateKey.PublicKey
	return privateKey, publicKey, err
}

func createUser(name string) *User {
	privateKey, publicKey, err := generateKeys()
	if err != nil {
		fmt.Println("Error while generating keys")
		return nil
	}
	newUser := &User{UUID: uuid.New().String(), name: name, privateKey: privateKey, publicKey: publicKey}
	return newUser
}

func createNewNode(payload string, ownerId string) *Node {
	user := usersStorage[ownerId]
	r, s, err := ecdsa.Sign(rand.Reader, user.privateKey, []byte(payload))
	if err != nil {
		fmt.Println("Failed to sign payload of new node")
		return nil
	}
	newNode := &Node{payload: payload, signature: append(r.Bytes(), s.Bytes()...), prev: nil}
	return newNode
}

func createNewCoin(ownerId string) *Node {
	if ownerId != GOODYS_ID {
		fmt.Println("Oncly goofy can create coins")
		return nil
	}
	node := createNewNode(CREATE_INSTRUCTION, ownerId)
	if node == nil {
		fmt.Println("node is nil")
		return nil
	}
	if ledger != nil {
		node.prev = ledger
	}
	ledger = node
	return ledger
}

func transferCoin(fromId string, toId string) {
	fromUser := usersStorage[fromId]
	toUser := usersStorage[toId]
	if fromUser == nil || toUser == nil {
		fmt.Println("Failed transaction! error: to user and from user ids' must be specified correctly.")
		return
	}
	newNode := createNewNode(TRANSFER_INSTRUCTION+fromId+":"+toId, ownerId)
	if newNode == nil {
		fmt.Println("Failed to create new node for transaction")
	}
	if ledger == nil {
		fmt.Println("No coins created yet. Failed to make a transaction.")
	}
	newNode.prev = ledger
	ledger = newNode

}

func main() {
	fmt.Println("Starting goofy coin mechanism")
	user := createUser(BOSS)
	if user == nil {
		fmt.Println("Failed to create a user")
	}
	GOOFYS_ID = user.UUID
	usersStorage[user.UUID] = user
	goofy := usersStorage[user.UUID]
	updatedNode := createNewCoin(goofy.UUID)
	if updatedNode == nil {
		fmt.Println("Error while creating new coin")
	}
	fmt.Println(ledger)
}
