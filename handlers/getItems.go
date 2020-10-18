package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/mount-joy/thelist-lambda/data"
	"github.com/mount-joy/thelist-lambda/db"
)

type getItems struct {
	db db.DB
}

func newGetItems() RouteHandler {
	return &getItems{
		db: db.NewDB(),
	}
}

func (g *getItems) match(request events.APIGatewayProxyRequest) bool {
	var re = regexp.MustCompile(`^/list/([\w-]+)/items/?$`)
	return re.MatchString(request.Path)
}

func (g *getItems) getItems(path *string) (*[]data.Item, error) {
	listID, err := getListID(path)
	if err != nil {
		return nil, err
	}

	return g.db.GetItemsOnList(listID)
}

func (g *getItems) handle(request events.APIGatewayProxyRequest) (interface{}, int) {
	items, err := g.getItems(&request.Path)

	if err != nil {
		log.Printf("Error: %s", err.Error())
		return nil, http.StatusInternalServerError
	}

	return items, http.StatusOK
}

func getListID(path *string) (*string, error) {
	if path == nil {
		return nil, errors.New("Path is nil")
	}

	parts := strings.SplitN(*path, "/", 4)
	if len(parts) < 4 {
		return nil, fmt.Errorf("Unable to match path: %s", *path)
	}
	return &parts[2], nil
}
