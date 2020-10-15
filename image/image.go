//图片处理包
package image

import (
	"github.com/disintegration/imaging"
)

//ResizeImage 压缩图片
func ResizeImage(sourceImageSrc string, dstImageSrc string, width int, height int) error {
	image, err := imaging.Open(sourceImageSrc)
	if err != nil {
		return err
	}
	dstImage := imaging.Resize(image, width, height, imaging.Lanczos)
	// Save the resulting image as JPEG.
	err = imaging.Save(dstImage, dstImageSrc)
	if err != nil {
		return err
	}

	return nil
}