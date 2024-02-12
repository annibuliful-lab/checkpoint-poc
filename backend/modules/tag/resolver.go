package tag

type TagResolver struct{}

var tagService = TagService{}

func (TagResolver) GetTags(data GetTagsInput) ([]Tag, error) {
	tags, _ := tagService.FindMany(data)
	return *tags, nil
}
