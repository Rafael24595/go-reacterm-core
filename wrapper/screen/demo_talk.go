package wrapper_screen

import (
	"time"

	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/model/chat"

	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/primitive/talk"
)

func NewDemoTalk() screen.Node {
	return talk.New().
		SetName("textinput - amet").
		SetOwner("human_001").
		AddMessage(
			chat.NewMessageWithTimestamp(
				"human_001",
				"Lorem ipsum dolor sit amet, consectetur adipiscing elit?",
				time.Now().Add(-15*time.Minute).UnixMilli(),
			),
			chat.NewMessageWithTimestamp(
				"human_002",
				"Morbi ac ex sit amet diam euismod vulputate ut eu leo.",
				time.Now().Add(-14*time.Minute).UnixMilli(),
			),
			chat.NewMessageWithTimestamp(
				"human_001",
				"Ok.",
				time.Now().Add(-12*time.Minute).UnixMilli(),
			),
			chat.NewMessageWithTimestamp(
				"human_001",
				"Nullam quis ante sodales, aliquet turpis ut, suscipit erat. Cras nec viverra dolor, non egestas erat. Vivamus ac pretium lectus. Proin id ligula scelerisque, condimentum elit sit amet, imperdiet magna.",
				time.Now().Add(-10*time.Minute).UnixMilli(),
			),
			chat.NewMessageWithTimestamp(
				"human_002",
				"Wow! Nunc imperdiet, turpis vel dictum pretium, sem nibh sodales est, nec pulvinar diam leo ac augue.",
				time.Now().Add(-8*time.Minute).UnixMilli(),
			),
			chat.NewMessageWithTimestamp(
				"human_002",
				"Quisque facilisis nisl nec ex feugiat, non tristique sem finibus.",
				time.Now().Add(-7*time.Minute).UnixMilli(),
			),
			chat.NewMessageWithTimestamp(
				"human_001",
				"Sed hendrerit elementum lorem, vel interdum velit. Vestibulum rhoncus rhoncus mi, in efficitur elit. Duis imperdiet dictum erat, vel laoreet lorem hendrerit eu.",
				time.Now().Add(-5*time.Minute).UnixMilli(),
			),
			chat.NewMessageWithTimestamp(
				"human_001",
				"Aenean lacinia porta dictum. Ut sed pulvinar purus, eget pretium tellus. In pretium finibus eros id pretium. Aliquam id interdum magna. Proin feugiat, turpis quis tincidunt elementum, neque justo efficitur elit, ac egestas ex lacus ac ante.",
				time.Now().Add(-4*time.Minute).UnixMilli(),
			),
			chat.NewMessageWithTimestamp(
				"human_002",
				"Entendido. Proin sollicitudin mi ac arcu dictum, eleifend varius tellus ultrices.",
				time.Now().Add(-2*time.Minute).UnixMilli(),
			),
			chat.NewMessageWithTimestamp(
				"human_001",
				"Donec id elit non mi porta gravida at eget metus. Nulli magna feugiat purus, ac porttitor elit sem id tellus. Aliquam erat volutpat.",
				time.Now().Add(-1*time.Minute).UnixMilli(),
			),
		).
		ToNode()
}
