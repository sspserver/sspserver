//
// @project GeniusRabbit rotator 2018 - 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018 - 2019
//

package models

// AdFileThumb of the file
type AdFileThumb struct {
	Path        string // Path to image or video
	Type        AdFileType
	Width       int
	Height      int
	ContentType string
}

// IsSuits thumb by size
func (th AdFileThumb) IsSuits(w, h, wmin, hmin int) bool {
	return th.Width <= w && th.Width >= wmin && th.Height <= h && th.Height >= hmin
}

// IsImage file type
func (th AdFileThumb) IsImage() bool {
	return th.Type.IsImage()
}

// IsVideo file type
func (th AdFileThumb) IsVideo() bool {
	return th.Type.IsVideo()
}

// AdFile information
type AdFile struct {
	ID          uint64
	Name        string
	Path        string // In case of HTML5, hare must be the path to directory on CDN
	Type        AdFileType
	ContentType string
	Width       int
	Height      int
	Thumbs      []AdFileThumb
}

// // AdFileByModel original
// func AdFileByModel(file *models.AdFile) *AdFile {
// 	var (
// 		newThumbs []AdFileThumb
// 		meta      = file.ObjectMeta()
// 	)

// 	// Prepare thumb list
// 	for _, thumb := range meta.Items {
// 		newThumbs = append(newThumbs, AdFileThumb{
// 			Path:        urlPathJoin(file.Path, thumb.Name),
// 			Type:        AdFileTypeByObjectType(thumb.Type),
// 			Width:       thumb.Width,
// 			Height:      thumb.Height,
// 			ContentType: thumb.ContentType,
// 		})
// 	}

// 	return &AdFile{
// 		ID:          file.ID,
// 		Name:        file.Name.String,
// 		Path:        urlPathJoin(file.Path, meta.Main.Name),
// 		Type:        AdFileTypeByObjectType(file.Type),
// 		ContentType: file.ContentType,
// 		Width:       meta.Main.Width,
// 		Height:      meta.Main.Height,
// 		Thumbs:      newThumbs,
// 	}
// }

// ThumbBy size borders and specific type
func (f *AdFile) ThumbBy(w, h, wmin, hmin int, adType AdFileType) (th *AdFileThumb) {
	for i := 0; i < len(f.Thumbs); i++ {
		if f.Thumbs[i].Type == adType && f.Thumbs[i].IsSuits(w, h, wmin, hmin) {
			return &f.Thumbs[i]
		}
	}
	return nil
}

// IsImage file type
func (f *AdFile) IsImage() bool {
	return f.Type.IsImage()
}

// IsVideo file type
func (f *AdFile) IsVideo() bool {
	return f.Type.IsVideo()
}

// func urlPathJoin(urlBase, name string) string {
// 	if strings.HasSuffix(urlBase, "/") != strings.HasPrefix(name, "/") {
// 		return urlBase + name
// 	}
// 	return strings.TrimRight(urlBase, "/") + "/" + strings.TrimLeft(name, "/")
// }
