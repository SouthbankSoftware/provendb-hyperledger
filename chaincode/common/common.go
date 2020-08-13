/*
 * @Author: guiguan
 * @Date:   2020-08-12T15:30:10+10:00
 * @Last modified by:   guiguan
 * @Last modified time: 2020-08-13T12:56:20+10:00
 */

package common

import "time"

const (
	// StateKeyData is the state key used for embedding data in Hyperledger Fabric
	StateKeyData = "data"
)

// EmbedDataReply represents the reply for the `EmbedData` function
type EmbedDataReply struct {
	TxnID      string    `json:"txnId"`
	CreateTime time.Time `json:"createTime"`
}
