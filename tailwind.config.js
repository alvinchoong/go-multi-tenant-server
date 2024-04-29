/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./cmd/server/**/*.go"],
  theme: {
    extend: {},
  },
  daisyui: {
    themes: ["light", "dark"],
  },
  plugins: [require("daisyui")],
};
