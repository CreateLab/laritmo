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
                page: 'var(--page)',
                surface: 'var(--surface)',
                sunken: 'var(--sunken)',
                ink: 'var(--ink)',
                'ink-muted': 'var(--ink-muted)',
                border: 'var(--border)',
                'border-strong': 'var(--border-strong)',
                action: {
                    DEFAULT: 'var(--action)',
                    hover: 'var(--action-hover)',
                    on: 'var(--on-action)',
                    soft: 'var(--action-soft)',
                    'soft-hover': 'var(--action-soft-hover)',
                },
                danger: {
                    soft: 'var(--danger-soft)',
                    'soft-hover': 'var(--danger-soft-hover)',
                    'soft-text': 'var(--danger-soft-text)',
                    bg: 'var(--danger-bg)',
                    text: 'var(--danger-text)',
                },
                accent: 'var(--accent)',
            },
            fontFamily: {
                sans: ['Roboto', 'system-ui', '-apple-system', 'sans-serif'],
                mono: ['"Roboto Mono"', 'ui-monospace', 'Consolas', 'monospace'],
            },
            maxWidth: { page: '2400px' },
            borderRadius: { sm: '6px', md: '8px', lg: '12px' },
            boxShadow: { sm: 'var(--shadow-sm)', md: 'var(--shadow-md)' },
        },
    },
    plugins: [],
}
