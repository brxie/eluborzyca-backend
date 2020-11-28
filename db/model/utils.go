package model

import "go.mongodb.org/mongo-driver/bson"

func toBSON(v interface{}) (*bson.M, error) {
	var (
		err error
		doc bson.M
	)
	data, err := bson.Marshal(v)
	if err != nil {
		return nil, err
	}

	err = bson.Unmarshal(data, &doc)
	return &doc, nil
}
