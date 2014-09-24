package main

// remove duplicates from a slice of strings
func uniq(a []string) []string {
        m := make(map[string]struct{})
        for _, s := range a {
                m[s] = struct{}{}
        }

        var result []string
        for key := range m {
                result = append(result, key)
        }
        return result
}
