package gostatsd

// Set is used for storing aggregated values for sets.
type Set struct {
	Values    map[string]struct{}
	Timestamp Nanotime // Last time value was updated
	Hostname  string   // Hostname of the source of the metric
	Tags      Tags     // The tags for the set
}

// NewSet initialises a new set.
func NewSet(timestamp Nanotime, values map[string]struct{}, hostname string, tags Tags) Set {
	return Set{Values: values, Timestamp: timestamp, Hostname: hostname, Tags: tags.Copy()}
}

// Sets stores a map of sets by tags.
type Sets map[string]map[string]Set

// MetricsName returns the name of the aggregated metrics collection.
func (s Sets) MetricsName() string {
	return "Sets"
}

// Delete deletes the metrics from the collection.
func (s Sets) Delete(k string) {
	delete(s, k)
}

// DeleteChild deletes the metrics from the collection for the given tags.
func (s Sets) DeleteChild(k, t string) {
	delete(s[k], t)
}

// HasChildren returns whether there are more children nested under the key.
func (s Sets) HasChildren(k string) bool {
	return len(s[k]) != 0
}

// Each iterates over each set.
func (s Sets) Each(f func(string, string, Set)) {
	for key, value := range s {
		for tags, set := range value {
			f(key, tags, set)
		}
	}
}

// copy returns a deep-ish copy of this Sets object
func (s Sets) copy() Sets {
	sNew := make(Sets, len(s))
	for name, value := range s {
		m := make(map[string]Set, len(value))
		for tagsKey, set := range value {
			m[tagsKey] = Set{
				Values: make(map[string]struct{}),
				// Timestamp: counter.Timestamp, // Timestamp is not required
				Hostname: set.Hostname,
				Tags:     set.Tags, // shallow copy is acceptable
			}
			for key := range set.Values {
				m[tagsKey].Values[key] = struct{}{}
			}
		}
		sNew[name] = m
	}
	return sNew
}
