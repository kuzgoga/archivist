package export

import "archivist/pkg/pipeline"

type ShortItem struct {
	Name    string
	Summary string
}

type TopicGroup struct {
	Topic string
	Items []ShortItem
}

type TagGroup struct {
	Tag    string
	Topics []TopicGroup
}

func groupInOrder(items []pipeline.CompleteItem) []TagGroup {
	var result []TagGroup

	tagIndex := make(map[string]int)
	topicIndex := make(map[string]map[string]int)

	for _, it := range items {
		ti, tagExists := tagIndex[it.Tag]
		if !tagExists {
			ti = len(result)
			tagIndex[it.Tag] = ti
			result = append(result, TagGroup{
				Tag:    it.Tag,
				Topics: nil,
			})
			topicIndex[it.Tag] = make(map[string]int)
		}

		topicsMap := topicIndex[it.Tag]
		to, topicExists := topicsMap[it.Topic]
		if !topicExists {
			to = len(result[ti].Topics)
			topicsMap[it.Topic] = to
			result[ti].Topics = append(result[ti].Topics, TopicGroup{
				Topic: it.Topic,
				Items: nil,
			})
		}
		newItem := ShortItem{
			Name:    it.Name,
			Summary: it.Summary,
		}
		result[ti].Topics[to].Items = append(result[ti].Topics[to].Items, newItem)
	}
	return result
}
