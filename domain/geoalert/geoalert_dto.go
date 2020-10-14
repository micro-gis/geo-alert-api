package geoalerts

type GeoAlert struct {
	Id string `json:"id"`
	UserId int64 `json:"user_id"`
	Title string `json:"title"`
	Description string `json:"description"`
	//TODO : replace by real geom implementation
	Geometry GeoPoint `json:"geometry"`
	//TODO : validate scope by enum
	Scope string `json:"scope"`
	TargetsUsersId []int64 `json:"targets_users_id"`
}

type GeoPoint struct {
	Type string
	Coordinates []float64
}