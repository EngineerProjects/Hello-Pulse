/** @type {import('tailwindcss').Config} */
module.exports = {
    content: [
      "./app/**/*.{js,ts,jsx,tsx}",
      "./components/**/*.{js,ts,jsx,tsx}",
    ],
    theme: {
      extend: {
        animation: {
          blob: "blob 7s infinite",
        },
        keyframes: {
          blob: {
            "0%": {
              transform: "translate(0px, 0px) scale(1)",
            },
            "33%": {
              transform: "translate(30px, -50px) scale(1.1)",
            },
            "66%": {
              transform: "translate(-20px, 20px) scale(0.9)",
            },
            "100%": {
              transform: "translate(0px, 0px) scale(1)",
            },
          },
        },
        fontFamily: {
          impact: ["Impact", "sans-serif"],
          outfit: ["Outfit", "sans-serif"],
        },
        colors: {
          helloblue: {
            50: "#edf4ff",
            100: "#dfeaff",
            500: "#3b82f6",
            600: "#2563eb",
          },
          hellopurple: {
            500: "#8b5cf6",
            600: "#7c3aed",
          }
        }
      },
    },
    plugins: [
      function({ addUtilities }) {
        const newUtilities = {
          '.animation-delay-2000': {
            'animation-delay': '2s',
          },
          '.animation-delay-4000': {
            'animation-delay': '4s',
          },
        }
        addUtilities(newUtilities)
      }
    ],
    darkMode: 'class',
  }