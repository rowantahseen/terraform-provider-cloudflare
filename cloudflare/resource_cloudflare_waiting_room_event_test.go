package cloudflare

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCloudflareWaitingRoomEvent_Create(t *testing.T) {
	t.Parallel()
	zoneID := os.Getenv("CLOUDFLARE_ZONE_ID")
	rnd := generateRandomResourceName()
	waitingRoomID := generateRandomResourceName()
	name := fmt.Sprintf("cloudflare_waiting_room_event.%s", rnd)
	waitingRoomEventName := fmt.Sprintf("waiting_room_event_%s", rnd)
	eventStartTime := time.Now()
	eventEndTime := eventStartTime.Add(5 * time.Minute)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudflareWaitingRoomEventDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudflareWaitingRoomEvent(rnd, waitingRoomEventName, zoneID, waitingRoomID, eventStartTime, eventEndTime),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "zone_id", zoneID),
					resource.TestCheckResourceAttr(name, "name", waitingRoomEventName),
					resource.TestCheckResourceAttr(name, "waiting_room_id", waitingRoomID),
					resource.TestCheckResourceAttr(name, "event_start_time", eventStartTime.Format(time.RFC3339)),
					resource.TestCheckResourceAttr(name, "event_end_time", eventEndTime.Format(time.RFC3339)),
					resource.TestCheckResourceAttr(name, "description", "my desc"),
					resource.TestCheckResourceAttr(name, "custom_page_html", "foobar"),
					resource.TestCheckResourceAttr(name, "disable_session_renewal", "true"),
					resource.TestCheckResourceAttr(name, "suspended", "true"),
					resource.TestCheckResourceAttr(name, "queueing_method", "fifo"),
					resource.TestCheckResourceAttr(name, "new_users_per_minute", "400"),
					resource.TestCheckResourceAttr(name, "total_active_users", "405"),
					resource.TestCheckResourceAttr(name, "session_duration", "10"),
					resource.TestCheckResourceAttr(name, "shuffle_at_event_start", "false"),
					resource.TestCheckNoResourceAttr(name, "prequeue_start_time"),
				),
			},
		},
	})
}

func testAccCheckCloudflareWaitingRoomEventDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*cloudflare.API)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudflare_waiting_room_event" {
			continue
		}

		_, err := client.WaitingRoomEvent(context.Background(), rs.Primary.Attributes["zone_id"], rs.Primary.Attributes["waiting_room_id"], rs.Primary.Attributes["id"])
		if err == nil {
			return fmt.Errorf("waiting room event still exists")
		}
	}

	return nil
}

func testAccCloudflareWaitingRoomEvent(resourceName, waitingRoomEventName, zoneID, waitingRoomID string, startTime, endTime time.Time) string {
	return fmt.Sprintf(`
resource "cloudflare_waiting_room_event" "%[1]s" {
  name                    = "%[2]s"
  zone_id                 = "%[3]s"
  waiting_room_id         = "%[4]s"
  event_start_time        = "%[5]s"
  event_end_time          = "%[6]s"
  total_active_users      = 405
  new_users_per_minute    = 400
  custom_page_html        = "foobar"
  queueing_method         = "fifo"
  shuffle_at_event_start  = false
  disable_session_renewal = true
  suspended               = true
  description             = "my desc"
  session_duration        = 10
}
`, resourceName, waitingRoomEventName, zoneID, waitingRoomID, startTime.Format(time.RFC3339), endTime.Format(time.RFC3339))
}
