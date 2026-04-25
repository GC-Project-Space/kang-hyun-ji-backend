package db

import "log"

func Seed() {
	seedCreatures()
	seedAchievements()
	seedTestUser()
}

func seedTestUser() {
	var count int
	DB.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&count)
	if count > 0 {
		return
	}
	_, err := DB.Exec(`INSERT INTO users (username, current_depth) VALUES (?, ?)`, "test_user", 0)
	if err != nil {
		log.Printf("seed test user failed: %v", err)
	}
}

func seedCreatures() {
	var count int
	DB.QueryRow(`SELECT COUNT(*) FROM sea_creatures`).Scan(&count)
	if count > 0 {
		return
	}

	creatures := []struct {
		name        string
		description string
		minDepth    int
		maxDepth    int
		rarity      string
		weight      int
		imagePath   string
	}{
		{"니모 (Clownfish)", "주황색 몸에 흰 줄무늬가 특징인 친숙한 물고기입니다.", 0, 150, "COMMON", 50, "/assets/images/creatures/clownfish.png"},
		{"해파리 (Jellyfish)", "투명한 몸으로 수면 근처를 유영하며 몽환적인 분위기를 자아냅니다.", 0, 300, "COMMON", 45, "/assets/images/creatures/jellyfish.png"},
		{"거북이 (Sea Turtle)", "느긋하게 바다를 여행하며 일기를 지켜주는 바다의 장수 상징입니다.", 0, 400, "UNCOMMON", 30, "/assets/images/creatures/sea-turtle.png"},
		{"돌고래 (Dolphin)", "지능이 높고 매끄러운 몸을 가진 바다의 친구입니다.", 0, 250, "UNCOMMON", 25, "/assets/images/creatures/dolphin.png"},
		{"불가사리 (Starfish)", "바다의 별이라고 불리며 산호초 사이에서 일기를 기다립니다.", 50, 450, "COMMON", 40, "/assets/images/creatures/starfish.png"},
		{"문어 (Octopus)", "영리한 지능을 가졌으며 자유자재로 몸의 색을 바꿉니다.", 100, 550, "UNCOMMON", 20, "/assets/images/creatures/octopus.png"},
		{"오징어 (Squid)", "빠르게 헤엄치며 위기 상황에는 먹물을 뿜어냅니다.", 150, 650, "COMMON", 35, "/assets/images/creatures/squid.png"},
		{"해마 (Seahorse)", "꼬리를 해초에 감고 조용히 머무르는 작은 생물입니다.", 50, 350, "UNCOMMON", 20, "/assets/images/creatures/seahorse.png"},
		{"고래 (Whale)", "거대한 몸집으로 바다의 경이로움을 느끼게 해주는 생물입니다.", 100, 800, "RARE", 10, "/assets/images/creatures/whale.png"},
		{"가오리 (Manta Ray)", "넓은 지느러미를 휘저으며 바닷속을 비행하듯 헤엄칩니다.", 200, 700, "RARE", 12, "/assets/images/creatures/manta-ray.png"},
		{"게 (Crab)", "단단한 껍질과 집게발로 바닥을 기어 다니는 생물입니다.", 0, 500, "COMMON", 40, "/assets/images/creatures/crab.png"},
		{"새우 (Shrimp)", "투명하고 작은 몸으로 부지런히 바닷속을 청소합니다.", 0, 600, "COMMON", 45, "/assets/images/creatures/shrimp.png"},
		{"복어 (Pufferfish)", "위험을 느끼면 몸을 동그랗게 부풀리는 귀여운 물고기입니다.", 100, 450, "UNCOMMON", 25, "/assets/images/creatures/pufferfish.png"},
		{"초롱아귀 (Anglerfish)", "머리 위의 빛나는 발광기로 먹이를 유혹하는 심해의 사냥꾼입니다.", 650, 1000, "RARE", 8, "/assets/images/creatures/anglerfish.png"},
		{"흡혈오징어 (Vampire Squid)", "망토 같은 지느러미를 가진 어두운 심해의 신비로운 생물입니다.", 750, 1000, "RARE", 7, "/assets/images/creatures/vampire-squid.png"},
		{"대왕오징어 (Giant Squid)", "심해 깊은 곳에 숨어 사는 전설적인 크기의 오징어입니다.", 600, 950, "LEGENDARY", 3, "/assets/images/creatures/giant-squid.png"},
		{"심해 등각류 (Isopod)", "바닥을 기어 다니며 심해의 영양분을 섭취하는 갑각류입니다.", 800, 1000, "UNCOMMON", 15, "/assets/images/creatures/deep-sea-isopod.png"},
		{"덤보문어 (Dumbo Octopus)", "귀 같은 지느러미를 펄럭이며 헤엄치는 매우 귀여운 심해어입니다.", 700, 1000, "LEGENDARY", 4, "/assets/images/creatures/dumbo-octopus.png"},
		{"클리오네 (Sea Angel)", "천사 같은 날개짓으로 심해를 부유하는 투명한 생물입니다.", 400, 900, "RARE", 6, "/assets/images/creatures/sea-angel.png"},
		{"엔젤피쉬 (Angelfish)", "화려한 줄무늬와 삼각형 지느러미가 아름다운 열대 물고기입니다.", 0, 200, "UNCOMMON", 22, "/assets/images/creatures/angelfish.png"},
		{"랜턴피쉬 (Lanternfish)", "몸에서 빛을 내며 심해를 유영하는 작은 발광 물고기입니다.", 500, 900, "RARE", 9, "/assets/images/creatures/lanternfish.png"},
		{"잠수함 (Submarine)", "바다 깊은 곳에서 발견된 신비로운 철제 물체입니다.", 300, 1000, "LEGENDARY", 2, "/assets/images/creatures/submarine.png"},
	}

	stmt, err := DB.Prepare(`
		INSERT INTO sea_creatures (name, description, min_depth, max_depth, rarity, spawn_weight, image_url)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		log.Printf("prepare creature insert failed: %v", err)
		return
	}
	defer stmt.Close()

	for _, c := range creatures {
		if _, err := stmt.Exec(c.name, c.description, c.minDepth, c.maxDepth, c.rarity, c.weight, c.imagePath); err != nil {
			log.Printf("insert creature %s failed: %v", c.name, err)
		}
	}
	log.Printf("seeded %d sea creatures", len(creatures))
}

func seedAchievements() {
	var count int
	DB.QueryRow(`SELECT COUNT(*) FROM achievements`).Scan(&count)
	if count > 0 {
		return
	}

	achievements := []struct {
		key           string
		name          string
		description   string
		conditionType string
		conditionVal  int
	}{
		// 수집 업적
		{"first_collect", "첫 부화", "첫 번째 생명체를 부화시켰습니다.", "COLLECTION_COUNT", 1},
		{"collect_5", "꼬마 탐험가", "생명체 5마리를 수집했습니다.", "COLLECTION_COUNT", 5},
		{"collect_10", "바다 탐험가", "생명체 10마리를 수집했습니다.", "COLLECTION_COUNT", 10},
		{"collect_30", "심해 탐험가", "생명체 30마리를 수집했습니다.", "COLLECTION_COUNT", 30},
		// 구간 도감 완성
		{"zone_shallow", "얕은 바다의 수집가", "0~250m 구간 도감을 완성했습니다.", "ZONE_COMPLETE", 0},
		{"zone_midwater", "중간 바다의 수집가", "250~500m 구간 도감을 완성했습니다.", "ZONE_COMPLETE", 250},
		{"zone_deep", "깊은 바다의 수집가", "500~750m 구간 도감을 완성했습니다.", "ZONE_COMPLETE", 500},
		{"zone_abyss", "심해의 수집가", "750~1000m 구간 도감을 완성했습니다.", "ZONE_COMPLETE", 750},
		{"full_collection", "완전한 도감", "모든 생명체를 수집했습니다.", "FULL_COLLECTION", 0},
		// 연속 작성
		{"streak_3", "3일 연속", "3일 연속으로 일기를 작성했습니다.", "STREAK", 3},
		{"streak_7", "일주일 연속", "7일 연속으로 일기를 작성했습니다.", "STREAK", 7},
		{"streak_30", "한 달 연속", "30일 연속으로 일기를 작성했습니다.", "STREAK", 30},
		// 누적 일기
		{"diary_10", "기록하는 잠수함", "일기 10개를 작성했습니다.", "DIARY_COUNT", 10},
		{"diary_50", "성실한 잠수함", "일기 50개를 작성했습니다.", "DIARY_COUNT", 50},
		// 수심 도달
		{"depth_250", "첫 번째 구역 진입", "수심 250m에 도달했습니다.", "DEPTH_REACHED", 250},
		{"depth_500", "두 번째 구역 진입", "수심 500m에 도달했습니다.", "DEPTH_REACHED", 500},
		{"depth_750", "세 번째 구역 진입", "수심 750m에 도달했습니다.", "DEPTH_REACHED", 750},
		{"depth_1000", "심해 정복", "수심 1000m에 도달했습니다.", "DEPTH_REACHED", 1000},
	}

	stmt, err := DB.Prepare(`
		INSERT INTO achievements (key, name, description, condition_type, condition_value)
		VALUES (?, ?, ?, ?, ?)
	`)
	if err != nil {
		log.Printf("prepare achievement insert failed: %v", err)
		return
	}
	defer stmt.Close()

	for _, a := range achievements {
		if _, err := stmt.Exec(a.key, a.name, a.description, a.conditionType, a.conditionVal); err != nil {
			log.Printf("insert achievement %s failed: %v", a.key, err)
		}
	}
	log.Printf("seeded %d achievements", len(achievements))
}
