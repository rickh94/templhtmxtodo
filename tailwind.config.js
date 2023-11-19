module.exports = {
  content: [
    "./**/*.{templ,go,html}",
    "./static/**/*.js",
    "./static/*.js",
    "./static/input.css"
  ],
  theme: {
    extend: {},
  },
  plugins: [
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography'),
  ],
}
