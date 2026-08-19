# -*- coding: utf-8 -*-
import json

Q = [
  ("The professor, together with his research assistants, were presenting the findings.",
   ["The professor, together with", "his research assistants,", "were presenting", "the findings"], 2, "hard",
   "C. Subjek utamanya 'the professor' (tunggal); 'together with his research assistants' cuma sisipan, jadi 'was presenting', bukan 'were'. Jebakannya: 'assistants' jamak di tengah.",
   {"noun": ["professor", "research", "assistants", "findings"], "preposition": ["together with"], "verb": ["presenting"]}),
  ("She avoided to answer the difficult question during the interview.",
   ["She avoided", "to answer the difficult", "question during", "the interview"], 1, "medium",
   "B. Verb 'avoid' diikuti GERUND, bukan infinitive, jadi 'avoided answering', bukan 'avoided to answer'. Beberapa verb punya pola gerund wajib.",
   {"verb": ["avoided", "to answer"], "noun": ["question", "interview"], "adjective": ["difficult"]}),
  ("Despite of the bad weather, the outdoor event proceeded as planned.",
   ["Despite of", "the bad weather,", "the outdoor event", "proceeded as planned"], 0, "medium",
   "A. 'Despite' nggak butuh 'of'. Yang benar 'Despite the bad weather'. Kalau mau pakai 'of', gunakan 'in spite of'.",
   {"preposition": ["Despite"], "adjective": ["bad"], "noun": ["weather", "event"]}),
  ("The manager asked that each employee submits a weekly report.",
   ["The manager asked", "that each employee", "submits a weekly", "report"], 2, "hard",
   "C. 'Asked that' menuntut subjunctive, jadi verb bentuk dasar 'submit', bukan 'submits'. Subjunctive nggak menyesuaikan subjek.",
   {"noun": ["manager", "employee", "weekly", "report"], "verb": ["asked", "submits"]}),
  ("If I was you, I would accept the job offer without hesitation.",
   ["If I was you,", "I would accept", "the job offer", "without hesitation"], 0, "medium",
   "A. Untuk situasi hipotetis gunakan 'were', bukan 'was': 'If I were you'. Ini subjunctive untuk kondisi yang nggak nyata.",
   {"verb": ["was", "would accept"], "noun": ["job", "offer"]}),
  ("One of the main reasons for the delay were technical issues.",
   ["One of the main reasons", "for the delay", "were technical issues", "(no error)"], 2, "medium",
   "C. Subjeknya 'one' (tunggal), jadi 'was technical issues'. 'Of the main reasons' cuma modifier jamak yang jadi jebakan.",
   {"determiner": ["One"], "noun": ["reasons", "delay", "technical", "issues"], "verb": ["were"]}),
  ("The students, as well as the teacher, was surprised by the announcement.",
   ["The students, as well as the teacher,", "was surprised", "by the announcement", "(no error)"], 1, "hard",
   "B. Subjek utamanya 'the students' (jamak), jadi 'were surprised'. 'As well as the teacher' cuma sisipan yang nggak mengubah jumlah subjek.",
   {"noun": ["students", "teacher", "announcement"], "preposition": ["as well as"], "verb": ["surprised"]}),
  ("She has been working on this project since three months.",
   ["She has been working", "on this project", "since three months", "(no error)"], 2, "medium",
   "C. 'Since' dipakai untuk titik waktu (since March), 'for' untuk durasi, jadi 'for three months'. Ini aturan preposisi waktu yang sering keliru.",
   {"preposition": ["since"], "noun": ["project", "months"], "verb": ["working"]}),
  ("The equipment are expensive to maintain.",
   ["The equipment", "are expensive", "to maintain", "(no error)"], 1, "medium",
   "B. 'Equipment' adalah uncountable noun yang selalu tunggal, jadi 'is expensive'. Kata yang terasa jamak begini sering bikin salah.",
   {"noun": ["equipment"], "verb": ["are", "maintain"], "adjective": ["expensive"]}),
  ("No sooner had he finished the exam when he realized he had made a mistake.",
   ["No sooner had he finished the exam", "when he realized", "he had made", "a mistake"], 1, "hard",
   "B. Pasangan tetap 'No sooner ... than'. 'When' nggak pernah dipakai setelah 'no sooner'; harus 'than he realized'.",
   {"adverb": ["No sooner"], "conjunction": ["when"], "noun": ["exam", "mistake"], "verb": ["finished", "realized", "made"]}),
]

out = []
for text, segs, ci, diff, exp, hl in Q:
    out.append({
        "section": "grammar_adv", "type": "error-identification", "difficulty": diff,
        "question_text": text,
        "options": segs, "correct_index": ci,
        "explanation": exp,
        "highlights": hl,
    })

with open("/home/nrizq/Documents/Codes/toefl-prep/generated/grammar_b.json", "w") as f:
    json.dump(out, f, ensure_ascii=False, indent=1)
print("grammar_b:", len(out))