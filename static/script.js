  //mobile menu side
  function openNav() {
    document.getElementById("menu").style.display = "flex";
  }
  function closeNav() {
    if (window.innerWidth <= 480) {
      document.getElementById("menu").style.display = "none";
    }
  }
  //mobile project side
  function openProjectSide() {
    let state = document.getElementById("side").style.display;
    if (state === "flex")
      document.getElementById("side").style.display = "none";
    else document.getElementById("side").style.display = "flex";
  }
  function closeProjectSide() {
    if (window.innerWidth <= 480) {
      document.getElementById("side").style.display = "none";
    }
  }
  //theme
  const toggleTheme = () => {
    const html = document.querySelector("html");
    const button = document.getElementById("themeButton");
    const isLight = html.classList.contains("light-theme");
    if (isLight) {
      html.classList.remove("light-theme");
      html.classList.add("dark-theme");
      button.textContent = "🌙";
      localStorage.setItem("theme", "dark");
    } else {
      html.classList.remove("dark-theme");
      html.classList.add("light-theme");
      button.textContent = "☀️";
      localStorage.setItem("theme", "light");
    }
  };
  const setThemeOnLoad = () => {
    const theme = localStorage.getItem("theme");
    const html = document.querySelector("html");
    const button = document.getElementById("themeButton");

    if (theme === "dark") {
      html.classList.remove("light-theme");
      html.classList.add("dark-theme");
      button.textContent = "🌙";
    } else {
      html.classList.remove("dark-theme");
      html.classList.add("light-theme");
      button.textContent = "☀️";
    }
  };
  const linkToGithub = (repo) => {
    const userTitle = document.querySelector(".user-title").textContent;
    const decodedUserTitle = decodeURIComponent(userTitle);
    const username = decodedUserTitle.trim();

    window.open(`https://github.com/${username}/${repo}`, "_blank");
  };
  document.addEventListener("DOMContentLoaded", setThemeOnLoad);
