package main

import (
	"time"

	firebase "firebase.google.com/go"
	"github.com/line/line-bot-sdk-go/linebot"
)

type app struct {
	bot          *linebot.Client
	client       *firebase.App
	sessionStore *sessionStore
	service      *service
}

type sessionStore struct {
	sessions sessions
	lifespan time.Duration
}

type userSession struct {
	orderID   string
	prevStep  int
	createdAt time.Time
	products  Products
}

type sessions map[string]*userSession

type service struct {
	menu          Menu
	locations     []Location
	businessHours businessHours
	detailTime    detailTime
}

type Location struct {
	Name string `firestore:"name,omitempty"`
}

type businessHours struct {
	today     string
	begin     detailTime `firestore:"begin,omitempty"`
	end       detailTime `firestore:"end,omitempty"`
	interval  int        `firestore:"interval,omitempty"`
	lastorder string     `firestore:"lastorder,omitempty"`
}

type detailTime struct {
	hour   int
	minute int
}

// map[{products Document の ID}] 個数
type Products map[string]int

func (ss *sessionStore) createSession(userID string) *userSession {
	ss.sessions[userID] = &userSession{prevStep: begin, createdAt: time.Now(), products: make(Products)}
	return ss.sessions[userID]
}

func (ss *sessionStore) deleteUserSession(userID string) {
	delete(ss.sessions, userID)
}

func (ss *sessionStore) checkSessionLifespan(userID string) (ok bool) {
	session := ss.sessions[userID]
	diff := time.Since(session.createdAt)
	if ok = diff <= ss.lifespan; ok {
		return true
	}
	return false
}

func (ss *sessionStore) searchSession(userID string) *userSession {
	if ss.sessions[userID] != nil {
		return ss.sessions[userID]
	}
	return nil
}
