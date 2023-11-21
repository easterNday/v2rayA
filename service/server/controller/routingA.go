package controller

import (
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/v2fly/v2ray-core/v5/common/platform/filesystem"
	"github.com/v2rayA/RoutingA"
	"github.com/v2rayA/v2ray-lib/router/routercommon"
	"github.com/v2rayA/v2rayA/common"
	"github.com/v2rayA/v2rayA/db/configure"
	"google.golang.org/protobuf/proto"
)

func GetGeoIPList(ctx *gin.Context) {
	// 获取 GeoIP 的 country_code 列表
	geoipBytes, err := filesystem.ReadAsset("geoip.dat")
	var geoipResult []string
	if err != nil {
		geoipResult = append(geoipResult, "Err")
	}
	var geoipList routercommon.GeoIPList
	if err := proto.Unmarshal(geoipBytes, &geoipList); err != nil {
		geoipResult = append(geoipResult, "Err")
	} else {
		for _, geoip := range geoipList.Entry {
			geoipResult = append(geoipResult, geoip.CountryCode)
		}
	}

	// 获取 GeoSite 的 country_code 列表
	geositeBytes, err := filesystem.ReadAsset("geosite.dat")
	var geositeResult []string
	if err != nil {
		geositeResult = append(geositeResult, "Err")
	}
	var geositeList routercommon.GeoSiteList
	if err := proto.Unmarshal(geositeBytes, &geositeList); err != nil {
		geositeResult = append(geositeResult, "Err")
	} else {
		for _, geosite := range geositeList.Entry {
			geositeResult = append(geositeResult, geosite.CountryCode)
		}
	}

	// 返回结果
	common.ResponseSuccess(ctx, gin.H{
		"GeoIP":   geoipResult,
		"GeoSite": geositeResult,
	})
}

func GetRoutingA(ctx *gin.Context) {
	common.ResponseSuccess(ctx, gin.H{
		"routingA": configure.GetRoutingA(),
	})
}
func PutRoutingA(ctx *gin.Context) {
	var data struct {
		RoutingA string `json:"routingA"`
	}
	err := ctx.ShouldBindJSON(&data)
	if err != nil {
		common.ResponseError(ctx, logError("bad request"))
		return
	}
	// remove hardcode replacement and try parsing
	lines := strings.Split(data.RoutingA, "\n")
	hardcodeReplacement := regexp.MustCompile(`\$\$.+?\$\$`)
	for i := range lines {
		hardcodes := hardcodeReplacement.FindAllString(lines[i], -1)
		for _, hardcode := range hardcodes {
			lines[i] = strings.Replace(lines[i], hardcode, "", 1)
		}
	}
	_, err = RoutingA.Parse(strings.Join(lines, "\n"))
	if err != nil {
		common.ResponseError(ctx, logError(err))
		return
	}
	err = configure.SetRoutingA(&data.RoutingA)
	if err != nil {
		common.ResponseError(ctx, logError(err))
		return
	}
	common.ResponseSuccess(ctx, nil)
}
