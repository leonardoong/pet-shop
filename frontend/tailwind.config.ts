import type { Config } from 'tailwindcss'

export default {
  content: ['./index.html', './src/**/*.{ts,tsx}'],
  theme: {
    extend: {
      colors: {
        primary: {
          50:  '#f4faf0',
          100: '#e5f4da',
          200: '#cce9b8',
          300: '#a8d688',
          400: '#7ebf57',
          500: '#5ea636',
          600: '#498528',
          700: '#3a6820',
          800: '#30531c',
          900: '#284619',
          950: '#12260a',
        },
        warm: {
          50:  '#fdf9f0',
          100: '#faf0dc',
        },
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'],
      },
    },
  },
  plugins: [],
} satisfies Config
