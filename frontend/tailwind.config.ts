import type { Config } from 'tailwindcss'
import typography from '@tailwindcss/typography'

export default {
  content: [
    './app/components/**/*.{vue,js,ts}',
    './app/layouts/**/*.vue',
    './app/pages/**/*.vue',
    './app/composables/**/*.{js,ts}',
    './app/app.vue',
  ],
  darkMode: 'class',
  theme: {
    extend: {
      fontFamily: {
        sans: ['Satoshi', 'system-ui', 'sans-serif'],
        mono: ['Space Mono', 'monospace'],
      },
      colors: {
        // Semantic theme colors (CSS variable based)
        th: {
          bg: 'rgb(var(--color-bg-primary) / <alpha-value>)',
          'bg-s': 'rgb(var(--color-bg-secondary) / <alpha-value>)',
          'bg-t': 'rgb(var(--color-bg-tertiary) / <alpha-value>)',
          'bg-e': 'rgb(var(--color-bg-elevated) / <alpha-value>)',
          text: 'rgb(var(--color-text-primary) / <alpha-value>)',
          'text-s': 'rgb(var(--color-text-secondary) / <alpha-value>)',
          'text-t': 'rgb(var(--color-text-tertiary) / <alpha-value>)',
          'text-m': 'rgb(var(--color-text-muted) / <alpha-value>)',
          card: 'rgb(var(--color-card-bg) / <alpha-value>)',
          input: 'rgb(var(--color-input-bg) / <alpha-value>)',
          nav: 'rgb(var(--color-nav-bg) / <alpha-value>)',
          overlay: 'rgb(var(--color-overlay) / <alpha-value>)',
        },
        // Static brand colors (unchanged between themes)
        surface: {
          900: '#030712',
          800: '#0A0F1E',
          700: '#111827',
          600: '#1E293B',
        },
        accent: {
          DEFAULT: '#0EA5E9',
          light: '#06B6D4',
          glow: '#22D3EE',
        },
        highlight: '#F59E0B',
        success: '#10B981',
        danger: '#F43F5E',
        warning: '#FBBF24',
        info: '#38BDF8',
      },
      animation: {
        'glow-pulse': 'glow-pulse 3s ease-in-out infinite',
        'float': 'float 6s ease-in-out infinite',
        'flow': 'flow 2s linear infinite',
        'fade-in-up': 'fade-in-up 0.6s ease-out forwards',
        'scale-in': 'scale-in 0.5s ease-out forwards',
      },
      keyframes: {
        'glow-pulse': {
          '0%, 100%': { opacity: '0.4', filter: 'blur(40px)' },
          '50%': { opacity: '0.8', filter: 'blur(60px)' },
        },
        'float': {
          '0%, 100%': { transform: 'translateY(0)' },
          '50%': { transform: 'translateY(-10px)' },
        },
        'flow': {
          '0%': { strokeDashoffset: '100' },
          '100%': { strokeDashoffset: '0' },
        },
        'fade-in-up': {
          '0%': { opacity: '0', transform: 'translateY(24px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' },
        },
        'scale-in': {
          '0%': { opacity: '0', transform: 'scale(0.8)' },
          '100%': { opacity: '1', transform: 'scale(1)' },
        },
      },
    },
  },
  plugins: [typography],
} satisfies Config
