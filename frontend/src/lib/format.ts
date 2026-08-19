import type { Pos, QuestionType, Section } from '../types';

export function formatDate(iso: string): string {
  const d = new Date(iso);
  return d.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
  });
}

export function formatDateTime(iso: string): string {
  const d = new Date(iso);
  return d.toLocaleString('en-US', {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  });
}

export function formatTime(iso: string): string {
  const d = new Date(iso);
  return d.toLocaleTimeString('en-US', {
    hour: '2-digit',
    minute: '2-digit',
    hour12: false,
  });
}

export function formatTimeRange(startIso: string, endIso?: string): string {
  const s = new Date(startIso);
  const date = s.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
  const start = formatTime(startIso);
  if (endIso) {
    const end = formatTime(endIso);
    return `${date} · ${start}-${end}`;
  }
  return `${date} · ${start}`;
}

export function formatDuration(ms: number | null | undefined): string {
  if (ms == null) return '-';
  const total = Math.round(ms / 1000);
  const m = Math.floor(total / 60);
  const s = total % 60;
  return `${m}:${String(s).padStart(2, '0')}`;
}

export function formatClock(totalSeconds: number): string {
  const m = Math.floor(totalSeconds / 60);
  const s = totalSeconds % 60;
  return `${m}:${String(s).padStart(2, '0')}`;
}

export function formatScore(scorePct: number): string {
  return `${scorePct}%`;
}

export const SECTION_LABEL: Record<Section, string> = {
  structure: 'Structure',
  vocabulary: 'Vocabulary',
  reading: 'Reading',
  grammar_adv: 'Advanced Grammar',
};

export const TYPE_LABEL: Record<QuestionType, string> = {
  'sentence-completion': 'Sentence Completion',
  'vocab-multiple-choice': 'Vocabulary Multiple Choice',
  'reading-comprehension': 'Reading Comprehension',
  'error-identification': 'Error Identification',
};

export const DEFAULT_TYPE: Record<Section, QuestionType> = {
  structure: 'sentence-completion',
  vocabulary: 'vocab-multiple-choice',
  reading: 'reading-comprehension',
  grammar_adv: 'error-identification',
};

export const ALL_SECTIONS: Section[] = ['structure', 'vocabulary', 'reading', 'grammar_adv'];

export const ALL_TYPES: QuestionType[] = [
  'sentence-completion',
  'vocab-multiple-choice',
  'reading-comprehension',
  'error-identification',
];

export const POS_LABEL: Record<Pos, string> = {
  verb: 'Verb',
  noun: 'Noun',
  pronoun: 'Pronoun',
  adjective: 'Adjective',
  adverb: 'Adverb',
  preposition: 'Preposition',
  conjunction: 'Conjunction',
  determiner: 'Determiner',
  other: 'Other',
};

export const POS_ID: Record<Pos, string> = {
  verb: 'Kata kerja: kata yang menunjukkan tindakan atau keadaan.',
  noun: 'Kata benda: orang, tempat, atau benda.',
  pronoun: 'Kata ganti: menggantikan kata benda.',
  adjective: 'Kata sifat: menerangkan kata benda.',
  adverb: 'Kata keterangan: menerangkan kata kerja, kata sifat, atau kalimat.',
  preposition: 'Kata depan: menunjukkan hubungan posisi, waktu, atau arah.',
  conjunction: 'Kata hubung: menghubungkan kata, frasa, atau klausa.',
  determiner: 'Kata penentu: artikel, kata tunjuk, atau kata jumlah sebelum kata benda.',
  other: 'Lainnya: partikel atau kata yang tidak termasuk kategori di atas.',
};