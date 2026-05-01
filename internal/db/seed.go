package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"sync"

	"github.com/saleemlawal/social/internal/store"
)

// Keep each DB op from queueing longer than store.QUERY_TIMEOUT_DURATION (waits for a conn count toward that deadline).
const seedParallelComments = 20

var tags = []string{
	"technology",
	"programming",
	"development",
	"software",
	"engineering",
	"design",
	"ux",
}

var commentsContent = []string{
	"Great post!",
	"Super duper OMG... So cool COOL!",
	"Totally agree with this.",
	"Interesting take!",
	"This changed my perspective.",
	"Love this content.",
	"Could not agree more.",
	"This is fire!",
	"Well said.",
	"Learned something new today.",
	"Absolutely brilliant!",
	"Never thought about it this way before.",
	"Sharing this with everyone I know.",
	"This deserves way more attention.",
	"You always deliver quality content.",
	"Mind blown honestly.",
	"Finally someone said it.",
	"This is exactly what I needed to read today.",
	"Saved this for later, so good.",
	"Keep up the great work!",
	"Underrated post right here.",
	"100% accurate.",
	"Facts only in this post.",
	"This hit different.",
	"I was literally just thinking about this.",
	"Short but powerful.",
	"This needs to be at the top.",
	"Couldn't have said it better myself.",
	"Genuinely helpful, thank you.",
	"Why is nobody talking about this more?",
	"This is the content I come here for.",
	"Dropping knowledge as always.",
	"Real talk right here.",
	"Came for the title, stayed for the content.",
	"Bold take, but I respect it.",
	"First time commenting but had to say something.",
	"This resonated deeply.",
	"Exactly my experience too.",
	"Wow, I did not expect that ending.",
	"You just described my life.",
}

func Seed(s store.Storage, db *sql.DB) {
	wg := sync.WaitGroup{}
	ctx := context.Background()

	users := generateUsers(100)
	tx, _ := db.BeginTx(ctx, nil)

	for _, user := range users {
		wg.Add(1)
		go func(user *store.User) {
			defer wg.Done()
			if err := s.Users.Create(ctx, tx, user); err != nil {
				_ = tx.Rollback()
				log.Println(err)
				return
			}
		}(user)
	}
	wg.Wait()
	tx.Commit()

	posts := generatePosts(200, users)
	for _, post := range posts {
		wg.Add(1)
		go func(post *store.Post) {
			defer wg.Done()
			if err := s.Posts.Create(ctx, post); err != nil {
				log.Println(err)
				return
			}
		}(post)
	}

	wg.Wait()

	comments := generateComments(10000, posts, users)
	sem := make(chan struct{}, seedParallelComments)
	for _, comment := range comments {
		wg.Add(1)
		sem <- struct{}{}
		go func(comment *store.Comment) {
			defer wg.Done()
			defer func() { <-sem }()
			if err := s.Comments.Create(ctx, comment); err != nil {
				log.Println(err)
			}
		}(comment)
	}
	wg.Wait()
	log.Printf("Seeded database with %d users, %d posts, and %d comments\n", len(users), len(posts), len(comments))
}

func generateComments(count int, posts []*store.Post, users []*store.User) []*store.Comment {
	comments := make([]*store.Comment, count)
	for i := range count {
		comments[i] = &store.Comment{
			PostID:  posts[rand.Intn(len(posts))].ID,
			UserID:  users[rand.Intn(len(users))].ID,
			Content: commentsContent[rand.Intn(len(commentsContent))],
			Likes:   rand.Intn(100),
		}
	}
	return comments
}

func generatePosts(count int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, count)
	for i := range count {
		posts[i] = &store.Post{
			Title:   fmt.Sprintf("Post %d", i),
			Content: fmt.Sprintf("Content %d", i),
			UserID:  users[rand.Intn(len(users))].ID,
			Tags: []string{
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))],
			},
		}
	}
	return posts
}

func generateUsers(count int) []*store.User {
	users := make([]*store.User, count)
	for i := range count {
		u := &store.User{
			Username: fmt.Sprintf("user%d", i),
			Email:    fmt.Sprintf("user%d@example.com", i),
			Role:     store.Role{Name: "User"},
		}
		if err := u.Password.Set("password"); err != nil {
			log.Fatalf("seed: user password: %v", err)
		}
		users[i] = u
	}
	return users
}
