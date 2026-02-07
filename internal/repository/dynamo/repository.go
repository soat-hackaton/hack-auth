package dynamo

import (
	"context"
	"errors"
	"fmt"
	"hack-auth/internal/domain"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type UserRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewUserRepository(client *dynamodb.Client, tableName string) *UserRepository {
	return &UserRepository{
		client:    client,
		tableName: tableName,
	}
}

func (r *UserRepository) Create(user *domain.User) error {
	// 1. Converte Domain -> Schema
	schema := toSchema(user)

	if schema.CreatedAt.IsZero() {
		schema.CreatedAt = time.Now()
	}

	// 2. Marshal do Schema
	av, err := attributevalue.MarshalMap(schema)
	if err != nil {
		return fmt.Errorf("failed to marshal schema: %w", err)
	}

	// 3. Persiste no DynamoDB
	_, err = r.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName:           aws.String(r.tableName),
		Item:                av,
		ConditionExpression: aws.String("attribute_not_exists(PK)"),
	})

	if err != nil {
		var condCheckFailed *types.ConditionalCheckFailedException
		if errors.As(err, &condCheckFailed) {
			return errors.New("email already exists")
		}
		return fmt.Errorf("failed to put item in dynamodb: %w", err)
	}

	return nil
}

func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
	pk := fmt.Sprintf("USER#%s", email)
	sk := "PROFILE"

	out, err := r.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get item: %w", err)
	}

	if out.Item == nil {
		return nil, errors.New("user not found")
	}

	// 1. Unmarshal no Schema
	var schema UserSchema
	err = attributevalue.UnmarshalMap(out.Item, &schema)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal schema: %w", err)
	}

	// 2. Converte Schema -> Domain e retorna
	return toDomain(schema), nil
}