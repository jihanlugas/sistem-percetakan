package photo

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/jihanlugas/sistem-percetakan/config"
	"github.com/jihanlugas/sistem-percetakan/constant"
	"github.com/jihanlugas/sistem-percetakan/model"
	"github.com/jihanlugas/sistem-percetakan/utils"
	"gorm.io/gorm"
	"strings"
)

type Repository interface {
	GetByID(conn *gorm.DB, id string) (tPhoto model.Photo, err error)
	Upload(conn *gorm.DB, base64Image string, refTable constant.RefTable) (tPhoto model.Photo, err error)
	Delete(conn *gorm.DB, id string) error
}

type repository struct {
}

func (r repository) GetByID(conn *gorm.DB, id string) (tPhoto model.Photo, err error) {
	return tPhoto, conn.Where("id = ?", id).First(&tPhoto).Error
}

func (r repository) Upload(conn *gorm.DB, base64Image string, refTable constant.RefTable) (tPhoto model.Photo, err error) {

	// Remove the base64 header (e.g., "data:image/png;base64,")
	if idx := strings.Index(base64Image, ","); idx != -1 {
		base64Image = base64Image[idx+1:]
	}

	// Decode the Base64 image
	data, err := base64.StdEncoding.DecodeString(base64Image)
	if err != nil {
		return tPhoto, err
	}

	img, format, err := utils.Base64ToImage(base64Image)
	if err != nil {
		return tPhoto, err
	}

	tPhotoinc, err := r.photoincGettouse(conn, refTable)
	if err != nil {
		return tPhoto, err
	}

	tPhotoID := utils.GetUniqueID()

	// save the image to a file
	err = r.saveLocal(fmt.Sprintf("%s/%s.%s", tPhotoinc.Folder, tPhotoID, format), data)
	if err != nil {
		return tPhoto, err
	}

	tPhoto = model.Photo{
		ID:          tPhotoID,
		ClientName:  tPhotoID,
		ServerName:  tPhotoID,
		RefTable:    string(refTable),
		Ext:         format,
		PhotoPath:   fmt.Sprintf("%s/%s.%s", tPhotoinc.Folder, tPhotoID, format),
		PhotoSize:   int64(len(data)),
		PhotoWidth:  int64(img.Bounds().Dx()),
		PhotoHeight: int64(img.Bounds().Dy()),
	}

	err = conn.Create(&tPhoto).Error

	return tPhoto, err
}

func (r repository) Delete(conn *gorm.DB, id string) error {
	var err error

	tPhoto, err := r.GetByID(conn, id)
	if err != nil {
		return err
	}

	// delete file resource
	err = r.deleteLocal(tPhoto.PhotoPath)
	if err != nil {
		return err
	}

	return conn.Delete(&tPhoto).Error
}

// everytime the func called add running + 1
func (r repository) photoincGettouse(conn *gorm.DB, refTable constant.RefTable) (tPhotoinc model.Photoinc, err error) {
	err = conn.Where("ref_table = ?", refTable).
		Order("folder_inc DESC").
		First(&tPhotoinc).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tPhotoinc, err = r.photoincNew(conn, refTable, 1)
		}
	} else {
		if tPhotoinc.Running >= config.PhotoincRunningLimit {
			tPhotoinc, err = r.photoincNew(conn, refTable, tPhotoinc.FolderInc+1)
		} else {
			err = r.photoincAddrunning(conn, tPhotoinc)
		}
	}

	return tPhotoinc, err
}

func (r repository) photoincNew(conn *gorm.DB, refTable constant.RefTable, folderInc int64) (tPhotoinc model.Photoinc, err error) {
	tPhotoinc = model.Photoinc{
		RefTable:  string(refTable),
		FolderInc: folderInc,
		Folder:    fmt.Sprintf("%s/%s/%s/%d", config.StorageDirectory, config.PhotoDirectory, refTable, folderInc),
		Running:   1,
	}

	err = r.createPhotoInc(conn, tPhotoinc)
	if err != nil {
		return tPhotoinc, err
	}

	err = utils.CreateFolder(tPhotoinc.Folder, 0777)
	if err != nil {
		return tPhotoinc, err
	}

	return tPhotoinc, err
}

func (r repository) photoincAddrunning(conn *gorm.DB, tPhotoinc model.Photoinc) error {
	tPhotoinc.Running = tPhotoinc.Running + 1
	return conn.Save(&tPhotoinc).Error
}

func (r repository) createPhotoInc(conn *gorm.DB, tPhotoinc model.Photoinc) error {
	return conn.Create(&tPhotoinc).Error
}

func (r repository) saveLocal(filepath string, data []byte) error {
	return utils.SaveFileLocal(filepath, data)
}

func (r repository) deleteLocal(filepath string) error {
	return utils.DeleteFileLocal(filepath)
}

func (r repository) saveAws() error {
	return errors.New("not implemented")
}

func NewRepository() Repository {
	return repository{}
}
