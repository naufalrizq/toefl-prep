export type Role = 'student' | 'admin';
export type Section = 'structure' | 'vocabulary' | 'reading' | 'grammar_adv';
export type QuestionType =
  | 'sentence-completion'
  | 'vocab-multiple-choice'
  | 'reading-comprehension'
  | 'error-identification';
export type Difficulty = 'easy' | 'medium' | 'hard';
export type Pos =
  | 'verb'
  | 'noun'
  | 'pronoun'
  | 'adjective'
  | 'adverb'
  | 'preposition'
  | 'conjunction'
  | 'determiner'
  | 'other';
export type ExamMode = 'per_question' | 'overall' | 'both';
export type AttemptStatus = 'in_progress' | 'submitted';

export interface User {
  id: number;
  email: string;
  role: Role;
}

export interface HighlightRegion {
  start: number;
  end: number;
  pos: Pos;
  label?: string;
}

export interface Question {
  id: number;
  section: Section;
  type: QuestionType;
  question_text: string;
  passage?: string;
  options: string[];
  correct_index: number;
  explanation: string;
  highlight_regions: HighlightRegion[];
  difficulty: Difficulty;
  active: boolean;
  created_at: string;
  updated_at: string;
}

export interface Draft {
  section?: Section;
  type?: QuestionType;
  difficulty?: Difficulty;
  question_text: string;
  passage?: string;
  options: string[];
  correct_index: number;
  explanation: string;
  highlights?: Record<string, string[]>;
}

export interface ImportResult {
  index: number;
  valid: boolean;
  error?: string;
  question?: Question;
}

export interface QuestionPage {
  items: Question[];
  page: number;
  limit: number;
  total: number;
}

export interface PartConfig {
  title: string;
  type: QuestionType;
  count: number;
}

export interface SectionConfig {
  parts: PartConfig[];
}

export interface ExamTemplateInput {
  title: string;
  section_filters: Partial<Record<Section, SectionConfig>>;
  shuffle?: boolean;
  mode: ExamMode;
  seconds_per_question?: number;
  total_minutes?: number;
}

export interface ExamTemplate extends ExamTemplateInput {
  id: number;
  published: boolean;
  active: boolean;
  created_at: string;
  updated_at: string;
}

export interface Attempt {
  id: number;
  user_id: number;
  exam_template_id: number;
  exam_title?: string;
  mode: 'per_question' | 'overall';
  started_at: string;
  finished_at?: string;
  status: AttemptStatus;
  score_pct?: number;
  deadline?: string;
  summary?: Record<string, unknown>;
}

export interface QuizSnapshot {
  question_text: string;
  passage?: string;
  options: string[];
  explanation?: string;
  highlight_regions: HighlightRegion[];
  section: Section;
  type: QuestionType;
  part?: string;
}

export interface QuizItem {
  id: number;
  question_snapshot: QuizSnapshot;
  flagged: boolean;
}

export interface StartResult {
  attempt: Attempt;
  items: QuizItem[];
}

export interface SectionReport {
  section: Section;
  correct: number;
  total: number;
  score_pct: number;
}

export interface ScoreReport {
  score_pct: number;
  correct: number;
  wrong: number;
  unanswered: number;
  total: number;
  sections: SectionReport[];
  worst_pos?: Pos;
  items?: Array<{
    id: number;
    is_correct: boolean;
    is_unanswered: boolean;
  }>;
}

export interface ReviewItem {
  id: number;
  section: Section;
  type: QuestionType;
  part?: string;
  question_text: string;
  passage?: string;
  options: string[];
  correct_index: number;
  chosen_index: number | null;
  is_correct: boolean;
  is_unanswered: boolean;
  explanation: string;
  highlight_regions: HighlightRegion[];
  time_taken_ms?: number | null;
  flagged?: boolean;
}

export interface Report {
  attempt: Attempt;
  report: ScoreReport;
  items: ReviewItem[];
}

export type Review = Report;

export interface DashboardStats {
  total_attempts: number;
  average_score: number;
  best_score: number;
  trend: number;
  series: Array<{ id: number; score_pct: number; started_at: string }>;
  sections: Partial<
    Record<
      Section,
      { correct: number; total: number; accuracy: number }
    >
  >;
  worst_pos: Partial<Record<Pos, number>>;
}