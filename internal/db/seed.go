package db

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"

	"github.com/saleemlawal/social/internal/store"
)

var tags = []string{
	"technology",
	"programming",
	"development",
	"software",
	"engineering",
	"design",
	"ux",
}

func Seed(s store.Storage) {
	wg := sync.WaitGroup{}
	ctx := context.Background()

	users := generateUsers(100)

	for _, user := range users {
		wg.Add(1)
		go func(user *store.User) {
			defer wg.Done()
			if err := s.Users.Create(ctx, user); err != nil {
				log.Println(err)
				return
			}
		}(user)
	}
	wg.Wait()

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

	comments := generateComments(500, posts, users)
	for _, comment := range comments {
		wg.Add(1)
		go func(comment *store.Comment) {
			defer wg.Done()
			if err := s.Comments.Create(ctx, comment); err != nil {
				log.Println(err)
				return
			}
		}(comment)
	}
	wg.Wait()
	log.Println("Seeded database with 100 users, 200 posts, and 500 comments")
}

func generateComments(count int, posts []*store.Post, users []*store.User) []*store.Comment {
	comments := make([]*store.Comment, count)
	for i := range count {
		comments[i] = &store.Comment{
			PostID: posts[rand.Intn(len(posts))].ID,
			UserID: users[rand.Intn(len(users))].ID,
			Content: fmt.Sprintf("Comment %d", i),
			Likes: rand.Intn(100),
		}
	}
	return comments
}


func generatePosts(count int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, count)
	for i := range count {
		posts[i] = &store.Post{
			Title: fmt.Sprintf("Post %d", i),
			Content: fmt.Sprintf("Content %d", i),
			UserID: users[rand.Intn(len(users))].ID,
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
		users[i] = &store.User{
			Username: fmt.Sprintf("user%d", i),
			Email: fmt.Sprintf("user%d@example.com", i),
			Password: "password",
		}
	}
	return users
}