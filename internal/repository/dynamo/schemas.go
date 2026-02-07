package dynamo

import "time"

// UserSchema representa a estrutura exata da tabela no DynamoDB
type UserSchema struct {
	PK        string    `dynamodbav:"PK"`
	SK        string    `dynamodbav:"SK"`
	Name      string    `dynamodbav:"name"`
	Email     string    `dynamodbav:"email"`
	Password  string    `dynamodbav:"password"`
	CreatedAt time.Time `dynamodbav:"created_at"`
}