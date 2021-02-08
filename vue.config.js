
module.exports = {
  "transpileDependencies": [
    "vuetify"
  ],
  "devServer": {
    "https": {
      key: process.env.KEY,
      cert: process.env.CERT,
    },
  },
  pages: {
    'home': {
      entry: './src/Home/main.ts',
      template: 'public/index.html',
      title: 'Home',
      chunks: ['chunk-vendors', 'chunk-common', 'home']
    },
    'login': {
      entry: './src/Login/main.ts',
      template: 'public/index.html',
      title: 'Login',
      chunks: ['chunk-vendors', 'chunk-common', 'login']
    },
    'register': {
      entry: './src/Register/main.ts',
      template: 'public/index.html',
      title: 'Register',
      chunks: ['chunk-vendors', 'chunk-common', 'register']
    }
  }
}