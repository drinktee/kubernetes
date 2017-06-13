package bcc

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/drinktee/bce-sdk-go/bce"
)

type Volume struct {
	Id            string             `json:"id"`
	Name          string             `json:"name"`
	DiskSizeInGB  int                `json:"diskSizeInGB"`
	PaymentTiming string             `json:"paymentTiming"`
	CreateTime    string             `json:"createTime"`
	ExpireTime    string             `json:"expireTime"`
	Status        VolumeStatus       `json:"status"`
	VolumeType    VolumeType         `json:"type"`
	StorageType   StorageType        `json:"storageType"`
	Desc          string             `json:"desc"`
	Attachments   []VolumeAttachment `json:"attachments"`
	ZoneName      string             `json:"zoneName"`
}

type VolumeStatus string

const (
	VOLUMESTATUS_CREATE             string = "Creating"
	VOLUMESTATUS_AVALIABLE          string = "Available"
	VOLUMESTATUS_ATTACHING          string = "Attaching"
	VOLUMESTATUS_NOTAVALIABLE       string = "NotAvailable"
	VOLUMESTATUS_INUSE              string = "InUse"
	VOLUMESTATUS_DETACHING          string = "Detaching"
	VOLUMESTATUS_DELETING           string = "Deleting"
	VOLUMESTATUS_DELETED            string = "Deleted"
	VOLUMESTATUS_SCALING            string = "Scaling"
	VOLUMESTATUS_EXPIRED            string = "Expired"
	VOLUMESTATUS_ERROR              string = "Error"
	VOLUMESTATUS_SNAPSHOTPROCESSING string = "SnapshotProcessing"
	VOLUMESTATUS_IMAGEPROCESSING    string = "ImageProcessing"
)

type VolumeType string

const (
	VOLUME_TYPE_SYSTEM    VolumeType = "System"
	VOLUME_TYPE_EPHEMERAL VolumeType = "Ephemeral"
	VOLUME_TYPE_CDS       VolumeType = "Cds"
)

type StorageType string

const (
	STORAGE_TYPE_STD1 StorageType = "std1"
	STORAGE_TYPE_HP1  StorageType = "hp1"
	STORAGE_TYPE_SATA StorageType = "sata"
	STORAGE_TYPE_SSD  StorageType = "ssd"
)

type VolumeAttachment struct {
	VolumeId   string `json:"volumeId"`
	InstanceId string `json:"instanceId"`
	// mount path
	Device string `json:"device"`
}

// DeleteVolume Delete a volume
func (c *Client) DeleteVolume(volumeId string) error {
	if volumeId == "" {
		return fmt.Errorf("DeleteVolume need a id")
	}
	req, err := bce.NewRequest("DELETE", c.GetURL("v2/volume"+"/"+volumeId, nil), nil)
	if err != nil {
		return err
	}
	_, err = c.SendRequest(req, nil)
	if err != nil {
		return err
	}
	return nil
}

type CreateVolumeArgs struct {
	PurchaseCount int          `json:"purchaseCount,omitempty"`
	CdsSizeInGB   int          `json:"cdsSizeInGB"`
	StorageType   StorageType  `json:"storageType"`
	Billing       *bce.Billing `json:"billing"`
	SnapshotId    string       `json:"snapshotId,omitempty"`
	ZoneName      string       `json:"zoneName,omitempty"`
}

type CreateVolumeResponse struct {
	VolumeIds []string `json:"volumeIds,omitempty"`
}

func (args *CreateVolumeArgs) validate() error {
	if args == nil {
		return fmt.Errorf("CreateVolumeArgs need args")
	}
	if args.StorageType == "" {
		return fmt.Errorf("CreateVolumeArgs need StorageType")
	}
	if args.Billing == nil {
		return fmt.Errorf("CreateVolumeArgs need Billing")
	}
	if args.CdsSizeInGB == 0 {
		return fmt.Errorf("CreateVolumeArgs need CdsSizeInGB")
	}
	return nil
}

// CreateVolumes create a volume
func (c *Client) CreateVolumes(args *CreateVolumeArgs) ([]string, error) {
	err := args.validate()
	if err != nil {
		return nil, err
	}
	params := map[string]string{
		"clientToken": c.GenerateClientToken(),
	}
	postContent, err := json.Marshal(args)
	if err != nil {
		return nil, err
	}
	req, err := bce.NewRequest("POST", c.GetURL("v2/volume", params), bytes.NewReader(postContent))
	if err != nil {
		return nil, err
	}
	resp, err := c.SendRequest(req, nil)
	if err != nil {
		return nil, err
	}
	bodyContent, err := resp.GetBodyContent()

	if err != nil {
		return nil, err
	}
	var blbsResp *CreateVolumeResponse
	err = json.Unmarshal(bodyContent, &blbsResp)

	if err != nil {
		return nil, err
	}
	return blbsResp.VolumeIds, nil

}

type GetVolumeListArgs struct {
	InstanceId string
	ZoneName   string
}
type GetVolumeListResponse struct {
	Volumes     []Volume `json:"volumes"`
	Marker      string   `json:"marker"`
	IsTruncated bool     `json:"isTruncated"`
	NextMarker  string   `json:"nextMarker"`
	MaxKeys     int      `json:"maxKeys"`
}

