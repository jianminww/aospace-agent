// Copyright (c) 2022 Institute of Software, Chinese Academy of Sciences (ISCAS)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pair

import (
	"agent/biz/docker"
	"agent/biz/model/device"
	"agent/biz/model/device_ability"
	"agent/biz/model/dto"
	"agent/biz/model/dto/pair/tryout"
	"agent/config"
	"fmt"

	"agent/utils/logger"
)

func ServiceTryout(req *tryout.TryoutCodeReq) (dto.BaseRspStr, error) {
	if device_ability.GetAbilityModel().RunInDocker {
		if docker.ContainersDownloading == docker.GetDockerStatus() {
			err := fmt.Errorf("docker images is downloading...")
			logger.AppLogger().Warnf("ServiceTryout,%v", err)
			return dto.BaseRspStr{Code: dto.AgentCodeDockerPulling, Message: err.Error(), Results: nil}, err
		}
	}

	return presetBoxInfo(req)
}

// 预置试用信息
func presetBoxInfo(req *tryout.TryoutCodeReq) (dto.BaseRspStr, error) {

	// 平台请求结构
	type platformReqStruct struct {
		Email   string `json:"email"`
		Code    string `json:"code"`
		Type    string `json:"type"`
		BoxInfo struct {
			BoxUUID   string            `json:"boxUUID"`
			Desc      string            `json:"desc"`
			Extra     map[string]string `json:"extra"`
			BoxPubKey string            `json:"boxPubKey"`
			AuthType  string            `json:"authType"`
		} `json:"boxInfo"`
	}
	// 平台响应结构
	type platformRspStruct struct {
		Code      string `json:"code"`
		Message   string `json:"message"`
		RequestId string `json:"requestId"`

		State   int32 `json:"state"` // 0-正常;1-禁用;2-已过期
		BoxInfo struct {
			AuthType     string `json:"authType"`
			SnNumber     string `json:"snNumber"`
			IsRegistered bool   `json:"isRegistered"`
		} `json:"boxInfo"`
	}

	// 请求平台
	parms := &platformReqStruct{}
	parms.Email = req.Email
	parms.Code = req.TryoutCode
	parms.Type = "pc_open"
	parms.BoxInfo.BoxUUID = device.GetDeviceInfo().BoxUuid
	parms.BoxInfo.Desc = "pc tryout"
	parms.BoxInfo.Extra = make(map[string]string)
	// parms.BoxInfo.BoxPubKey = strings.ReplaceAll(string(device.GetBoxPubKey()), "\n", "")
	parms.BoxInfo.BoxPubKey = string(device.GetDevicePubKey())
	parms.BoxInfo.AuthType = "box_pub_key"
	url := device.GetApiBaseUrl() + config.Config.Platform.PresetBoxInfo.Path
	logger.AppLogger().Debugf("presetBoxInfo, url:%+v, parms:%+v", url, parms)

	//TODO: Use pubkey as SN Number
	device.UpdateSnNumber(parms.BoxInfo.BoxUUID)
	device.UpdateApplyEmail(req.Email)
	return dto.BaseRspStr{Code: dto.AgentCodeOkStr, Message: "OK", Results: nil}, nil

}
