package wrapper_screen

import (
	"time"

	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/model/chat"

	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/primitive/talk"
)

func NewTestTalk() screen.Node {
	return talk.New().
		SetName("textinput - amet").
		SetOwner("human_001").
		AddMessage(
			chat.Message{
				Time:    time.Now().Add(-15 * time.Minute).UnixMilli(),
				Owner:   "human_001",
				Message: "Lorem ipsum dolor sit amet, consectetur adipiscing elit?",
			},
			chat.Message{
				Time:    time.Now().Add(-14 * time.Minute).UnixMilli(),
				Owner:   "human_002",
				Message: "Morbi ac ex sit amet diam euismod vulputate ut eu leo.",
			},
			chat.Message{
				Time:    time.Now().Add(-12 * time.Minute).UnixMilli(),
				Owner:   "human_001",
				Message: "Ok.",
			},
			chat.Message{
				Time:    time.Now().Add(-10 * time.Minute).UnixMilli(),
				Owner:   "human_001",
				Message: "Nullam quis ante sodales, aliquet turpis ut, suscipit erat. Cras nec viverra dolor, non egestas erat. Vivamus ac pretium lectus. Proin id ligula scelerisque, condimentum elit sit amet, imperdiet magna.",
			},
			chat.Message{
				Time:    time.Now().Add(-8 * time.Minute).UnixMilli(),
				Owner:   "human_002",
				Message: "Wow! Nunc imperdiet, turpis vel dictum pretium, sem nibh sodales est, nec pulvinar diam leo ac augue.",
			},
			chat.Message{
				Time:    time.Now().Add(-7 * time.Minute).UnixMilli(),
				Owner:   "human_002",
				Message: "Quisque facilisis nisl nec ex feugiat, non tristique sem finibus.",
			},
			chat.Message{
				Time:    time.Now().Add(-5 * time.Minute).UnixMilli(),
				Owner:   "human_001",
				Message: "Sed hendrerit elementum lorem, vel interdum velit. Vestibulum rhoncus rhoncus mi, in efficitur elit. Duis imperdiet dictum erat, vel laoreet lorem hendrerit eu.",
			},
			chat.Message{
				Time:    time.Now().Add(-4 * time.Minute).UnixMilli(),
				Owner:   "human_001",
				Message: "Aenean lacinia porta dictum. Ut sed pulvinar purus, eget pretium tellus. In pretium finibus eros id pretium. Aliquam id interdum magna. Proin feugiat, turpis quis tincidunt elementum, neque justo efficitur elit, ac egestas ex lacus ac ante.",
			},
			chat.Message{
				Time:    time.Now().Add(-2 * time.Minute).UnixMilli(),
				Owner:   "human_002",
				Message: "Entendido. Proin sollicitudin mi ac arcu dictum, eleifend varius tellus ultrices.",
			},
			chat.Message{
				Time:    time.Now().Add(-1 * time.Minute).UnixMilli(),
				Owner:   "human_001",
				Message: "Donec id elit non mi porta gravida at eget metus. Nulli magna feugiat purus, ac porttitor elit sem id tellus. Aliquam erat volutpat.",
			},
		).
		ToNode()
}
