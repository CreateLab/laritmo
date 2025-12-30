/** @type {import('tailwindcss').Config} */
export default {
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
            },
        },
    },
    plugins: [],
}