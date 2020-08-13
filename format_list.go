package youtube

type FormatList []Format

func (list FormatList) FindByQuality(quality string) *Format {
	for i := range list {
		if list[i].Quality == quality {
			return &list[i]
		}
	}
	return nil
}

func (list FormatList) FindByItag(itagNo int) *Format {
	for i := range list {
		if list[i].ItagNo == itagNo {
			return &list[i]
		}
	}
	return nil
}
