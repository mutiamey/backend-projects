package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Article defines the structure for blog posts stored in JSON files
type Article struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Date    string `json:"date"`
	Content string `json:"content"`
}

// PageData passes structured data into the HTML template renderer
type PageData struct {
	Title    string
	Articles []Article
	Article  Article
}

const articlesDir = "articles"

// HTML template with styles matching the wireframe designs
const htmlTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <style>
        body { 
            font-family: system-ui, -apple-system, sans-serif; 
            max-width: 480px; 
            margin: 40px auto; 
            padding: 24px; 
            border: 2px solid #000; 
            border-radius: 12px; 
            background: #fff; 
        }
        h1 { margin-top: 0; font-size: 28px; font-weight: 800; }
        .article-list { list-style: none; padding: 0; margin: 20px 0 0 0; }
        .article-item { 
            display: flex; 
            justify-content: space-between; 
            align-items: center; 
            margin-bottom: 14px; 
        }
        .article-item a { 
            text-decoration: none; 
            color: #000; 
            font-weight: 700; 
            font-size: 18px; 
        }
        .date { color: #888; font-size: 16px; font-weight: 600; }
        .btn { 
            text-decoration: none; 
            padding: 8px 18px; 
            border: 2px solid #000; 
            border-radius: 6px; 
            color: #000; 
            font-weight: bold; 
            font-size: 16px; 
            cursor: pointer; 
            background: #fff; 
            display: inline-block; 
            margin-top: 10px;
        }
        .btn:hover { background: #f0f0f0; }
        .form-group { margin-bottom: 16px; }
        input[type="text"], textarea { 
            width: 100%; 
            padding: 10px; 
            box-sizing: border-box; 
            border: 2px solid #000; 
            border-radius: 6px; 
            font-size: 15px; 
            font-weight: 600;
        }
        textarea { height: 200px; font-family: inherit; }
        .top-bar { display: flex; justify-content: space-between; align-items: center; }
        .add-link { text-decoration: none; color: #000; font-weight: bold; font-size: 16px; }
        .action-link { text-decoration: none; color: #888; font-size: 15px; font-weight: 600; margin-left: 10px; }
        .action-link:hover { color: #000; }
    </style>
</head>
<body>
    {{if eq .Title "Home"}}
        <h1>Personal Blog</h1>
        <ul class="article-list">
            {{range .Articles}}
                <li class="article-item">
                    <a href="/article/{{.ID}}">{{.Title}}</a>
                    <span class="date">{{.Date}}</span>
                </li>
            {{end}}
        </ul>
    {{else if eq .Title "Article"}}
        <h1>{{.Article.Title}}</h1>
        <div class="date" style="margin-bottom: 20px;">{{.Article.Date}}</div>
        <p style="line-height: 1.5; white-space: pre-wrap; font-size: 16px; font-weight: 600; color: #000;">{{.Article.Content}}</p>
    {{else if eq .Title "Admin Dashboard"}}
        <div class="top-bar">
            <h1>Personal Blog</h1>
            <a href="/new" class="add-link">+ Add</a>
        </div>
        <ul class="article-list">
            {{range .Articles}}
                <li class="article-item">
                    <span style="font-size: 18px; font-weight: 700;">{{.Title}}</span>
                    <div>
                        <a href="/edit/{{.ID}}" class="action-link">Edit</a>
                        <a href="/delete/{{.ID}}" class="action-link" onclick="return confirm('Delete article?')">Delete</a>
                    </div>
                </li>
            {{end}}
        </ul>
    {{else if eq .Title "New Article"}}
        <h1>New Article</h1>
        <form action="/new" method="POST">
            <div class="form-group">
                <input type="text" name="title" placeholder="Article Title" required>
            </div>
            <div class="form-group">
                <input type="text" name="date" placeholder="Publishing Date" required value="{{.Article.Date}}">
            </div>
            <div class="form-group">
                <textarea name="content" placeholder="Content" required></textarea>
            </div>
            <button type="submit" class="btn">Publish</button>
        </form>
    {{else if eq .Title "Update Article"}}
        <h1>Update Article</h1>
        <form action="/edit/{{.Article.ID}}" method="POST">
            <div class="form-group">
                <input type="text" name="title" value="{{.Article.Title}}" placeholder="Article Title" required>
            </div>
            <div class="form-group">
                <input type="text" name="date" value="{{.Article.Date}}" placeholder="Publishing Date" required>
            </div>
            <div class="form-group">
                <textarea name="content" placeholder="Content" required>{{.Article.Content}}</textarea>
            </div>
            <button type="submit" class="btn">Update</button>
        </form>
    {{end}}
</body>
</html>
`

func main() {
	// Initialize local folder and seed dummy data matching the wireframe mockup
	os.MkdirAll(articlesDir, 0755)
	seedInitialArticles()

	// Guest section routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	})
	http.HandleFunc("/home", handleHome)
	http.HandleFunc("/article/", handleArticle)

	// Admin section routes (protected with HTTP basic auth)
	http.HandleFunc("/admin", basicAuth(handleAdmin))
	http.HandleFunc("/new", basicAuth(handleNewArticle))
	http.HandleFunc("/edit/", basicAuth(handleEditArticle))
	http.HandleFunc("/delete/", basicAuth(handleDeleteArticle))

	fmt.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

// basicAuth protects administrative pages requiring credentials
func basicAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok || username != "admin" || password != "password" {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted Admin Section"`)
			http.Error(w, "Unauthorized access", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	articles := loadAllArticles()
	renderTemplate(w, PageData{Title: "Home", Articles: articles})
}

func handleArticle(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/article/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	article, err := loadArticleByID(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	renderTemplate(w, PageData{Title: "Article", Article: article})
}

func handleAdmin(w http.ResponseWriter, r *http.Request) {
	articles := loadAllArticles()
	renderTemplate(w, PageData{Title: "Admin Dashboard", Articles: articles})
}

func handleNewArticle(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		today := time.Now().Format("2006-01-02")
		renderTemplate(w, PageData{Title: "New Article", Article: Article{Date: today}})
		return
	}

	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		date := r.FormValue("date")
		content := r.FormValue("content")

		articles := loadAllArticles()

		// Cari ID terbesar di antara seluruh artikel agar tidak ada ID yang bertabrakan/tertimpa
		maxID := 0
		for _, a := range articles {
			if a.ID > maxID {
				maxID = a.ID
			}
		}
		newID := maxID + 1

		article := Article{
			ID:      newID,
			Title:   title,
			Date:    formatDate(date),
			Content: content,
		}

		saveArticle(article)
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	}
}

func handleEditArticle(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/edit/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if r.Method == http.MethodGet {
		article, err := loadArticleByID(id)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		renderTemplate(w, PageData{Title: "Update Article", Article: article})
		return
	}

	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		date := r.FormValue("date")
		content := r.FormValue("content")

		article := Article{
			ID:      id,
			Title:   title,
			Date:    date,
			Content: content,
		}

		saveArticle(article)
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	}
}

func handleDeleteArticle(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/delete/")
	filePath := filepath.Join(articlesDir, fmt.Sprintf("%s.json", idStr))
	os.Remove(filePath)
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func renderTemplate(w http.ResponseWriter, data PageData) {
	tmpl := template.Must(template.New("blog").Parse(htmlTemplate))
	tmpl.Execute(w, data)
}

func saveArticle(article Article) {
	filePath := filepath.Join(articlesDir, fmt.Sprintf("%d.json", article.ID))
	data, _ := json.MarshalIndent(article, "", "  ")
	os.WriteFile(filePath, data, 0644)
}

func loadArticleByID(id int) (Article, error) {
	filePath := filepath.Join(articlesDir, fmt.Sprintf("%d.json", id))
	data, err := os.ReadFile(filePath)
	if err != nil {
		return Article{}, err
	}

	var article Article
	json.Unmarshal(data, &article)
	return article, nil
}

func loadAllArticles() []Article {
	files, err := os.ReadDir(articlesDir)
	if err != nil {
		return []Article{}
	}

	var articles []Article
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			idStr := strings.TrimSuffix(file.Name(), ".json")
			id, err := strconv.Atoi(idStr)
			if err == nil {
				article, err := loadArticleByID(id)
				if err == nil {
					articles = append(articles, article)
				}
			}
		}
	}

	// Sort articles ascending by ID (ID 1 first, ID 9 last - matching the wireframe)
	for i := 0; i < len(articles); i++ {
		for j := i + 1; j < len(articles); j++ {
			if articles[i].ID > articles[j].ID { // Ubah '<' menjadi '>'
				articles[i], articles[j] = articles[j], articles[i]
			}
		}
	}

	return articles
}

// Seed initial mockup data matching roadmap.sh wireframe exactly
func seedInitialArticles() {
	files, _ := os.ReadDir(articlesDir)
	if len(files) > 0 {
		return
	}

	dummies := []Article{
		{1, "My first article", "August 7, 2024", "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."},
		{2, "Second article", "August 4, 2024", "Lorem ipsum dolor sit amet, consectetur adipiscing elit."},
		{3, "Third article", "August 1, 2024", "Lorem ipsum dolor sit amet, consectetur adipiscing elit."},
		{4, "Fourth article", "July 30, 2024", "Lorem ipsum dolor sit amet, consectetur adipiscing elit."},
		{5, "Fifth article", "July 21, 2024", "Lorem ipsum dolor sit amet, consectetur adipiscing elit."},
		{6, "Sixth article", "July 15, 2024", "Lorem ipsum dolor sit amet, consectetur adipiscing elit."},
		{7, "Seventh article", "July 8, 2024", "Lorem ipsum dolor sit amet, consectetur adipiscing elit."},
		{8, "Eighth article", "July 4, 2024", "Lorem ipsum dolor sit amet, consectetur adipiscing elit."},
		{9, "Ninth Aritcle", "July 1, 2024", "Lorem ipsum dolor sit amet, consectetur adipiscing elit."},
	}

	for _, d := range dummies {
		saveArticle(d)
	}
}

// formatDate converts input date string (YYYY-MM-DD) into readable format (e.g. August 11, 2026)
func formatDate(inputDate string) string {
	t, err := time.Parse("2006-01-02", inputDate)
	if err != nil {
		return inputDate
	}
	return t.Format("August 2, 2006")
}
