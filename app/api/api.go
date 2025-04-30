package api

import (
	"crypto/rsa"
	"fmt"
	"gin-server/app/models"
	"gin-server/app/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var appleService = &utils.AppleService{
	Config: models.AppleConfig{},
	HttpClient: &http.Client{
		Timeout: 10 * time.Second,
	},
	PublicKeys: make(map[string]*rsa.PublicKey),
	KeysExpiry: time.Time{},
}

type AppleApi struct{}

// AppleAuth 授权
func (appleApi *AppleApi) AppleAuth(ctx *gin.Context) {
	var request models.AppleToken

	appleService.Config = models.AppleConfig{
		ClientID: "",
		TeamID:   "",
	}

	appleUser, err := appleService.VerifyAppleToken(request.IdentityToken)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "apple auth failed",
		})
		return
	}

	fmt.Println(appleUser)
}
