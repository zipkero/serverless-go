package user

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type User struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

var (
	ErrorFailedToFetchRecord     = "failed to fetch record"
	ErrorFailedToFetchRecords    = "failed to fetch records"
	ErrorFailedToUnmarshalRecrd  = "failed to unmarshal record"
	ErrorFailedToUnmarshalRecrds = "failed to unmarshal records"
)

func FetchUser(email string, tableName string, dynamoClient dynamodbiface.DynamoDBAPI) (*User, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String(tableName),
	}

	result, err := dynamoClient.GetItem(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	item := new(User)
	err = dynamodbattribute.UnmarshalMap(result.Item, item)
	if err != nil {
		return nil, errors.New(ErrorFailedToUnmarshalRecrd)
	}
	return item, nil
}

func FetchUsers() ([]*User, error) {

}

func CreateUser(user *User) {

}

func UpdateUser() {

}

func DeleteUser() {

}
