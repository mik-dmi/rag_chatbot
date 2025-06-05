package db

import (
	"context"
	"fmt"
	"log"

	"github.com/mik-dmi/rag_chatbot/backend/internal/store"
)

var usernames = []string{
	"alpha_wolf", "beta_bear", "charlie_fox", "delta_eagle", "echo_lion",
	"foxtrot_snake", "golf_panther", "hotel_hawk", "india_tiger", "juliet_shark",
	"kilo_cobra", "lima_wolf", "mike_raven", "november_lynx", "oscar_ferret",
	"papa_otter", "quebec_jaguar", "romeo_gecko", "sierra_puma", "tango_badger",
	"uniform_pelican", "victor_mongoose", "whiskey_owl", "xray_leopard", "yankee_cougar",
	"zulu_falcon", "ninja_squirrel", "cyber_penguin", "ghost_hyena", "blaze_dragon",
	"storm_rider", "pixel_pirate", "shadow_puma", "frost_sniper", "ember_knight",
	"lava_ghost", "jungle_scout", "venom_scythe", "lunar_archer", "solar_blade",
	"phantom_falcon", "arctic_raven", "electric_zebra", "neon_bison", "omega_spider",
	"rapid_iguana", "vortex_gryphon", "iron_fang", "steel_talon", "crimson_mantis",
}

func Seed(postgresStore store.PostgreStorage) {
	ctx := context.Background()

	users := generateUsers(100)

	for _, user := range users {
		if err := postgresStore.Users.CreateUser(ctx, user); err != nil {
			log.Println("Error creating user:", err)
			return
		}
	}

	log.Println("Seeding complete")

}

func generateUsers(num int) []*store.PostgreUser {
	users := make([]*store.PostgreUser, num)

	for i := 0; i < num; i++ {
		users[i] = &store.PostgreUser{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@example.com",
		}
	}
	return users
}
