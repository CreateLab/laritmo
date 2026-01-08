/** @type {import('tailwindcss').Config} */
export default {
    darkMode: 'class',
    content: [
        "./index.html",
        "./src/**/*.{vue,js,ts,jsx,tsx}",
    ],
    theme: {
        extend: {
            colors: {
                'forest-beige': '#F5F5DC',
                'forest-green': '#8FBC8F',
                'forest-dark': '#2F4F2F',
                'forest-mint': '#98D8C8',
                'forest-green-dark': '#6B9B6B',
                'forest-mint-dark': '#5A8A7A',
                'dark': {
                    'bg': '#1a1a1a',
                    'surface': '#2d2d2d',
                    'border': '#404040',
                    'text': '#e0e0e0',
                    'text-secondary': '#a0a0a0',
                },
            },
        },
    },
    plugins: [],
}