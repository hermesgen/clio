/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './assets/ssg/**/*.html',
    './assets/ssg/**/*.tmpl',
    './assets/template/**/*.tmpl',
    './assets/static/css/prose.css',
  ],
  theme: {
    extend: {},
  },
  plugins: [
    // Removed require('@tailwindcss/typography'),
  ],
}
