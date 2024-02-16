package service

import (
	"fmt"
	"strings"

	"github.com/litterGuy/imaotai/config"
	"github.com/litterGuy/imaotai/db"
	"github.com/litterGuy/imaotai/models"
	"github.com/litterGuy/imaotai/reqfunc"
)

// Reservation 预约
func Reservation(conf *config.Config) (string, error) {
	var mtversion *models.MtVersion
	err := db.Gormdb.First(&mtversion).Error
	if err != nil {
		return "", err
	}
	var session *models.Session
	err = db.Gormdb.First(&session).Error
	if err != nil {
		return "", err
	}
	var sessionitems []*models.SessionItem
	err = db.Gormdb.Find(&sessionitems).Error
	if err != nil {
		return "", err
	}
	result := make([]string, 0)
	for i, account := range conf.Account {
		result = append(result, fmt.Sprintf("------------帐号%v:%s--------------\n\n", i, account.Phone))
		for _, item := range sessionitems {
			shopid, err := GetShopId(item.ItemCode, *account)
			if err != nil {
				return "", err
			}
			if len(shopid) == 0 {
				continue
			}
			rt, err := reqfunc.Reservation(reqfunc.UserInfo{
				UserId: account.UserId,
				Lat:    fmt.Sprintf("%v", account.Lat),
				Lng:    fmt.Sprintf("%v", account.Lng),
				Token:  account.Token,
			}, item.ItemCode, shopid, fmt.Sprintf("%v", session.SessionId), mtversion.Version)
			if err != nil {
				return "", err
			}
			result = append(result, rt)
		}
		result = append(result, "------------------------------\n\n")
	}

	return strings.Join(result, "\n"), nil
}
