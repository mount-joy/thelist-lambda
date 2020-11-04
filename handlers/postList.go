package handlers

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/aws/aws-lambda-go/events"
	"github.com/mount-joy/thelist-lambda/db"
)

type postList struct {
	db db.DB
}

func newPostList() RouteHandler {
	return &postList{
		db: db.DynamoDB(),
	}
}

func (g *postList) match(request events.APIGatewayProxyRequest) bool {
	return request.Path == "/lists" && request.HTTPMethod == "POST"
}

func (g *postList) handle(request events.APIGatewayProxyRequest) (interface{}, int) {
	log.Printf("Running %s", request.Path)
	listID, err := g.createList(request.Path)

	if err != nil {
		log.Printf("Error: %s", err.Error())
		return nil, http.StatusInternalServerError
	}

	list := map[string]string{"listID": listID}

	return list, http.StatusOK
}

func (g *postList) createList(path string) (string, error) {
	listID := generateListID()
	return listID, g.db.CreateList(&listID)
}

func generateListID() (string) {
	return uuid.New().String()
}
