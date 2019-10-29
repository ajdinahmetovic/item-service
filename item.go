package main

import (
	"context"

	"github.com/ajdinahmetovic/item-service/db"
	"github.com/ajdinahmetovic/item-service/proto/v1"
)

//CreateItem func
func (s *Server) CreateItem(ctx context.Context, request *proto.CreateItemReq) (*proto.CreateItemRes, error) {
	item := request.GetItem()
	id, err := db.AddItem(&db.Item{
		ID:          int(item.ID),
		Title:       item.Title,
		Description: item.Description,
		UserID:      int(item.UserID),
	})
	if err != nil {
		return nil, err
	}

	return &proto.CreateItemRes{Message: "Item added", ID: int32(id)}, nil
}

//GetItem func
func (s *Server) GetItem(ctx context.Context, request *proto.GetItemReq) (*proto.GetItemRes, error) {
	queryItem := db.Item{
		ID:          int(request.GetID()),
		Title:       request.GetTitle(),
		Description: request.GetDescription(),
		UserID:      int(request.GetUserID()),
	}

	items, err := db.FindItem(&queryItem)
	if err != nil {
		return &proto.GetItemRes{Item: nil, Message: "Failed to query database"}, err
	}

	res := make([]*proto.Item, 0)
	for _, i := range items {
		res = append(res, &proto.Item{
			ID:          int32(i.ID),
			Title:       i.Title,
			Description: i.Description,
			UserID:      int32(i.UserID),
		})
	}
	return &proto.GetItemRes{Item: res, Message: "Items found"}, nil
}

//UpdateItem func
func (s *Server) UpdateItem(ctx context.Context, request *proto.UpdateItemReq) (*proto.UpdateItemRes, error) {
	item := request.GetItem()
	err := db.UpdateItem(&db.Item{
		ID:          int(item.ID),
		Title:       item.Title,
		Description: item.Description,
		UserID:      int(item.UserID),
	})
	if err != nil {
		return &proto.UpdateItemRes{Item: nil, Message: "Update in DB failed"}, err
	}
	return &proto.UpdateItemRes{Item: nil, Message: "Update successful"}, nil
}

//DeleteItem func
func (s *Server) DeleteItem(ctx context.Context, request *proto.DeleteItemReq) (*proto.DeleteItemRes, error) {
	id := request.GetID()
	err := db.DeleteItem(int(id))
	if err != nil {
		return &proto.DeleteItemRes{Message: "Faield to delete"}, nil
	}
	return &proto.DeleteItemRes{Message: "Deleted successful"}, nil
}
