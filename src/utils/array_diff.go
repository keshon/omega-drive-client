package utils

/*
	Array comparsion funcs to get the difference
	The implementation is taken from here: https://stackoverflow.com/a/45428032/10352443

	Innokentiy Sokolov
	https://github.com/keshon

	2022-03-24
*/

func StringArrayDifference(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}

	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}
