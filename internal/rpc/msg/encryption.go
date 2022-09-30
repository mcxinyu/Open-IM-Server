package msg

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	"Open_IM/pkg/proto/encryption"
	sdk_ws "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	"errors"
	"strings"
)

func EncryptionContent(msgData *sdk_ws.MsgData, recvIDList []string, operationID string) (error, []*encryption.MsgContent) {
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImEncryptionName, operationID)
	if etcdConn == nil {
		errMsg := operationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(operationID, errMsg)
		return errors.New(errMsg), nil
	}
	client := encryption.NewEncryptionClient(etcdConn)
	encryptContentReq := encryption.EncryptContentReq{OperationID: operationID, RecvIDList: recvIDList,
		MsgContent: &encryption.MsgContent{SendID: msgData.SendID, RecvID: msgData.RecvID, GroupID: msgData.GroupID,
			SessionType: msgData.SessionType, MsgFrom: msgData.MsgFrom, ContentType: msgData.ContentType, Content: msgData.Content, KeyVersion: msgData.KeyVersion}}
	RpcResp, err := client.EncryptContent(context.Background(), &encryptContentReq)

	if err != nil {
		log.Error(operationID, "EncryptContent failed ", err.Error(), encryptContentReq.String())
		return utils.Wrap(err, ""), nil
	}
	if RpcResp.CommonResp.ErrCode != 0 {
		log.Error(operationID, "EncryptContent failed ", RpcResp.CommonResp, encryptContentReq.String())
		return errors.New(RpcResp.CommonResp.ErrMsg), nil
	}
	if len(RpcResp.MsgContentList) != len(recvIDList) {
		log.Error(operationID, "len(RpcResp.MsgContentList) != len(recvIDList) ", len(RpcResp.MsgContentList), len(recvIDList))
		return errors.New("len(RpcResp.MsgContentList) != len(recvIDList)"), nil
	}
	return nil, RpcResp.MsgContentList
}
