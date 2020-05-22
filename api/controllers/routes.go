package controllers

import "github.com/tqvthu/social_friend_restful/api/middlewares"

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	// Relationship
	s.Router.HandleFunc("/friend/connect", middlewares.SetMiddlewareJSON(s.MakeFriend)).Methods("POST")
	s.Router.HandleFunc("/friendList", middlewares.SetMiddlewareJSON(s.GetFriendList)).Queries("email", "{email}").Methods("GET")
	s.Router.HandleFunc("/friend/common", middlewares.SetMiddlewareJSON(s.GetCommonFriends)).Methods("POST")
	s.Router.HandleFunc("/friend/subscribe", middlewares.SetMiddlewareJSON(s.SubscribeFriend)).Methods("POST")
	s.Router.HandleFunc("/friend/block", middlewares.SetMiddlewareJSON(s.Block)).Methods("POST")
	s.Router.HandleFunc("/friend/recipients", middlewares.SetMiddlewareJSON(s.GetRecipients)).Methods("POST")



}