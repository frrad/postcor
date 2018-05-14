package client

import (
	"fmt"
)

func (c *PClient) FavoriteFEWorkouts() (string, error) {
	uid := c.GetUserId()

	startDate, endDate := "2018-05-01", "2018-05-13"

	url := fmt.Sprintf("https://na.preva.com/exerciser-api/exerciser-account/id/%d/fitness-activity?local-start-date=%s&local-end-date=%s", uid, startDate, endDate)

	fmt.Println(url)

	return c.GetPage(url)
}
