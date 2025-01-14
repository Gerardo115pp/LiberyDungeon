package video_moment_models

type Video struct {
	VideoUUID    string `json:"video_uuid"`
	VideoCluster string `json:"video_cluster"`
}

type VideoMoment struct {
	ID         int    `json:"id"`
	VideoUUID  string `json:"video_uuid"`
	MomentTime int    `json:"moment_time"`
}

type VideoMoments struct {
	Video
	Moments []VideoMoment `json:"moments"`
}
