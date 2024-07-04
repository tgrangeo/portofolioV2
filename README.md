# Portfolio V2

Portfolio V2 is a personal portfolio website built with Go (Golang) and Flutter, showcasing my projects and providing information about me.

## Features

- Displays a list of projects from GitHub repositories.
- Allows users to view project details and README files.
- Includes an "About Me" section with a biography and skills.

## Technologies Used

- **Backend:** Go (Golang), HTTP Server, GitHub API integration.
- **Frontend:** Htmx, HTML/CSS.
- **Markdown Parsing:** gomarkdown/markdown library.
- **Others:** GitHub Api, dotenv for environment variables.

## Setup

### Prerequisites

- Go (Golang) installed on your machine.
- GitHub account with personal access token for API access.

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/tgrangeo/portofolioV2.git
   cd portofolioV2
   ```

2. Set up environment variables:

   Create a `.env` file in the root directory and add the following:

   ```plaintext
   GITHUB_TOKEN=your_github_personal_access_token
   ```

   Replace `your_github_personal_access_token` with your actual GitHub personal access token.

3. Build and run the Go backend:

   ```bash
   ./bin/air
   ```

4. Access the application:

   Open a web browser and go to `http://localhost:3000` to view the portfolio.

## Usage

- **Projects:** Navigate to the "Projects" section to see a list of GitHub repositories.
- **Project Details:** Click on a project to view its details and README file.
- **About Me:** Visit the "About Me" section to learn more about the owner of the portfolio.

## Acknowledgments

- **gomarkdown/markdown:** Used for parsing Markdown files.
- **joho/godotenv:** Used for loading environment variables from `.env` files.
