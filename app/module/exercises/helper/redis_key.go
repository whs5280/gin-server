package helper

import "fmt"

func GetAuthorizeKey(userId string) string {
	return fmt.Sprintf("%s:%s", "user_token", userId)
}