// GetVolumeList get all volumes
func (c *Client) GetVolumeList(args *GetVolumeListArgs) ([]Volume, error) {
	if args == nil {
		args = &GetVolumeListArgs{}
	}
	params := map[string]string{
		"zoneName":   args.ZoneName,
		"instanceId": args.InstanceId,
	}
	req, err := bce.NewRequest("GET", c.GetURL("v2/volume", params), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.SendRequest(req, nil)
	bodyContent, err := resp.GetBodyContent()

	if err != nil {
		return nil, err
	}
	var blbsResp *GetVolumeListResponse
	err = json.Unmarshal(bodyContent, &blbsResp)

	if err != nil {
		return nil, err
	}
	return blbsResp.Volumes, nil
}

type DescribeVolumeResponse struct {
	Volume *Volume `json:"volume"`
}

// DescribeVolume describe a volume
// More info see https://cloud.baidu.com/doc/BCC/API.html#.E6.9F.A5.E8.AF.A2.E7.A3.81.E7.9B.98.E8.AF.A6.E6.83.85
func (c *Client) DescribeVolume(id string) (*Volume, error) {
	if id == "" {
		return nil, fmt.Errorf("DescribeVolume need a id")
	}
	req, err := bce.NewRequest("GET", c.GetURL("v2/volume"+"/"+id, nil), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.SendRequest(req, nil)
	if err != nil {
		return nil, err
	}
	bodyContent, err := resp.GetBodyContent()
	if err != nil {
		return nil, err
	}
	var ins DescribeVolumeResponse
	err = json.Unmarshal(bodyContent, &ins)

	if err != nil {
		return nil, err
	}
	return ins.Volume, nil
}

// AttachCDSVolumeArgs describe attachcds args
type AttachCDSVolumeArgs struct {
	VolumeId   string `json:"-"`
	InstanceId string `json:"instanceId"`
}
type AttachCDSVolumeResponse struct {
	VolumeAttachment *VolumeAttachment `json:"volumeAttachment"`
}

func (args *AttachCDSVolumeArgs) validate() error {
	if args == nil {
		return fmt.Errorf("AttachCDSVolumeArgs need args")
	}
	if args.VolumeId == "" {
		return fmt.Errorf("AttachCDSVolumeArgs need VolumeId")
	}
	if args.InstanceId == "" {
		return fmt.Errorf("AttachCDSVolumeArgs need InstanceId")
	}
	return nil
}

// AttachCDSVolume attach a cds to vm
func (c *Client) AttachCDSVolume(args *AttachCDSVolumeArgs) (*VolumeAttachment, error) {
	err := args.validate()
	if err != nil {
		return nil, err
	}
	params := map[string]string{
		"attach":      "",
		"clientToken": c.GenerateClientToken(),
	}
	postContent, err := json.Marshal(args)
	if err != nil {
		return nil, err
	}
	req, err := bce.NewRequest("PUT", c.GetURL("v2/volume"+"/"+args.VolumeId, params), bytes.NewReader(postContent))
	if err != nil {
		return nil, err
	}
	resp, err := c.SendRequest(req, nil)
	if err != nil {
		return nil, err
	}
	bodyContent, err := resp.GetBodyContent()
	if err != nil {
		return nil, err
	}
	var blbsResp AttachCDSVolumeResponse
	err = json.Unmarshal(bodyContent, &blbsResp)

	if err != nil {
		return nil, err
	}
	return blbsResp.VolumeAttachment, nil
}

// DetachCDSVolume detach a cds
// TODO: if a volume is ddetaching, need to wait
func (c *Client) DetachCDSVolume(args *AttachCDSVolumeArgs) error {
	err := args.validate()
	if err != nil {
		return err
	}
	params := map[string]string{
		"detach":      "",
		"clientToken": c.GenerateClientToken(),
	}
	postContent, err := json.Marshal(args)
	if err != nil {
		return err
	}
	req, err := bce.NewRequest("PUT", c.GetURL("v2/volume"+"/"+args.VolumeId, params), bytes.NewReader(postContent))
	if err != nil {
		return err
	}
	_, err = c.SendRequest(req, nil)
	if err != nil {
		return err
	}
	return nil
}

// DeleteCDS delete a cds
func (c *Client) DeleteCDS(volumeID string) error {
	if volumeID == "" {
		return fmt.Errorf("DeleteCDS need a volumeId")
	}
	params := map[string]string{
		"clientToken": c.GenerateClientToken(),
	}
	req, err := bce.NewRequest("DELETE", c.GetURL("v2/volume"+"/"+volumeID, params), nil)
	if err != nil {
		return err
	}
	_, err = c.SendRequest(req, nil)
	if err != nil {
		return err
	}
	return nil
}

// RollbackVolume rollback a volume
// TODO
func (c *Client) RollbackVolume() {

}

// PurchaseReservedVolume purchaseReserved a volume
// TODO
func (c *Client) PurchaseReservedVolume() {

}
