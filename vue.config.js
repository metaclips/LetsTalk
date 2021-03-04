module.exports = {
  "transpileDependencies": [
    "vuetify"
  ],
  "devServer": {
    "https": {}
  },
  "pages": {
    "home": {
      "entry": "./src/Home/main.ts",
      "template": "public/index.html",
      "title": "Home",
      "chunks": [
        "chunk-vendors",
        "chunk-common",
        "home"
      ]
    },
    "admin": {
      "entry": "./src/Admin/main.ts",
      "template": "public/index.html",
      "title": "Admin",
      "chunks": [
        "chunk-vendors",
        "chunk-common",
        "admin"
      ]
    },
    "login": {
      "entry": "./src/Login/main.ts",
      "template": "public/index.html",
      "title": "Login",
      "chunks": [
        "chunk-vendors",
        "chunk-common",
        "login"
      ]
    },
  }
}