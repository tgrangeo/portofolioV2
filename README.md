**🚀 Portfolio V2**

Portfolio V2 is a personal portfolio website built with Go (Golang) and Flutter, designed to showcase my projects and provide details about me.
🛠️ Roadmap

Add GitHub repo quick link
Home page: Select GitHub profile
Improve link color

    Make project page responsive

**✨ Features**

    📂 Displays a list of projects from GitHub repositories.
    📝 Allows users to view project details and README files.
    👤 Includes an "About Me" section with biography and skills.


    Next: 
        adding ai auto generate resume from last readmes to know better him (cool in case the user don t use it s readme)

**🔧 Technologies Used**

    Backend: Go (Golang), HTTP Server, GitHub API integration.
    Frontend: Htmx, HTML/CSS.
    Markdown Parsing: gomarkdown/markdown library.
    Others: GitHub API, dotenv for environment variables.

**⚙️ Setup**
📝 Prerequisites

    Go (Golang) installed on your machine.
    GitHub account with a personal access token for API access.

**🛠️ Installation**

    Clone the repository:

    bash

git clone https://github.com/tgrangeo/portofolioV2.git
cd portofolioV2

Set up environment variables:

Create a .env file in the root directory and add the following:

plaintext

GITHUB_TOKEN=your_github_personal_access_token

Replace your_github_personal_access_token with your actual GitHub personal access token.

Build and run the Go backend:

bash

    ./bin/air

    Access the application:

    Open a web browser and go to http://localhost:3000 to view the portfolio.

**🚀 Usage**

    Projects: Navigate to the "Projects" section to see a list of GitHub repositories.
    Project Details: Click on a project to view its details and README file.
    About Me: Visit the "About Me" section to learn more about the portfolio owner.

**🙏 Acknowledgments**

    gomarkdown/markdown: For Markdown parsing.
    joho/godotenv: For loading environment variables from .env files.