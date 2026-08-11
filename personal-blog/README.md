# Personal Blog

A web application built with Go standard library (`net/http` and `html/template`) to publish and manage personal blog articles.

This project is built to complete the **[Personal Blog](https://roadmap.sh/projects/personal-blog)** from **roadmap.sh**.

## Features

- **Guest Section**:
  - Home Page (`/home`): View list of published articles.
  - Article Page (`/article/:id`): Read full content of a selected article.
- **Admin Section** (Protected with HTTP Basic Authentication):
  - Dashboard (`/admin`): Overview of all articles with options to edit or delete.
  - Add Article Page (`/new`): Form to create and publish new articles.
  - Edit Article Page (`/edit/:id`): Form to update existing articles.
- Local JSON file-based storage inside the `articles/` directory.
- Pre-seeded mockup data matching the official roadmap.sh wireframe layout.
- Built strictly with the **Go Standard Library** (zero external dependencies).

## How to Run

1. Open your terminal and navigate to this folder:

```bash
go run main.go
# open localhost in your browser and admin
```