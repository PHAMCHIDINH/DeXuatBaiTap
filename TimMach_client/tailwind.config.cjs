/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './index.html',
    './src/**/*.{ts,tsx}',
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          DEFAULT: '#2563eb',
          foreground: '#ffffff',
        },
        success: '#16a34a',
        warning: '#f59e0b',
        danger: '#dc2626',
      },
      borderRadius: {
        lg: '12px',
      },
    },
  },
  plugins: [],
};
