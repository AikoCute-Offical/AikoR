package dev

import "github.com/AikoCute-Offical/AikoR/api"

// ReportUserTraffic reports the user traffic
func (c *APIClient) ReportUserTraffic(userTraffic *[]api.UserTraffic) error {
	path := "/api/v1/server/Aiko/push"

	// json structure: {uid1: [u, d], uid2: [u, d], uid1: [u, d], uid3: [u, d]}
	data := make(map[int][]int64, len(*userTraffic))
	for _, traffic := range *userTraffic {
		data[traffic.UID] = []int64{traffic.Upload, traffic.Download}
	}

	res, err := c.client.R().SetBody(data).ForceContentType("application/json").Post(path)
	_, err = c.parseResponse(res, path, err)
	if err != nil {
		return err
	}

	return nil
}
