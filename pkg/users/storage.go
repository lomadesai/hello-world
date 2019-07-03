package users

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const (
	PartitionKeyName = "UserId"
	SortKeyName      = "DateOfBirth"
)

type UserStore struct {
	tableName string
	db        *dynamodb.DynamoDB
}

type DynamoItem struct {
	Id          string
	FirstName   string
	LastName    string
	DateOfBirth string
	Email       string
	PhoneNumber string
}

func InitUserStore(tableName string, region string) (*UserStore, error) {
	store := &UserStore{tableName: tableName}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		fmt.Println("Failed to create aws session", err)
		return nil, err
	}

	store.db = dynamodb.New(sess)
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String(PartitionKeyName),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String(SortKeyName),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String(PartitionKeyName),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String(SortKeyName),
				KeyType:       aws.String("RANGE"),
			},
		},
		BillingMode: aws.String(dynamodb.BillingModePayPerRequest),
		SSESpecification: &dynamodb.SSESpecification{
			Enabled: aws.Bool(true),
			SSEType: aws.String(dynamodb.SSETypeKms),
		},
		TableName: aws.String(tableName),
	}

	_, err = store.db.CreateTable(input)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			switch awsErr.Code() {
			case dynamodb.ErrCodeTableAlreadyExistsException:
				fallthrough
			case dynamodb.ErrCodeResourceInUseException:
				fmt.Println("Dynamodb table already exists", tableName)
			default:
				fmt.Print("failed to create dynamodb table", err.Error())
				return nil, err
			}
		}
	}
	return store, nil
}

func (s *UserStore) UpsertUser(id string, user *User) error {
	input := &dynamodb.PutItemInput{
		Item:         s.buildAttributeMap(id, user),
		ReturnValues: aws.String("NONE"),
		TableName:    aws.String(s.tableName),
	}

	_, err := s.db.PutItem(input)
	if err != nil {
		fmt.Println("failed to upsert user item to dynamodb", err.Error())
		return err
	}
	return nil
}

func (s *UserStore) buildAttributeMap(id string, user *User) map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		PartitionKeyName: {
			S: aws.String(id),
		},
		SortKeyName: {
			S: aws.String(user.DateOfBirth),
		},
		"FirstName": {
			S: aws.String(user.FirstName),
		},
		"LastName": {
			S: aws.String(user.LastName),
		},
		"Email": {
			S: aws.String(user.Email),
		},
		"PhoneNumber": {
			S: aws.String(user.PhoneNumber),
		},
	}
}
