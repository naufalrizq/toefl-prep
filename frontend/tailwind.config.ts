import type { Config } from 'tailwindcss';

export default {
  content: ['./index.html', './src/**/*.{ts,tsx}'],
  theme: {
    extend: {
      colors: {
        paper: 'var(--paper)',
        card: 'var(--card)',
        'card-muted': 'var(--card-muted)',
        ink: 'var(--ink)',
        'ink-muted': 'var(--ink-muted)',
        'ink-faint': 'var(--ink-faint)',
        border: 'var(--border)',
        'border-strong': 'var(--border-strong)',
        primary: 'var(--primary)',
        'on-primary': 'var(--on-primary)',
        'primary-soft': 'var(--primary-soft)',
        success: 'var(--success)',
        'success-soft': 'var(--success-soft)',
        danger: 'var(--danger)',
        'danger-soft': 'var(--danger-soft)',
        warning: 'var(--warning)',
        'pos-verb': 'var(--pos-verb)',
        'pos-noun': 'var(--pos-noun)',
        'pos-pronoun': 'var(--pos-pronoun)',
        'pos-adjective': 'var(--pos-adjective)',
        'pos-adverb': 'var(--pos-adverb)',
        'pos-preposition': 'var(--pos-preposition)',
        'pos-conjunction': 'var(--pos-conjunction)',
        'pos-determiner': 'var(--pos-determiner)',
        'pos-other': 'var(--pos-other)',
        'on-pos': 'var(--on-pos)',
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'],
        serif: ['Lora', 'Georgia', 'serif'],
      },
      boxShadow: {
        sm: '0 1px 2px rgba(0,0,0,0.05)',
        md: '0 4px 6px rgba(0,0,0,0.1)',
        lg: '0 10px 15px rgba(0,0,0,0.1)',
      },
      borderRadius: {
        sm: '8px',
        md: '12px',
      },
      maxWidth: {
        reading: '720px',
        app: '1200px',
      },
      transitionTimingFunction: {
        DEFAULT: 'ease',
      },
      transitionDuration: {
        DEFAULT: '200ms',
      },
    },
  },
  plugins: [],
} satisfies Config;